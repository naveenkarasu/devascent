package tui

import (
	"testing"

	"devascent/internal/save"
	"devascent/internal/ticket"
)

// enterStep1 (from the graduation screen) opens the board as the persistent home,
// seeding the first sprint and writing it to the save slot.
func TestStep1_EnterSeedsAndPersists(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	m := Model{screen: screenStep0Complete, lang: "python"}
	got, _ := mustModel(m.enterStep1())
	if got.screen != screenBoard {
		t.Fatalf("enterStep1 should open the board, got screen %d", got.screen)
	}
	if !got.step1Home || got.boardSprint == nil {
		t.Fatal("enterStep1 should mark step1Home and seed a sprint")
	}
	st, err := save.LoadLang("python")
	if err != nil || st == nil || st.Sprint == nil {
		t.Fatalf("enterStep1 should persist the sprint to the slot, got %v, %v", st, err)
	}
	if st.Stage != "step1" {
		t.Errorf("saved stage = %q, want step1", st.Stage)
	}
}

// A ticket move on the real career home persists and survives a save→resume.
func TestStep1_MovePersistsAndResumes(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	m, _ := mustModel(Model{screen: screenStep0Complete, lang: "python"}.enterStep1())

	tk := m.boardSprint.Find("PXF-202") // a To Do ticket
	if tk == nil {
		t.Fatal("seed should contain PXF-202")
	}
	_ = tk.MoveTo(ticket.InProgress)
	m.persist()

	// resume from the saved slot
	st, _ := save.LoadLang("python")
	resumed, _ := mustModel((Model{lang: "python", resume: st}).applyResume())
	if resumed.screen != screenBoard {
		t.Fatalf("resume should land on the board, got screen %d", resumed.screen)
	}
	if !resumed.step1Home {
		t.Error("resume should mark the board as the persistent home")
	}
	rt := resumed.boardSprint.Find("PXF-202")
	if rt == nil || rt.Status != ticket.InProgress {
		t.Fatalf("the moved status should survive resume, got %+v", rt)
	}
}

// The cheat preview (step1Home == false) must NOT write to the save.
func TestStep1_CheatPreviewDoesNotPersist(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	m := Model{screen: screenBoard, lang: "python", step1Home: false}
	m.boardProject, m.boardSprint = seedSprint1()
	m.persist()
	if st, _ := save.LoadLang("python"); st != nil {
		t.Fatal("a cheat-preview board must not touch the save slot")
	}
}
