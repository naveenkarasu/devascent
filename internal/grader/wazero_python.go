package grader

// WazeroPython is the SANDBOXED grader backend (NFR-1): it runs the player's
// Python inside a wazero WASI sandbox (pure Go, no cgo) instead of a host
// subprocess. The wasm module sees ONLY a per-run temp dir — no host FS or
// network. It reuses BuildPyDriver + ParseHarnessOutput, so grading is identical
// to NativePython; only execution differs.
//
// The 25MB CPython-WASI artifact is downloaded on first use (checksum-pinned)
// into the user cache dir, and a persisted compilation cache makes subsequent
// submissions ~250ms (measured; see Spike Results).

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"github.com/tetratelabs/wazero/sys"
)

const (
	// Pinned CPython-WASI artifact (VMware Wasm Language Runtimes, Python 3.12).
	pyWasmURL    = "https://github.com/vmware-labs/webassembly-language-runtimes/releases/download/python%2F3.12.0%2B20231211-040d5a6/python-3.12.0.wasm"
	pyWasmSHA256 = "e5dc5a398b07b54ea8fdb503bf68fb583d533f10ec3f930963e02b9505f7a763"
)

// WazeroPython runs player Python in a wazero WASI sandbox.
//
// The runtime and the (expensive, 25MB) compiled module are built ONCE in
// ensureReady and reused across every Run; each Run only instantiates a fresh,
// isolated module instance with its own temp-dir mount. Recompiling per-Run was
// a ~3s/grade cold cost — compiling once amortizes it to a single startup hit.
type WazeroPython struct {
	Timeout  time.Duration
	WasmPath string // explicit path to python.wasm; if empty, download-on-first-use

	mu       sync.Mutex
	rt       wazero.Runtime
	compiled wazero.CompiledModule
	cache    wazero.CompilationCache
	cacheDir string
}

// NewWazeroPython returns a sandboxed backend. The wasm artifact is resolved
// lazily (downloaded on first Run if not already cached). DEVASCENT_PYTHON_WASM
// overrides the artifact path (used by tests to point at a local file).
func NewWazeroPython() *WazeroPython {
	return &WazeroPython{Timeout: 15 * time.Second, WasmPath: os.Getenv("DEVASCENT_PYTHON_WASM")}
}

func cacheBaseDir() string {
	base, err := os.UserCacheDir()
	if err != nil {
		base = os.TempDir()
	}
	return filepath.Join(base, "DevAscent", "wasm")
}

// ensureReady builds the shared runtime + compiled module ONCE: it loads the
// wasm bytes (downloading + verifying on first use), opens a persisted
// compilation cache, instantiates WASI, and compiles the 25MB module (the
// expensive step). Subsequent Runs reuse rt+compiled and only instantiate.
// Safe for concurrent callers.
func (g *WazeroPython) ensureReady(ctx context.Context) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.compiled != nil {
		return nil
	}
	dir := cacheBaseDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	path := g.WasmPath
	if path == "" {
		path = filepath.Join(dir, "python.wasm")
	}
	if _, err := os.Stat(path); err != nil {
		if g.WasmPath != "" { // explicit path must exist
			return fmt.Errorf("python.wasm not found at %s: %w", path, err)
		}
		if err := downloadVerify(ctx, pyWasmURL, path, pyWasmSHA256); err != nil {
			return fmt.Errorf("download python.wasm: %w", err)
		}
	}
	wasm, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	g.cacheDir = filepath.Join(dir, "wzcache")
	_ = os.MkdirAll(g.cacheDir, 0o755)
	g.cache, _ = wazero.NewCompilationCacheWithDir(g.cacheDir)

	// WithCloseOnContextDone is REQUIRED for NFR-1: without it wazero never
	// interrupts a CPU-bound guest (e.g. `while True: pass`), so a per-call
	// deadline only fires if the module returns on its own. The deadline is
	// supplied per-Run via the context passed to InstantiateModule.
	rt := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig().
		WithCompilationCache(g.cache).
		WithCloseOnContextDone(true))
	wasi_snapshot_preview1.MustInstantiate(ctx, rt)
	compiled, err := rt.CompileModule(ctx, wasm)
	if err != nil {
		rt.Close(ctx)
		return fmt.Errorf("compile python.wasm: %w", err)
	}
	g.rt = rt
	g.compiled = compiled
	return nil
}

