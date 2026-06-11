package grader

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"devascent/internal/toolchain"
)

// LocalToolchain grades by shelling out to the player's OWN installed toolchain
// (ADR-0007: DevAscent bundles no runtime). It dispatches per language to a
// langAdapter; the Python adapter reuses BuildPyDriver/ParseHarnessOutput so
// Python grading is byte-identical to the native/wazero backends.
//
// One shared *toolchain.Detector supplies the resolved PATH + exe resolution, so
// grader child processes find the same toolchains the picker detected. Execution
// safety is accident-protection (the player's own code on their own machine):
// context timeout + scoped throwaway workdir + the resolved PATH. (Process-tree
// kill for runaway grandchildren is a documented hardening follow-up; the
// function-call path the bench uses doesn't spawn children.)
type LocalToolchain struct {
	det      *toolchain.Detector
	timeout  time.Duration
	adapters map[string]langAdapter
}

// NewLocalToolchain builds the BYO-toolchain grader. Python grades function-call
// tests (pythonAdapter); the other languages grade the oracle checks (stdout /
// compiles / compile-error) via nativeAdapters.
func NewLocalToolchain(det *toolchain.Detector) *LocalToolchain {
	adapters := map[string]langAdapter{"python": pythonAdapter{}}
	for lang, a := range nativeAdapters() {
		adapters[lang] = a
	}
	return &LocalToolchain{
		det: det,
		// Generous enough for a real compile step (javac/rustc seconds; a cold
		// dotnet restore can be longer) while still bounding runaway player code.
		timeout:  30 * time.Second,
		adapters: adapters,
	}
}

// Run is the function-call (Check=tests) path — unchanged semantics.
func (g *LocalToolchain) Run(lang, source, funcName string, tests []TestCase, shape Shape) (Verdict, error) {
	return g.Grade(GradeRequest{Lang: lang, Source: source, Check: CheckTests, FuncName: funcName, Tests: tests, Shape: shape})
}

// Grade runs the request via the language's adapter in a scoped temp dir.
func (g *LocalToolchain) Grade(req GradeRequest) (Verdict, error) {
	if req.Check == CheckNone {
		return Verdict{Err: "reveal-only (no grading for this exercise)"}, nil
	}
	a, ok := g.adapters[req.Lang]
	if !ok {
		return Verdict{Err: req.Lang + " grading is not available yet (reference-only)"}, nil
	}
	dir, err := os.MkdirTemp("", "devascent-grade-")
	if err != nil {
		return Verdict{}, err
	}
	defer os.RemoveAll(dir)
	return a.grade(req, runner{dir: dir, det: g.det, timeout: g.timeout})
}

// langAdapter encapsulates one language's compile/run/interpret logic. The
// LocalToolchain owns the sandbox (temp dir, timeout, PATH); the adapter owns the
// language specifics.
type langAdapter interface {
	grade(req GradeRequest, r runner) (Verdict, error)
}

// runner is the scoped executor handed to adapters: a throwaway working dir, the
// detector's resolved PATH, a bounded timeout, and exe resolution.
type runner struct {
	dir     string
	det     *toolchain.Detector
	timeout time.Duration
}

type execOut struct {
	stdout   string
	stderr   string
	exit     int
	timedOut bool
}

func (r runner) write(name, content string) error {
	return os.WriteFile(filepath.Join(r.dir, name), []byte(content), 0o600)
}

// resolve returns the absolute path of the first available alternative (e.g.
// "python","python3"), falling back to the first name so exec yields a clear
// "not found" error if truly absent.
func (r runner) resolve(names ...string) string {
	for _, n := range names {
		if p, ok := r.det.Resolve(n); ok {
			return p
		}
	}
	return names[0]
}

// env is the parent environment with PATH overridden to the detector's resolved
// dirs (replacing, not appending — duplicate PATH keys are platform-ambiguous).
func (r runner) env() []string {
	out := make([]string, 0, len(os.Environ())+1)
	for _, e := range os.Environ() {
		if len(e) >= 5 && strings.EqualFold(e[:5], "PATH=") {
			continue
		}
		out = append(out, e)
	}
	return append(out, r.det.PathEnv())
}

// run executes name+args in the scoped dir with a timeout.
func (r runner) run(name string, args ...string) execOut {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = r.dir
	cmd.Env = r.env()
	var so, se bytes.Buffer
	cmd.Stdout = &so
	cmd.Stderr = &se
	err := cmd.Run()
	o := execOut{stdout: so.String(), stderr: se.String()}
	if ctx.Err() == context.DeadlineExceeded {
		o.timedOut, o.exit = true, -1
		return o
	}
	if err == nil {
		return o
	}
	if ee, ok := err.(*exec.ExitError); ok {
		o.exit = ee.ExitCode()
		return o
	}
	o.exit = -1
	if o.stderr == "" {
		o.stderr = err.Error()
	}
	return o
}

