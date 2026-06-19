package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"devascent/internal/mentor"
	"devascent/internal/ticket"
)

// Step-1 day cycle (S3–S6). Ending a day no longer jumps straight to the next:
//
//	[working board] --[e]--> [cooldown: evening recap + live countdown, the team
//	works, board read-only] --(timer ends / [s] skip)--> [new day "join standup"]
//	--> [morning standup: yesterday / today / blockers] --[enter]--> [working board]
//
// The cooldown is wall-clock (Sprint.CooldownEndsAt) so the beat survives quitting
// mid-day, mirroring the economy nudge recharge. The team's delegated work advances
// at the day flip (advanceTeamWork), never on a keypress.

// cooldownTickMsg drives the once-a-second countdown refresh while in cooldown.
type cooldownTickMsg struct{}

func cooldownTick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg { return cooldownTickMsg{} })
}

// ── end of day → cooldown ─────────────────────────────────────────────────────

// enterCooldown ends the working day: the team's day plays out over the cooldown
// beat. Sets the wall-clock deadline and shows the recap + countdown.
func (m Model) enterCooldown() (tea.Model, tea.Cmd) {
	sp := m.boardSprint
	if sp == nil {
		return m, nil
	}
	sp.Phase = ticket.PhaseCooldown
	sp.CooldownEndsAt = time.Now().Add(time.Duration(sp.CooldownDuration()) * time.Second).UTC().Format(time.RFC3339)
	m.screen = screenCooldown
	m.persist()
	return m, cooldownTick()
}

// cooldownRemaining is how long the cooldown beat has left (0 once elapsed). Pure
// (takes now) so the math is testable without a clock.
func cooldownRemaining(sp *ticket.Sprint, now time.Time) time.Duration {
	if sp == nil || sp.CooldownEndsAt == "" {
		return 0
	}
	end, err := time.Parse(time.RFC3339, sp.CooldownEndsAt)
	if err != nil {
		return 0
	}
	if d := end.Sub(now); d > 0 {
		return d
	}
	return 0
}

// onCooldownTick refreshes the countdown; when the beat elapses it flips to the
// next day. Stale ticks (we already left cooldown) stop the loop.
func (m Model) onCooldownTick() (tea.Model, tea.Cmd) {
	if m.boardSprint == nil || m.boardSprint.Phase != ticket.PhaseCooldown {
		return m, nil
	}
	if cooldownRemaining(m.boardSprint, time.Now()) <= 0 {
		return m.finishCooldown()
	}
	return m, cooldownTick()
}

// finishCooldown flips to the next day: the team's delegated work advances and
// the board waits at "join standup". The player stays on whatever screen they're
// on (the cooldown screen shows the morning prompt; a read-only board shows the
// join banner) — they explicitly join the standup to start working.
func (m Model) finishCooldown() (tea.Model, tea.Cmd) {
	m.boardSprint.CooldownEndsAt = ""
	return m.advanceDay()
}

// advanceDay flips to the next day: Day++, the team's delegated work progresses,
// and the sprint waits at the morning standup (Phase=standup). It does NOT change
// the screen — callers decide where the player lands.
func (m Model) advanceDay() (tea.Model, tea.Cmd) {
	m.boardSprint.Day++
	day := m.boardSprint.Day
	advanceTeamWork(m.boardSprint, m.playerLvl, day)
	m.boardSprint.Phase = ticket.PhaseStandup
	m.boardCol = defaultFocusCol(m.boardSprint, day)
	m.boardRow = 0
	m.persist()
	return m, nil
}

// handleCooldownKey: adjust/skip the beat, browse the board read-only, or — once
// the day is ready — join the morning standup.
func (m Model) handleCooldownKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	sp := m.boardSprint
	if sp == nil {
		m.screen = screenBoard
		return m, nil
	}
	switch msg.String() {
	case "enter":
		if sp.Phase == ticket.PhaseStandup {
			return m.enterStandup()
		}
	case "s": // skip the wait
		if sp.Phase == ticket.PhaseCooldown {
			return m.finishCooldown()
		}
	case "+", "=":
		if sp.Phase == ticket.PhaseCooldown {
			m.adjustCooldown(15)
		}
	case "-", "_":
		if sp.Phase == ticket.PhaseCooldown {
			m.adjustCooldown(-15)
		}
	case "b", "esc": // browse the board (read-only) while the day plays out
		m.screen = screenBoard
	}
	return m, nil
}

