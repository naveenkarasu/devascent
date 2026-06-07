package content

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"devascent/internal/grader"
)

// resolveWasm returns a path to a local python.wasm for the wazero gate, or ""
// to signal "skip", taken from the DEVASCENT_PYTHON_WASM env var. We never trigger
// the 25MB download from a test — the gate is opt-in and expects the artifact to
// already be on disk.
func resolveWasm() string {
	if p := os.Getenv("DEVASCENT_PYTHON_WASM"); p != "" {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

// TestWazeroBackendMatchesGate is the NFR-1 equivalence gate: it runs the SAME
// 267 bench reference solutions through the sandboxed WazeroPython backend and
// asserts they all pass — proving the wazero backend grades identically to
// NativePython (same BuildPyDriver/ParseHarnessOutput seam, just sandboxed
// execution). Opt-in: needs a local python.wasm (set DEVASCENT_PYTHON_WASM or have
// the spike artifact on disk). Slow (~70s) so it's skipped by default in -short.
func TestWazeroBackendMatchesGate(t *testing.T) {
	if testing.Short() {
		t.Skip("wazero gate is slow; skipped in -short")
	}
	wasm := resolveWasm()
	if wasm == "" {
		t.Skip("python.wasm not found; set DEVASCENT_PYTHON_WASM to run the wazero gate")
	}
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Problems) < 100 {
		t.Fatalf("expected 100+ bench problems, got %d", len(c.Problems))
	}
	g := grader.NewWazeroPython()
	g.WasmPath = wasm
	fails := 0
	for _, p := range c.Problems {
		if p.Solution == "" {
			t.Errorf("%s: no reference solution", p.ID)
			continue
		}
		v, err := g.Run("python", p.Solution, p.FuncName, p.Tests, p.GraderShape())
		if err != nil {
			t.Errorf("%s: wazero grader error: %v", p.ID, err)
			continue
		}
		if !v.Passed {
			fails++
			t.Errorf("%s (%s): reference solution did NOT pass under wazero: %+v", p.ID, p.Difficulty, v.Results)
		}
	}
	if fails == 0 {
		t.Logf("wazero backend graded all %d bench problems clean (NFR-1 equivalence holds)", len(c.Problems))
	}
}

// TestWazeroTimeout is the NFR-1 bounded-execution gate: a CPU-bound infinite
// loop must be PREEMPTED at the deadline (proving WithCloseOnContextDone works),
// not hang the grader. This is the property the correctness gate can't exercise,
// since every reference solution terminates. Uses a short timeout so the test is
// fast; if preemption is broken, this test hangs to the go-test timeout instead
// of returning "time limit exceeded".
func TestWazeroTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("wazero gate is slow; skipped in -short")
	}
	wasm := resolveWasm()
	if wasm == "" {
		t.Skip("python.wasm not found; set DEVASCENT_PYTHON_WASM to run the wazero gate")
	}
	g := grader.NewWazeroPython()
	g.WasmPath = wasm
	g.Timeout = 5 * time.Second
	src := "def spin():\n    while True:\n        pass\n"
	done := make(chan grader.Verdict, 1)
	go func() {
		v, _ := g.Run("python", src, "spin", []grader.TestCase{{Name: "t", Input: []any{}}}, grader.Shape{})
		done <- v
	}()
	select {
	case v := <-done:
		if !strings.Contains(strings.ToLower(v.Err), "time limit") {
			t.Errorf("expected time-limit verdict, got Err=%q passed=%v results=%+v", v.Err, v.Passed, v.Results)
		}
	case <-time.After(30 * time.Second):
		t.Fatal("grader did not preempt infinite loop within 30s — WithCloseOnContextDone not enforcing NFR-1")
	}

	// CRITICAL: the timeout must close only the offending instance, NOT the
	// shared runtime. Grade a normal problem on the SAME grader afterward — if
	// preemption tore down the runtime, every subsequent submission would fail.
	v, err := g.Run("python", "def add(a, b):\n    return a + b\n", "add",
		[]grader.TestCase{{Name: "t", Input: []any{2, 3}, Expected: 5}}, grader.Shape{})
	if err != nil {
		t.Fatalf("grade after timeout errored (shared runtime torn down?): %v", err)
	}
	if !v.Passed {
		t.Fatalf("grade after timeout did not pass (shared runtime damaged?): %+v", v)
	}
}

