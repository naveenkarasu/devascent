package grader

import (
	"os/exec"
	"testing"
)

func skipIfNoPython(t *testing.T, g NativePython) {
	if _, err := exec.LookPath(g.Exe); err != nil {
		t.Skipf("python (%s) not found on PATH", g.Exe)
	}
}

func TestNativePython_CorrectAndWrong(t *testing.T) {
	g := NewNativePython()
	skipIfNoPython(t, g)

	tests := []TestCase{
		{Name: "t1", Input: []any{2, 3}, Expected: 5},
		{Name: "t2", Input: []any{-1, 1}, Expected: 0},
		{Name: "t3", Input: []any{0, 0}, Expected: 0},
	}

	// Correct implementation passes.
	v, err := g.Run("python", "def add(a, b):\n    return a + b\n", "add", tests, Shape{})
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	if !v.Passed {
		t.Fatalf("expected pass, got %+v", v)
	}

	// Wrong implementation fails (and we still get per-case results).
	v2, _ := g.Run("python", "def add(a, b):\n    return a - b\n", "add", tests, Shape{})
	if v2.Passed {
		t.Fatalf("expected fail for wrong impl, got pass: %+v", v2)
	}
	if len(v2.Results) != 3 {
		t.Fatalf("expected 3 case results, got %d", len(v2.Results))
	}
}

func TestNativePython_RuntimeErrorReported(t *testing.T) {
	g := NewNativePython()
	skipIfNoPython(t, g)

	tests := []TestCase{{Name: "t1", Input: []any{1}, Expected: 1}}
	v, _ := g.Run("python", "def f(x):\n    return x[0]\n", "f", tests, Shape{}) // x[0] on int -> TypeError
	if v.Passed {
		t.Fatalf("expected failure on runtime error")
	}
	if len(v.Results) == 1 && v.Results[0].Err == "" {
		t.Fatalf("expected per-case error captured, got %+v", v.Results)
	}
}
