package content

import "testing"

// allowedTime/allowedSpace are the fixed complexity ladders the write-up MCQ
// draws options from. Authoring outside the ladder breaks distractor
// generation, so the corpus is gated here.
var allowedTime = map[string]bool{
	"O(1)": true, "O(log n)": true, "O(n)": true, "O(n log n)": true,
	"O(n^2)": true, "O(n^3)": true, "O(2^n)": true, "O(n!)": true,
	"O(V+E)": true, "O(m*n)": true,
}

var allowedSpace = map[string]bool{
	"O(1)": true, "O(log n)": true, "O(n)": true, "O(n^2)": true,
	"O(V+E)": true, "O(m*n)": true,
}

func TestBenchComplexityMetadataComplete(t *testing.T) {
	cat, err := Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(cat.Problems) == 0 {
		t.Fatal("no bench problems loaded")
	}
	for _, p := range cat.Problems {
		if p.TimeComplexity == "" || p.SpaceComplexity == "" {
			t.Errorf("%s: missing complexity metadata (time=%q space=%q)", p.ID, p.TimeComplexity, p.SpaceComplexity)
			continue
		}
		if !allowedTime[p.TimeComplexity] {
			t.Errorf("%s: time_complexity %q not in the allowed ladder", p.ID, p.TimeComplexity)
		}
		if !allowedSpace[p.SpaceComplexity] {
			t.Errorf("%s: space_complexity %q not in the allowed ladder", p.ID, p.SpaceComplexity)
		}
	}
}
