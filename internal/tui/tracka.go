package tui

// Track A surfaces: the hint picker (A2), the write-up gate (A1), the
// graduation-gate view (A3), and the mentor picker (A4). State lives on Model
// (app.go); these are the handlers and renderers.

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"devascent/internal/economy"
	"devascent/internal/engine"
	"devascent/internal/mentor"
	"devascent/internal/save"
)

var (
	mentorOnce sync.Once
	mentorSvc  *mentor.Service
)

// tuiMentor lazily builds the shared mentor service from mentor.json.
func tuiMentor() *mentor.Service {
	mentorOnce.Do(func() {
		cfg, _ := mentor.LoadConfig()
		mentorSvc = mentor.NewService(cfg)
	})
	return mentorSvc
}

// ── shared helpers ───────────────────────────────────────────────────────────

// ensureWallet grants the starting stash on first bench contact and accrues
// nudge recharges.
func (m *Model) ensureWallet() {
	if !m.wallet.Init {
		st := save.State{}
		m.wallet = economy.Load(&st, time.Now())
	} else {
		m.wallet.Recharge(time.Now())
	}
	if m.solveRecords == nil {
		m.solveRecords = map[string]save.SolveRecord{}
	}
	if m.nudgeUsed == nil {
		m.nudgeUsed = map[string]int{}
	}
}

// recordFail accumulates pity-rule bookkeeping for an unsolved problem.
func (m *Model) recordFail(id string) {
	m.ensureWallet()
	rec := m.solveRecords[id]
	rec.FailedRuns++
	if rec.FirstTryAt == "" {
		rec.FirstTryAt = time.Now().UTC().Format(time.RFC3339)
	}
	m.solveRecords[id] = rec
	m.persist()
}

func (m Model) walletLine() string {
	w := m.wallet
	w.Recharge(time.Now())
	line := fmt.Sprintf("⬡ %d tokens · ◉ %d/%d nudges", w.Tokens, w.NudgeCharges, economy.NudgeMax)
	if next := w.NextRecharge(time.Now()); next > 0 {
		line += fmt.Sprintf(" (next %dm)", int(next.Minutes())+1)
	}
	return line
}

// provisionalIDs lists solved problems whose write-up is still pending.
func (m Model) provisionalIDs() []string {
	var out []string
	for _, p := range m.cat.Problems { // catalog order keeps the queue stable
		if m.solvedSet[p.ID] && !m.solveRecords[p.ID].WriteupDone {
			out = append(out, p.ID)
		}
	}
	return out
}

func (m Model) gateProgress() engine.GateProgress {
	recs := m.solveRecords
	return engine.Blind75Progress(m.cat.Problems, m.solvedSet, func(id string) bool {
		return recs[id].WriteupDone
	})
}

// payMilestones awards MilestoneAward per gate category that just crossed its
// minimum (idempotent via m.milestones).
func (m *Model) payMilestones(before, after engine.GateProgress) int {
	paid := map[string]bool{}
	for _, c := range m.milestones {
		paid[c] = true
	}
	was := map[string]bool{}
	for _, c := range before.Categories {
		was[c.Category] = c.Done >= c.Required
	}
	total := 0
	for _, c := range after.Categories {
		if c.Done >= c.Required && !was[c.Category] && !paid[c.Category] {
			total += economy.MilestoneAward
			m.milestones = append(m.milestones, c.Category)
		}
	}
	return total
}

func (m Model) mentorLabel() string {
	if cfg := tuiMentor().Config(); cfg.Backend != "" {
		return cfg.Backend
	}
	return "built-in playbook"
}

// ── A2: hint picker ──────────────────────────────────────────────────────────