// adjustCooldown changes the configured cooldown length (remembered for next time)
// and shifts the live deadline by the same amount.
func (m *Model) adjustCooldown(delta int) {
	sp := m.boardSprint
	want := sp.CooldownDuration() + delta
	if want < ticket.CooldownMinSecs {
		want = ticket.CooldownMinSecs
	}
	if want > ticket.CooldownMaxSecs {
		want = ticket.CooldownMaxSecs
	}
	sp.CooldownSecs = want
	rem := int(cooldownRemaining(sp, time.Now()).Seconds()) + delta
	if rem < 1 {
		rem = 1
	}
	sp.CooldownEndsAt = time.Now().Add(time.Duration(rem) * time.Second).UTC().Format(time.RFC3339)
	m.persist()
}

// ── morning standup ───────────────────────────────────────────────────────────

// standupMsg carries the async AI-rendered standup (empty Text → templated).
type standupMsg struct{ resp mentor.Response }

func (m Model) enterStandup() (tea.Model, tea.Cmd) {
	m.screen = screenStandup
	m.standupText = ""
	m.standupBusy = false
	// When an AI backend is connected, the whole standup is rendered in ONE
	// batched call (the team "speaks" in the player's own AI's voice); offline
	// it stays the templated view. Hint() falls back to a template on any failure.
	if tuiMentor().AIEnabled() {
		day := m.boardSprint.Day
		m.standupBusy = true
		return m, standupCmd(day, m.standupFacts(day))
	}
	return m, nil
}

// standupCmd renders the standup off the UI loop from the public status facts —
// never tickets, tests, or solutions.
func standupCmd(day int, facts []string) tea.Cmd {
	req := mentor.Request{Kind: mentor.KindStandup, Persona: "Sam, the Engineering Manager", Day: day, Status: facts}
	return func() tea.Msg {
		return standupMsg{resp: tuiMentor().Hint(context.Background(), req)}
	}
}

// standupFacts builds the plain-text (no ANSI) per-person status the AI renders.
func (m Model) standupFacts(day int) []string {
	sp := m.boardSprint
	facts := []string{"You: " + m.yourStandupLine(day)}
	seen := map[string]bool{}
	for _, t := range sp.Tickets {
		if t.Assignee == "" || t.Assignee == "you" || assignKind(t.Assignee, m.playerLvl) != "delegate" || seen[t.Assignee] {
			continue
		}
		switch {
		case t.Status == ticket.InProgress:
			seen[t.Assignee] = true
			facts = append(facts, fmt.Sprintf("%s: working on %s, on track for D%d, no blockers.", t.Assignee, t.Key, t.DueDay))
		case t.Status == ticket.Done && t.ResolvedDay == day-1:
			seen[t.Assignee] = true
			facts = append(facts, fmt.Sprintf("%s: delivered %s yesterday, free for more.", t.Assignee, t.Key))
		}
	}
	var arrived []string
	for _, t := range sp.Tickets {
		if t.AssignedDay == day && t.Status != ticket.Backlog {
			arrived = append(arrived, t.Key)
		}
	}
	if len(arrived) > 0 {
		facts = append(facts, "New tickets on the board today: "+strings.Join(arrived, ", ")+".")
	}
	for _, t := range sp.Tickets {
		if t.Visible(day) && t.Overdue(day) {
			facts = append(facts, fmt.Sprintf("Manager concern: %s is overdue (was due D%d).", t.Key, t.DueDay))
		}
	}
	return facts
}

func (m Model) handleStandupKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter": // start the working day
		m.boardSprint.Phase = ticket.PhaseWorking
		m.boardCol = defaultFocusCol(m.boardSprint, m.boardSprint.Day)
		m.boardRow = 0
		m.screen = screenBoard
		m.persist()
	case "esc": // peek at the board first (still read-only until you start the day)
		m.screen = screenBoard
	}
	return m, nil
}

