package grader

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"devascent/internal/toolchain"
)

// RepoRequest grades a multi-file working copy ("repo") by running a test command
// inside it — the step up from single-source Grade that the Step-1 ticket board
// needs (a ticket's work is a small repo, not one function). Files maps a
// repo-relative path (forward slashes) to its content; Command is argv, where
// Command[0] is the program (a bare name is resolved against the detector; pass an
// absolute path for tools the detector doesn't track, e.g. git).
//
// SPIKE SCOPE (#52): Files and Command sit at ONE trust level here — fine for
// proving feasibility, but for real tickets the command + hidden tests must be
// TICKET-OWNED and player edits confined to declared paths, otherwise a player
// could set Command to ["true"] and pass everything. Enforcing that ownership
// split is the ticket engine's job (#56), not this primitive's.
type RepoRequest struct {
	Files   map[string]string
	Command []string
	Timeout time.Duration // optional; falls back to the grader's default
}

// GradeRepo materializes Files into a scoped throwaway dir, runs Command, and
// returns the shared Verdict (Passed iff the command exits 0 — the regression-
// test pattern: the hidden suite fails on unpatched code and passes once fixed).
// It reuses the LocalToolchain runner, so it inherits the resolved PATH, the
// spawn.Hide no-console-flash behavior, and the bounded timeout.
func (g *LocalToolchain) GradeRepo(req RepoRequest) (Verdict, error) {
	if len(req.Command) == 0 {
		return Verdict{Err: "no test command"}, nil
	}
	dir, err := os.MkdirTemp("", "devascent-repo-")
	if err != nil {
		return Verdict{}, err
	}
	defer os.RemoveAll(dir)
	if err := writeRepo(dir, req.Files); err != nil {
		return Verdict{}, err
	}
	to := req.Timeout
	if to <= 0 {
		to = g.timeout
	}
	r := runner{dir: dir, det: g.det, timeout: to}
	prog := req.Command[0]
	if !filepath.IsAbs(prog) {
		prog = r.resolve(prog)
	}
	return repoVerdict(r.run(prog, req.Command[1:]...)), nil
}

// writeRepo writes a repo-relative file map into dir, creating parent dirs so a
// nested path like "pkg/x.py" works (runner.write is flat and would error).
func writeRepo(dir string, files map[string]string) error {
	for rel, content := range files {
		p := filepath.Join(dir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
			return err
		}
	}
	return nil
}

// repoVerdict maps a command's exit into the shared Verdict vocabulary: a single
// "tests" CaseResult, the player's stdout kept for debugging, and (on failure)
// the first stderr line as the error — the same shape the bench UI already
// renders, so the ticket workbench (#61) needs no new render path.
func repoVerdict(o execOut) Verdict {
	if o.timedOut {
		return Verdict{Err: "time limit exceeded"}
	}
	pass := o.exit == 0
	v := Verdict{
		Passed:  pass,
		Results: []CaseResult{{Name: "tests", Passed: pass}},
		Stdout:  strings.TrimSpace(o.stdout),
	}
	if !pass {
		v.Err = firstNonEmpty(firstLine(o.stderr), firstLine(o.stdout))
		v.Results[0].Err = v.Err
	}
	return v
}

// gitAt runs git in dir (a persistent repo, unlike GradeRepo's ephemeral dir) for
// the git-state grading primitive — driving AND inspecting state for git-domain
// tickets. git is invoked by ABSOLUTE path because runner.env() replaces PATH with
// the detector's toolchain dirs, which don't include git; and the throwaway repo's
// author identity + no-signing are injected per-invocation (this is an isolated
// temp repo, not the player's real commits).
func (g *LocalToolchain) gitAt(dir string, args ...string) execOut {
	r := runner{dir: dir, det: g.det, timeout: g.timeout}
	full := append([]string{
		"-c", "user.email=ci@devascent.local",
		"-c", "user.name=DevAscent",
		"-c", "commit.gpgsign=false",
	}, args...)
	return r.run(resolveGit(g.det), full...)
}

// resolveGit returns an absolute path to git: the detector first (in case it ever
// tracks git), then the process PATH, then the common Git-for-Windows location,
// finally the bare name so exec yields a clear "not found".
func resolveGit(det *toolchain.Detector) string {
	if det != nil {
		if p, ok := det.Resolve("git"); ok {
			return p
		}
	}
	if p, err := exec.LookPath("git"); err == nil {
		return p
	}
	if p := `C:\Program Files\Git\bin\git.exe`; fileExists(p) {
		return p
	}
	return "git"
}

func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}
