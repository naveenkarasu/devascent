package engine

import "devascent/internal/content"

// Step 0 completion milestone targets (scaled from the gate spec's shapes —
// count + coverage + difficulty floor; tunable). NOT the literal 49/75 Blind-75
// gate. Moved from internal/tui; the TUI now aliases these.
const (
	Step0BankTarget = 15 // distinct problems banked
	Step0CatTarget  = 6  // distinct categories covered
	Step0HardTarget = 2  // hard problems solved
)

// BenchStats counts distinct banked problems, categories covered, and hards from
// the solved set, resolving each id through the catalog index (ids with no entry
// are skipped, e.g. after a catalog change).
func BenchStats(solvedSet map[string]bool, probByID map[string]content.Problem) (banked, cats, hard int) {
	seenCat := map[string]bool{}
	for id := range solvedSet {
		p, ok := probByID[id]
		if !ok {
			continue
		}
		banked++
		if p.Category != "" && !seenCat[p.Category] {
			seenCat[p.Category] = true
			cats++
		}
		if p.Difficulty == "hard" {
			hard++
		}
	}
	return
}

// Step0Met reports whether the Step 0 milestone (bank + category + hard floors)
// is reached.
func Step0Met(solvedSet map[string]bool, probByID map[string]content.Problem) bool {
	b, c, h := BenchStats(solvedSet, probByID)
	return b >= Step0BankTarget && c >= Step0CatTarget && h >= Step0HardTarget
}

// Step0Profile computes the competency summary shown at completion. Problem-
// Solving and Language Proficiency are computed; Code Quality and Speed need
// write-ups/timing (not built) and are reported as pending. allProblems is the
// full catalog (for category-breadth denominator).
func Step0Profile(solvedSet map[string]bool, probByID map[string]content.Problem, allProblems []content.Problem) (problemSolving, langProf int, track string) {
	b, c, h := BenchStats(solvedSet, probByID)
	totalCats := map[string]bool{}
	for _, p := range allProblems {
		if p.Category != "" {
			totalCats[p.Category] = true
		}
	}
	clamp := func(f float64) float64 {
		if f > 1 {
			return 1
		}
		return f
	}
	bankRatio := clamp(float64(b) / 25.0)
	hardCov := clamp(float64(h) / 5.0)
	catBreadth := 0.0
	if len(totalCats) > 0 {
		catBreadth = float64(c) / float64(len(totalCats))
	}
	problemSolving = int(100 * (0.5*bankRatio + 0.3*hardCov + 0.2*catBreadth))
	langProf = int(clamp(float64(b)/25.0) * 100)
	track = "startup (provisional)"
	if b >= 25 && h >= 4 {
		track = "FAANG (provisional)"
	}
	return
}
