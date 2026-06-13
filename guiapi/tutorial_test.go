package guiapi

import (
	"testing"

	"devascent/internal/save"
)

// Tutorial frontier persistence: advances persist into the language slot with
// the TUI's resume fields, revisits never regress, restarts resume, completion
// marks the run done — and a run already past the tutorial keeps its Stage.
func TestTutorialFrontierPersistence(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	e, err := New()
	if err != nil {
		t.Fatal(err)
	}
	tut := e.StartTutorial("go")
	n := tut.Count()
	if n < 2 {
		t.Fatalf("need 2+ go lessons, got %d", n)
	}
	if pos := tut.Resume(); pos.Lesson != 0 || pos.Stage != 0 || pos.Done {
		t.Fatalf("fresh resume = %+v", pos)
	}

	// Advance within lesson 0, then into lesson 1.
	if pos := tut.Advance(0, 1); pos.Lesson != 0 || pos.Stage != 1 {
		t.Fatalf("advance(0,1) = %+v", pos)
	}
	if pos := tut.Advance(1, 0); pos.Lesson != 1 || pos.Stage != 0 {
		t.Fatalf("advance(1,0) = %+v", pos)
	}
	// Revisit never regresses.
	if pos := tut.Advance(0, 0); pos.Lesson != 1 || pos.Stage != 0 {
		t.Fatalf("frontier regressed: %+v", pos)
	}

	// The slot carries the TUI's resume fields.
	st, err := save.LoadLang("go")
	if err != nil || st == nil {
		t.Fatalf("go slot not written: %v, %v", st, err)
	}
	if st.Stage != "tutorial" || st.LessonIdx != 1 || st.StageIdx != 0 {
		t.Fatalf("slot fields wrong: stage=%q lesson=%d stageIdx=%d", st.Stage, st.LessonIdx, st.StageIdx)
	}

	// Restart: a fresh Engine resumes the frontier from disk.
	e2, err := New()
	if err != nil {
		t.Fatal(err)
	}
	tut2 := e2.StartTutorial("go")
	if pos := tut2.Resume(); pos.Lesson != 1 || pos.Stage != 0 || pos.Done {
		t.Fatalf("resume after restart = %+v", pos)
	}

	// Completion: advancing past the last lesson marks the run done.
	if pos := tut2.Advance(n, 0); !pos.Done {
		t.Fatalf("completion not reported: %+v", pos)
	}
	if st, _ := save.LoadLang("go"); st == nil || st.Stage != "done" {
		t.Fatalf("completed run stage = %v", st)
	}

	// A run already past the tutorial keeps its Stage (bookmarks still move).
	if err := save.SaveLang("rust", save.State{Stage: "bench", Placement: "test-out"}); err != nil {
		t.Fatal(err)
	}
	e3, err := New()
	if err != nil {
		t.Fatal(err)
	}
	tut3 := e3.StartTutorial("rust")
	tut3.Advance(0, 1)
	st3, _ := save.LoadLang("rust")
	if st3 == nil || st3.Stage != "bench" {
		t.Fatalf("bench run regressed to %q", st3.Stage)
	}
	if st3.LessonIdx != 0 || st3.StageIdx != 1 {
		t.Fatalf("bookmark not recorded: %+v", st3)
	}
}
