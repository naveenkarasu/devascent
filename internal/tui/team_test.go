package tui

import (
	"strings"
	"testing"

	"devascent/internal/ticket"
)

// Delegating a ticket: the teammate starts it now, and it's delivered at the due
// day via the day cycle — never before.
func TestDelegate_StartsThenDeliversAtDueDay(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	_, sp := seedSprint1()
	m := Model{screen: screenBoard, lang: "python", boardSprint: sp, playerLvl: 1}

	tk := &ticket.Ticket{Key: "PXF-900", Type: ticket.Task, Title: "delegated thing",
		Status: ticket.ToDo, Priority: ticket.PMinor, Assignee: "Maya",
		AssignedDay: sp.Day, DueDay: sp.Day + 2}
	sp.Tickets = append(sp.Tickets, tk)

	applyAssignment(tk, m.playerLvl, sp.Day)
	if tk.Status != ticket.InProgress {
		t.Fatalf("delegated ticket should start (In Progress), got %s", tk.Status)
	}
	if len(tk.Comments) == 0 || tk.Comments[0].Author != "Maya" {
		t.Fatalf("teammate should acknowledge, got %+v", tk.Comments)
	}

	// advance one day — still before the due day → not delivered yet
	m, _ = mustModel(m.advanceDay())
	if tk.Status == ticket.Done {
		t.Fatal("delegated ticket must not finish before its due day")
	}
	// advance to/past the due day → delivered
	m, _ = mustModel(m.advanceDay())
	if tk.Status != ticket.Done || tk.ResolvedDay != m.boardSprint.Day {
		t.Fatalf("delegated ticket should be delivered at the due day, got %s (resolved D%d, day D%d)", tk.Status, tk.ResolvedDay, m.boardSprint.Day)
	}
}

// Escalating a ticket: the senior replies with guidance immediately.
func TestEscalate_GivesGuidance(t *testing.T) {
	_, sp := seedSprint1()
	playerLvl := 1 // a junior escalating to a principal
	tk := &ticket.Ticket{Key: "PXF-901", Type: ticket.Bug, Title: "need help",
		Status: ticket.ToDo, Assignee: "Raj", AssignedDay: sp.Day, DueDay: sp.Day + 1}

	if assignKind(tk.Assignee, playerLvl) != "escalate" {
		t.Fatalf("assigning to Raj from level 1 should be an escalation")
	}
	applyAssignment(tk, playerLvl, sp.Day)
	if tk.Status != ticket.Done {
		t.Fatalf("an escalation should resolve with advice, got %s", tk.Status)
	}
	if len(tk.Comments) == 0 || !strings.Contains(strings.ToLower(tk.Comments[0].Body), "test") {
		t.Fatalf("senior should leave actionable guidance, got %+v", tk.Comments)
	}
}

// Your own tickets never auto-progress on the day cycle.
func TestAdvanceDay_DoesNotTouchYourTickets(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	_, sp := seedSprint1()
	m := Model{screen: screenBoard, lang: "python", boardSprint: sp, playerLvl: 1}
	yours := sp.Find("PXF-201") // assigned to you, In ToDo
	start := yours.Status
	for i := 0; i < 5; i++ {
		m, _ = mustModel(m.advanceDay())
	}
	if got := m.boardSprint.Find("PXF-201").Status; got != start {
		t.Fatalf("your own ticket must not auto-progress, was %s now %s", start, got)
	}
}
