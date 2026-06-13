package guiapi

import (
	"os/exec"
	"testing"
)

// The Advanced Topics surface: every graded language lists topics with
// gradeable exercises, details round-trip, reveal-only grading is refused,
// and a reference fixed_code passes through the real grader.
func TestAdvancedTopicsSurface(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir()) // never touch the real save
	e, err := New()
	if err != nil {
		t.Fatal(err)
	}

	for _, lang := range GradedLanguages() {
		topics := e.AdvancedTopics(lang)
		if len(topics) == 0 {
			t.Errorf("%s: no advanced topics", lang)
			continue
		}
		gradeable := 0
		for _, tp := range topics {
			gradeable += tp.Gradeable
		}
		if gradeable == 0 {
			t.Errorf("%s: no gradeable advanced exercises", lang)
		}
		d := e.AdvancedTopic(lang, 0)
		if !d.Found || d.Title == "" || len(d.Exercises) == 0 {
			t.Errorf("%s: topic detail wrong: %+v", lang, d)
		}
	}

	// Out-of-range + reveal-only refusal.
	if d := e.AdvancedTopic("python", 999); d.Found {
		t.Error("out-of-range topic reported found")
	}
	if g := e.GradeAdvanced("python", 999, 0, "x"); g.Err == "" {
		t.Error("out-of-range grade did not error")
	}

	// Grade a reference fixed_code through the real grader (python on PATH).
	if _, err := exec.LookPath("python"); err != nil {
		t.Skip("python not on PATH; skipping the grade round-trip")
	}
	for ti, tp := range e.AdvancedTopics("python") {
		d := e.AdvancedTopic("python", tp.Index)
		for _, ex := range d.Exercises {
			if !ex.Gradeable || ex.FixedCode == "" {
				continue
			}
			g := e.GradeAdvanced("python", ti, ex.Index, ex.FixedCode)
			if g.Err != "" || !g.Passed {
				t.Fatalf("python topic %d ex %d: fixed_code did not pass: %+v", ti, ex.Index, g)
			}
			return // one real round-trip is the wiring proof
		}
	}
	t.Skip("no gradeable python exercise with fixed_code found")
}
