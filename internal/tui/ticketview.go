package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"devascent/internal/grader"
	"devascent/internal/ticket"
)

// Step-1 ticket detail + work loop (#60 view, #61 workflow). A workable ticket
// (assigned @you with a Grading) moves through the columns by real work:
//
//	To Do --[s] start--> In Progress --[e] edit / [r] grade--> (pass) In Review
//	      --[a] answer the reviewer--> Done  (engine's grade-gated Resolve)
//
// Moves mutate the in-memory sprint (shared *Ticket pointers), so the card travels
// columns immediately; save-backed persistence is wired with the career home (#63).

// repoGradeMsg carries the async result of grading a ticket via GradeRepo.
type repoGradeMsg struct{ v grader.Verdict }

// openTicket switches to the detail view for t, seeding the work buffer and
// remembering the board to return to.
func (m Model) openTicket(t *ticket.Ticket) (tea.Model, tea.Cmd) {
	m.detailTicket = t
	m.detailReturn = screenBoard
	m.workCode = startCode(t)
	m.workVerdict = nil
	m.workStatus = ""
	m.reviewQ = ""
	m.inputActive = false
	m.screen = screenTicket
	return m, nil
}

func (m Model) handleTicketKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	t := m.detailTicket
	if t == nil || msg.String() == "esc" {
		m.screen = m.detailReturn
		return m, nil
	}
	// Between days (cooldown / standup-pending) the ticket is read-only.
	if m.boardSprint != nil && m.boardSprint.Phase != ticket.PhaseWorking {
		m.workStatus = dimStyle.Render("Read-only until you start the day — [esc] back.")
		return m, nil
	}
	if msg.String() == "E" { // edit this ticket's fields (assignee, priority, …)
		return m.openEditTicket(t)
	}
	if t.Assignee != "you" {
		return m, nil // someone else's ticket — read-only (but you can still [E]dit it)
	}
	if t.Grading == nil { // ungraded chore (onboarding) — self-attest done
		if msg.String() == "d" && t.Status != ticket.Done {
			_ = t.MarkDone()
			if m.boardSprint != nil {
				t.ResolvedDay = m.boardSprint.Day
			}
			m.workStatus = okStyle.Render("✓ " + t.Key + " marked done.")
			m.persist()
		}
		return m, nil
	}
	switch t.Status {
	case ticket.ToDo:
		if msg.String() == "s" { // start work → In Progress + open the editor
			_ = t.MoveTo(ticket.InProgress)
			m.workCode = startCode(t)
			m.workVerdict = nil
			m.workStatus = "You picked it up. The editor is opening — fix it, save & close, then [r]."
			m.persist()
			return m, editorCmd(m.workCode, m.editorChoice, t.Title, langExt(learnLang(t)))
		}
	case ticket.InProgress:
		switch msg.String() {
		case "e":
			return m, editorCmd(m.workCode, m.editorChoice, t.Title, langExt(learnLang(t)))
		case "r":
			m.workStatus = "Running the hidden tests…"
			return m, m.gradeTicketCmd(t)
		}
	case ticket.InReview:
		if msg.String() == "a" { // open the answer field (modal text input)
			m.inputActive = true
			m.input = ""
		}
	}
	return m, nil
}

// gradeTicketCmd grades the player's work buffer with the ticket-owned command +
// hidden tests via GradeRepo (the trust split is enforced by Grading.Assemble).
func (m Model) gradeTicketCmd(t *ticket.Ticket) tea.Cmd {
	g, code, grd := m.g, m.workCode, t.Grading
	return func() tea.Msg {
		lt, ok := g.(*grader.LocalToolchain)
		if !ok {
			return repoGradeMsg{grader.Verdict{Err: "repo grading needs the local-toolchain grader"}}
		}
		files, cmd := grd.Assemble(map[string]string{workEditPath(grd): code})
		v, err := lt.GradeRepo(grader.RepoRequest{Files: files, Command: cmd})
		if err != nil {
			v = grader.Verdict{Err: err.Error()}
		}
		return repoGradeMsg{v}
	}
}

