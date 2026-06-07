package toolchain

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// testDetector builds a detector with injected lookup + exec so the probe engine
// can be tested without any real toolchain installed.
func testDetector(look lookFn, run execFn) *Detector {
	return &Detector{
		pathDirs: []string{"/fake/bin"},
		look:     look,
		run:      run,
		cache:    map[string]Probe{},
	}
}

// foundAll resolves every base name to a fake absolute path.
func foundAll(name string) (string, bool) { return "/fake/bin/" + name, true }

// foundExcept resolves everything except the named executables.
func foundExcept(missing ...string) lookFn {
	return func(name string) (string, bool) {
		for _, m := range missing {
			if name == m {
				return "", false
			}
		}
		return "/fake/bin/" + name, true
	}
}

func base(name string) string { return filepath.Base(name) }

func TestPresence_MissingExe(t *testing.T) {
	d := testDetector(foundExcept("rustc"), func(context.Context, string, []string, string, ...string) execResult {
		t.Fatalf("exec should not run when an exe is missing")
		return execResult{}
	})
	p := d.Presence("rust")
	if p.Status != Missing {
		t.Fatalf("rust should be Missing, got %s", p.Status)
	}
	if p.Depth != DepthPresence || !strings.Contains(p.Reason, "rustc") {
		t.Fatalf("bad probe: %+v", p)
	}
}

func TestPresence_FoundParsesVersion(t *testing.T) {
	d := testDetector(foundAll, func(_ context.Context, _ string, _ []string, name string, _ ...string) execResult {
		if base(name) != "python" {
			t.Fatalf("unexpected exe %q", name)
		}
		return execResult{stdout: "Python 3.13.1\n"}
	})
	p := d.Presence("python")
	if p.Status != Available || p.Depth != DepthPresence {
		t.Fatalf("python should be provisionally Available, got %+v", p)
	}
	if p.Version != "3.13.1" {
		t.Fatalf("version = %q, want 3.13.1", p.Version)
	}
}

func TestPresence_JavaVersionFromStderr(t *testing.T) {
	// javac prints its version to stderr — the probe must read both streams.
	d := testDetector(foundAll, func(context.Context, string, []string, string, ...string) execResult {
		return execResult{stderr: "javac 21.0.2\n"}
	})
	if v := d.Presence("java").Version; v != "21.0.2" {
		t.Fatalf("java version = %q, want 21.0.2", v)
	}
}

func TestCapability_Pass(t *testing.T) {
	d := testDetector(foundAll, func(_ context.Context, _ string, _ []string, name string, args ...string) execResult {
		switch base(name) {
		case "python":
			if len(args) > 0 && args[0] == "--version" {
				return execResult{stdout: "Python 3.13.1\n"}
			}
			return execResult{stdout: sentinel + "\n"}
		}
		return execResult{}
	})
	p := d.Capability(context.Background(), "python")
	if p.Status != Available || p.Depth != DepthCapability {
		t.Fatalf("python capability should pass, got %+v", p)
	}
}

func TestCapability_CompileFails(t *testing.T) {
	// Java: javac (compile) exits nonzero → Broken with a compile reason.
	d := testDetector(foundAll, func(_ context.Context, _ string, _ []string, name string, _ ...string) execResult {
		if base(name) == "javac" {
			return execResult{stderr: "Canary.java:1: error: bad\n", exit: 1}
		}
		return execResult{stdout: sentinel}
	})
	p := d.Capability(context.Background(), "java")
	if p.Status != Broken {
		t.Fatalf("java should be Broken, got %s", p.Status)
	}
	if !strings.Contains(p.Reason, "compile") || !strings.Contains(p.Reason, "javac") {
		t.Fatalf("reason should name the compile step: %q", p.Reason)
	}
}

func TestCapability_NoSentinel(t *testing.T) {
	d := testDetector(foundAll, func(context.Context, string, []string, string, ...string) execResult {
		return execResult{stdout: "something else\n"} // ran clean but wrong output
	})
	p := d.Capability(context.Background(), "javascript")
	if p.Status != Broken || !strings.Contains(p.Reason, "expected output") {
		t.Fatalf("missing sentinel should be Broken, got %+v", p)
	}
}

