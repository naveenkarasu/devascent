// Package engine holds DevAscent's UI-neutral session logic — the pure
// decision/selection/scoring functions that BOTH the terminal (Bubble Tea) and
// the GUI (Wails) frontends drive. It imports content/grader but never a UI
// toolkit, so a frontend is a thin presentation layer over these functions.
//
// Extracted from internal/tui via the strangler pattern: functions moved here
// verbatim, with same-signature delegates left in the TUI. This file holds the
// content-selection logic (intake, dev set, bench tiering); bench scoring lives
// in bench.go and the language starter dispatch in starter.go.
package engine

import (
	"math/rand"
	"sort"

	"devascent/internal/content"
)

// IntakeSize is the number of items in the Entrance Test intake.
const IntakeSize = 10

// bandDifficulties maps the self-report level to the allowed difficulties.
// never → easy only; a-little → easy+medium; regularly → medium+hard.
func bandDifficulties(level string) map[int]bool {
	switch level {
	case "never":
		return map[int]bool{1: true}
	case "regularly":
		return map[int]bool{2: true, 3: true}
	default: // a-little / unknown
		return map[int]bool{1: true, 2: true}
	}
}

// SelectIntake builds the 10-item intake from items WITHIN the self-report's
// difficulty band (never upward-fallback — "never" stays easy). It aims for a
// coding/machine/spec mix (so routing signals exist), tops up from any in-band
// item to reach 10, orders ascending difficulty, and puts a coding item first
// (the "floor" that routing reads). No repeats. `exclude` (used by redo) pushes
// previously-seen item IDs to the BACK of each measure pool so fresh items are
// preferred — but they still backfill if the band is too small to avoid repeats.
func SelectIntake(pool []content.Diagnostic, level string, exclude map[string]bool, rng *rand.Rand) []content.Diagnostic {
	band := bandDifficulties(level)
	byMeasure := map[string][]content.Diagnostic{}
	for _, d := range pool {
		if band[d.Difficulty] {
			byMeasure[d.Measures] = append(byMeasure[d.Measures], d)
		}
	}
	for m := range byMeasure {
		g := byMeasure[m]
		rng.Shuffle(len(g), func(i, j int) { g[i], g[j] = g[j], g[i] })
		if len(exclude) > 0 { // fresh-first: non-excluded items ahead of excluded
			var fresh, stale []content.Diagnostic
			for _, d := range g {
				if exclude[d.ID] {
					stale = append(stale, d)
				} else {
					fresh = append(fresh, d)
				}
			}
			byMeasure[m] = append(fresh, stale...)
		}
	}
	used := map[string]bool{}
	picked := make([]content.Diagnostic, 0, IntakeSize)
	take := func(m string, k int) {
		for _, d := range byMeasure[m] {
			if len(picked) >= IntakeSize || k <= 0 {
				break
			}
			if used[d.ID] {
				continue
			}
			picked = append(picked, d)
			used[d.ID] = true
			k--
		}
	}
	take("coding", 6)
	take("machine", 3)
	take("spec", 1)
	if len(picked) < IntakeSize { // top up from any in-band item, coding first
		for _, m := range []string{"coding", "machine", "spec"} {
			take(m, IntakeSize)
		}
	}
	sort.SliceStable(picked, func(i, j int) bool { return picked[i].Difficulty < picked[j].Difficulty })
	for i := range picked { // rotate the first coding item to the front (the floor)
		if picked[i].Measures == "coding" {
			c := picked[i]
			copy(picked[1:i+1], picked[:i])
			picked[0] = c
			break
		}
	}
	return picked
}

// SelectDevSet picks up to n dev tasks across DISTINCT categories (random
// category order, random task within each). If fewer than n categories exist,
// it fills from the remaining tasks. Variety without repeats within a run.
func SelectDevSet(all []content.DevTask, n int, rng *rand.Rand) []content.DevTask {
	byCat := map[string][]content.DevTask{}
	var cats []string
	for _, t := range all {
		if _, ok := byCat[t.Category]; !ok {
			cats = append(cats, t.Category)
		}
		byCat[t.Category] = append(byCat[t.Category], t)
	}
	rng.Shuffle(len(cats), func(i, j int) { cats[i], cats[j] = cats[j], cats[i] })
	set := make([]content.DevTask, 0, n)
	used := map[string]bool{}
	for _, c := range cats {
		if len(set) >= n {
			break
		}
		v := byCat[c]
		t := v[rng.Intn(len(v))]
		set = append(set, t)
		used[t.ID] = true
	}
	if len(set) < n { // few categories → fill with random unused tasks
		var rest []content.DevTask
		for _, t := range all {
			if !used[t.ID] {
				rest = append(rest, t)
			}
		}
		rng.Shuffle(len(rest), func(i, j int) { rest[i], rest[j] = rest[j], rest[i] })
		for _, t := range rest {
			if len(set) >= n {
				break
			}
			set = append(set, t)
		}
	}
	return set
}