// TestWazeroDownloadPath exercises the first-run download (fetch from GitHub →
// SHA-256 verify → atomic rename) that every shipped user hits but no other test
// touches. It removes the cached artifact and runs a grade with NO WasmPath set,
// forcing NewWazeroPython() down the real downloadVerify path. Opt-in (network +
// 25MB): set DEVASCENT_WASM_DOWNLOAD=1. A checksum mismatch or bad URL fails here.
func TestWazeroDownloadPath(t *testing.T) {
	if os.Getenv("DEVASCENT_WASM_DOWNLOAD") == "" {
		t.Skip("set DEVASCENT_WASM_DOWNLOAD=1 to exercise the 25MB first-run download")
	}
	base, err := os.UserCacheDir()
	if err != nil {
		t.Fatal(err)
	}
	cached := filepath.Join(base, "DevAscent", "wasm", "python.wasm")
	if err := os.Remove(cached); err != nil && !os.IsNotExist(err) {
		t.Fatalf("clearing cached wasm: %v", err)
	}
	t.Logf("removed %s; forcing fresh download", cached)

	g := grader.NewWazeroPython() // no WasmPath → download-on-first-run
	g.WasmPath = ""
	v, err := g.Run("python", "def add(a, b):\n    return a + b\n", "add",
		[]grader.TestCase{{Name: "t", Input: []any{2, 3}, Expected: 5}}, grader.Shape{})
	if err != nil {
		t.Fatalf("download-path grade errored: %v", err)
	}
	if !v.Passed {
		t.Fatalf("download-path grade did not pass: %+v", v)
	}
	if _, err := os.Stat(cached); err != nil {
		t.Fatalf("expected downloaded artifact at %s: %v", cached, err)
	}
	t.Logf("download + SHA-256 verify + grade OK; artifact at %s", cached)
}

// TestWazeroScaffoldBypassesRejected mirrors TestScaffoldBypassesRejected on the
// sandboxed backend: identity "solutions" for the deep-copy/round-trip scaffolds
// must still FAIL under wazero (the boundary assertions are in the shared driver,
// not the backend, so this confirms the seam carries through).
func TestWazeroScaffoldBypassesRejected(t *testing.T) {
	if testing.Short() {
		t.Skip("wazero gate is slow; skipped in -short")
	}
	wasm := resolveWasm()
	if wasm == "" {
		t.Skip("python.wasm not found; set DEVASCENT_PYTHON_WASM to run the wazero gate")
	}
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	byID := map[string]Problem{}
	for _, p := range c.Problems {
		byID[p.ID] = p
	}
	overrides := map[string]string{
		"nc-encode-decode-strings": "\ndef encode(strs):\n    return strs\ndef decode(s):\n    return s\n",
		"nc-serialize-tree":        "\ndef serialize(root):\n    return root\ndef deserialize(data):\n    return data\n",
		"nc-clone-graph":           "\ndef clone(node):\n    return node\n",
		"nc-copy-random-list":      "\ndef copy_random_list(head):\n    return head\n",
	}
	g := grader.NewWazeroPython()
	g.WasmPath = wasm
	for id, override := range overrides {
		p, ok := byID[id]
		if !ok {
			t.Errorf("%s: problem not found", id)
			continue
		}
		wrong := p.Solution + override
		v, err := g.Run("python", wrong, p.FuncName, p.Tests, p.GraderShape())
		if err != nil {
			t.Errorf("%s: wazero grader error: %v", id, err)
			continue
		}
		if v.Passed {
			t.Errorf("%s: identity-bypass PASSED under wazero — scaffold not enforced", id)
		}
	}
}
