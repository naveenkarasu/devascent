package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"devascent/internal/save"
)

// With more than one language slot, [c] opens the profile picker; choosing a
// slot resumes THAT language's run (not just the latest).
func TestProfilePicker_ResumesChosenSlot(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	if err := save.SaveLang("python", save.State{Stage: "done", Placement: "test-out"}); err != nil {
		t.Fatal(err)
	}
	if err := save.SaveLang("go", save.State{Stage: "tutorial", LessonIdx: 2}); err != nil {
		t.Fatal(err)
	}
	m := New()
	if len(m.profiles) != 2 {
		t.Fatalf("profiles = %d, want 2", len(m.profiles))
	}
	nm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("c")})
	m = nm.(Model)
	if m.screen != screenProfilePick {
		t.Fatalf("want profile picker, got screen %d", m.screen)
	}
	// pick the python slot (find its row; order is most-recent-first)
	want := 0
	for i, p := range m.profiles {
		if p.Lang == "python" {
			want = i
		}
	}
	for m.profIdx != want {
		nm, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
		m = nm.(Model)
	}
	nm, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = nm.(Model)
	if m.lang != "python" || m.placement != "test-out" {
		t.Fatalf("resumed wrong slot: lang=%q placement=%q", m.lang, m.placement)
	}
}

// Solves banked from the GUI (SolvedIDs in the shared slot, no TUI bench run)
// are restored into the TUI's solved set on resume.
func TestResume_RestoresGUIBankedSolves(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	if err := save.SaveLang("python", save.State{Stage: "done", Placement: "test-out",
		SolvedIDs: []string{"nc-two-sum", "nc-trapping-rain-water"}}); err != nil {
		t.Fatal(err)
	}
	m := New()
	nm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("c")})
	m = nm.(Model)
	if !m.solvedSet["nc-two-sum"] || !m.solvedSet["nc-trapping-rain-water"] {
		t.Fatalf("GUI-banked solves not restored: %v", m.solvedSet)
	}
}

// [t] on a lesson's reading stage jumps to the final hands-on stage.
func TestLesson_TestMeJumpsToFinalStage(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	if err := save.SaveLang("python", save.State{Stage: "tutorial", LessonIdx: 0, StageIdx: 0}); err != nil {
		t.Fatal(err)
	}
	m := New()
	nm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("c")})
	m = nm.(Model)
	if m.screen != screenLesson || m.stageIdx != 0 {
		t.Fatalf("want lesson stage 0, got screen %d stage %d", m.screen, m.stageIdx)
	}
	nm, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("t")})
	m = nm.(Model)
	if m.stageIdx != len(m.les.stages)-1 {
		t.Fatalf("test-me did not jump: stage %d of %d", m.stageIdx, len(m.les.stages))
	}
}

// The bench menu's dev-literacy practice runs the drills and returns to the
// bench menu (a routed gate run still hands off to the bench).
func TestDevLiteracyRevisit_ReturnsToBenchMenu(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	if err := save.SaveLang("python", save.State{Stage: "done", Placement: "test-out"}); err != nil {
		t.Fatal(err)
	}
	m := New()
	nm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("c")})
	m = nm.(Model)
	// enter the bench browse menu, pick the dev-literacy option
	nm, _ = m.startBench()
	m = nm.(Model)
	devPos := -1
	for i, o := range m.benchMenu {
		if o.kind == "devlit" {
			devPos = i
		}
	}
	if devPos < 0 {
		t.Fatal("no devlit option on the bench menu")
	}
	m.benchMenuIdx = devPos
	nm, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = nm.(Model)
	if m.screen != screenDevLiteracy || !m.devRevisit {
		t.Fatalf("devlit revisit not started: screen %d revisit %v", m.screen, m.devRevisit)
	}
	// completion routes back to the bench menu, not the handoff
	nm, _ = m.devDone()
	m = nm.(Model)
	if m.screen != screenBenchMenu || m.devRevisit {
		t.Fatalf("revisit completion: screen %d revisit %v", m.screen, m.devRevisit)
	}
}
