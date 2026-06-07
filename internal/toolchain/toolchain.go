// Package toolchain detects whether a language's runtime/compiler is installed
// and actually works on THIS machine (ADR-0007: DevAscent bundles no runtimes; the
// player brings their own). Detection drives the capability gate in the TUI:
// available → the language is selectable; missing/broken → it routes to the
// install guide.
//
// Two-phase probe (see the design spec):
//   - Presence:   fast — locate the executable(s) + parse --version (~ms).
//   - Capability: deep — write a tiny canary and run the REAL compile+run
//     pipeline the exercises use, requiring the sentinel on stdout (~ms–seconds).
//
// A language is authoritatively Available only after Capability passes; Presence
// gives a provisional answer so the UI can render instantly.
package toolchain

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

// sentinel is what every capability canary must print on stdout to prove the
// full compile+run pipeline works.
const sentinel = "DEVASCENT_OK"

// Status is a language's detected state on this machine.
type Status int

const (
	Unknown   Status = iota // not probed yet
	Available               // toolchain present (provisional after Presence, authoritative after Capability)
	Missing                 // a required executable was not found — definitive
	Broken                  // present but the canary failed (e.g. JRE without JDK, broken linker)
)

func (s Status) String() string {
	switch s {
	case Available:
		return "available"
	case Missing:
		return "missing"
	case Broken:
		return "broken"
	default:
		return "unknown"
	}
}

// Depth records how far the probe got, so the UI can show "…verifying" for a
// presence-only result vs trust a capability-verified one.
type Depth int

const (
	DepthNone Depth = iota
	DepthPresence
	DepthCapability
)

// Probe is the result of detecting one language.
type Probe struct {
	Lang    string
	Status  Status
	Version string // best-effort parsed version ("3.13.1"); may be empty
	Reason  string // why Missing/Broken — shown on the install screen
	Depth   Depth
}

// OK reports whether the language is usable (authoritatively or provisionally).
func (p Probe) OK() bool { return p.Status == Available }

// execResult is the outcome of running one command.
type execResult struct {
	stdout string
	stderr string
	exit   int  // process exit code; -1 if it failed to start
	timed  bool // true if killed by the context deadline
}

// execFn runs name+args in dir with env and returns the result. Injected so the
// detector is unit-testable without real toolchains.
type execFn func(ctx context.Context, dir string, env []string, name string, args ...string) execResult

// lookFn resolves an executable base name to an absolute path. Injected for tests.
type lookFn func(name string) (string, bool)

// Detector probes languages, resolving the augmented PATH once and caching
// results. Safe for concurrent use.
type Detector struct {
	pathDirs []string // resolved PATH (with login-shell / known-dir augmentation)
	look     lookFn
	run      execFn

	stub map[string]Probe // non-nil in NewStub: canned results, no real probing

	mu    sync.Mutex
	cache map[string]Probe
}

// NewStub builds a Detector that returns canned probe results and performs NO
// real probing — for tests in other packages (e.g. the TUI) that need
// deterministic toolchain availability without depending on what's installed.
// Seeded languages return their Probe; any other language reports Missing.
func NewStub(seed map[string]Probe) *Detector {
	return &Detector{stub: seed, cache: map[string]Probe{}}
}

// stubProbe returns the canned probe for a language when the detector is a stub.
func (d *Detector) stubProbe(lang string) (Probe, bool) {
	if d.stub == nil {
		return Probe{}, false
	}
	if p, ok := d.stub[lang]; ok {
		p.Lang = lang
		return p, true
	}
	return Probe{Lang: lang, Status: Missing, Depth: DepthPresence, Reason: "not installed"}, true
}

// New builds a Detector with the real PATH resolution + exec, ready for use.
func New() *Detector {
	dirs := defaultPathDirs()
	d := &Detector{pathDirs: dirs, cache: map[string]Probe{}}
	d.look = func(name string) (string, bool) { return lookPathIn(dirs, name) }
	d.run = realExec
	return d
}

// PathEnv returns the resolved PATH as a single PATH=… entry, for the grader to
// hand to its child processes so they find the same toolchains the picker did.
func (d *Detector) PathEnv() string {
	return "PATH=" + strings.Join(d.pathDirs, string(os.PathListSeparator))
}

// Resolve returns the absolute path of an executable using the detector's
// augmented PATH (the grader must resolve exe names itself: exec.LookPath uses
// the process PATH, not a child's custom Env). Returns ("", false) if not found.
func (d *Detector) Resolve(name string) (string, bool) {
	if d.look == nil {
		return "", false
	}
	return d.look(name)
}

