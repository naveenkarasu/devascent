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

// TestAdvancedExercisesGradeNativeRoundTrip is the multi-language "solve &
// verify" gate (ADR-0007): it grades the AUTHORED solutions for every gradeable
// Advanced-Topics exercise through the player's REAL toolchain and reports a
// per-language scorecard. This is the deterministic, system-driven version of
// "play the game and solve" — the model solution is the player's solve, verified
// by the actual compiler/interpreter, not by WASM.
//
// Per check mode:
//   - tests        (Python): FixedCode passes the tests; BrokenCode fails them.
//   - compile-error (Rust):  BrokenCode fails to compile WITH the signal code;
//     FixedCode compiles cleanly. (Both directions — deterministic.)
//   - compiles     (TS):     FixedCode compiles; BrokenCode does not. (Both.)
//   - stdout       (Go/Java/JS): FixedCode prints the signal. (Fixed only —
//     broken is often a race/non-determinism that can coincidentally match.)
//
// Slow (compiles/runs ~100 exercises across 6 languages) → opt-in, -short-skip.
func TestAdvancedExercisesGradeNativeRoundTrip(t *testing.T) {
	if testing.Short() {
		t.Skip("native advanced round-trip is slow; skipped in -short")
	}
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	det := toolchain.New()
	g := grader.NewLocalToolchain(det)

	type tally struct {
		gradeable    int
		fixedOK      int // authored solution graded as correct
		brokenOK     int // broken graded as incorrect (only where deterministic)
		brokenChecks int // how many broken-direction checks we ran
		issues       []string
	}
	stats := map[string]*tally{}
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

	for _, topic := range c.AdvancedTopics {
		for i, ex := range topic.Exercises {
			if ex.Check == "" || ex.Check == "none" {
				continue
			}
			lang := topic.Lang
			if stats[lang] == nil {
				stats[lang] = &tally{}
			}
			st := stats[lang]
			st.gradeable++
			if !checkAvail(lang) {
				st.issues = append(st.issues, fmt.Sprintf("%s ex%d: %s toolchain unavailable — skipped", topic.Title, i+1, lang))
				continue
			}
			label := fmt.Sprintf("%s ex%d (%s)", topic.Title, i+1, ex.Check)

			switch ex.Check {
			case "tests":
				fv, _ := g.Run(lang, ex.FixedCode, ex.FuncName, ex.Tests, ex.GraderShape())
				if fv.Passed {
					st.fixedOK++
				} else {
					st.issues = append(st.issues, label+": FixedCode did NOT pass tests")
				}
				st.brokenChecks++
				bv, _ := g.Run(lang, ex.BrokenCode, ex.FuncName, ex.Tests, ex.GraderShape())
				if !bv.Passed {
					st.brokenOK++
				} else {
					st.issues = append(st.issues, label+": BrokenCode unexpectedly PASSED tests")
				}
			case "compile-error":
				// Broken must fail to compile with the signal; fixed must compile.
				bv, _ := g.Grade(grader.GradeRequest{Lang: lang, Check: grader.CheckCompileError, Source: ex.BrokenCode, Signal: ex.Signal})
				st.brokenChecks++
				if bv.Passed {
					st.brokenOK++
				} else {
					st.issues = append(st.issues, label+": BrokenCode did NOT fail with "+ex.Signal)
				}
				fv, _ := g.Grade(grader.GradeRequest{Lang: lang, Check: grader.CheckCompiles, Source: ex.FixedCode})
				if fv.Passed {
					st.fixedOK++
				} else {
					st.issues = append(st.issues, label+": FixedCode did NOT compile")
				}
			case "compiles":
				fv, _ := g.Grade(grader.GradeRequest{Lang: lang, Check: grader.CheckCompiles, Source: ex.FixedCode})
				if fv.Passed {
					st.fixedOK++
				} else {
					st.issues = append(st.issues, label+": FixedCode did NOT compile")
				}
				st.brokenChecks++
				bv, _ := g.Grade(grader.GradeRequest{Lang: lang, Check: grader.CheckCompiles, Source: ex.BrokenCode})
				if !bv.Passed {
					st.brokenOK++
				} else {
					st.issues = append(st.issues, label+": BrokenCode unexpectedly compiled")
				}
			case "stdout":
				// Fixed must print the signal. Broken is intentionally NOT asserted:
				// many are concurrency/race bugs whose output is non-deterministic.
				fv, _ := g.Grade(grader.GradeRequest{Lang: lang, Check: grader.CheckStdout, Source: ex.FixedCode, Signal: ex.Signal})
				if fv.Passed {
					st.fixedOK++
				} else {
					got := ""
					if len(fv.Results) > 0 {
						got = fv.Results[0].Got
					}
					st.issues = append(st.issues, fmt.Sprintf("%s: FixedCode stdout != %q (got %q / err %q)", label, ex.Signal, got, fv.Err))
				}
			}
		}
	}

	// Report the per-language scorecard.
	langs := make([]string, 0, len(stats))
	for l := range stats {
		langs = append(langs, l)
	}
	sort.Strings(langs)
	t.Log("Advanced-Topics native grading scorecard (authored solutions through the real toolchain):")
	for _, l := range langs {
		st := stats[l]
		t.Logf("  %-11s gradeable=%-3d fixed-correct=%-3d broken-correct=%d/%d",
			l, st.gradeable, st.fixedOK, st.brokenOK, st.brokenChecks)
		for _, iss := range st.issues {
			t.Logf("      ⚠ %s", iss)
		}
		// A whole language grading 0 of its authored solutions = a systemic
		// grader/adapter problem worth failing on.
		if st.gradeable > 0 && st.fixedOK == 0 && avail[l] {
			t.Errorf("%s: 0/%d authored solutions graded correctly — systemic problem", l, st.gradeable)
		}
	}
}
