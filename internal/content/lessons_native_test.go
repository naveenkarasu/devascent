package content

import (
	"context"
	"fmt"
	"sort"
	"testing"
	"time"

	"devascent/internal/grader"
	"devascent/internal/toolchain"
)

// TestLessonsGradeNativeRoundTrip is the Tutorial-Island analogue of the
// Advanced-Topics native gate: for every lesson task that ships a reference
// Solution, it grades that solution in the lesson's language through the player's
// REAL toolchain and asserts it passes. This proves the per-language lessons
// (their Tests, in particular) actually grade green in that language — catching
// the static-typing pitfalls (heterogeneous test types, null/None returns,
// map-as-argument) that "it compiles" alone would miss.
//
// Python lessons carry no Solution (the original is already proven) and are
// skipped. Unavailable toolchains are reported and skipped, never failed. Slow
// (a compile/run per task) → opt-in, -short-skip.
func TestLessonsGradeNativeRoundTrip(t *testing.T) {
	if testing.Short() {
		t.Skip("native lesson round-trip is slow; skipped in -short")
	}
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	det := toolchain.New()
	g := grader.NewLocalToolchain(det)

	avail := map[string]bool{}
	checkAvail := func(lang string) bool {
		if v, ok := avail[lang]; ok {
			return v
		}
		ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
		v := det.Capability(ctx, lang).Status == toolchain.Available
		cancel()
		avail[lang] = v
		return v
	}

	type tally struct {
		graded, ok int
		issues     []string
	}
	stats := map[string]*tally{}

	for _, l := range c.Lessons {
		lang := normLang(l.Lang)
		for si, st := range l.Stages {
			if st.Task == nil || st.Task.Solution == "" {
				continue
			}
			tk := st.Task
			if stats[lang] == nil {
				stats[lang] = &tally{}
			}
			s := stats[lang]
			if !checkAvail(lang) {
				s.issues = append(s.issues, fmt.Sprintf("%s stage%d (%s): toolchain unavailable — skipped", l.ID, si+1, lang))
				continue
			}
			s.graded++
			v, gerr := g.Run(lang, tk.Solution, tk.FuncName, tk.Tests, grader.Shape{})
			if v.Passed {
				s.ok++
				continue
			}
			detail := v.Err
			if detail == "" && gerr != nil {
				detail = gerr.Error()
			}
			for _, r := range v.Results {
				if !r.Passed {
					detail += fmt.Sprintf(" [%s got=%q want=%q]", r.Name, r.Got, r.Expected)
				}
			}
			s.issues = append(s.issues, fmt.Sprintf("%s stage%d %s: solution did NOT pass — %s", l.ID, si+1, tk.FuncName, detail))
		}
	}

	langs := make([]string, 0, len(stats))
	for l := range stats {
		langs = append(langs, l)
	}
	sort.Strings(langs)
	t.Log("Tutorial-Island native grading scorecard (reference solutions through the real toolchain):")
	for _, l := range langs {
		s := stats[l]
		t.Logf("  %-11s graded=%-3d passed=%-3d", l, s.graded, s.ok)
		for _, iss := range s.issues {
			t.Logf("      ⚠ %s", iss)
		}
		if s.graded > 0 && s.ok < s.graded && avail[l] {
			t.Errorf("%s: %d/%d lesson solutions failed to grade green", l, s.graded-s.ok, s.graded)
		}
	}
}
