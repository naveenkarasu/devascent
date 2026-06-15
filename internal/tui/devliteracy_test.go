package tui

import (
	"testing"

	"devascent/internal/content"
)

func devLiteracyModel() Model {
	m := Model{screen: screenDevLiteracy, inputActive: true}
	m.dev = []content.DevTask{{ID: "a", Title: "A"}, {ID: "b", Title: "B"}}
	m.enterDevTask()
	return m
}

// [tab] skips a task you're stuck on — advances without quitting.
func TestDevLiteracy_SkipAdvances(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	nm, _ := devLiteracyModel().skipDevTask()
	got := nm.(Model)
	if got.quitting {
		t.Fatal("skip must not quit the game")
	}
	if got.devIdx != 1 || got.curDev.ID != "b" {
		t.Fatalf("skip should move to task 2 (id=b), got idx=%d id=%q", got.devIdx, got.curDev.ID)
	}
}

// [esc] backs out to the menu (the bench hand-off) — it must NOT quit the game.
func TestDevLiteracy_EscLeavesToMenu(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	nm, _ := devLiteracyModel().leaveDevLiteracy()
	got := nm.(Model)
	if got.quitting {
		t.Fatal("esc must not quit dev-literacy")
	}
	if got.screen != screenHandoff {
		t.Fatalf("esc should land on the hand-off menu (with [b] bench), got screen %d", got.screen)
	}
	if got.inputActive {
		t.Fatal("leaving should clear modal input")
	}
}

// Entered as bench practice, leaving returns to the bench browse, not the hand-off.
func TestDevLiteracy_EscRevisitReturnsToBench(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	m := devLiteracyModel()
	m.devRevisit = true
	if c, err := content.Load(); err == nil {
		m.cat = c // startBench needs the problem catalog
	}
	nm, _ := m.leaveDevLiteracy()
	got := nm.(Model)
	if got.quitting {
		t.Fatal("esc must not quit")
	}
	if got.screen != screenBenchMenu {
		t.Fatalf("practice revisit should return to the bench menu, got screen %d", got.screen)
	}
}
