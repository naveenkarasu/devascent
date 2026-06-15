package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

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

// [n] opens the form; filling it (type via ←/→, title via typing) and [enter]
// files a new To Do ticket assigned to you, with the cursor on it.
func TestNewTicket_CreatesToDoCard(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	proj, sp := seedSprint1()
	m := Model{screen: screenBoard, lang: "python", boardProject: proj, boardSprint: sp, playerLvl: 1}
	before := len(sp.Column(ticket.ToDo))

	opened, _ := m.handleBoardKey(mkKey("n"))
	m = opened.(Model)
	if m.screen != screenNewTicket {
		t.Fatalf("n should open the create form, got screen %d", m.screen)
	}

	// type the title (focus starts on Title)
	for _, ch := range "Investigate flaky checkout test" {
		m, _ = mustModel(m.handleFormKey(runeKey(ch)))
	}
	// move to the Type field (form navigation is on key TYPE, not runes)
	guard := 0
	for ticketField(m) != fldType && guard < 20 {
		m, _ = mustModel(m.handleFormKey(tea.KeyMsg{Type: tea.KeyDown}))
		guard++
	}
	for formTypes[m.ntType] != ticket.Bug && guard < 40 {
		m, _ = mustModel(m.handleFormKey(tea.KeyMsg{Type: tea.KeyRight}))
		guard++
	}
	// submit
	created, _ := m.handleFormKey(tea.KeyMsg{Type: tea.KeyEnter})
	got := created.(Model)

	if got.screen != screenBoard {
		t.Fatalf("creating should return to the board, got screen %d", got.screen)
	}
	if len(got.boardSprint.Column(ticket.ToDo)) != before+1 {
		t.Fatalf("To Do should gain one card, want %d", before+1)
	}
	nt := got.selectedTicket()
	if nt == nil || nt.Title != "Investigate flaky checkout test" || nt.Type != ticket.Bug || nt.Assignee != "you" || nt.Status != ticket.ToDo {
		t.Fatalf("new ticket has wrong fields: %+v", nt)
	}
	if !strings.HasPrefix(nt.Key, "PXF-") {
		t.Errorf("new ticket should have a fresh PXF key, got %q", nt.Key)
	}
}

func ticketField(m Model) int { return m.ntFocus }

func runeKey(r rune) tea.KeyMsg { return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}} }

// Delegation vs escalation depends on the player's level relative to the teammate.
func TestAssignKind(t *testing.T) {
	if got := assignKind("you", 1); got != "self" {
		t.Errorf("you → self, got %s", got)
	}
	if got := assignKind("Maya", 1); got != "delegate" { // junior, peer/below a level-1 player
		t.Errorf("Maya@lvl1 → delegate, got %s", got)
	}
	if got := assignKind("Raj", 1); got != "escalate" { // principal, above a level-1 player
		t.Errorf("Raj@lvl1 → escalate, got %s", got)
	}
	if got := assignKind("Raj", 4); got != "delegate" { // as a principal, Raj is a peer
		t.Errorf("Raj@lvl4 → delegate, got %s", got)
	}
	if got := assignKind("Sam", 4); got != "escalate" { // manager is still above
		t.Errorf("Sam@lvl4 → escalate, got %s", got)
	}
}

// [E] edit reassigns and re-prioritizes an existing ticket (and recomputes its SLA).
func TestForm_EditReassignsAndReprioritizes(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	_, sp := seedSprint1()
	m := Model{screen: screenBoard, lang: "python", boardSprint: sp, playerLvl: 3}
	tk := sp.Find("PXF-201")

	m, _ = mustModel(m.openEditTicket(tk))
	if m.screen != screenNewTicket || m.ntEditKey != "PXF-201" {
		t.Fatalf("openEditTicket should load the form for PXF-201, got screen %d key %q", m.screen, m.ntEditKey)
	}
	m.ntAssignee = idxOf(assigneeOptions(), "Priya")
	m.ntPri = idxOf(formPris, ticket.PCritical)

	out, _ := m.submitForm()
	got := out.(Model)
	// Reassigning to a teammate (delegation) opens discuss-&-agree before they start.
	if got.screen != screenDiscuss {
		t.Fatalf("reassigning to a teammate should open discuss-&-agree, got screen %d", got.screen)
	}
	e := got.boardSprint.Find("PXF-201")
	if e.Assignee != "Priya" || e.Priority != ticket.PCritical {
		t.Fatalf("edit didn't apply: assignee=%q priority=%q", e.Assignee, e.Priority)
	}
	if e.DueDay != ticket.DueDayFor(ticket.PCritical, e.AssignedDay) {
		t.Errorf("SLA should be recomputed on priority change, got DueDay %d", e.DueDay)
	}
	// Agreeing starts the teammate on it and returns to the board.
	agreed, _ := got.handleDiscussKey(mkKey("enter"))
	g2 := agreed.(Model)
	if e.Status != ticket.InProgress {
		t.Errorf("after agreeing, %s should be In Progress, got %s", e.Key, e.Status)
	}
	if g2.screen != screenBoard {
		t.Errorf("agree should return to the board, got screen %d", g2.screen)
	}
}
