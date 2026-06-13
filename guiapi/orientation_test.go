package guiapi

import (
	"os/exec"
	"testing"

	"devascent/internal/content"
)

// The orientation session machine selects a ladder, advances through every item,
// and terminates with a valid placement.
func TestOrientationSession(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir()) // never touch the real save
	e, err := New()
	if err != nil {
		t.Fatal(err)
	}
	o := e.StartOrientation("python", "a-little")
	if first := o.Step(); first.Done || first.Total == 0 {
		t.Fatalf("empty orientation: %+v", first)
	}
	for steps := 0; ; steps++ {
		if steps > 50 {
			t.Fatal("orientation did not terminate")
		}
		s := o.Step()
		if s.Done {
			break
		}
		switch s.Kind {
		case "code":
			// code items grade without advancing now — commit explicitly
			o.SubmitCode("def stub():\n    return None\n")
			o.AdvanceOrientation()
		case "spec":
			o.SubmitSpec("")
		default:
			o.SubmitChoice(0)
		}
	}
	final := o.Step()
	switch final.Placement {
	case "test-out", "dev-literacy", "tutorial-full":
		// ok
	default:
		t.Errorf("invalid placement %q (score %d/%d)", final.Placement, final.Score, final.Total)
	}
}

// langExe maps a language to the executable its grader needs on PATH.
func langExe(lang string) string {
	switch lang {
	case "go":
		return "go"
	case "javascript", "typescript":
		return "node"
	case "rust":
		return "rustc"
	case "java":
		return "javac"
	case "csharp":
		return "dotnet"
	}
	return ""
}

// Tutorial lessons load per language and their reference solutions grade clean
// through the real grader (proves the GUI tutorial → grader wiring end-to-end).
// Non-python lessons carry the reference Solution; python's originals are proven
// elsewhere. Grades the first language whose toolchain is available.
func TestTutorialGradesReferenceSolution(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir()) // never touch the real save
	e, err := New()
	if err != nil {
		t.Fatal(err)
	}
	cat, err := content.Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, lang := range []string{"go", "javascript", "typescript", "rust", "java", "csharp"} {
		exe := langExe(lang)
		if exe == "" {
			continue
		}
		if _, err := exec.LookPath(exe); err != nil {
			continue
		}
		tut := e.StartTutorial(lang)
		graded := 0
		for li, l := range cat.LessonsForLang(lang) {
			for si, st := range l.Stages {
				if st.Task == nil || st.Task.Solution == "" {
					continue
				}
				gr := tut.GradeStage(li, si, st.Task.Solution)
				if gr.Err != "" {
					t.Fatalf("%s lesson %d stage %d grade error: %s", lang, li, si, gr.Err)
				}
				if !gr.Passed {
					t.Errorf("%s lesson %d (%s) stage %d: reference solution did not pass: %+v", lang, li, l.ID, si, gr)
				}
				graded++
			}
		}
		if graded > 0 {
			t.Logf("%s: graded %d lesson tasks via reference solutions", lang, graded)
			return
		}
	}
	t.Skip("no language with lessons+solutions and an available toolchain")
}