// Languages returns the language keys the detector knows how to probe.
func (d *Detector) Languages() []string {
	out := make([]string, len(specs))
	for i, s := range specs {
		out[i] = s.lang
	}
	return out
}

// Get returns the best cached probe for a language (Unknown if never probed).
func (d *Detector) Get(lang string) Probe {
	if p, ok := d.stubProbe(lang); ok {
		return p
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if p, ok := d.cache[lang]; ok {
		return p
	}
	return Probe{Lang: lang, Status: Unknown}
}

// Invalidate drops the cached result for a language (the "I just installed it"
// re-check path).
func (d *Detector) Invalidate(lang string) {
	d.mu.Lock()
	delete(d.cache, lang)
	d.mu.Unlock()
}

// Presence is the fast probe: are the required executables on PATH, and what
// version? It returns Missing definitively, or a provisional Available
// (Depth=Presence) that Capability later confirms or downgrades to Broken.
func (d *Detector) Presence(lang string) Probe {
	if p, ok := d.stubProbe(lang); ok {
		return p
	}
	sp, ok := specFor(lang)
	if !ok {
		return Probe{Lang: lang, Status: Unknown, Reason: "unknown language"}
	}
	for _, group := range sp.exes {
		if _, found := d.resolveGroup(group); !found {
			p := Probe{
				Lang:   lang,
				Status: Missing,
				Depth:  DepthPresence,
				Reason: missingReason(group),
			}
			d.store(p)
			return p
		}
	}
	// All required executables found — parse a version from the primary one.
	primary, _ := d.resolveGroup(sp.exes[0])
	ver := ""
	if len(sp.versionArgs) > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		res := d.run(ctx, "", d.env(), primary, sp.versionArgs...)
		cancel()
		ver = parseVersion(res.stdout + "\n" + res.stderr) // javac prints to stderr
	}
	p := Probe{Lang: lang, Status: Available, Version: ver, Depth: DepthPresence}
	d.store(p)
	return p
}

// Capability is the deep probe: run the real compile+run canary. This is the
// authoritative "available." Cached; call Invalidate to force a re-run.
func (d *Detector) Capability(ctx context.Context, lang string) Probe {
	if p, ok := d.stubProbe(lang); ok {
		return p
	}
	if cached := d.Get(lang); cached.Depth == DepthCapability {
		return cached
	}
	sp, ok := specFor(lang)
	if !ok {
		return Probe{Lang: lang, Status: Unknown, Reason: "unknown language"}
	}
	// Presence first — never run a canary for a toolchain we can't even find.
	if pres := d.Presence(lang); pres.Status == Missing {
		return pres
	}

	dir, err := os.MkdirTemp("", "devascent-probe-"+lang+"-")
	if err != nil {
		return Probe{Lang: lang, Status: Broken, Depth: DepthCapability, Reason: "could not create probe dir: " + err.Error()}
	}
	defer os.RemoveAll(dir)

	cs := sp.canary(runtime.GOOS, dir, d)
	for name, content := range cs.files {
		if werr := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600); werr != nil {
			return Probe{Lang: lang, Status: Broken, Depth: DepthCapability, Reason: "could not write canary: " + werr.Error()}
		}
	}

	env := d.env()
	var last execResult
	for i, st := range cs.steps {
		exe := st.name
		if !isPathy(exe) {
			if resolved, found := d.look(exe); found {
				exe = resolved
			}
		}
		last = d.run(ctx, dir, env, exe, st.args...)
		if last.timed {
			return d.broken(lang, "timed out running "+st.name)
		}
		if last.exit != 0 {
			stage := "compile"
			if i == len(cs.steps)-1 {
				stage = "run"
			}
			return d.broken(lang, stage+" step failed ("+st.name+"): "+firstLine(last.stderr+last.stdout))
		}
	}
	if sp.verify != nil {
		if ok, reason := sp.verify(last); !ok {
			return d.broken(lang, reason)
		}
	} else if !strings.Contains(last.stdout, sentinel) {
		return d.broken(lang, "canary ran but did not produce the expected output")
	}
	// Keep the version from Presence if we have it.
	ver := d.Get(lang).Version
	p := Probe{Lang: lang, Status: Available, Version: ver, Depth: DepthCapability}
	d.store(p)
	return p
}

