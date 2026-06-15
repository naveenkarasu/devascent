package tui

import (
	"strings"
	"testing"

	"devascent/internal/ticket"
)

func TestNextKey(t *testing.T) {
	_, sp := seedSprint1()
	got := nextKey(sp, &ticket.Project{Key: "PXF"})
	if !strings.HasPrefix(got, "PXF-") {
		t.Errorf("nextKey = %q, want a PXF- key", got)
	}
	for _, tk := range sp.Tickets { // must be fresh (not collide with any existing key)
		if tk.Key == got {
			t.Fatalf("nextKey returned an existing key %q", got)
		}
	}
	if e := nextKey(&ticket.Sprint{}, nil); e != "PXF-1" {
		t.Errorf("nextKey on empty = %q, want PXF-1", e)
	}
}

// [n] opens the form; creating files a new To Do ticket assigned to you and puts
// the cursor on it.
func TestNewTicket_CreatesToDoCard(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	proj, sp := seedSprint1()
	m := Model{screen: screenBoard, lang: "python", boardProject: proj, boardSprint: sp}
	before := len(sp.Column(ticket.ToDo))

	opened, _ := m.handleBoardKey(mkKey("n"))
	m = opened.(Model)
	if m.screen != screenNewTicket || !m.inputActive {
		t.Fatalf("n should open the file-a-ticket form (modal), got screen %d active %v", m.screen, m.inputActive)
	}

	m.ntType = 1 // Bug
	created, _ := m.createTicket("Investigate flaky checkout test")
	got := created.(Model)
	if got.screen != screenBoard {
		t.Fatalf("creating should return to the board, got screen %d", got.screen)
	}
	td := got.boardSprint.Column(ticket.ToDo)
	if len(td) != before+1 {
		t.Fatalf("To Do should gain one card: got %d, want %d", len(td), before+1)
	}
	// the cursor should land on the freshly-created card
	nt := got.selectedTicket()
	if nt == nil || nt.Title != "Investigate flaky checkout test" || nt.Type != ticket.Bug || nt.Assignee != "you" || nt.Status != ticket.ToDo {
		t.Fatalf("new ticket has wrong fields: %+v", nt)
	}
	if !strings.HasPrefix(nt.Key, "PXF-") {
		t.Errorf("new ticket should have a fresh PXF key, got %q", nt.Key)
	}
}
