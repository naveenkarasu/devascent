package engine

// Write-up gate (Track A1): after a bench solve passes, the player proves
// understanding with a complexity MCQ + a short free-text approach note. The
// MCQ is graded deterministically against the problem's authored canonical
// complexity; the free text is required but never content-graded (it feeds
// the A4 mentor). Options are stable per problem (hash-seeded shuffle) so a
// resumed run shows the same question.

import (
	"hash/fnv"
	"strings"

	"devascent/internal/content"
)

// complexityChain is the ordered Big-O ladder distractors are drawn from.
var complexityChain = []string{
	"O(1)", "O(log n)", "O(n)", "O(n log n)", "O(n^2)", "O(n^3)", "O(2^n)", "O(n!)",
}

// offChainDistractors supplies distractors for the two ladder entries that
// sit outside the single-variable chain (near-synonyms like O(n^2) for
// O(m*n) are deliberately excluded — a distractor must be wrong, not a
// notation quibble).
var offChainDistractors = map[string][]string{
	"O(V+E)": {"O(n log n)", "O(n^2)", "O(2^n)"},
	"O(m*n)": {"O(n)", "O(n log n)", "O(2^n)"},
}

// MCQ is one deterministic multiple-choice question.
type MCQ struct {
	Question string
	Options  []string
	Correct  int // index into Options
}

// MinWriteupLen is the minimum free-text length for the approach note.
const MinWriteupLen = 20

// ComplexityMCQ builds the write-up question for p. ok is false when the
// problem has no authored complexity (write-up then needs only the text).
func ComplexityMCQ(p content.Problem) (MCQ, bool) {
	correct := p.TimeComplexity
	if correct == "" {
		return MCQ{}, false
	}
	distract := offChainDistractors[correct]
	if distract == nil {
		idx := -1
		for i, c := range complexityChain {
			if c == correct {
				idx = i
				break
			}
		}
		if idx < 0 {
			return MCQ{}, false // not on the ladder; corpus test prevents this
		}
		// Nearest neighbours are the tempting wrong answers: i-1, i+1, i-2, i+2…
		for d := 1; len(distract) < 3 && d < len(complexityChain); d++ {
			if i := idx - d; i >= 0 {
				distract = append(distract, complexityChain[i])
			}
			if i := idx + d; i < len(complexityChain) && len(distract) < 3 {
				distract = append(distract, complexityChain[i])
			}
		}
	}
	opts := append([]string{correct}, distract...)
	// Stable per-problem shuffle (Fisher-Yates driven by FNV of the ID).
	h := fnv.New32a()
	h.Write([]byte(p.ID))
	seed := h.Sum32()
	for i := len(opts) - 1; i > 0; i-- {
		seed = seed*1664525 + 1013904223
		j := int(seed % uint32(i+1))
		opts[i], opts[j] = opts[j], opts[i]
	}
	correctIdx := 0
	for i, o := range opts {
		if o == correct {
			correctIdx = i
			break
		}
	}
	return MCQ{
		Question: "What is the worst-case time complexity of the intended solution?",
		Options:  opts,
		Correct:  correctIdx,
	}, true
}

// WriteupTextOK is the deterministic acceptance rule for the free-text
// approach note: long enough after trimming, and not one repeated rune.
func WriteupTextOK(text string) bool {
	t := strings.TrimSpace(text)
	if len(t) < MinWriteupLen {
		return false
	}
	first := rune(0)
	varied := false
	for _, r := range t {
		if first == 0 {
			first = r
		} else if r != first && r != ' ' {
			varied = true
			break
		}
	}
	return varied
}
