package content

import (
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"

	"devascent/internal/grader"
)

// TestGoBenchInferenceCoverage measures how much of the real 267-problem bench
// the Go function-call harness can handle: for each problem it generates the Go
// starter stub + harness and checks it COMPILES. A high rate means the type
// inference + codegen cover the bench; failures are the problems that need a
// stored signature override (or special handling, e.g. op-list "design" problems).
//
// This is the Go spike's headline metric. Slow (a `go build` per problem) →
// opt-in, -short-skip. Parallelized across CPUs.
func TestGoBenchInferenceCoverage(t *testing.T) {
	if testing.Short() {
		t.Skip("Go bench coverage gate is slow; skipped in -short")
	}
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go not on PATH")
	}
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	type job struct {
		id     string
		driver string
	}
	var jobs []job
	for _, p := range c.Problems {
		if p.FuncName == "" || len(p.Tests) == 0 {
			continue
		}
		stub := grader.GoStarter(p.FuncName, p.Solution, p.Tests, p.GraderShape())
		driver, derr := grader.BuildGoDriver(stub, p.FuncName, p.Tests, p.GraderShape())
		if derr != nil {
			jobs = append(jobs, job{p.ID, ""}) // generation failure counts as not-covered
			continue
		}
		jobs = append(jobs, job{p.ID, driver})
	}

	var mu sync.Mutex
	var fails []string
	var wg sync.WaitGroup
	sem := make(chan struct{}, 6)
	for _, j := range jobs {
		wg.Add(1)
		sem <- struct{}{}
		go func(j job) {
			defer wg.Done()
			defer func() { <-sem }()
			if j.driver == "" || !goCompiles(j.driver) {
				mu.Lock()
				fails = append(fails, j.id)
				mu.Unlock()
			}
		}(j)
	}
	wg.Wait()

	total := len(jobs)
	ok := total - len(fails)
	pct := 0.0
	if total > 0 {
		pct = 100 * float64(ok) / float64(total)
	}
	t.Logf("Go function-call coverage: %d/%d bench problems compile cleanly (%.1f%%)", ok, total, pct)
	const showN = 40
	if len(fails) > 0 {
		shown := fails
		if len(shown) > showN {
			shown = shown[:showN]
		}
		t.Logf("non-compiling (need a stored signature / special handling), first %d of %d: %v", len(shown), len(fails), shown)
	}
}

// goCompiles writes the program to a temp dir and runs `go build`.
func goCompiles(src string) bool {
	dir, err := os.MkdirTemp("", "devascent-gocov-")
	if err != nil {
		return false
	}
	defer os.RemoveAll(dir)
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(src), 0o600); err != nil {
		return false
	}
	cmd := exec.Command("go", "build", "-o", filepath.Join(dir, "out"), "main.go")
	cmd.Dir = dir
	return cmd.Run() == nil
}
