package guiapi

import "testing"

// findCodeStep walks an orientation to the first code item, answering any
// choice/spec items along the way, and returns its 1-based index (0 = none).
func driveToCodeItem(o *Orientation) bool {
	for i := 0; i < 40; i++ {
		st := o.Step()
		if st.Done {
			return false
		}
		if st.Kind == "code" {
			return true
		}
		if st.Kind == "choice" {
			o.SubmitChoice(0)
		} else {
			o.SubmitSpec("answer")
		}
	}
	return false
}

func TestOrientation_CodeSubmitDoesNotAdvance(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	e, err := New()
	if err != nil {
		t.Fatal(err)
	}
	o := e.StartOrientation("python", "regularly")
	if !driveToCodeItem(o) {
		t.Skip("no code item in the python intake")
	}
	before := o.Step()

	// A failing submit must NOT advance — same item stays on screen for a retry.
	out := o.SubmitCode("def wrong(:\n  syntax error")
	if out.Advanced {
		t.Fatal("failing code submit advanced the session")
	}
	if out.Passed {
		t.Fatal("broken code reported as passed")
	}
	if o.Step().Index != before.Index {
		t.Fatalf("session moved off the item: %d → %d", before.Index, o.Step().Index)
	}

	// Re-running is allowed and still doesn't advance on another failure.
	out = o.SubmitCode("still broken")
	if out.Advanced || o.Step().Index != before.Index {
		t.Fatal("second failing submit advanced")
	}

	// Only an explicit Advance moves on (here: a Skip — last grade failed).
	adv := o.AdvanceOrientation()
	if !adv.Advanced {
		t.Fatal("AdvanceOrientation did not advance")
	}
	if o.Step().Index == before.Index && !o.Step().Done {
		t.Fatal("still on the same item after advancing")
	}
}

func TestOrientation_CodePassThenContinueScores(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	e, err := New()
	if err != nil {
		t.Fatal(err)
	}
	o := e.StartOrientation("python", "regularly")
	if !driveToCodeItem(o) {
		t.Skip("no code item")
	}
	st := o.Step()
	d := o.diag[o.idx]
	if d.Task == nil || d.Task.Solution == "" {
		t.Skip("code item has no reference solution to submit")
	}
	codingBefore := o.codingOK

	out := o.SubmitCode(d.Task.Solution)
	if !out.Passed || out.Advanced {
		t.Fatalf("reference solution: passed=%v advanced=%v", out.Passed, out.Advanced)
	}
	// Still on the same item until Continue.
	if o.Step().Index != st.Index {
		t.Fatal("pass advanced without Continue")
	}
	o.AdvanceOrientation()
	if d.Measures == "coding" && o.codingOK != codingBefore+1 {
		t.Fatalf("a passed coding item did not score: %d → %d", codingBefore, o.codingOK)
	}
}