// applyTicketGrade reacts to a grade: a pass moves In Progress → In Review and the
// reviewer posts a question; a fail keeps the ticket in place.
func (m Model) applyTicketGrade(v grader.Verdict) (tea.Model, tea.Cmd) {
	m.workVerdict = &v
	t := m.detailTicket
	switch {
	case v.Err != "":
		m.workStatus = errStyle.Render("✗ " + oneline(v.Err))
	case v.Passed:
		if t != nil && t.Status == ticket.InProgress {
			_ = t.MoveTo(ticket.InReview)
			m.reviewQ = reviewQuestion(t)
			t.Comments = append(t.Comments, ticket.Comment{Author: "Sam (reviewer)", Body: m.reviewQ})
			m.workStatus = okStyle.Render("✓ Hidden tests pass!") + " Moved to review — [a] answer Sam's question."
			m.persist()
		}
	default:
		m.workStatus = "Some tests failed — [e] edit and [r] run again."
	}
	return m, nil
}

// answerReview records the player's answer and resolves the ticket to Done (the
// engine gates this on the grade having passed, which it has by now).
func (m Model) answerReview(ans string) (tea.Model, tea.Cmd) {
	m.inputActive = false
	m.input = ""
	t := m.detailTicket
	if t == nil {
		return m, nil
	}
	t.Comments = append(t.Comments, ticket.Comment{Author: "you", Body: ans})
	if err := t.Resolve(true); err != nil {
		m.workStatus = errStyle.Render("✗ " + err.Error())
		return m, nil
	}
	if m.boardSprint != nil {
		t.ResolvedDay = m.boardSprint.Day
	}
	m.workStatus = okStyle.Render("✓ Approved — " + t.Key + " is Done.")
	m.persist()
	return m, nil
}

func (m Model) ticketDetailView() string {
	t := m.detailTicket
	if t == nil {
		return boxStyle.Render("No ticket open.\n\n" + dimStyle.Render("[esc] back"))
	}
	w := 76
	if m.width > 12 && m.width-4 < w {
		w = m.width - 4
	}

	var b strings.Builder

	head := titleStyle.Render(t.Key) + "  " + typeTag(t)
	if t.Priority != "" {
		if t.Priority.IsHigh() {
			head += "  " + errStyle.Render(t.Priority.Label())
		} else {
			head += "  " + dimStyle.Render(t.Priority.Label())
		}
	}
	head += dimStyle.Render("   " + t.Status.Label())
	if t.Assignee == "you" {
		head += dimStyle.Render(" · ") + codeStyle.Render("@you")
	}
	b.WriteString(head + "\n")
	b.WriteString(t.Title + "\n\n")

	// meta line: reporter · points · labels · assigned/due (SLA)
	day := 0
	if m.boardSprint != nil {
		day = m.boardSprint.Day
	}
	var meta []string
	if t.Reporter != "" {
		meta = append(meta, "filed by "+t.Reporter)
	}
	if t.Points > 0 {
		meta = append(meta, fmt.Sprintf("%dpt", t.Points))
	}
	if len(t.Labels) > 0 {
		meta = append(meta, strings.Join(t.Labels, ","))
	}
	meta = append(meta, fmt.Sprintf("assigned D%d", t.AssignedDay))
	if t.DueDay > 0 {
		if t.Overdue(day) {
			meta = append(meta, "OVERDUE (was due D"+fmt.Sprint(t.DueDay)+")")
		} else {
			meta = append(meta, fmt.Sprintf("due D%d", t.DueDay))
		}
	}
	metaStr := dimStyle.Render(strings.Join(meta, " · "))
	if t.Overdue(day) {
		metaStr = errStyle.Render(strings.Join(meta, " · ")) // whole line red when overdue
	}
	b.WriteString(metaStr + "\n\n")

	if t.Desc != "" {
		b.WriteString(dimStyle.Render("DESCRIPTION") + "\n")
		b.WriteString(wrap(t.Desc, w) + "\n\n")
	}
	if len(t.Acceptance) > 0 {
		b.WriteString(dimStyle.Render("ACCEPTANCE") + "\n")
		for _, c := range t.Acceptance {
			mark := "  ▢ "
			if c.Met {
				mark = "  " + okStyle.Render("☑") + " "
			}
			b.WriteString(mark + c.Text + "\n")
		}
		b.WriteString("\n")
	}
	if len(t.Subtasks) > 0 {
		b.WriteString(dimStyle.Render("SUBTASKS") + "\n")
		for _, s := range t.Subtasks {
			mark := "  ▢ "
			if s.Done {
				mark = "  " + okStyle.Render("☑") + " "
			}
			b.WriteString(mark + s.Title + "\n")
		}
		b.WriteString("\n")
	}
	if strings.TrimSpace(t.Learn) != "" {
		b.WriteString(dimStyle.Render("LEARN") + "\n")
		b.WriteString(indentLines(highlightCode(t.Learn, learnLang(t)), "  ") + "\n\n")
	}
	if len(t.Attachments) > 0 {
		b.WriteString(dimStyle.Render("ATTACHMENTS") + "\n")
		for _, a := range t.Attachments {
			b.WriteString("  📎 " + a.Name + "\n")
			if a.Body != "" {
				b.WriteString(dimStyle.Render(indentLines(a.Body, "    ")) + "\n")
			}
		}
		b.WriteString("\n")
	}
	if len(t.Links) > 0 {
		var ls []string
		for _, l := range t.Links {
			ls = append(ls, string(l.Kind)+" "+l.Key)
		}
		b.WriteString(dimStyle.Render("LINKS  "+strings.Join(ls, " · ")) + "\n\n")
	}
	if len(t.Watchers) > 0 {
		b.WriteString(dimStyle.Render("watched by "+strings.Join(t.Watchers, ", ")) + "\n\n")
	}
	if len(t.Comments) > 0 {
		b.WriteString(dimStyle.Render("COMMENTS") + "\n")
		for _, c := range t.Comments {
			b.WriteString("  " + okStyle.Render(c.Author) + "\n")
			b.WriteString(wrap(c.Body, w-2) + "\n")
		}
		b.WriteString("\n")
	}

	// ── work section ──
	sepW := w
	if sepW > 46 {
		sepW = 46
	}
	b.WriteString(dimStyle.Render(strings.Repeat("─", sepW)) + "\n")
	if v := renderVerdict(m.workVerdict); v != "" {
		b.WriteString(v + "\n")
	}
	if m.workStatus != "" {
		b.WriteString(m.workStatus + "\n")
	}
	if t.Status == ticket.InReview && m.inputActive {
		b.WriteString("\n  your answer: " + m.input + "▏\n")
	}
	b.WriteString("\n" + dimStyle.Render(ticketFooter(t, m.inputActive)))
	return b.String()
}

