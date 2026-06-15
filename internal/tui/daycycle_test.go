package tui

import (
	"testing"
	"time"

	"devascent/internal/ticket"
)

// cooldownRemaining is pure wall-clock math: future deadline → >0, past/empty → 0.
func TestCooldownRemaining(t *testing.T) {
	sp := &ticket.Sprint{}
	if cooldownRemaining(sp, time.Now()) != 0 {
		t.Error("no deadline → 0")
	}
	sp.CooldownEndsAt = time.Now().Add(30 * time.Second).UTC().Format(time.RFC3339)
	if d := cooldownRemaining(sp, time.Now()); d <= 0 || d > 31*time.Second {
		t.Errorf("future deadline → ~30s, got %v", d)
	}
	sp.CooldownEndsAt = time.Now().Add(-time.Minute).UTC().Format(time.RFC3339)
	if cooldownRemaining(sp, time.Now()) != 0 {
		t.Error("past deadline → 0")
	}
}

// [e] enters cooldown (deadline in the future, countdown started); a tick before
// the deadline keeps counting; once elapsed a tick flips to the next day.
func TestCooldown_EnterTickAndElapse(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	m := step1Model()
	out, cmd := m.enterCooldown()
	m = out.(Model)
	if m.screen != screenCooldown || m.boardSprint.Phase != ticket.PhaseCooldown {
		t.Fatalf("enterCooldown wrong: screen %d phase %q", m.screen, m.boardSprint.Phase)
	}
	if cmd == nil {
		t.Error("enterCooldown should start the countdown tick")
	}
	if cooldownRemaining(m.boardSprint, time.Now()) <= 0 {
		t.Error("the deadline should be in the future")
	}
	// a tick before the deadline keeps counting
	out, _ = m.onCooldownTick()
	m = out.(Model)
	if m.boardSprint.Phase != ticket.PhaseCooldown {
		t.Error("a tick before the deadline should stay in cooldown")
	}
	// force the deadline into the past → the next tick flips the day
	day := m.boardSprint.Day
	m.boardSprint.CooldownEndsAt = time.Now().Add(-time.Second).UTC().Format(time.RFC3339)
	out, _ = m.onCooldownTick()
	m = out.(Model)
	if m.boardSprint.Phase != ticket.PhaseStandup || m.boardSprint.Day != day+1 {
		t.Fatalf("an elapsed tick should flip to standup-pending day+1, got phase %q day %d", m.boardSprint.Phase, m.boardSprint.Day)
	}
}

// [s] during cooldown skips the wait and flips the day immediately.
func TestCooldown_SkipFlipsDay(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	m := step1Model()
	m, _ = mustModel(m.enterCooldown())
	day := m.boardSprint.Day
	out, _ := m.handleCooldownKey(mkKey("s"))
	m = out.(Model)
	if m.boardSprint.Phase != ticket.PhaseStandup || m.boardSprint.Day != day+1 {
		t.Fatalf("skip should flip to standup-pending day+1, got phase %q day %d", m.boardSprint.Phase, m.boardSprint.Day)
	}
}

// During cooldown the board locks mutations but still allows navigation, and esc
// returns to the cooldown screen.
func TestReadOnlyBoard_LocksMutationsAllowsNav(t *testing.T) {
	m := step1Model()
	m.boardSprint.Phase = ticket.PhaseCooldown
	m.boardCol = 1
	out, _ := m.handleBoardKey(mkKey("n")) // file-a-ticket is locked
	m = out.(Model)
	if m.screen == screenNewTicket {
		t.Fatal("new-ticket must be locked during cooldown")
	}
	out, _ = m.handleBoardKey(mkKey("a")) // nav still works
	m = out.(Model)
	if m.boardCol != 0 {
		t.Errorf("read-only nav should still move, col=%d", m.boardCol)
	}
	out, _ = m.handleBoardKey(mkKey("esc"))
	m = out.(Model)
	if m.screen != screenCooldown {
		t.Errorf("esc during cooldown should return to the recap screen, got %d", m.screen)
	}
}

// Once the new day is pending, [j] on the board joins the morning standup.
func TestReadOnlyBoard_StandupJoin(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	m := step1Model()
	m.boardSprint.Phase = ticket.PhaseStandup
	out, _ := m.handleBoardKey(mkKey("j"))
	m = out.(Model)
	if m.screen != screenStandup {
		t.Fatalf("[j] during standup-pending should join the standup, got %d", m.screen)
	}
}

// Filing a ticket assigned to a teammate opens discuss-&-agree (the ticket is
// created but not started); agreeing starts them with the plan as a comment.
func TestDiscuss_CreateDelegateThenAgree(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	proj, sp := seedSprint1()
	m := Model{screen: screenNewTicket, lang: "python", boardProject: proj, boardSprint: sp, playerLvl: 1}
	m.ntTitle = "Tidy the asset tags"
	m.ntAssignee = idxOf(assigneeOptions(), "Maya") // a junior → delegate from level 1
	before := len(sp.Tickets)

	out, _ := m.submitForm()
	m = out.(Model)
	if m.screen != screenDiscuss {
		t.Fatalf("delegating a new ticket should open discuss, got screen %d", m.screen)
	}
	if len(m.boardSprint.Tickets) != before+1 {
		t.Fatalf("the ticket should be created pending agreement")
	}
	tk := m.discussTk
	if tk == nil || tk.Assignee != "Maya" || tk.Status != ticket.ToDo {
		t.Fatalf("pending delegated ticket wrong: %+v", tk)
	}
	if m.discussPlan == "" {
		t.Error("the teammate should propose a plan")
	}

	plan := m.discussPlan
	out, _ = m.handleDiscussKey(mkKey("enter"))
	m = out.(Model)
	if tk.Status != ticket.InProgress {
		t.Fatalf("agree should start the ticket, got %s", tk.Status)
	}
	if len(tk.Comments) == 0 || tk.Comments[len(tk.Comments)-1].Body != plan {
		t.Fatalf("agree should post the plan as the teammate's comment, got %+v", tk.Comments)
	}
	if m.screen != screenBoard {
		t.Errorf("agree should return to the board, got %d", m.screen)
	}
}

// Cancelling discuss reverts the assignment — you keep the ticket, unstarted.
func TestDiscuss_CancelReverts(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	_, sp := seedSprint1()
	tk := &ticket.Ticket{Key: "PXF-951", Status: ticket.ToDo, Assignee: "Maya"}
	sp.Tickets = append(sp.Tickets, tk)
	m := Model{screen: screenDiscuss, lang: "python", boardSprint: sp, playerLvl: 1,
		discussTk: tk, discussAss0: "you"}
	out, _ := m.handleDiscussKey(mkKey("esc"))
	m = out.(Model)
	if tk.Assignee != "you" {
		t.Fatalf("cancel should revert the assignee to you, got %q", tk.Assignee)
	}
	if tk.Status != ticket.ToDo {
		t.Errorf("cancel must not start the ticket, got %s", tk.Status)
	}
	if m.screen != screenBoard {
		t.Errorf("cancel should return to the board, got %d", m.screen)
	}
}
