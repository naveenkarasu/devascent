package grader

import (
	"os/exec"
	"strings"
	"testing"

	"devascent/internal/toolchain"
)

// The "Off By One" ticket (PXF-101): a 1-indexed paginator. Three fixtures — a
// buggy implementation, the fix, and an emptied file — exercise the repo grader.
const (
	buggyPaginate = `def paginate(items, page, size):
    start = page * size  # BUG: page 1 should start at index 0, not size
    return items[start : start + size]
`
	fixedPaginate = `def paginate(items, page, size):
    start = (page - 1) * size
    return items[start : start + size]
`
	// emptyPaginate has no paginate(): the hidden check's import fails → nonzero
	// exit → FAIL. Proves you cannot pass by deleting the implementation.
	emptyPaginate = "# (player removed the implementation)\n"

	// checkPaginate is the ticket-owned hidden test: plain asserts + nonzero exit
	// on failure (no pytest/unittest discovery, so it runs on a bare interpreter).
	checkPaginate = `from paginate import paginate

items = list(range(1, 11))  # 1..10
assert paginate(items, 1, 3) == [1, 2, 3], paginate(items, 1, 3)
assert paginate(items, 2, 3) == [4, 5, 6], paginate(items, 2, 3)
assert paginate(items, 4, 3) == [10], paginate(items, 4, 3)
print("OK")
`
)

func newLocalForRepoTest(t *testing.T) *LocalToolchain {
	t.Helper()
	det := toolchain.New()
	if det.Presence("python").Status == toolchain.Missing {
		t.Skip("python not installed; skipping repo-grader tests")
	}
	return NewLocalToolchain(det)
}

// TestGradeRepo_OffByOne_TestCommand proves a ticket = repo + test command: the
// hidden check fails on the buggy code, passes once fixed, and cannot be passed
// by emptying the file.
func TestGradeRepo_OffByOne_TestCommand(t *testing.T) {
	g := newLocalForRepoTest(t)
	cmd := []string{"python", "check.py"}

	buggy, err := g.GradeRepo(RepoRequest{
		Files:   map[string]string{"paginate.py": buggyPaginate, "check.py": checkPaginate},
		Command: cmd,
	})
	if err != nil {
		t.Fatal(err)
	}
	if buggy.Passed {
		t.Fatalf("buggy paginator must FAIL the hidden check, got %+v", buggy)
	}

	fixed, err := g.GradeRepo(RepoRequest{
		Files:   map[string]string{"paginate.py": fixedPaginate, "check.py": checkPaginate},
		Command: cmd,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !fixed.Passed {
		t.Fatalf("fixed paginator must PASS, got %+v", fixed)
	}
	if fixed.Stdout != "OK" {
		t.Errorf("expected the check's stdout 'OK' to be surfaced, got %q", fixed.Stdout)
	}

	empty, _ := g.GradeRepo(RepoRequest{
		Files:   map[string]string{"paginate.py": emptyPaginate, "check.py": checkPaginate},
		Command: cmd,
	})
	if empty.Passed {
		t.Fatalf("emptying the implementation must NOT pass (anti-cheat)")
	}
}

// TestGradeRepo_GitState proves the git-state grading primitive: the fix is
// committed and inspectable at HEAD (commit count + the patched line present,
// the buggy line gone) — the basis for git-domain tickets.
func TestGradeRepo_GitState(t *testing.T) {
	g := newLocalForRepoTest(t)
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed; skipping git-state test")
	}
	dir := t.TempDir()

	if err := writeRepo(dir, map[string]string{"paginate.py": buggyPaginate}); err != nil {
		t.Fatal(err)
	}
	mustGit(t, g, dir, "init", "-q")
	mustGit(t, g, dir, "add", ".")
	mustGit(t, g, dir, "commit", "-q", "-m", "initial (buggy)")

	// The player applies and commits the fix.
	if err := writeRepo(dir, map[string]string{"paginate.py": fixedPaginate}); err != nil {
		t.Fatal(err)
	}
	mustGit(t, g, dir, "add", ".")
	mustGit(t, g, dir, "commit", "-q", "-m", "fix off-by-one")

	if n := strings.TrimSpace(g.gitAt(dir, "rev-list", "--count", "HEAD").stdout); n != "2" {
		t.Fatalf("expected 2 commits, got %q", n)
	}
	head := g.gitAt(dir, "show", "HEAD:paginate.py").stdout
	if !strings.Contains(head, "(page - 1) * size") {
		t.Fatalf("fixed line not present at HEAD:\n%s", head)
	}
	if strings.Contains(head, "start = page * size") {
		t.Fatalf("buggy line still present at HEAD")
	}
}

func mustGit(t *testing.T, g *LocalToolchain, dir string, args ...string) {
	t.Helper()
	if o := g.gitAt(dir, args...); o.exit != 0 {
		t.Fatalf("git %v failed (exit %d): %s", args, o.exit, firstNonEmpty(o.stderr, o.stdout))
	}
}