func renderVerdict(v *grader.Verdict) string {
	switch {
	case v == nil:
		return ""
	case v.Err != "":
		return errStyle.Render("✗ " + oneline(v.Err))
	case v.Passed:
		return okStyle.Render("✓ hidden tests pass")
	default:
		return errStyle.Render("✗ hidden tests failed")
	}
}

// ticketFooter is the status-appropriate key hints for a workable ticket.
func ticketFooter(t *ticket.Ticket, answering bool) string {
	if t.Assignee != "you" {
		return "[E] edit   ·   [esc] back to the board"
	}
	if t.Grading == nil { // ungraded chore / onboarding
		if t.Status == ticket.Done {
			return "✓ done   ·   [E] edit   ·   [esc] back"
		}
		return "[d] mark done   ·   [E] edit   ·   [esc] back"
	}
	switch t.Status {
	case ticket.ToDo:
		return "[s] start work   ·   [E] edit   ·   [esc] back"
	case ticket.InProgress:
		return "[e] edit code   ·   [r] run tests   ·   [E] edit fields   ·   [esc] back"
	case ticket.InReview:
		if answering {
			return "[enter] submit answer   ·   [esc] cancel"
		}
		return "[a] answer the reviewer   ·   [E] edit   ·   [esc] back"
	case ticket.Done:
		return "✓ done   ·   [E] edit   ·   [esc] back"
	}
	return "[esc] back"
}

// workEditPath is the single file the player edits (first EditPaths, else the
// first StartFile).
func workEditPath(g *ticket.Grading) string {
	if g == nil {
		return "solution"
	}
	if len(g.EditPaths) > 0 {
		return g.EditPaths[0]
	}
	for p := range g.StartFiles {
		return p
	}
	return "solution"
}

// startCode is the player's initial buffer for a ticket (its starting file).
func startCode(t *ticket.Ticket) string {
	if t == nil || t.Grading == nil {
		return ""
	}
	return t.Grading.StartFiles[workEditPath(t.Grading)]
}

func reviewQuestion(t *ticket.Ticket) string {
	return "Before I approve: in one line — what was the root cause, and why doesn't your fix break the last (partial) page?"
}

// learnLang is the language the Learn snippet is highlighted as (and the editor
// extension for the work buffer).
func learnLang(t *ticket.Ticket) string {
	if t.Grading != nil && t.Grading.Lang != "" {
		return t.Grading.Lang
	}
	return "python"
}

// wrap soft-wraps s to width w (terminal-aware, never below a sane floor).
func wrap(s string, w int) string {
	if w < 8 {
		w = 8
	}
	return lipgloss.NewStyle().Width(w).Render(s)
}