func (m Model) handleHintKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.hintBusy {
		return m, nil // a mentor call is in flight
	}
	m.ensureWallet()
	id := m.curProblem.ID
	rec := m.solveRecords[id]
	switch msg.String() {
	case "esc", "h":
		m.hintMode = false
		m.hintArm = 0
		return m, nil
	case "1":
		m.hintArm = 0
		if !m.wallet.SpendNudge(time.Now()) {
			next := m.wallet.NextRecharge(time.Now())
			m.hintNote = fmt.Sprintf("no nudges left — next recharges in %dm", int(next.Minutes())+1)
			return m, nil
		}
		m.hintText = mentor.Nudge(m.curProblem.Category, m.nudgeUsed[id])
		m.nudgeUsed[id]++
		m.hintNote = "nudge (free)"
		m.persist()
		return m, nil
	case "2", "3":
		tier := economy.TierStrategy
		if msg.String() == "3" {
			tier = economy.TierWalkthrough
		}
		cost := economy.HintCost(tier)
		pity := false
		switch {
		case rec.HintTier >= tier:
			cost = 0 // already paid on this problem; re-showing is free
		case tier == economy.TierStrategy && economy.PityEligible(rec, time.Now()):
			cost, pity = 0, true
		default:
			if m.wallet.Tokens < cost {
				m.hintNote = fmt.Sprintf("not enough tokens (need %d) — bank a clean solve or finish a write-up", cost)
				return m, nil
			}
			if m.hintArm != tier { // two-step confirm before spending
				m.hintArm = tier
				m.hintNote = fmt.Sprintf("press [%s] again to spend %d token(s) — this solve will count less toward mastery", msg.String(), cost)
				return m, nil
			}
		}
		m.hintArm = 0
		if pity {
			rec.PityUsed = true
		} else if cost > 0 {
			m.wallet.Spend(cost)
		}
		if tier > rec.HintTier {
			rec.HintTier = tier
		}
		m.solveRecords[id] = rec
		m.persist() // debit committed BEFORE the mentor call
		req := mentor.Request{
			Kind: mentor.KindStrategy, Lang: m.lang, Title: m.curProblem.Title,
			Prompt: m.curProblem.Prompt, Category: m.curProblem.Category,
			Difficulty: m.curProblem.Difficulty, PlayerCode: m.task.code,
			FailedRuns: rec.FailedRuns,
		}
		if tier == economy.TierWalkthrough {
			req.Kind = mentor.KindWalkthrough
		}
		if pity {
			m.hintNote = "free hint — you've earned it for persistence"
		}
		if tuiMentor().AIEnabled() {
			m.hintBusy = true
			m.hintText = ""
			return m, hintCmd(req, tier, cost)
		}
		resp := tuiMentor().Hint(context.Background(), req)
		m.hintText = resp.Text
		if m.hintNote == "" || !pity {
			m.hintNote = "from the playbook"
		}
		return m, nil
	}
	return m, nil
}

// hintCmd runs the mentor call off the UI loop (up to 45s).
func hintCmd(req mentor.Request, tier, cost int) tea.Cmd {
	return func() tea.Msg {
		resp := tuiMentor().Hint(context.Background(), req)
		return hintMsg{resp: resp, tier: tier, cost: cost}
	}
}

func (m Model) renderHintPanel() string {
	if !m.hintMode && m.hintText == "" {
		return ""
	}
	var b strings.Builder
	b.WriteString("\n")
	if m.hintMode {
		b.WriteString(titleStyle.Render("Hints — "+m.walletLine()) + "\n")
		b.WriteString("  [1] Nudge — free (uses a nudge charge)\n")
		b.WriteString(fmt.Sprintf("  [2] Strategy — %d ⬡ (mastery ×%.1f)\n", economy.StrategyCost, economy.MasteryWeight(economy.TierStrategy)))
		b.WriteString(fmt.Sprintf("  [3] Walkthrough — %d ⬡ (mastery ×%.1f)\n", economy.WalkthroughCost, economy.MasteryWeight(economy.TierWalkthrough)))
		b.WriteString(dimStyle.Render("  [esc] close") + "\n")
	}
	if m.hintBusy {
		b.WriteString(dimStyle.Render("…the mentor is thinking") + "\n")
	}
	if m.hintNote != "" {
		b.WriteString(dimStyle.Render(m.hintNote) + "\n")
	}
	if m.hintText != "" {
		b.WriteString(okStyle.Render("Mentor: ") + m.hintText + "\n")
	}
	return b.String()
}

// ── A1: write-up gate ────────────────────────────────────────────────────────

// openWriteups queues problems for the write-up screen. fromBench resumes the
// bench flow when the queue empties; otherwise it returns to the browse menu.
func (m Model) openWriteups(ids []string, fromBench bool) (tea.Model, tea.Cmd) {
	if len(ids) == 0 {
		return m.startBench()
	}
	m.wuQueue = ids
	m.wuQIdx = 0
	m.wuFromBench = fromBench
	return m.enterWriteup()
}

