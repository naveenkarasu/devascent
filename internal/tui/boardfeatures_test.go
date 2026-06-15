package tui

import (
	"strings"
	"testing"

	"devascent/internal/ticket"
)

// B4: [f] cycles the quick filter; the columns then show only matching cards.
func TestBoard_QuickFilter(t *testing.T) {
	m := step1Model()
	m.boardSprint.Day = 5 // reveal everything; PXF-201 (due D2) is now overdue
	all := len(m.columnCards(ticket.ToDo))
	// cycle to "mine" — every seed ticket is @you, so the count is unchanged
	nm, _ := m.handleBoardKey(mkKey("f"))
	m = nm.(Model)
	if boardFilters[m.boardFilter].name != "mine" {
		t.Fatalf("f should select 'mine', got %q", boardFilters[m.boardFilter].name)
	}
	if len(m.columnCards(ticket.ToDo)) != all {
		t.Errorf("all seed tickets are @you, so 'mine' shouldn't change the count")
	}
	// cycle to "overdue" — fewer cards
	nm, _ = m.handleBoardKey(mkKey("f"))
	m = nm.(Model)
	if boardFilters[m.boardFilter].name != "overdue" {
		t.Fatalf("expected 'overdue', got %q", boardFilters[m.boardFilter].name)
	}
	for _, c := range m.columnCards(ticket.ToDo) {
		if !c.Overdue(m.boardSprint.Day) {
			t.Errorf("overdue filter leaked a non-overdue card: %s", c.Key)
		}
	}
}

// B3: WIP limits show in the column header and warn when exceeded.
func TestBoard_WIPLimit(t *testing.T) {
	if wipLimit(ticket.InProgress) != 3 || wipLimit(ticket.InReview) != 2 || wipLimit(ticket.ToDo) != 0 {
		t.Fatal("unexpected WIP limits")
	}
	m := step1Model()
	m.boardSprint.Day = 5
	// pile 4 tickets into In Progress (limit 3) → header should warn
	n := 0
	for _, tk := range m.boardSprint.Tickets {
		if tk.Grading != nil && n < 4 {
			tk.Status = ticket.InProgress
			n++
		}
	}
	out := m.boardView()
	if !strings.Contains(out, "In Progress (4/3)") {
		t.Errorf("WIP badge should show 4/3, got:\n%s", out)
	}
}

// B2/B5: [g] cycles grouping; the grouped overview renders lanes (epic shows progress).
func TestBoard_GroupingOverview(t *testing.T) {
	m := step1Model()
	m.boardSprint.Day = 5
	nm, _ := m.handleBoardKey(mkKey("g"))
	m = nm.(Model)
	if boardGroupings[m.boardGroup] != "epic" {
		t.Fatalf("g should group by epic, got %q", boardGroupings[m.boardGroup])
	}
	out := m.boardView()
	if !strings.Contains(out, "grouped by epic") {
		t.Errorf("grouped view should render, got:\n%s", out)
	}
	// cycle to priority
	nm, _ = m.handleBoardKey(mkKey("g"))
	if boardGroupings[nm.(Model).boardGroup] != "priority" {
		t.Errorf("expected priority grouping")
	}
}

// B1: [b] opens the backlog view listing Backlog-status tickets.
func TestBoard_BacklogView(t *testing.T) {
	m := step1Model()
	nm, _ := m.handleBoardKey(mkKey("b"))
	m = nm.(Model)
	if !m.boardBacklog {
		t.Fatal("b should open the backlog")
	}
	out := m.boardView()
	for _, want := range []string{"Backlog", "PXF-300", "OAuth", "Stripe billing"} {
		if !strings.Contains(out, want) {
			t.Errorf("backlog view missing %q", want)
		}
	}
}

// B6: [v] opens the analytics overlay (committed/done/burndown).
func TestBoard_Analytics(t *testing.T) {
	m := step1Model()
	m.boardSprint.Day = 3
	nm, _ := m.handleBoardKey(mkKey("v"))
	m = nm.(Model)
	if !m.boardAnalytic {
		t.Fatal("v should open analytics")
	}
	out := m.boardView()
	for _, want := range []string{"analytics", "Committed", "Burndown", "Velocity"} {
		if !strings.Contains(out, want) {
			t.Errorf("analytics view missing %q", want)
		}
	}
}

// esc is layered: it closes an open overlay before leaving the board.
func TestBoard_LayeredEsc(t *testing.T) {
	m := step1Model()
	m.boardReturn = screenStep0Complete
	m.boardBacklog = true
	nm, _ := m.handleBoardKey(mkKey("esc"))
	m = nm.(Model)
	if m.boardBacklog || m.screen == screenStep0Complete {
		t.Fatal("esc should close the backlog overlay, not leave the board")
	}
	nm, _ = m.handleBoardKey(mkKey("esc"))
	if nm.(Model).screen != screenStep0Complete {
		t.Fatal("a second esc should leave the board")
	}
}
