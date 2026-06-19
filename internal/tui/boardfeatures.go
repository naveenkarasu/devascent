package tui

import (
	"fmt"
	"strings"

	"devascent/internal/ticket"
)

// Wave B — Jira-parity board features: quick filters (B4), WIP limits (B3),
// swimlane grouping + epics (B2/B5), a backlog view (B1), and sprint analytics
// (B6). Toggled with f / g / b / v on the board.

// ── B4: quick filters ──
var boardFilters = []struct {
	name string
	pred func(t *ticket.Ticket, day int) bool
}{
	{"all", func(t *ticket.Ticket, day int) bool { return true }},
	{"mine", func(t *ticket.Ticket, day int) bool { return t.Assignee == "you" }},
	{"overdue", func(t *ticket.Ticket, day int) bool { return t.Overdue(day) }},
	{"high-pri", func(t *ticket.Ticket, day int) bool { return t.Priority.IsHigh() }},
}

// columnCards is the visible tickets in a status after the active quick-filter —
// the single source the board, nav, and selection all read.
func (m Model) columnCards(col ticket.Status) []*ticket.Ticket {
	day := m.boardSprint.Day
	pred := boardFilters[m.boardFilter].pred
	var out []*ticket.Ticket
	for _, t := range m.boardSprint.ColumnVisible(col, day) {
		if pred(t, day) {
			out = append(out, t)
		}
	}
	return out
}

// ── B3: WIP limits ──
func wipLimit(s ticket.Status) int {
	switch s {
	case ticket.InProgress:
		return 3
	case ticket.InReview:
		return 2
	}
	return 0 // no limit
}

// ── B2/B5: swimlane groupings ──
var boardGroupings = []string{"none", "epic", "priority", "assignee"}

func groupKey(t *ticket.Ticket, mode string) string {
	switch mode {
	case "epic":
		if t.Epic == "" {
			return "(no epic)"
		}
		return t.Epic
	case "priority":
		if t.Priority == "" {
			return "(none)"
		}
		return t.Priority.Label()
	case "assignee":
		if t.Assignee == "" {
			return "(unassigned)"
		}
		return t.Assignee
	}
	return ""
}

// boardGroupedView is the swimlane overview: visible+filtered tickets grouped by
// epic / priority / assignee, each group a labelled lane (epics show progress).
func (m Model) boardGroupedView() string {
	mode := boardGroupings[m.boardGroup]
	day := m.boardSprint.Day
	groups := map[string][]*ticket.Ticket{}
	var order []string
	for _, t := range m.boardSprint.Tickets {
		if !t.Visible(day) || !boardFilters[m.boardFilter].pred(t, day) {
			continue
		}
		k := groupKey(t, mode)
		if _, ok := groups[k]; !ok {
			order = append(order, k)
		}
		groups[k] = append(groups[k], t)
	}

	var b strings.Builder
	b.WriteString(m.boardHeader() + "\n\n")
	for _, g := range order {
		hdr := headFocus.Render("▸ " + g)
		if mode == "epic" {
			hdr += "  " + epicProgress(groups[g])
		}
		b.WriteString(hdr + "\n")
		for _, t := range groups[g] {
			b.WriteString("   " + statusDot(t.Status) + " " + codeStyle.Render(t.Key) + " " +
				dimStyle.Render(fmt.Sprintf("%-11s", t.Status.Label())) + " " + trunc(t.Title, 40) + "\n")
		}
		b.WriteString("\n")
	}
	b.WriteString(dimStyle.Render(fmt.Sprintf("grouped by %s · [g] regroup · [f] filter (%s) · [esc] board",
		mode, boardFilters[m.boardFilter].name)))
	return b.String()
}

func statusDot(s ticket.Status) string {
	switch s {
	case ticket.Done:
		return okStyle.Render("●")
	case ticket.InProgress, ticket.InReview:
		return headFocus.Render("●")
	default:
		return dimStyle.Render("○")
	}
}

func epicProgress(ts []*ticket.Ticket) string {
	done := 0
	for _, t := range ts {
		if t.Status == ticket.Done {
			done++
		}
	}
	return dimStyle.Render(fmt.Sprintf("(%d/%d done)", done, len(ts)))
}

// ── B1: backlog ──
func (m Model) boardBacklogView() string {
	var bl []*ticket.Ticket
	for _, t := range m.boardSprint.Tickets {
		if t.Status == ticket.Backlog {
			bl = append(bl, t)
		}
	}
	var b strings.Builder
	b.WriteString(titleStyle.Render("Backlog") + dimStyle.Render(fmt.Sprintf("   %d item(s) not yet pulled into the sprint", len(bl))) + "\n\n")
	if len(bl) == 0 {
		b.WriteString(dimStyle.Render("  (empty)\n"))
	}
	for _, t := range bl {
		b.WriteString("   " + codeStyle.Render(t.Key) + " " + typeTag(t) + "  " + trunc(t.Title, 44) + "\n")
	}
	b.WriteString("\n" + dimStyle.Render("[b] back to the board · [esc] board"))
	return b.String()
}

// ── B6: analytics ──
func (m Model) boardAnalyticsView() string {
	sp := m.boardSprint
	committed, done := sp.Committed(), sp.DonePoints()
	doneCount := 0
	for _, t := range sp.Tickets {
		if t.Status == ticket.Done {
			doneCount++
		}
	}
	var b strings.Builder
	b.WriteString(titleStyle.Render(fmt.Sprintf("Sprint %d — analytics (Day %d)", sp.Number, sp.Day)) + "\n\n")
	b.WriteString(fmt.Sprintf("  Committed    %d pts\n", committed))
	b.WriteString(fmt.Sprintf("  Done         %d pts\n", done))
	b.WriteString(fmt.Sprintf("  Remaining    %d pts\n", committed-done))
	b.WriteString(fmt.Sprintf("  Throughput   %d tickets done\n", doneCount))
	if od := sp.OverdueCount(sp.Day); od > 0 {
		b.WriteString("  " + errStyle.Render(fmt.Sprintf("Overdue      %d", od)) + "\n")
	}
	if sp.Day > 0 {
		b.WriteString(fmt.Sprintf("  Velocity     %.1f pts/day\n", float64(done)/float64(sp.Day)))
	}
	b.WriteString("\n  Burndown\n  " + burndownBar(committed, done) + "\n")
	b.WriteString("\n" + dimStyle.Render("[v] or [esc] close"))
	return b.String()
}

func burndownBar(committed, done int) string {
	if committed <= 0 {
		return dimStyle.Render("—")
	}
	const width = 24
	filled := done * width / committed
	if filled > width {
		filled = width
	}
	return okStyle.Render(strings.Repeat("█", filled)) +
		dimStyle.Render(strings.Repeat("░", width-filled)) +
		fmt.Sprintf("  %d/%d pts", done, committed)
}