func (m Model) enterWriteup() (tea.Model, tea.Cmd) {
	id := m.wuQueue[m.wuQIdx]
	p := m.probByID[id]
	m.wuMCQ, m.wuHasMCQ = engine.ComplexityMCQ(p)
	m.wuSel, m.wuPhase, m.wuErr = 0, 0, ""
	m.input = ""
	m.inputActive = false
	if !m.wuHasMCQ { // no authored complexity → straight to the note
		m.wuPhase = 1
		m.inputActive = true
	}
	m.screen = screenWriteup
	m.ctx = ctxNone
	m.task = nil
	return m, nil
}

func (m Model) handleWriteupKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "s": // keep it provisional; the bench-menu entry brings you back
		return m.nextWriteup(false)
	case "up", "k":
		if m.wuPhase == 0 && m.wuSel > 0 {
			m.wuSel--
		}
	case "down", "j":
		if m.wuPhase == 0 && m.wuSel < len(m.wuMCQ.Options)-1 {
			m.wuSel++
		}
	case "enter":
		if m.wuPhase != 0 {
			return m, nil
		}
		if m.wuSel == m.wuMCQ.Correct {
			m.wuErr = ""
			m.wuPhase = 1
			m.input = ""
			m.inputActive = true
		} else {
			m.wuErr = "Not quite — think about what the solution actually does on a big input."
		}
	}
	return m, nil
}

// submitWriteupText accepts the approach note (called from submitInput).
func (m Model) submitWriteupText(text string) (tea.Model, tea.Cmd) {
	if !engine.WriteupTextOK(text) {
		m.wuErr = fmt.Sprintf("A bit more — describe your approach in at least %d characters.", engine.MinWriteupLen)
		return m, nil
	}
	m.inputActive = false
	m.ensureWallet()
	id := m.wuQueue[m.wuQIdx]
	rec := m.solveRecords[id]
	before := m.gateProgress()
	rec.WriteupDone = true
	rec.MCQCorrect = true
	rec.WriteupText = text
	m.solveRecords[id] = rec
	after := m.gateProgress()
	award := economy.WriteupAward + m.payMilestones(before, after)
	m.wallet.Award(award)
	m.wuAward += award
	m.persist()
	return m.nextWriteup(true)
}

// nextWriteup advances the queue; done resumes the bench or the menu.
func (m Model) nextWriteup(accepted bool) (tea.Model, tea.Cmd) {
	m.inputActive = false
	if accepted {
		m.status = fmt.Sprintf("✓ Banked in full (+%d ⬡).", m.wuAward)
	} else {
		m.status = "Kept provisional — find it under “Write-ups pending” when ready."
	}
	m.wuAward = 0
	m.wuQIdx++
	if m.wuQIdx < len(m.wuQueue) {
		return m.enterWriteup()
	}
	m.wuQueue = nil
	if m.wuFromBench {
		return m.benchContinue()
	}
	return m.startBench()
}

func (m Model) renderWriteup() string {
	id := m.wuQueue[m.wuQIdx]
	p := m.probByID[id]
	var b strings.Builder
	b.WriteString(titleStyle.Render("Explain it — "+p.Title) + "\n")
	queue := ""
	if len(m.wuQueue) > 1 {
		queue = fmt.Sprintf("  ·  %d of %d", m.wuQIdx+1, len(m.wuQueue))
	}
	b.WriteString(dimStyle.Render("Tests passed — explaining banks it in full (+1 ⬡)."+queue) + "\n\n")
	if m.wuAward > 0 && m.wuPhase == 0 {
		b.WriteString(okStyle.Render(fmt.Sprintf("Clean solve: +%d ⬡", m.wuAward)) + "\n\n")
	}
	if m.wuPhase == 0 && m.wuHasMCQ {
		b.WriteString(m.wuMCQ.Question + "\n\n")
		for i, o := range m.wuMCQ.Options {
			cursor := "  "
			line := o
			if i == m.wuSel {
				cursor = "› "
				line = okStyle.Render(line)
			}
			b.WriteString(cursor + line + "\n")
		}
		if m.wuErr != "" {
			b.WriteString("\n" + errStyle.Render(m.wuErr) + "\n")
		}
		b.WriteString("\n" + dimStyle.Render("[↑/↓] choose   ·   [enter] answer   ·   [s] later (stays provisional)"))
	} else {
		b.WriteString("In a sentence or two: how does your solution work?\n\n")
		b.WriteString("> " + m.input + "▌\n")
		if m.wuErr != "" {
			b.WriteString("\n" + errStyle.Render(m.wuErr) + "\n")
		}
		b.WriteString("\n" + dimStyle.Render("[enter] submit   ·   [esc] back"))
	}
	return b.String()
}