func (d *Detector) broken(lang, reason string) Probe {
	p := Probe{Lang: lang, Status: Broken, Depth: DepthCapability, Reason: reason}
	d.store(p)
	return p
}

// store caches a probe, never letting a shallower probe overwrite a deeper one.
func (d *Detector) store(p Probe) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if cur, ok := d.cache[p.Lang]; ok && cur.Depth > p.Depth {
		return
	}
	d.cache[p.Lang] = p
}

// resolveGroup finds the first available alternative in a group (e.g.
// {"g++","clang++"}) and returns its absolute path.
func (d *Detector) resolveGroup(group []string) (string, bool) {
	for _, name := range group {
		if path, ok := d.look(name); ok {
			return path, true
		}
	}
	return "", false
}

// env builds the child-process environment: inherit the parent's and override
// only PATH with the resolved dirs. We deliberately do NOT scope HOME/cache dirs
// for the detection canary: the workdir is already isolated (cmd.Dir is the
// throwaway temp dir, and compiled outputs go there via explicit -o paths), and
// overriding HOME/CARGO_HOME breaks toolchains that read their install/config
// from there — notably rustup, which fails with "could not choose a version of
// rustc" when it can't find its toolchain store. A hello-world canary doesn't
// meaningfully pollute the real caches, and reusing them keeps probes fast and
// makes them behave exactly as real play will. (Cache scoping for actual GRADING
// runs is a separate, per-language decision — and must preserve RUSTUP_HOME.)
func (d *Detector) env() []string {
	env := append([]string{}, os.Environ()...)
	return setEnv(env, "PATH", strings.Join(d.pathDirs, string(os.PathListSeparator)))
}

// ---- language probe specs -------------------------------------------------

// canarySpec is the files + ordered command steps for a capability probe. The
// last step's stdout must contain the sentinel.
type canarySpec struct {
	files map[string]string
	steps []step
}

type step struct {
	name string // exe (base name → resolved via PATH, or an absolute path)
	args []string
}

// langSpec describes how to detect one language.
type langSpec struct {
	lang        string
	label       string
	exes        [][]string // presence: each group must be satisfied; alternatives within a group
	versionArgs []string   // args to read a version from the primary exe
	canary      func(goos, dir string, d *Detector) canarySpec
	// verify, if non-nil, validates the final step's output INSTEAD of the default
	// "stdout contains the sentinel" check — for canaries that prove capability a
	// different way (e.g. C# checks `dotnet --list-sdks` lists ≥1 SDK).
	verify func(last execResult) (ok bool, reason string)
}

// bin returns the platform-correct freshly-built binary name/path in dir.
func bin(goos, dir, base string) string {
	if goos == "windows" {
		base += ".exe"
	}
	return filepath.Join(dir, base)
}

