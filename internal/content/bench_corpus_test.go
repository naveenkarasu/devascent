package content

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"devascent/internal/grader"
	"devascent/internal/toolchain"
)

// TestBenchSolutionCorpus is the per-language native grading gate: for every
// solution file in testdata/bench/<lang>/<problemID>.<ext>, it grades that
// real, correct solution through the REAL native grader (grader.New) and asserts
// it passes. This is the thing scaffold-compile coverage cannot prove — that a
// correct solution actually grades clean end-to-end (serialization + comparison)
// in each language's own toolchain.
//
// The corpus grows incrementally (Go first, then the other languages); a
// language with no solution files yet is skipped, as is one whose toolchain is
// absent. Slow (compiles/runs every solution) → opt-in, -short-skip.
func TestBenchSolutionCorpus(t *testing.T) {
	if testing.Short() {
		t.Skip("native bench corpus gate is slow; skipped in -short")
	}
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	byID := map[string]Problem{}
	for _, p := range c.Problems {
		byID[p.ID] = p
	}

	exts := map[string]string{
		"go": ".go", "java": ".java", "csharp": ".cs",
		"javascript": ".js", "typescript": ".ts", "rust": ".rs", "python": ".py",
	}
	det := toolchain.New()
	g := grader.New(det)

	for lang, ext := range exts {
		dir := filepath.Join("testdata", "bench", lang)
		entries, err := os.ReadDir(dir)
		if err != nil || len(entries) == 0 {
			continue // no corpus for this language yet
		}
		pass, total := 0, 0
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ext) {
				continue
			}
			id := strings.TrimSuffix(e.Name(), ext)
			p, ok := byID[id]
			if !ok {
				t.Errorf("%s/%s: no bench problem with id %q", lang, e.Name(), id)
				continue
			}
			src, err := os.ReadFile(filepath.Join(dir, e.Name()))
			if err != nil {
				t.Errorf("%s/%s: %v", lang, e.Name(), err)
				continue
			}
			total++
			v, err := g.Run(lang, string(src), p.FuncName, p.Tests, p.GraderShape())
			if err != nil {
				t.Errorf("%s %s: grader error: %v", lang, id, err)
				continue
			}
			if !v.Passed {
				t.Errorf("%s %s: corpus solution did NOT grade clean: err=%q results=%+v", lang, id, v.Err, v.Results)
				continue
			}
			pass++
		}
		if total > 0 {
			t.Logf("%-11s bench corpus: %d/%d solutions grade clean natively", lang, pass, total)
		}
	}
}