// ── A3: graduation gate view ─────────────────────────────────────────────────

func (m Model) renderGate() string {
	g := m.gateProgress()
	var b strings.Builder
	b.WriteString(titleStyle.Render("Graduation gate — Blind 75") + "\n")
	b.WriteString(dimStyle.Render("A problem counts once it's solved AND explained (write-up done).") + "\n\n")
	head := fmt.Sprintf("Fully banked: %d / %d", g.Full, g.Target)
	if g.Provisional > 0 {
		head += fmt.Sprintf("   (+%d solved, write-up pending)", g.Provisional)
	}
	if g.Met {
		b.WriteString(okStyle.Render(head+"   —   GATE MET. Apprenticeship complete.") + "\n\n")
	} else {
		b.WriteString(head + "\n\n")
	}
	for _, c := range g.Categories {
		mark := dimStyle
		if c.Done >= c.Required {
			mark = okStyle
		}
		b.WriteString(fmt.Sprintf("  %s %-32s %d/%d (of %d)\n",
			mark.Render(pickMark(c.Done >= c.Required)), c.Category, c.Done, c.Required, c.Available))
	}
	b.WriteString("\nMandatory classics:\n")
	for _, it := range g.Mandatory {
		title := it.Title
		if title == "" {
			title = it.Slug
		}
		if it.Done {
			b.WriteString("  " + okStyle.Render("✓ "+title) + "\n")
		} else {
			b.WriteString("  ○ " + title + "\n")
		}
	}
	b.WriteString("\n" + dimStyle.Render("[esc] back to the bench menu"))
	return b.String()
}

func pickMark(ok bool) string {
	if ok {
		return "✓"
	}
	return "·"
}

// ── A4: mentor picker ────────────────────────────────────────────────────────

func (m Model) openMentorPicker() (tea.Model, tea.Cmd) {
	m.mentorRows = tuiMentor().Statuses()
	m.mentorIdx = 0
	m.mentorNote = ""
	m.screen = screenMentor
	return m, nil
}

func (m Model) handleMentorKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.mentorBusy {
		return m, nil
	}
	switch msg.String() {
	case "esc", "m":
		return m.startBench()
	case "up", "k":
		if m.mentorIdx > 0 {
			m.mentorIdx--
		}
	case "down", "j":
		if m.mentorIdx < len(m.mentorRows)-1 {
			m.mentorIdx++
		}
	case "enter":
		row := m.mentorRows[m.mentorIdx]
		if !row.Present {
			m.mentorNote = "✗ not detected — install it (or start it) first"
			return m, nil
		}
		m.mentorBusy = true
		m.mentorNote = "probing " + row.Name + "…"
		id := row.ID
		return m, func() tea.Msg {
			return mentorSelectMsg{id: id, err: tuiMentor().Select(context.Background(), id)}
		}
	}
	return m, nil
}

func (m Model) renderMentor() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("AI mentor — bring your own") + "\n")
	b.WriteString(dimStyle.Render("Paid hints can be answered by an AI you already have. The game ships none,\nstores no keys, and always falls back to the built-in playbook offline.") + "\n\n")
	for i, r := range m.mentorRows {
		cursor := "  "
		mark := "·"
		if r.Present {
			mark = "✓"
		}
		line := fmt.Sprintf("%s %-42s %s", mark, r.Name, r.Info)
		if r.Selected {
			line += "   ← current"
		}
		if i == m.mentorIdx {
			cursor = "› "
			line = okStyle.Render(line)
		}
		b.WriteString(cursor + line + "\n")
	}
	b.WriteString("\n")
	if m.mentorBusy {
		b.WriteString(dimStyle.Render("…probing (sends a one-line test prompt)") + "\n")
	}
	if m.mentorNote != "" {
		b.WriteString(m.mentorNote + "\n")
	}
	b.WriteString(dimStyle.Render("[↑/↓] choose   ·   [enter] probe & use   ·   [esc] back"))
	return b.String()
}
