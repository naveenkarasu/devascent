package content

import (
	"os/exec"
	"testing"

	"devascent/internal/grader"
)

// TestAdvancedTopicsLoad: the Stage-2 spike content loads and is well-formed.
func TestAdvancedTopicsLoad(t *testing.T) {
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	// Every general-purpose language must have a populated Advanced Topics track.
	for _, lang := range []string{"python", "java", "cpp", "csharp", "javascript", "typescript", "rust", "go"} {
		if n := len(c.AdvancedTopicsByLang(lang)); n < 5 {
			t.Errorf("language %q has only %d advanced topics (expected >= 5)", lang, n)
		}
	}
	for _, at := range c.AdvancedTopics {
		if at.Lang == "" || at.Group == "" || at.Title == "" {
			t.Errorf("advanced topic missing lang/group/title: %+v", at.Group)
		}
		if at.Tag == "E" || at.Tag == "gotcha" {
			if len(at.Exercises) == 0 {
				t.Errorf("%s/%s is exercisable but has no exercises", at.Lang, at.Group)
			}
		}
		for _, e := range at.Exercises {
			if e.Prompt == "" || e.BrokenCode == "" || e.FixedCode == "" || e.Bug == "" {
				t.Errorf("%s/%s: exercise missing a required field (prompt/broken/fixed/bug)", at.Lang, at.Group)
			}
		}
	}
}

// TestAdvancedExercisesGradeRoundTrip is the LOAD-BEARING spike proof: every
// Check=="tests" exercise must (1) grade its FixedCode as PASS through the real
// grader, and (2) grade its BrokenCode as FAIL — proving the schema carries enough
// to auto-grade AND that the exercise actually discriminates correct from buggy.
// This is the forward-compat claim, demonstrated today on the only graded language.
func TestAdvancedExercisesGradeRoundTrip(t *testing.T) {
	if _, err := exec.LookPath("python"); err != nil {
		t.Skip("python not found on PATH")
	}
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	g := grader.NewNativePython()
	graded := 0
	for _, at := range c.AdvancedTopics {
		if at.Lang != "python" {
			continue
		}
		for _, e := range at.Exercises {
			if e.Check != "tests" {
				continue
			}
			graded++
			// (1) FixedCode must PASS.
			v, err := g.Run("python", e.FixedCode, e.FuncName, e.Tests, e.GraderShape())
			if err != nil {
				t.Errorf("%s/%s %q: grader error on FixedCode: %v", at.Lang, at.Group, e.Prompt[:20], err)
				continue
			}
			if !v.Passed {
				t.Errorf("%s/%s: FixedCode did NOT pass its tests: %+v", at.Lang, at.Group, v.Results)
			}
			// (2) BrokenCode must FAIL (otherwise the exercise teaches nothing).
			vb, err := g.Run("python", e.BrokenCode, e.FuncName, e.Tests, e.GraderShape())
			if err == nil && vb.Passed {
				t.Errorf("%s/%s: BrokenCode unexpectedly PASSED — exercise does not discriminate", at.Lang, at.Group)
			}
		}
	}
	if graded == 0 {
		t.Error("no Check==tests exercises found to prove the grading round-trip")
	} else {
		t.Logf("grading round-trip proven on %d Python advanced exercise(s)", graded)
	}
}
