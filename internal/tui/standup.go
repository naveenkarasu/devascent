package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"devascent/internal/ticket"
)

// Step-1 day cycle (#72/#73): [e] ends the day → a standup recap with manager Sam's
// accountability note → [enter] advances to the next day (revealing newly-assigned
// tickets). Sam's note is templated offline; an AI overview can replace standupNote
// when a mentor backend is configured (follow-up).

func (m Model) enterStandup() (tea.Model, tea.Cmd) {
	m.screen = screenStandup
	return m, nil
}

func (m Model) handleStandupKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.screen = screenBoard
	case "enter":
		return m.advanceDay()
	}
	return m, nil
}

// advanceDay ends the day: Day++, re-focus, persist, and the board reveals any
// tickets assigned to the new day. Incomplete tickets carry over.
func (m Model) advanceDay() (tea.Model, tea.Cmd) {
	m.boardSprint.Day++
	m.boardCol = defaultFocusCol(m.boardSprint, m.boardSprint.Day)
	m.boardRow = 0
	m.screen = screenBoard
	m.persist()
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

// standupNote is manager Sam's end-of-day review (templated): praise for what
// shipped, then accountability on overdue and not-yet-started-but-due-soon work.
func (m Model) standupNote() []string {
	sp, day := m.boardSprint, m.boardSprint.Day
	var lines []string

	if d := m.doneToday(); d > 0 {
		lines = append(lines, okStyle.Render(fmt.Sprintf("Nice — %d ticket(s) shipped today. 👏", d)))
	}
	for _, t := range sp.Tickets {
		if t.Visible(day) && t.Overdue(day) {
			lines = append(lines, errStyle.Render(fmt.Sprintf("⚠ %s is overdue (was due D%d) — what's the plan?", t.Key, t.DueDay)))
		}
	}
	for _, t := range sp.Tickets {
		if t.Visible(day) && t.Assignee == "you" && t.NotStarted() && t.Grading != nil &&
			t.DueDay > 0 && !t.Overdue(day) && t.DueDay-day <= 1 {
			lines = append(lines, fmt.Sprintf("%s is due D%d and hasn't been started — can you pick it up next?", t.Key, t.DueDay))
		}
	}
	if len(lines) == 0 {
		lines = append(lines, dimStyle.Render("All on track. See you tomorrow."))
	}
	return lines
}

func (m Model) standupView() string {
	sp := m.boardSprint
	var b strings.Builder
	b.WriteString(titleStyle.Render(fmt.Sprintf("Standup — end of Day %d", sp.Day)) + "\n\n")

	remaining := 0
	for _, st := range []ticket.Status{ticket.ToDo, ticket.InProgress, ticket.InReview} {
		remaining += len(sp.ColumnVisible(st, sp.Day))
	}
	b.WriteString(fmt.Sprintf("  Shipped today   %d\n", m.doneToday()))
	b.WriteString(fmt.Sprintf("  Points done     %d / %d committed\n", sp.DonePoints(), sp.Committed()))
	b.WriteString(fmt.Sprintf("  Still open      %d\n", remaining))
	if od := sp.OverdueCount(sp.Day); od > 0 {
		b.WriteString("  " + errStyle.Render(fmt.Sprintf("Overdue         %d", od)) + "\n")
	}
	arrivingNext := 0
	for _, t := range sp.Tickets {
		if t.AssignedDay == sp.Day+1 {
			arrivingNext++
		}
	}
	if arrivingNext > 0 {
		b.WriteString(dimStyle.Render(fmt.Sprintf("  Arriving Day %d  %d new ticket(s)\n", sp.Day+1, arrivingNext)))
	}

	b.WriteString("\n" + dimStyle.Render("Sam (manager):") + "\n")
	for _, ln := range m.standupNote() {
		b.WriteString("  " + ln + "\n")
	}

	b.WriteString("\n" + dimStyle.Render(fmt.Sprintf("[enter] start Day %d   ·   [esc] back to the board", sp.Day+1)))
	return b.String()
}
