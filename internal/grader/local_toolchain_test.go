package grader

import (
	"strings"
	"testing"

	"devascent/internal/toolchain"
)

// newLocalForTest builds a LocalToolchain on the real detector, skipping if
// Python isn't installed (these tests shell out to the player's interpreter).
func newLocalForTest(t *testing.T) *LocalToolchain {
	t.Helper()
	det := toolchain.New()
	if det.Presence("python").Status == toolchain.Missing {
		t.Skip("python not installed; skipping LocalToolchain tests")
	}
	return NewLocalToolchain(det)
}

func TestLocalToolchain_TestsPass(t *testing.T) {
	g := newLocalForTest(t)
	tests := []TestCase{
		{Name: "a", Input: []any{2, 3}, Expected: 5},
		{Name: "b", Input: []any{-1, 1}, Expected: 0},
	}
	v, err := g.Run("python", "def add(a, b):\n    return a + b\n", "add", tests, Shape{})
	if err != nil {
		t.Fatal(err)
	}
	if !v.Passed {
		t.Fatalf("add should pass, got %+v", v)
	}
}

func TestLocalToolchain_TestsFail(t *testing.T) {
	g := newLocalForTest(t)
	tests := []TestCase{{Name: "a", Input: []any{2, 3}, Expected: 5}}
	v, _ := g.Run("python", "def add(a, b):\n    return a - b\n", "add", tests, Shape{})
	if v.Passed {
		t.Fatalf("a-b should NOT pass the add tests")
	}
}

func TestLocalToolchain_RuntimeError(t *testing.T) {
	g := newLocalForTest(t)
	tests := []TestCase{{Name: "a", Input: []any{5}, Expected: 5}}
	v, _ := g.Run("python", "def f(x):\n    return x[0]\n", "f", tests, Shape{}) // int is not subscriptable
	if v.Passed {
		t.Fatalf("expected failure from a runtime error")
	}
}

func TestLocalToolchain_Shape_LinkedList(t *testing.T) {
	g := newLocalForTest(t)
	// reverse a linked list built from an array arg, returning a node dumped to array
	src := "def reverse(head):\n    prev = None\n    while head:\n        nxt = head.next\n        head.next = prev\n        prev = head\n        head = nxt\n    return prev\n"
	tests := []TestCase{{Name: "a", Input: []any{[]any{1, 2, 3}}, Expected: []any{3, 2, 1}}}
	v, err := g.Run("python", src, "reverse", tests, Shape{Kind: "linkedlist", ArgKinds: []string{"node"}, RetKind: "node"})
	if err != nil {
		t.Fatal(err)
	}
	if !v.Passed {
		t.Fatalf("linked-list reverse should pass, got %+v", v)
	}
}

func TestLocalToolchain_Stdout(t *testing.T) {
	g := newLocalForTest(t)
	v, err := g.Grade(GradeRequest{
		Lang: "python", Check: CheckStdout,
		Source: "print('hello world')\n", Signal: "hello world",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !v.Passed {
		t.Fatalf("stdout should match, got %+v", v)
	}
	// mismatch
	v2, _ := g.Grade(GradeRequest{Lang: "python", Check: CheckStdout, Source: "print('nope')\n", Signal: "hello world"})
	if v2.Passed {
		t.Fatalf("mismatched stdout should fail")
	}
}

func TestLocalToolchain_Compiles(t *testing.T) {
	g := newLocalForTest(t)
	good, _ := g.Grade(GradeRequest{Lang: "python", Check: CheckCompiles, Source: "x = 1 + 1\n"})
	if !good.Passed {
		t.Fatalf("valid python should compile, got %+v", good)
	}
	bad, _ := g.Grade(GradeRequest{Lang: "python", Check: CheckCompiles, Source: "def (:\n"}) // syntax error
	if bad.Passed {
		t.Fatalf("invalid python should fail to compile")
	}
}

func TestLocalToolchain_UnavailableLanguageIsReferenceOnly(t *testing.T) {
	g := newLocalForTest(t)
	// cpp has no registered adapter (deliberately omitted), so it stays
	// reference-only. (rust/go/etc. now DO have adapters — see native_adapters.go.)
	v, _ := g.Grade(GradeRequest{Lang: "cpp", Check: CheckCompileError, Source: "int main(){}", Signal: "error"})
	if v.Passed {
		t.Fatalf("cpp has no adapter; should not pass")
	}
	if !strings.Contains(v.Err, "reference-only") {
		t.Fatalf("unavailable language should report reference-only, got %q", v.Err)
	}
}

func TestLocalToolchain_NoneCheckIsRevealOnly(t *testing.T) {
	g := newLocalForTest(t)
	v, _ := g.Grade(GradeRequest{Lang: "python", Check: CheckNone, Source: "whatever"})
	if v.Passed || !strings.Contains(v.Err, "reveal-only") {
		t.Fatalf("CheckNone should be reveal-only, got %+v", v)
	}
}

// compileErrorVerdict is exercised directly (no language adapter uses it yet).
func TestCompileErrorVerdict(t *testing.T) {
	pass := compileErrorVerdict(execOut{exit: 1, stderr: "error[E0382]: borrow of moved value"}, "E0382")
	if !pass.Passed {
		t.Fatalf("matching E0382 should pass")
	}
	cleanCompile := compileErrorVerdict(execOut{exit: 0}, "E0382")
	if cleanCompile.Passed {
		t.Fatalf("a clean compile should fail a compile-error check")
	}
	wrongCode := compileErrorVerdict(execOut{exit: 1, stderr: "error[E0499]: cannot borrow"}, "E0382")
	if wrongCode.Passed {
		t.Fatalf("wrong error code should fail")
	}
}
