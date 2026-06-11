package guiapi

import (
	"os/exec"
	"testing"

	"devascent/internal/save"
)

const twoSumPy = "def two_sum(nums, target):\n" +
	"    seen = {}\n" +
	"    for i, n in enumerate(nums):\n" +
	"        if target - n in seen:\n" +
	"            return [seen[target - n], i]\n" +
	"        seen[n] = i\n" +
	"    return []\n"

// Banking: a pass banks the problem into the shared save, the scorecard and the
// browse list reflect it, a re-pass is not newly banked, and a fresh Engine
// (= app restart) restores it from disk.
func TestBankingPersistsAcrossEngines(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	if _, err := exec.LookPath("python"); err != nil {
		t.Skip("python not on PATH; skipping the banking round-trip")
	}

	e, err := New()
	if err != nil {
		t.Fatal(err)
	}
	if p := e.Progress("python"); p.Banked != 0 || p.Step0Met {
		t.Fatalf("fresh save not empty: %+v", p)
	}

	g := e.Grade("python", "nc-two-sum", twoSumPy)
	if !g.Passed || !g.Banked || !g.NewlyBanked || g.SaveErr != "" {
		t.Fatalf("first pass did not bank cleanly: %+v", g)
	}
	if p := e.Progress("python"); p.Banked != 1 {
		t.Fatalf("banked count after pass = %d, want 1", p.Banked)
	}
	solved := false
	for _, s := range e.Problems("python") {
		if s.ID == "nc-two-sum" {
			solved = s.Solved
		}
	}
	if !solved {
		t.Fatal("browse list does not mark nc-two-sum solved")
	}

	// Re-pass: still banked, not newly.
	if g2 := e.Grade("python", "nc-two-sum", twoSumPy); !g2.Banked || g2.NewlyBanked {
		t.Fatalf("re-pass banking wrong: %+v", g2)
	}

	// Next skips the solved problem.
	if next := e.NextProblem("python", "nc-two-sum"); next == "" || next == "nc-two-sum" {
		t.Fatalf("NextProblem after a banked solve = %q", next)
	}

	// Restart: a fresh Engine restores the banked set from disk.
	e2, err := New()
	if err != nil {
		t.Fatal(err)
	}
	if p := e2.Progress("python"); p.Banked != 1 {
		t.Fatalf("banked count after restart = %d, want 1", p.Banked)
	}

	// Per-language isolation: the python bank never leaks into another slot,
	// and the profile picker lists exactly the slots that exist.
	if p := e2.Progress("go"); p.Banked != 0 {
		t.Fatalf("go slot inherited python's bank: %+v", p)
	}
	profs := e2.Profiles()
	if len(profs) != 1 || profs[0].Lang != "python" || profs[0].Banked != 1 {
		t.Fatalf("profiles wrong: %+v", profs)
	}
	for _, s := range e2.Problems("go") {
		if s.ID == "nc-two-sum" && s.Solved {
			t.Fatal("go browse list shows python's solve")
		}
	}
}

// A completed orientation persists its placement block into the shared save.
func TestOrientationPersistsPlacement(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	e, err := New()
	if err != nil {
		t.Fatal(err)
	}
	o := e.StartOrientation("python", "never") // clamp: never → tutorial-full
	for steps := 0; !o.Step().Done; steps++ {
		if steps > 50 {
			t.Fatal("orientation did not terminate")
		}
		switch o.Step().Kind {
		case "code":
			o.SubmitCode("def stub():\n    return None\n")
		case "spec":
			o.SubmitSpec("")
		default:
			o.SubmitChoice(0)
		}
	}
	st, err := save.LoadLang("python")
	if err != nil || st == nil {
		t.Fatalf("save not written: st=%v err=%v", st, err)
	}
	if st.Placement != "tutorial-full" || st.Level != "never" {
		t.Fatalf("persisted placement block wrong: placement=%q level=%q", st.Placement, st.Level)
	}
	if p := e.Progress("python"); p.Placement != "tutorial-full" {
		t.Fatalf("Progress placement = %q", p.Placement)
	}
}