// ---- Python adapter -------------------------------------------------------

type pythonAdapter struct{}

func (pythonAdapter) grade(req GradeRequest, r runner) (Verdict, error) {
	py := r.resolve("python", "python3")
	switch req.Check {
	case CheckTests, "":
		driver, err := BuildPyDriver(req.Source, req.FuncName, req.Tests, req.Shape)
		if err != nil {
			return Verdict{}, err
		}
		if err := r.write("run.py", driver); err != nil {
			return Verdict{}, err
		}
		o := r.run(py, "run.py")
		if o.timedOut {
			return Verdict{Err: "time limit exceeded"}, nil
		}
		// ParseHarnessOutput finds the marker on stdout; include stderr so a
		// traceback (no marker) surfaces as the error message (parity with the
		// native backend's CombinedOutput).
		return ParseHarnessOutput(o.stdout+"\n"+o.stderr, req.Tests), nil
	case CheckStdout:
		if err := r.write("main.py", req.Source); err != nil {
			return Verdict{}, err
		}
		o := r.run(py, "main.py")
		if o.timedOut {
			return Verdict{Err: "time limit exceeded"}, nil
		}
		return stdoutVerdict(o, req.Signal), nil
	case CheckCompiles:
		if err := r.write("main.py", req.Source); err != nil {
			return Verdict{}, err
		}
		o := r.run(py, "-c", "import py_compile,sys; py_compile.compile(sys.argv[1], doraise=True)", "main.py")
		return compilesVerdict(o), nil
	case CheckCompileError:
		return Verdict{Err: "compile-error checks are not applicable to Python"}, nil
	default:
		return Verdict{Err: "unsupported check: " + string(req.Check)}, nil
	}
}

// ---- shared check interpreters (reused by future language adapters) --------

// normalizeNL normalizes line endings (Windows CRLF / lone CR → LF) and trims
// surrounding whitespace, so stdout comparison is platform-independent — e.g.
// Java's System.out.println emits \r\n on Windows but the authored signal uses
// \n. Without this, correct multi-line programs fail only on Windows (and only
// for real Windows players, never on Linux CI).
func normalizeNL(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return strings.TrimSpace(s)
}

// stdoutVerdict compares captured stdout to the expected Signal (line-ending
// normalized and trimmed).
func stdoutVerdict(o execOut, want string) Verdict {
	if o.exit != 0 && strings.TrimSpace(o.stdout) == "" {
		return Verdict{Err: firstNonEmpty(firstLine(o.stderr), "program exited "+strconv.Itoa(o.exit))}
	}
	got := normalizeNL(o.stdout)
	want = normalizeNL(want)
	pass := got == want
	return Verdict{Passed: pass, Results: []CaseResult{{Name: "stdout", Passed: pass, Got: got, Expected: want}}}
}

// compilesVerdict passes iff compilation succeeded (exit 0).
func compilesVerdict(o execOut) Verdict {
	if o.exit == 0 {
		return Verdict{Passed: true, Results: []CaseResult{{Name: "compiles", Passed: true}}}
	}
	return Verdict{Passed: false, Results: []CaseResult{{Name: "compiles", Passed: false, Err: firstLine(o.stderr)}}}
}

// compileErrorVerdict passes iff compilation FAILED with the expected diagnostic
// (the trybuild/compiletest "compiler-as-oracle" pattern). Reused by the future
// Rust/C++/C#/Java adapters; no language uses it yet.
func compileErrorVerdict(o execOut, signal string) Verdict {
	combined := o.stdout + "\n" + o.stderr
	if o.exit == 0 {
		return Verdict{Passed: false, Results: []CaseResult{{Name: "compile-error", Passed: false, Err: "expected a compile error" + signalSuffix(signal) + ", but it compiled cleanly"}}}
	}
	if signal != "" && !strings.Contains(combined, signal) {
		return Verdict{Passed: false, Results: []CaseResult{{Name: "compile-error", Passed: false, Expected: signal, Got: firstLine(combined)}}}
	}
	return Verdict{Passed: true, Results: []CaseResult{{Name: "compile-error", Passed: true, Expected: signal}}}
}

func signalSuffix(s string) string {
	if s == "" {
		return ""
	}
	return " (" + s + ")"
}

func firstNonEmpty(a, b string) string {
	if strings.TrimSpace(a) != "" {
		return a
	}
	return b
}

func firstLine(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		s = s[:i]
	}
	if len(s) > 240 {
		s = s[:240] + "…"
	}
	return s
}
