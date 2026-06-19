package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"devascent/internal/ticket"
)

func demoModel() Model {
	proj, sp := demoBoard()
	return Model{screen: screenBoard, boardProject: proj, boardSprint: sp, boardCol: defaultFocusCol(sp, sp.Day)}
}

func TestBoardView_RendersDemo(t *testing.T) {
	out := demoModel().boardView()
	for _, want := range []string{
		"Pixel Forge", "Sprint 3", "committed 13", "done 5pt",
		"To Do", "In Progress", "In Review", "Done",
		"PXF-104", "PXF-101", "PXF-099", "PXF-097",
		"@you", "esc back",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("board view missing %q", want)
		}
	}
	// Eyeball the layout (plain box-drawing; color shows in a real terminal).
	t.Log("\n" + out)
}

func TestDefaultFocusCol_PrefersInProgress(t *testing.T) {
	_, sp := demoBoard()
	if idx := defaultFocusCol(sp, sp.Day); ticket.BoardColumns[idx] != ticket.InProgress {
		t.Fatalf("focus should default to In Progress, got %s", ticket.BoardColumns[idx])
	}
}

func TestTrunc(t *testing.T) {
	if got := trunc("hello world", 5); got != "hell…" {
		t.Errorf("trunc long = %q, want %q", got, "hell…")
	}
	if got := trunc("hi", 5); got != "hi" {
		t.Errorf("trunc short = %q", got)
	}
}

func mkKey(s string) tea.KeyMsg {
	switch s {
	case "esc":
		return tea.KeyMsg{Type: tea.KeyEsc}
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter}
	default:
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
	}
}

func boardNav(m Model, s string) Model {
	nm, _ := m.handleBoardKey(mkKey(s))
	return nm.(Model)
}

func TestBoard_ColumnNavigation(t *testing.T) {
	m := demoModel() // focus defaults to In Progress (col 1)
	if ticket.BoardColumns[m.boardCol] != ticket.InProgress {
		t.Fatalf("setup: expected focus on In Progress, got %d", m.boardCol)
	}
	m = boardNav(m, "d") // → In Review
	m = boardNav(m, "d") // → Done
	if ticket.BoardColumns[m.boardCol] != ticket.Done {
		t.Fatalf("after d,d focus should be Done, got %d", m.boardCol)
	}
	m = boardNav(m, "d") // clamp at the last column
	if ticket.BoardColumns[m.boardCol] != ticket.Done {
		t.Fatalf("d past the end must clamp at Done, got %d", m.boardCol)
	}
	m = boardNav(m, "a") // ← In Review
	if ticket.BoardColumns[m.boardCol] != ticket.InReview {
		t.Fatalf("a should move back to In Review, got %d", m.boardCol)
	}
}

func TestBoard_CardNavigationClamps(t *testing.T) {
	m := demoModel()
	for ticket.BoardColumns[m.boardCol] != ticket.Done { // move to Done (3 cards)
		m = boardNav(m, "d")
	}
	if m.boardRow != 0 {
		t.Fatalf("changing column should reset the card cursor, got row %d", m.boardRow)
	}
	m = boardNav(m, "s")
	m = boardNav(m, "s")
	if m.boardRow != 2 {
		t.Fatalf("two s in a 3-card column should select row 2, got %d", m.boardRow)
	}
	m = boardNav(m, "s") // clamp at the last card
	if m.boardRow != 2 {
		t.Fatalf("s past the last card must clamp, got %d", m.boardRow)
	}
	m = boardNav(m, "w")
	if m.boardRow != 1 {
		t.Fatalf("w should move the cursor up, got %d", m.boardRow)
	}
}

func TestBoard_HelpToggleAndEsc(t *testing.T) {
	m := demoModel()
	m = boardNav(m, "?")
	if !m.boardHelp {
		t.Fatal("? should open the help overlay")
	}
	if !strings.Contains(m.boardView(), "Board — keys") {
		t.Fatal("help view should render the key list")
	}
	m = boardNav(m, "esc") // closes help, does NOT leave the board
	if m.boardHelp {
		t.Fatal("esc should close the help overlay")
	}
	m.boardReturn = screenBenchMenu
	m = boardNav(m, "esc") // now leaves the board
	if m.screen != screenBenchMenu {
		t.Fatalf("esc should leave the board to boardReturn, got screen %d", m.screen)
	}
}

func TestBoard_NarrowFallback(t *testing.T) {
	m := demoModel()
	m.width = 70 // below boardWideMin → single-column + tab strip
	out := m.boardView()
	if !strings.Contains(out, "In Prog") {
		t.Error("narrow view should render the abbreviated tab strip")
	}
	if !strings.Contains(out, "PXF-101") {
		t.Error("narrow view should render the focused column's cards")
	}
	if !strings.Contains(out, "‹") {
		t.Error("narrow view should mark the focused tab")
	}
}