func downloadVerify(ctx context.Context, url, dest, wantSHA string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d", resp.StatusCode)
	}
	tmp := dest + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	h := sha256.New()
	if _, err := io.Copy(io.MultiWriter(f, h), resp.Body); err != nil {
		f.Close()
		return err
	}
	f.Close()
	if got := hex.EncodeToString(h.Sum(nil)); got != wantSHA {
		os.Remove(tmp)
		return fmt.Errorf("checksum mismatch: got %s want %s", got, wantSHA)
	}
	return os.Rename(tmp, dest)
}

// Warm pays the one-time cost up front: it ensures python.wasm is downloaded and
// runs one trivial grade so the disk compilation cache is populated. After Warm,
// the player's first real submission hits the ~200ms warm path instead of the
// ~3s cold compile (and never blocks on the 25MB fetch). Errors are swallowed —
// a failed warm just means the first real Run pays the cost as before.
func (g *WazeroPython) Warm() {
	_, _ = g.Run("python", "def __warm():\n    return 0\n", "__warm",
		[]TestCase{{Name: "warm", Input: []any{}, Expected: 0}}, Shape{})
}

func (g *WazeroPython) Run(lang, source, funcName string, tests []TestCase, shape Shape) (Verdict, error) {
	if lang != "python" {
		return Verdict{Err: "unsupported language: " + lang}, nil
	}
	ctx := context.Background()
	if err := g.ensureReady(ctx); err != nil {
		return Verdict{}, err
	}
	driver, err := BuildPyDriver(source, funcName, tests, shape)
	if err != nil {
		return Verdict{}, err
	}

	dir, err := os.MkdirTemp("", "devascent-wz-")
	if err != nil {
		return Verdict{}, err
	}
	defer os.RemoveAll(dir)
	if err := os.WriteFile(filepath.Join(dir, "run.py"), []byte(driver), 0o600); err != nil {
		return Verdict{}, err
	}

	rctx, cancel := context.WithTimeout(ctx, g.Timeout)
	defer cancel()

	var stdout, stderr writeBuf
	// Sandbox: a fresh anonymous instance mounting ONLY this run's temp dir (no
	// host FS, no network). The runtime + compiled module are shared (built once
	// in ensureReady); each Run is isolated by its own mount + module instance.
	modCfg := wazero.NewModuleConfig().
		WithName("").
		WithArgs("python", "/app/run.py").
		WithStdout(&stdout).WithStderr(&stderr).
		WithFSConfig(wazero.NewFSConfig().WithDirMount(dir, "/app"))

	mod, err := g.rt.InstantiateModule(rctx, g.compiled, modCfg)
	if mod != nil {
		mod.Close(rctx)
	}
	if rctx.Err() == context.DeadlineExceeded {
		return Verdict{Err: "time limit exceeded"}, nil
	}
	if se, ok := err.(*sys.ExitError); ok && se.ExitCode() != 0 {
		// non-zero exit: fall through; stderr captured in output
	}
	out := stdout.String()
	if out == "" {
		out = stderr.String()
	}
	return ParseHarnessOutput(out, tests), nil
}

// Grade satisfies the Grader interface. The sandboxed WASM Python backend
// supports only the tests check; under ADR-0007 it is being retired in favor of
// LocalToolchain, which owns the compile-error/compiles/stdout modes.
func (g *WazeroPython) Grade(req GradeRequest) (Verdict, error) {
	switch req.Check {
	case CheckTests, "":
		return g.Run(req.Lang, req.Source, req.FuncName, req.Tests, req.Shape)
	default:
		return Verdict{Err: "WazeroPython supports only the tests check, got: " + string(req.Check)}, nil
	}
}

// Close releases the shared runtime. Optional — the OS reclaims everything on
// exit — but lets long-lived hosts free the compiled module explicitly.
func (g *WazeroPython) Close() error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.rt != nil {
		return g.rt.Close(context.Background())
	}
	return nil
}

// writeBuf is a tiny io.Writer accumulator (avoids importing bytes here).
type writeBuf struct{ b []byte }

func (w *writeBuf) Write(p []byte) (int, error) { w.b = append(w.b, p...); return len(p), nil }
func (w *writeBuf) String() string              { return string(w.b) }