// doneToday counts tickets resolved on the current day.
func (m Model) doneToday() int {
	n := 0
	for _, t := range m.boardSprint.Tickets {
		if t.Status == ticket.Done && t.ResolvedDay == m.boardSprint.Day {
			n++
		}
	}
	return n
}

// standupNote is manager Sam's accountability note: praise for recent deliveries,
// then overdue and not-yet-started-but-due-soon work. Templated; an AI overview
// replaces this when a mentor backend is wired (S9).
func (m Model) standupNote() []string {
	sp, day := m.boardSprint, m.boardSprint.Day
	var lines []string
	for _, t := range sp.Tickets {
		if t.Visible(day) && t.Overdue(day) {
			lines = append(lines, errStyle.Render(fmt.Sprintf("⚠ %s is overdue (was due D%d) — what's the plan?", t.Key, t.DueDay)))
		}
	}
	for _, t := range sp.Tickets {
		if t.Visible(day) && t.Assignee == "you" && t.NotStarted() && t.Grading != nil &&
			t.DueDay > 0 && !t.Overdue(day) && t.DueDay-day <= 1 {
			lines = append(lines, fmt.Sprintf("%s is due D%d and hasn't been started — can you pick it up today?", t.Key, t.DueDay))
		}
	}
	if len(lines) == 0 {
		lines = append(lines, dimStyle.Render("Good place to be — let's keep it moving. Have a good day."))
	}
	return lines
}

// ── views ─────────────────────────────────────────────────────────────────────

// cooldownView is the between-days beat: the evening recap (what YOU did), the
// team working in the background, and a live countdown — or, once the beat is
// over, the prompt to join the morning standup.
func (m Model) cooldownView() string {
	sp := m.boardSprint
	if sp == nil {
		return boxStyle.Render("No sprint loaded.\n\n" + dimStyle.Render("[esc] back"))
	}
	recapDay := sp.Day
	if sp.Phase == ticket.PhaseStandup {
		recapDay = sp.Day - 1
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render(fmt.Sprintf("Day %d — wrapping up", recapDay)) + "\n\n")

	b.WriteString(dimStyle.Render("Your day") + "\n")
	yours := m.yourRecap(recapDay)
	if len(yours) == 0 {
		b.WriteString("  " + dimStyle.Render("A quiet one — nothing banked today.") + "\n")
	} else {
		for _, ln := range yours {
			b.WriteString("  " + ln + "\n")
		}
	}
	b.WriteString("\n")

	if sp.Phase == ticket.PhaseStandup {
		b.WriteString(okStyle.Render("☀ A new day. The team has updates.") + "\n\n")
		b.WriteString(dimStyle.Render("[enter] join the standup   ·   [b] peek at the board"))
		return b.String()
	}

	b.WriteString(dimStyle.Render("The rest of the day plays out…") + "\n")
	for _, ln := range m.teamBeat() {
		b.WriteString("  " + ln + "\n")
	}
	b.WriteString("\n")
	rem := cooldownRemaining(sp, time.Now())
	b.WriteString("  " + okStyle.Render(fmt.Sprintf("Day %d starts in %s", sp.Day+1, fmtCountdown(rem))) + "\n")
	b.WriteString("\n" + dimStyle.Render(fmt.Sprintf("[s] skip   ·   [+/-] adjust (%ds)   ·   [b] browse board (read-only)", sp.CooldownDuration())))
	return b.String()
}

// yourRecap lists what YOU shipped on the given day plus your still-open work.
func (m Model) yourRecap(day int) []string {
	sp := m.boardSprint
	var out []string
	for _, t := range sp.Tickets {
		if t.Assignee == "you" && t.Status == ticket.Done && t.ResolvedDay == day {
			out = append(out, okStyle.Render("✓ shipped ")+t.Key+" — "+t.Title)
		}
	}
	for _, t := range sp.Tickets {
		if t.Assignee == "you" && (t.Status == ticket.InProgress || t.Status == ticket.InReview) {
			out = append(out, dimStyle.Render("• "+t.Key+" ("+strings.ToLower(t.Status.Label())+")"))
		}
	}
	return out
}

