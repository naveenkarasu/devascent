package tui

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"

	"devascent/internal/content"
	"devascent/internal/engine"
)

// The content-selection logic now lives in internal/engine (UI-neutral, shared
// with the GUI). The functions below are thin delegates that preserve the TUI's
// original names/signatures so call sites and tests are unchanged — the strangler
// seam. Presentation-only helpers (benchOption / benchMenuOptions) stay here.

const intakeSize = engine.IntakeSize

func selectIntake(pool []content.Diagnostic, level string, exclude map[string]bool, rng *rand.Rand) []content.Diagnostic {
	return engine.SelectIntake(pool, level, exclude, rng)
}

func selectDevSet(all []content.DevTask, n int, rng *rand.Rand) []content.DevTask {
	return engine.SelectDevSet(all, n, rng)
}

func diagsByIDs(all []content.Diagnostic, ids []string) ([]content.Diagnostic, bool) {
	return engine.DiagsByIDs(all, ids)
}

func devsByIDs(all []content.DevTask, ids []string) ([]content.DevTask, bool) {
	return engine.DevsByIDs(all, ids)
}

func dedupeBySlug(ps []content.Problem, rng *rand.Rand) []content.Problem {
	return engine.DedupeBySlug(ps, rng)
}

func filterProblems(problems []content.Problem, kind, value string) []content.Problem {
	return engine.FilterProblems(problems, kind, value)
}

func selectBench(problems []content.Problem, placement string, rng *rand.Rand) []content.Problem {
	return engine.SelectBench(problems, placement, rng)
}

func benchIDs(ps []content.Problem) []string {
	return engine.BenchIDs(ps)
}

func problemsByIDs(all []content.Problem, ids []string) ([]content.Problem, bool) {
	return engine.ProblemsByIDs(all, ids)
}

func diagIDs(ds []content.Diagnostic) []string {
	return engine.DiagIDs(ds)
}

func devTaskIDs(ts []content.DevTask) []string {
	return engine.DevTaskIDs(ts)
}

// benchOption is one row in the bench browse menu.
type benchOption struct {
	label string
	kind  string // "all" | "category" | "list"
	value string
}

// benchMenuOptions builds the browse menu: All, then each category (with count),
// then any curated list that has members. Presentation-only — stays in the TUI.
func benchMenuOptions(problems []content.Problem) []benchOption {
	catCount := map[string]int{}
	var cats []string
	diffCount := map[string]int{}
	listSlugs := map[string]map[string]bool{} // list -> set of distinct canonical slugs
	for _, p := range problems {
		if _, ok := catCount[p.Category]; !ok {
			cats = append(cats, p.Category)
		}
		catCount[p.Category]++
		diffCount[p.Difficulty]++
		for _, l := range p.Lists {
			if listSlugs[l] == nil {
				listSlugs[l] = map[string]bool{}
			}
			listSlugs[l][p.CanonicalSlug()] = true
		}
	}
	sort.Strings(cats)
	opts := []benchOption{{label: fmt.Sprintf("All problems (%d) — tiered to your level", len(problems)), kind: "all"}}
	for _, c := range cats {
		opts = append(opts, benchOption{label: fmt.Sprintf("%s (%d)", c, catCount[c]), kind: "category", value: c})
	}
	// by difficulty (so hards are reachable — the milestone needs them)
	for _, d := range []string{"easy", "medium", "hard"} {
		if n := diffCount[d]; n > 0 {
			opts = append(opts, benchOption{label: fmt.Sprintf("%s problems only (%d)", strings.Title(d), n), kind: "difficulty", value: d})
		}
	}
	// curated lists — counted by DISTINCT canonical slug (variants don't inflate)
	for _, lst := range []struct {
		key, label string
		total      int
	}{
		{"blind75", "Blind 75", 75},
		{"neetcode150", "NeetCode 150", 150},
	} {
		if n := len(listSlugs[lst.key]); n > 0 {
			opts = append(opts, benchOption{label: fmt.Sprintf("%s (%d/%d)", lst.label, n, lst.total), kind: "list", value: lst.key})
		}
	}
	return opts
}