func TestCapability_Timeout(t *testing.T) {
	d := testDetector(foundAll, func(context.Context, string, []string, string, ...string) execResult {
		return execResult{timed: true, exit: -1}
	})
	p := d.Capability(context.Background(), "go")
	if p.Status != Broken || !strings.Contains(p.Reason, "timed out") {
		t.Fatalf("timeout should be Broken, got %+v", p)
	}
}

func TestCapability_MissingShortCircuits(t *testing.T) {
	d := testDetector(foundExcept("go"), func(context.Context, string, []string, string, ...string) execResult {
		t.Fatalf("should not exec a canary when go is missing")
		return execResult{}
	})
	if p := d.Capability(context.Background(), "go"); p.Status != Missing {
		t.Fatalf("missing go should short-circuit to Missing, got %s", p.Status)
	}
}

func TestCache_DepthNotOverwrittenAndInvalidate(t *testing.T) {
	d := testDetector(foundAll, func(_ context.Context, _ string, _ []string, name string, args ...string) execResult {
		if len(args) > 0 && args[0] == "--version" {
			return execResult{stdout: "v1.0.0"}
		}
		return execResult{stdout: sentinel}
	})
	cap := d.Capability(context.Background(), "python")
	if cap.Depth != DepthCapability {
		t.Fatalf("expected capability depth")
	}
	// A later Presence must NOT downgrade the cached capability result.
	d.Presence("python")
	if got := d.Get("python"); got.Depth != DepthCapability {
		t.Fatalf("presence overwrote capability cache: %+v", got)
	}
	d.Invalidate("python")
	if got := d.Get("python"); got.Status != Unknown {
		t.Fatalf("invalidate should clear cache, got %+v", got)
	}
}

func TestLookPathIn_RealFS(t *testing.T) {
	dir := t.TempDir()
	name := "mytool"
	if runtime.GOOS == "windows" {
		name = "mytool.exe"
	}
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	got, ok := lookPathIn([]string{dir}, "mytool")
	if !ok {
		t.Fatalf("mytool should be found in %s", dir)
	}
	if base(got) != name {
		t.Fatalf("resolved %q, want basename %q", got, name)
	}
	if _, ok := lookPathIn([]string{dir}, "ghosttool"); ok {
		t.Fatalf("ghosttool should not be found")
	}
}

func TestMergeDirs_Dedup(t *testing.T) {
	got := mergeDirs([]string{"/a", "/b"}, []string{"/b", "/c", ""})
	want := []string{"/a", "/b", "/c"}
	if strings.Join(got, ":") != strings.Join(want, ":") {
		t.Fatalf("mergeDirs = %v, want %v", got, want)
	}
}

func TestParseVersion(t *testing.T) {
	cases := map[string]string{
		"Python 3.13.1":         "3.13.1",
		"go version go1.26.3":   "1.26.3",
		"javac 21.0.2":          "21.0.2",
		"no version here":       "",
		"rustc 1.87.0 (abc123)": "1.87.0",
	}
	for in, want := range cases {
		if got := parseVersion(in); got != want {
			t.Fatalf("parseVersion(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestLanguagesAndPathEnv(t *testing.T) {
	d := testDetector(foundAll, nil)
	if len(d.Languages()) != len(specs) {
		t.Fatalf("Languages() count mismatch")
	}
	if !strings.HasPrefix(d.PathEnv(), "PATH=") {
		t.Fatalf("PathEnv should start with PATH=, got %q", d.PathEnv())
	}
}

// --- real-toolchain integration (skips when the toolchain isn't installed) ---
// These validate the actual canary command chains end-to-end for languages
// present on the build machine, catching a wrong compile/run chain.

func realCapability(t *testing.T, lang string) {
	t.Helper()
	d := New()
	if d.Presence(lang).Status == Missing {
		t.Skipf("%s not installed; skipping real capability probe", lang)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	p := d.Capability(ctx, lang)
	if p.Status != Available {
		t.Fatalf("%s present but capability not Available: %+v", lang, p)
	}
	if p.Depth != DepthCapability {
		t.Fatalf("%s should be capability-verified, got depth %d", lang, p.Depth)
	}
}

func TestRealCapability_Python(t *testing.T) { realCapability(t, "python") }
func TestRealCapability_Go(t *testing.T)     { realCapability(t, "go") }
func TestRealCapability_Node(t *testing.T)   { realCapability(t, "javascript") }
func TestRealCapability_Rust(t *testing.T)   { realCapability(t, "rust") }
