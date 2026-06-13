package guiapi

import (
	"strings"
	"testing"
)

// pythonDef catches a Python template leaking onto a non-Python surface — the
// release-blocking bug where the entrance test and tutorial seeded Python
// starters for every language.
func assertNativeStarter(t *testing.T, surface, lang, starter string) {
	t.Helper()
	if starter == "" {
		t.Errorf("%s/%s: empty starter", surface, lang)
		return
	}
	if lang != "python" && strings.Contains(starter, "def ") {
		t.Errorf("%s/%s: starter is Python:\n%s", surface, lang, starter)
	}
}

func TestReferenceCpp_StartersAreCpp(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	e, err := New()
	if err != nil {
		t.Fatal(err)
	}
	// Even reference-only C++ must seed C++ syntax everywhere an editor opens.
	p := e.cat.Problems[0]
	d := e.Problem(p.ID, "cpp")
	assertNativeStarter(t, "bench", "cpp", d.Starter)
	o := e.StartOrientation("cpp", "never")
	for i := 0; i < 40; i++ {
		st := o.Step()
		if st.Done {
			break
		}
		if st.Kind == "code" {
			assertNativeStarter(t, "orientation", "cpp", st.Starter)
		}
		o.SubmitChoice(0)
	}
}

func TestOrientationAndTutorial_StartersAreLanguageNative(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	e, err := New()
	if err != nil {
		t.Fatal(err)
	}
	for _, lang := range GradedLanguages() {
		// Entrance test: walk the intake far enough to see every code item
		// (submit wrong choice/spec answers; code items are skipped via a
		// failing grade-free path — we only need Step()'s starter).
		o := e.StartOrientation(lang, "never")
		codeSeen := 0
		for i := 0; i < 40; i++ {
			st := o.Step()
			if st.Done {
				break
			}
			switch st.Kind {
			case "code":
				codeSeen++
				assertNativeStarter(t, "orientation", lang, st.Starter)
				o.SubmitChoice(0) // anything that advances; grading not under test
			case "choice":
				o.SubmitChoice(0)
			default:
				o.SubmitSpec("answer")
			}
		}
		if codeSeen == 0 {
			t.Errorf("orientation/%s: no code item reached", lang)
		}

		// Tutorial Island: every staged task must serve a native starter.
		tut := e.StartTutorial(lang)
		tasks := 0
		for i := 0; i < tut.Count(); i++ {
			for _, stg := range tut.Lesson(i).Stages {
				if stg.HasTask {
					tasks++
					assertNativeStarter(t, "tutorial", lang, stg.Starter)
				}
			}
		}
		if tasks == 0 {
			t.Errorf("tutorial/%s: no staged tasks found", lang)
		}
	}
}

func TestLanguages_IncludesReferenceOnlyCpp(t *testing.T) {
	langs := Languages()
	if len(langs) != len(GradedLanguages())+1 {
		t.Fatalf("languages = %d", len(langs))
	}
	last := langs[len(langs)-1]
	if last.ID != "cpp" || last.Graded || last.Label != "C++" {
		t.Fatalf("cpp row: %+v", last)
	}
	for _, l := range langs[:len(langs)-1] {
		if !l.Graded || l.Label == "" {
			t.Fatalf("graded row: %+v", l)
		}
	}
}
