// Package grader runs player code against tests and returns a structured verdict.
//
// The Grader interface is the seam from ADR-0005 (sidecar grader): the TUI talks
// only to this interface, so the execution backend is swappable. The shipped
// backend MUST be the WASM/wazero sandbox (NFR-1); NativePython is DEV-ONLY.
package grader

import (
	"os"

	"devascent/internal/toolchain"
)

// New returns the active grader backend, selected by DEVASCENT_GRADER:
//
//	"native"  → NativePython   (DEV-ONLY, un-sandboxed host Python; tests/oracle)
//	"wazero"  → WazeroPython   (sandboxed WASM Python; retained during transition)
//	(unset)   → LocalToolchain (ADR-0007 default: BYO — grades via the player's
//	                            own installed toolchain; det supplies the PATH)
//
// ADR-0007 retired the bundled WASM runtime: the default now shells out to the
// player's installed toolchain. WazeroPython stays reachable behind the env var
// during the transition so the 267-gate can dual-run against it. All Python
// backends share BuildPyDriver/ParseHarnessOutput, so grading stays identical.
func New(det *toolchain.Detector) Grader {
	switch os.Getenv("DEVASCENT_GRADER") {
	case "native":
		return NewNativePython()
	case "wazero":
		return NewWazeroPython()
	default:
		return NewLocalToolchain(det)
	}
}

// TestCase: call funcName(*Input) and compare the return to Expected.
type TestCase struct {
	Name     string `json:"name" yaml:"name"`
	Input    []any  `json:"input" yaml:"input"`
	Expected any    `json:"expected" yaml:"expected"`
	Hidden   bool   `json:"hidden" yaml:"hidden,omitempty"`
}

// CaseResult is the outcome of a single test.
type CaseResult struct {
	Name     string
	Passed   bool
	Got      string
	Expected string
	Err      string
}

// Verdict is the result of grading a submission. Same shape the Step-0 grader
// will emit, so Step -1 and Step 0 share one grading vocabulary.
type Verdict struct {
	Passed  bool
	Results []CaseResult
	Err     string // harness/runtime/timeout error (empty on normal completion)
}

// Shape describes how to marshal structured (non-JSON-native) arguments and
// return values for data-structure problems. The ZERO value means plain JSON I/O
// (the only mode Step -1 and most bench problems use) — identical to the old
// behavior. For node problems the harness injects a prelude defining the class
// (ListNode/TreeNode) + __build/__dump, builds flagged args from arrays before
// the call, and dumps a flagged return back to an array for comparison.
type Shape struct {
	Kind     string   // "" | "linkedlist" | "tree" — selects the injected prelude
	ArgKinds []string // per-arg: "node" (build from array) else raw; short/empty => all raw
	RetKind  string   // "node" (dump return to array) else raw
}

// Check selects what a GradeRequest verifies. CheckTests is the function-call
// path (Step 0 + most surfaces); the others are the Stage-2 "compiler-as-oracle"
// modes (ADR-0007 / Runtime Detection design spec).
type Check string

const (
	CheckTests        Check = "tests"         // call funcName against Tests, compare
	CheckCompileError Check = "compile-error" // compile must FAIL with Signal in stderr (e.g. Rust E0382)
	CheckCompiles     Check = "compiles"      // compile must SUCCEED (no run)
	CheckStdout       Check = "stdout"        // run, compare stdout to Signal
	CheckNone         Check = "none"          // reveal-only — never executed
)

// GradeRequest is the general grading request. The zero Check (== "") is treated
// as CheckTests, so existing callers via Run keep their exact semantics.
type GradeRequest struct {
	Lang     string
	Source   string
	Check    Check
	FuncName string     // tests
	Tests    []TestCase // tests
	Shape    Shape      // tests
	Signal   string     // compile-error: expected diagnostic; stdout: expected output
	Stdin    string     // optional
}

// Grader runs player source against a check in the given language.
//
//   - Run is the function-call (Check=tests) path — IDENTICAL signature to before,
//     so every existing bench/diagnostic caller is unchanged.
//   - Grade is the general path (compile-error|compiles|stdout|tests|none) used by
//     the Stage-2 surfaces.
//
// Backends implement Run as a thin wrapper over Grade(GradeRequest{Check:tests}).
type Grader interface {
	Run(lang, source, funcName string, tests []TestCase, shape Shape) (Verdict, error)
	Grade(req GradeRequest) (Verdict, error)
}

// Warmer is implemented by backends with eager startup cost (download + compile)
// worth paying off the interactive path. The TUI calls Warm() in a background
// command at launch so the first real submission hits a warm cache. Backends
// without startup cost (NativePython) simply don't implement it.
type Warmer interface {
	Warm()
}