var specs = []langSpec{
	{
		lang: "python", label: "Python",
		exes: [][]string{{"python", "python3"}}, versionArgs: []string{"--version"},
		canary: func(goos, dir string, d *Detector) canarySpec {
			py, _ := d.resolveGroup([]string{"python", "python3"})
			if py == "" {
				py = "python"
			}
			return canarySpec{
				files: map[string]string{"canary.py": "import sys\nprint('" + sentinel + "' if sys.version_info[0] == 3 else 'OLD')\n"},
				steps: []step{{py, []string{"canary.py"}}},
			}
		},
	},
	{
		lang: "javascript", label: "JavaScript",
		exes: [][]string{{"node"}}, versionArgs: []string{"--version"},
		canary: func(goos, dir string, d *Detector) canarySpec {
			return canarySpec{
				files: map[string]string{"canary.js": "console.log('" + sentinel + "')\n"},
				steps: []step{{"node", []string{"canary.js"}}},
			}
		},
	},
	{
		lang: "typescript", label: "TypeScript",
		exes: [][]string{{"tsc"}, {"node"}}, versionArgs: []string{"--version"},
		canary: func(goos, dir string, d *Detector) canarySpec {
			return canarySpec{
				files: map[string]string{"canary.ts": "const m: string = '" + sentinel + "'\nconsole.log(m)\n"},
				steps: []step{
					{"tsc", []string{"canary.ts"}}, // emits canary.js in dir
					{"node", []string{"canary.js"}},
				},
			}
		},
	},
	{
		lang: "go", label: "Go",
		exes: [][]string{{"go"}}, versionArgs: []string{"version"},
		canary: func(goos, dir string, d *Detector) canarySpec {
			return canarySpec{
				files: map[string]string{"canary.go": "package main\nimport \"fmt\"\nfunc main() { fmt.Println(\"" + sentinel + "\") }\n"},
				steps: []step{{"go", []string{"run", "canary.go"}}},
			}
		},
	},
	{
		lang: "java", label: "Java",
		exes: [][]string{{"javac"}, {"java"}}, versionArgs: []string{"-version"},
		canary: func(goos, dir string, d *Detector) canarySpec {
			return canarySpec{
				files: map[string]string{"Canary.java": "public class Canary { public static void main(String[] a) { System.out.println(\"" + sentinel + "\"); } }\n"},
				steps: []step{
					{"javac", []string{"Canary.java"}},
					{"java", []string{"Canary"}},
				},
			}
		},
	},
	{
		lang: "csharp", label: "C#",
		exes: [][]string{{"dotnet"}}, versionArgs: []string{"--version"},
		// NOTE: a full compile+run canary needs an SDK-version-matched TargetFramework,
		// which is fragile to pin blind. C# is reveal-only (Axis B) for a long while,
		// so we verify the meaningful adequacy question instead: is an SDK installed
		// (not just the runtime)? `dotnet --list-sdks` lists ≥1 line iff an SDK exists.
		// Replace with a real project canary when C# grading is authored.
		canary: func(goos, dir string, d *Detector) canarySpec {
			return canarySpec{
				steps: []step{{"dotnet", []string{"--list-sdks"}}},
			}
		},
		// `dotnet --list-sdks` prints one line per installed SDK (empty if only the
		// runtime is present). A non-empty list = a usable C# compiler.
		verify: func(last execResult) (bool, string) {
			if strings.TrimSpace(last.stdout) == "" {
				return false, "the .NET runtime is present but no SDK was found — install the .NET SDK (not just the runtime) to compile C#"
			}
			return true, ""
		},
	},
	{
		lang: "cpp", label: "C++",
		exes: [][]string{{"g++", "clang++"}}, versionArgs: []string{"--version"},
		canary: func(goos, dir string, d *Detector) canarySpec {
			cxx, _ := d.resolveGroup([]string{"g++", "clang++"})
			if cxx == "" {
				cxx = "g++"
			}
			out := bin(goos, dir, "canary")
			return canarySpec{
				files: map[string]string{"canary.cpp": "#include <iostream>\nint main() { std::cout << \"" + sentinel + "\" << std::endl; return 0; }\n"},
				steps: []step{
					{cxx, []string{"canary.cpp", "-o", out}},
					{out, nil},
				},
			}
		},
	},
	{
		lang: "rust", label: "Rust",
		exes: [][]string{{"rustc"}}, versionArgs: []string{"--version"},
		canary: func(goos, dir string, d *Detector) canarySpec {
			out := bin(goos, dir, "canary")
			return canarySpec{
				files: map[string]string{"canary.rs": "fn main() { println!(\"" + sentinel + "\"); }\n"},
				steps: []step{
					{"rustc", []string{"canary.rs", "-o", out}},
					{out, nil},
				},
			}
		},
	},
}

func specFor(lang string) (langSpec, bool) {
	for _, s := range specs {
		if s.lang == lang {
			return s, true
		}
	}
	return langSpec{}, false
}

// ---- helpers --------------------------------------------------------------

var versionRe = regexp.MustCompile(`\d+\.\d+(?:\.\d+)?`)

func parseVersion(s string) string { return versionRe.FindString(s) }

func missingReason(group []string) string {
	if len(group) == 1 {
		return group[0] + " was not found on your PATH"
	}
	return "none of [" + strings.Join(group, ", ") + "] was found on your PATH"
}

func firstLine(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		s = s[:i]
	}
	if len(s) > 200 {
		s = s[:200] + "…"
	}
	if s == "" {
		return "no error output"
	}
	return s
}