// teamBeat narrates the delegated work in flight during the cooldown.
func (m Model) teamBeat() []string {
	sp := m.boardSprint
	var out []string
	for _, t := range sp.Tickets {
		if t.Assignee == "" || t.Assignee == "you" || assignKind(t.Assignee, m.playerLvl) != "delegate" {
			continue
		}
		if t.Status == ticket.InProgress {
			out = append(out, fmt.Sprintf("%s is heads-down on %s.", t.Assignee, t.Key))
		}
	}
	if len(out) == 0 {
		out = append(out, dimStyle.Render("The office is quiet — no delegated work in flight."))
	}
	return out
}

func fmtCountdown(d time.Duration) string {
	s := int(d.Seconds() + 0.5)
	if s < 0 {
		s = 0
	}
	return fmt.Sprintf("%d:%02d", s/60, s%60)
}

// standupView is the morning standup: per-person yesterday/today/blockers, the
// work that arrived today, and Sam's focus for the day.
func (m Model) standupView() string {
	sp := m.boardSprint
	if sp == nil {
		return boxStyle.Render("No sprint loaded.\n\n" + dimStyle.Render("[esc] back"))
	}
	day := sp.Day
	var b strings.Builder
	b.WriteString(titleStyle.Render(fmt.Sprintf("☀ Standup — Day %d", day)) + "\n")
	b.WriteString(dimStyle.Render("Pixel Forge daily — yesterday · today · blockers") + "\n\n")

	switch {
	case m.standupBusy:
		b.WriteString(dimStyle.Render("…the team is gathering") + "\n")
	case m.standupText != "": // AI-rendered standup
		b.WriteString(m.standupText + "\n")
	default: // templated standup
		for _, line := range m.standupLines(day) {
			b.WriteString("  " + line + "\n")
		}
		var arrived []string
		for _, t := range sp.Tickets {
			if t.AssignedDay == day && t.Status != ticket.Backlog {
				arrived = append(arrived, t.Key)
			}
		}
		if len(arrived) > 0 {
			b.WriteString("\n" + okStyle.Render("New on the board today: ") + strings.Join(arrived, ", ") + "\n")
		}
		b.WriteString("\n" + dimStyle.Render("Sam (manager):") + "\n")
		for _, ln := range m.standupNote() {
			b.WriteString("  " + ln + "\n")
		}
	}

	b.WriteString("\n" + dimStyle.Render("[enter] start the day   ·   [esc] peek at the board"))
	return b.String()
}

// standupLines is one short yesterday/today/blockers line for you and for each
// teammate with active or just-delivered delegated work.
func (m Model) standupLines(day int) []string {
	sp := m.boardSprint
	out := []string{okStyle.Render("You") + " — " + m.yourStandupLine(day)}
	seen := map[string]bool{}
	for _, t := range sp.Tickets {
		if t.Assignee == "" || t.Assignee == "you" || assignKind(t.Assignee, m.playerLvl) != "delegate" || seen[t.Assignee] {
			continue
		}
		switch {
		case t.Status == ticket.InProgress:
			seen[t.Assignee] = true
			out = append(out, codeStyle.Render(t.Assignee)+fmt.Sprintf(" — on %s, on track for D%d. No blockers.", t.Key, t.DueDay))
		case t.Status == ticket.Done && t.ResolvedDay == day-1:
			seen[t.Assignee] = true
			out = append(out, codeStyle.Render(t.Assignee)+fmt.Sprintf(" — delivered %s yesterday. Free for more.", t.Key))
		}
	}
	return out
}

// yourStandupLine summarises your own committed work for the day.
func (m Model) yourStandupLine(day int) string {
	sp := m.boardSprint
	var inProg, todo []string
	for _, t := range sp.Tickets {
		if t.Assignee != "you" || !t.Visible(day) {
			continue
		}
		switch t.Status {
		case ticket.InProgress, ticket.InReview:
			inProg = append(inProg, t.Key)
		case ticket.ToDo:
			todo = append(todo, t.Key)
		}
	}
	switch {
	case len(inProg) > 0:
		return "continuing " + strings.Join(inProg, ", ") + "."
	case len(todo) > 0:
		return "picking up " + todo[0] + " today."
	default:
		return "clear board — grab something new."
	}
}
