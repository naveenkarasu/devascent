package guiapi

import (
	"os/exec"
	"strings"
	"testing"
)

// Smoke test for the GUI seam: load, list, open a problem with a starter, and
// grade a correct + a wrong solution through the real grader.
func TestGuiApiSmoke(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir()) // never touch the real save
	e, err := New()
	if err != nil {
		t.Fatal(err)
	}
	if len(e.Problems("python")) < 100 {
		t.Fatalf("expected 100+ bench problems, got %d", len(e.Problems("python")))
	}

	d := e.Problem("nc-two-sum", "python")
	if !d.Found || d.FuncName != "two_sum" || strings.TrimSpace(d.Starter) == "" {
		t.Fatalf("problem detail wrong: %+v", d)
	}

	if _, err := exec.LookPath("python"); err != nil {
		t.Skip("python not on PATH; skipping the grade round-trip")
	}
	correct := "def two_sum(nums, target):\n    seen = {}\n    for i, n in enumerate(nums):\n        if target - n in seen:\n            return [seen[target - n], i]\n        seen[n] = i\n    return []\n"
	if g := e.Grade("python", "nc-two-sum", correct); !g.Passed || g.CasesFailed != 0 {
		t.Errorf("correct solution did not pass: %+v", g)
	}
	wrong := "def two_sum(nums, target):\n    return []\n"
	if g := e.Grade("python", "nc-two-sum", wrong); g.Passed {
		t.Errorf("wrong solution unexpectedly passed: %+v", g)
	}
}