// isPathy reports whether name is an absolute path or contains a separator (so
// it should be used directly rather than resolved via PATH).
func isPathy(name string) bool {
	return filepath.IsAbs(name) || strings.ContainsAny(name, `/\`)
}

func setEnv(env []string, key, val string) []string {
	prefix := key + "="
	for i, e := range env {
		if strings.HasPrefix(e, prefix) || (runtime.GOOS == "windows" && strings.HasPrefix(strings.ToUpper(e), strings.ToUpper(prefix))) {
			env[i] = prefix + val
			return env
		}
	}
	return append(env, prefix+val)
}

// realExec runs a command with output capture, honoring the context deadline.
func realExec(ctx context.Context, dir string, env []string, name string, args ...string) execResult {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	cmd.Env = env
	var so, se bytes.Buffer
	cmd.Stdout = &so
	cmd.Stderr = &se
	err := cmd.Run()
	res := execResult{stdout: so.String(), stderr: se.String()}
	if ctx.Err() == context.DeadlineExceeded {
		res.timed = true
		res.exit = -1
		return res
	}
	if err == nil {
		res.exit = 0
		return res
	}
	if ee, ok := err.(*exec.ExitError); ok {
		res.exit = ee.ExitCode()
		return res
	}
	// Failed to start (e.g. exe not found).
	res.exit = -1
	if res.stderr == "" {
		res.stderr = err.Error()
	}
	return res
}

// lookPathIn resolves name against the given dirs (honoring PATHEXT on Windows).
func lookPathIn(dirs []string, name string) (string, bool) {
	if isPathy(name) {
		if isExecutable(name) {
			return name, true
		}
		return "", false
	}
	exts := []string{""}
	if runtime.GOOS == "windows" {
		exts = windowsExts()
	}
	for _, d := range dirs {
		if d == "" {
			continue
		}
		for _, ext := range exts {
			p := filepath.Join(d, name+ext)
			if isExecutable(p) {
				return p, true
			}
		}
	}
	return "", false
}

func windowsExts() []string {
	pe := os.Getenv("PATHEXT")
	if pe == "" {
		pe = ".COM;.EXE;.BAT;.CMD"
	}
	parts := strings.Split(pe, ";")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, strings.ToLower(p))
		}
	}
	return out
}

func isExecutable(path string) bool {
	fi, err := os.Stat(path)
	if err != nil || fi.IsDir() {
		return false
	}
	if runtime.GOOS == "windows" {
		return true // extension already filtered by PATHEXT
	}
	return fi.Mode()&0o111 != 0
}

func splitPath(p string) []string {
	if p == "" {
		return nil
	}
	return strings.Split(p, string(os.PathListSeparator))
}

// mergeDirs appends extras not already present, preserving order.
func mergeDirs(base, extra []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(base)+len(extra))
	for _, list := range [][]string{base, extra} {
		for _, d := range list {
			if d == "" || seen[d] {
				continue
			}
			seen[d] = true
			out = append(out, d)
		}
	}
	return out
}

// defaultPathDirs resolves the augmented PATH. A GUI/TUI launched from a desktop
// shortcut does NOT inherit the login-shell PATH (Homebrew, rustup, nvm, …), so
// on unix we capture it from the login shell; on Windows we add known install
// dirs. (The stale-PATH-after-install case is handled by [R] re-check / relaunch.)
func defaultPathDirs() []string {
	dirs := splitPath(os.Getenv("PATH"))
	if runtime.GOOS == "windows" {
		return mergeDirs(dirs, windowsKnownDirs())
	}
	if lp, err := loginShellPath(); err == nil {
		dirs = mergeDirs(dirs, splitPath(lp))
	}
	return dirs
}

// loginShellPath asks the user's login shell for its PATH (the one with their
// toolchain dirs). Bounded by a short timeout; failure is non-fatal.
func loginShellPath() (string, error) {
	sh := os.Getenv("SHELL")
	if sh == "" {
		sh = "/bin/sh"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, sh, "-lc", `printf %s "$PATH"`).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// windowsKnownDirs returns common per-user install locations to augment PATH
// (existence-filtered), covering toolchains installed after the session began.
func windowsKnownDirs() []string {
	var c []string
	add := func(p string) {
		if p != "" {
			c = append(c, p)
		}
	}
	home := os.Getenv("USERPROFILE")
	local := os.Getenv("LOCALAPPDATA")
	pf := os.Getenv("ProgramFiles")
	if home != "" {
		add(filepath.Join(home, ".cargo", "bin"))
		add(filepath.Join(home, "go", "bin"))
	}
	if local != "" {
		add(filepath.Join(local, "Programs", "Python"))
		add(filepath.Join(local, "Microsoft", "dotnet"))
	}
	if pf != "" {
		add(filepath.Join(pf, "Go", "bin"))
		add(filepath.Join(pf, "dotnet"))
		add(filepath.Join(pf, "nodejs"))
	}
	var out []string
	for _, p := range c {
		if fi, err := os.Stat(p); err == nil && fi.IsDir() {
			out = append(out, p)
		}
	}
	return out
}
