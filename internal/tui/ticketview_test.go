package tui

import (
	"strings"
	"testing"

	"devascent/internal/ticket"
)

// enter on the cursor card opens the detail view; esc returns to the board.
func TestTicketDetail_OpenAndBack(t *testing.T) {
	m := demoModel() // focus In Progress, row 0 → PXF-101 (the enriched demo ticket)
	sel := m.selectedTicket()
	if sel == nil || sel.Key != "PXF-101" {
		t.Fatalf("setup: cursor should be on PXF-101, got %v", sel)
	}

	nm, _ := m.handleBoardKey(mkKey("enter"))
	got := nm.(Model)
	if got.screen != screenTicket || got.detailTicket == nil {
		t.Fatalf("enter should open the ticket detail, got screen %d", got.screen)
	}

	out := got.ticketDetailView()
	for _, want := range []string{
		"PXF-101", "Off-by-one in the paginator",
		"DESCRIPTION", "ACCEPTANCE", "page 1 returns items 1..size",
		"LEARN", "paginate", "COMMENTS", "Sam (manager)", "[esc] back",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("detail view missing %q", want)
		}
	}

	back, _ := got.handleTicketKey(mkKey("esc"))
	if back.(Model).screen != screenBoard {
		t.Fatalf("esc should return to the board, got screen %d", back.(Model).screen)
	}
}

// Opening an empty column's "selection" is a safe no-op (no card under cursor).
func TestTicketDetail_NoSelectionNoOpen(t *testing.T) {
	m := demoModel()
	// point the cursor past the end of a column
	m.boardRow = 99
	if m.selectedTicket() != nil {
		t.Fatal("out-of-range cursor should select no ticket")
	}
	nm, _ := m.handleBoardKey(mkKey("enter"))
	if nm.(Model).screen != screenBoard {
		t.Fatal("enter with no selection must stay on the board")
	}
}

func TestLearnLang_Default(t *testing.T) {
	if got := learnLang(&ticket.Ticket{}); got != "python" {
		t.Errorf("learnLang default = %q, want python", got)
	}
	g := &ticket.Ticket{Grading: &ticket.Grading{Lang: "go"}}
	if got := learnLang(g); got != "go" {
		t.Errorf("learnLang from grading = %q, want go", got)
	}
}