// DiagsByIDs / DevsByIDs rebuild a selection from saved IDs; they return ok=true
// ONLY if every id resolves, so a partial/stale save triggers a clean re-select.
func DiagsByIDs(all []content.Diagnostic, ids []string) ([]content.Diagnostic, bool) {
	idx := map[string]content.Diagnostic{}
	for _, d := range all {
		idx[d.ID] = d
	}
	out := make([]content.Diagnostic, 0, len(ids))
	for _, id := range ids {
		d, ok := idx[id]
		if !ok {
			return nil, false
		}
		out = append(out, d)
	}
	return out, len(out) > 0
}

func DevsByIDs(all []content.DevTask, ids []string) ([]content.DevTask, bool) {
	idx := map[string]content.DevTask{}
	for _, t := range all {
		idx[t.ID] = t
	}
	out := make([]content.DevTask, 0, len(ids))
	for _, id := range ids {
		t, ok := idx[id]
		if !ok {
			return nil, false
		}
		out = append(out, t)
	}
	return out, len(out) > 0
}

// DedupeBySlug keeps one RANDOM variant per canonical slug (so a curated list
// serves the right count of distinct problems, but the specific phrasing varies
// run-to-run for freshness).
func DedupeBySlug(ps []content.Problem, rng *rand.Rand) []content.Problem {
	bySlug := map[string][]content.Problem{}
	var order []string
	for _, p := range ps {
		s := p.CanonicalSlug()
		if _, ok := bySlug[s]; !ok {
			order = append(order, s)
		}
		bySlug[s] = append(bySlug[s], p)
	}
	out := make([]content.Problem, 0, len(order))
	for _, s := range order {
		v := bySlug[s]
		out = append(out, v[rng.Intn(len(v))])
	}
	return out
}

// FilterProblems returns the subset matching a menu option.
func FilterProblems(problems []content.Problem, kind, value string) []content.Problem {
	if kind == "all" {
		return problems
	}
	var out []content.Problem
	for _, p := range problems {
		switch kind {
		case "category":
			if p.Category == value {
				out = append(out, p)
			}
		case "difficulty":
			if p.Difficulty == value {
				out = append(out, p)
			}
		case "list":
			for _, l := range p.Lists {
				if l == value {
					out = append(out, p)
					break
				}
			}
		}
	}
	return out
}

// tierOrder maps the Step -1 placement to the difficulty order the bench serves.
func tierOrder(placement string) []string {
	switch placement {
	case "test-out":
		return []string{"medium", "hard", "easy"}
	default: // tutorial-full / dev-literacy / unknown → ease in
		return []string{"easy", "medium", "hard"}
	}
}

// SelectBench orders the pool by the player's tier (from placement), shuffled
// within each difficulty — so a stronger player starts harder, everyone ramps,
// no repeats. (Streak-based re-tiering is a later refinement.)
func SelectBench(problems []content.Problem, placement string, rng *rand.Rand) []content.Problem {
	byDiff := map[string][]content.Problem{}
	for _, p := range problems {
		byDiff[p.Difficulty] = append(byDiff[p.Difficulty], p)
	}
	seen := map[string]bool{}
	out := make([]content.Problem, 0, len(problems))
	for _, tier := range tierOrder(placement) {
		g := append([]content.Problem(nil), byDiff[tier]...)
		rng.Shuffle(len(g), func(i, j int) { g[i], g[j] = g[j], g[i] })
		for _, p := range g {
			out = append(out, p)
			seen[p.Difficulty] = true
		}
	}
	for _, p := range problems { // safety: any difficulty not in the tier list
		if !seen[p.Difficulty] {
			out = append(out, p)
		}
	}
	return out
}

func BenchIDs(ps []content.Problem) []string {
	out := make([]string, len(ps))
	for i, p := range ps {
		out[i] = p.ID
	}
	return out
}

func ProblemsByIDs(all []content.Problem, ids []string) ([]content.Problem, bool) {
	idx := map[string]content.Problem{}
	for _, p := range all {
		idx[p.ID] = p
	}
	out := make([]content.Problem, 0, len(ids))
	for _, id := range ids {
		p, ok := idx[id]
		if !ok {
			return nil, false
		}
		out = append(out, p)
	}
	return out, len(out) > 0
}

func DiagIDs(ds []content.Diagnostic) []string {
	out := make([]string, len(ds))
	for i, d := range ds {
		out[i] = d.ID
	}
	return out
}

func DevTaskIDs(ts []content.DevTask) []string {
	out := make([]string, len(ts))
	for i, t := range ts {
		out[i] = t.ID
	}
	return out
}
