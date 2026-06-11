package guiapi

import "testing"

// Every bench category resolves to a primer in every graded language (the
// Learn drawer must never come up empty), and the browse list carries the
// curated-list tags the filter chips run on.
func TestPrimersAndListTags(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir()) // never touch the real save
	e, err := New()
	if err != nil {
		t.Fatal(err)
	}
	probs := e.Problems("python")
	cats := map[string]bool{}
	blind, nc := 0, 0
	for _, p := range probs {
		if p.Category != "" {
			cats[p.Category] = true
		}
		for _, l := range p.Lists {
			switch l {
			case "blind75":
				blind++
			case "neetcode150":
				nc++
			}
		}
	}
	if blind < 75 || nc < 100 {
		t.Fatalf("list tags missing from browse rows: blind75=%d neetcode150=%d", blind, nc)
	}
	for cat := range cats {
		for _, lang := range GradedLanguages() {
			pv := e.PrimerFor(cat, lang)
			if !pv.Found || len(pv.Sections) == 0 || pv.Summary == "" {
				t.Errorf("primer missing/empty for %q in %s", cat, lang)
			}
		}
	}
}
