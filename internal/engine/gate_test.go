package engine

import (
	"testing"

	"devascent/internal/content"
)

func b75(id, title, cat string) content.Problem {
	return content.Problem{ID: id, Title: title, Category: cat, Lists: []string{"blind75"}}
}

func TestBlind75Progress_CountsBySlugAndWriteup(t *testing.T) {
	probs := []content.Problem{
		b75("p1", "Two Sum", "Arrays"),
		b75("p1b", "Two Sum", "Arrays"), // variant: same slug
		b75("p2", "Valid Parentheses", "Stack"),
		{ID: "p3", Title: "Not In List", Category: "Arrays"}, // untagged → ignored
	}
	solved := map[string]bool{"p1": true, "p1b": true, "p2": true}
	full := map[string]bool{"p1": true} // p1b/p2 solved but no write-up
	g := Blind75Progress(probs, solved, func(id string) bool { return full[id] })

	if g.Full != 1 {
		t.Fatalf("Full = %d, want 1 (variants share a slug)", g.Full)
	}
	if g.Provisional != 1 {
		t.Fatalf("Provisional = %d, want 1 (p2)", g.Provisional)
	}
	if g.Met {
		t.Fatal("gate met with 1 full solve")
	}
	// Two Sum is fully banked → mandatory row done; Valid Parentheses is only provisional.
	var ts, vp bool
	for _, m := range g.Mandatory {
		if m.Slug == "two-sum" {
			ts = m.Done
		}
		if m.Slug == "valid-parentheses" {
			vp = m.Done
		}
	}
	if !ts || vp {
		t.Fatalf("mandatory: two-sum=%v (want true) valid-parentheses=%v (want false)", ts, vp)
	}
}

func TestBlind75Progress_RealCatalog(t *testing.T) {
	cat, err := content.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	g := Blind75Progress(cat.Problems, map[string]bool{}, func(string) bool { return false })

	avail := 0
	requiredSum := 0
	for _, c := range g.Categories {
		avail += c.Available
		requiredSum += c.Required
	}
	if avail < 70 {
		t.Fatalf("blind75 slugs = %d, expected ~74", avail)
	}
	if requiredSum >= Blind75Target {
		t.Fatalf("category minimums sum to %d ≥ target %d — no free choice left", requiredSum, Blind75Target)
	}
	if avail < Blind75Target {
		t.Fatalf("target %d unreachable with %d slugs", Blind75Target, avail)
	}
	// Every mandatory slug must exist in the tagged catalog.
	for _, m := range g.Mandatory {
		if m.Title == "" {
			t.Errorf("mandatory slug %q not found in blind75 catalog", m.Slug)
		}
	}
}

func TestComplexityMCQ_DeterministicAndCorrect(t *testing.T) {
	p := content.Problem{ID: "x-test", TimeComplexity: "O(n)"}
	q1, ok := ComplexityMCQ(p)
	if !ok {
		t.Fatal("MCQ not built")
	}
	q2, _ := ComplexityMCQ(p)
	if len(q1.Options) != 4 {
		t.Fatalf("options = %d, want 4", len(q1.Options))
	}
	if q1.Options[q1.Correct] != "O(n)" {
		t.Fatalf("correct option = %q", q1.Options[q1.Correct])
	}
	for i := range q1.Options {
		if q1.Options[i] != q2.Options[i] {
			t.Fatal("MCQ not deterministic")
		}
	}
	seen := map[string]bool{}
	for _, o := range q1.Options {
		if seen[o] {
			t.Fatalf("duplicate option %q", o)
		}
		seen[o] = true
	}
}

func TestComplexityMCQ_OffChainAndMissing(t *testing.T) {
	q, ok := ComplexityMCQ(content.Problem{ID: "g", TimeComplexity: "O(V+E)"})
	if !ok || q.Options[q.Correct] != "O(V+E)" {
		t.Fatalf("off-chain MCQ: ok=%v q=%+v", ok, q)
	}
	for _, o := range q.Options {
		if o == "O(n)" {
			t.Fatal("O(n) is a near-synonym distractor for O(V+E)")
		}
	}
	if _, ok := ComplexityMCQ(content.Problem{ID: "none"}); ok {
		t.Fatal("MCQ built without authored complexity")
	}
}

func TestComplexityMCQ_EveryBenchProblem(t *testing.T) {
	cat, err := content.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	for _, p := range cat.Problems {
		q, ok := ComplexityMCQ(p)
		if !ok {
			t.Errorf("%s: no MCQ", p.ID)
			continue
		}
		if q.Options[q.Correct] != p.TimeComplexity {
			t.Errorf("%s: correct index points at %q, want %q", p.ID, q.Options[q.Correct], p.TimeComplexity)
		}
	}
}

func TestWriteupTextOK(t *testing.T) {
	if WriteupTextOK("   short  ") {
		t.Fatal("accepted too-short text")
	}
	if WriteupTextOK("aaaaaaaaaaaaaaaaaaaaaaaaaaaa") {
		t.Fatal("accepted single-rune spam")
	}
	if !WriteupTextOK("Used a hash map to track seen values in one pass.") {
		t.Fatal("rejected a normal write-up")
	}
}
