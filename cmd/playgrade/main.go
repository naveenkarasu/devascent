// Command playgrade is a headless "play as a student" harness: it grades
// arbitrary student-written source against a DevAscent Advanced-Topics exercise
// through the REAL game grader (grader.New — exactly what the TUI constructs),
// without opening the interactive editor. Dev/QA tool; not part of the shipped
// product.
//
//	go run ./cmd/playgrade -lang rust -list
//	go run ./cmd/playgrade -lang rust -topic "Lifetimes" -ex 1 -code my.rs -expect pass
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"devascent/internal/content"
	"devascent/internal/grader"
	"devascent/internal/toolchain"
)

func main() {
	lang := flag.String("lang", "", "session language (go|java|javascript|typescript|rust|python|csharp)")
	list := flag.Bool("list", false, "list gradeable Advanced-Topics exercises for -lang (JSON)")
	topic := flag.String("topic", "", "topic title (substring match)")
	ex := flag.Int("ex", 0, "1-based exercise index within the topic")
	codeFile := flag.String("code", "", "path to a file holding the student's source")
	expect := flag.String("expect", "", "pass|fail — what the student INTENDED, for grading-correctness check")
	benchList := flag.Bool("bench-list", false, "list bench (DSA) problems (JSON); -diff filters by difficulty")
	bench := flag.Bool("bench", false, "grade -code against bench problem -id")
	id := flag.String("id", "", "bench problem ID (with -bench)")
	diff := flag.String("diff", "", "filter bench list by difficulty (easy|medium|hard)")
	withSolution := flag.Bool("with-solution", false, "include the python reference Solution in -bench-list output")
	flag.Parse()

	c, err := content.Load()
	if err != nil {
		die(err)
	}

	if *benchList {
		type bitem struct {
			ID         string `json:"id"`
			Difficulty string `json:"difficulty"`
			Category   string `json:"category"`
			FuncName   string `json:"funcName"`
			Prompt     string `json:"prompt"`
			Solution   string `json:"solution,omitempty"` // python reference (to translate)
		}
		var out []bitem
		for _, p := range c.Problems {
			if *diff != "" && p.Difficulty != *diff {
				continue
			}
			sol := ""
			if *withSolution {
				sol = p.Solution
			}
			out = append(out, bitem{ID: p.ID, Difficulty: p.Difficulty, Category: p.Category, FuncName: p.FuncName, Prompt: p.Prompt, Solution: sol})
		}
		b, _ := json.MarshalIndent(out, "", "  ")
		fmt.Println(string(b))
		return
	}

	if *bench {
		var p *content.Problem
		for i := range c.Problems {
			if c.Problems[i].ID == *id {
				p = &c.Problems[i]
				break
			}
		}
		if p == nil {
			die(fmt.Errorf("bench problem id %q not found", *id))
		}
		src, err := os.ReadFile(*codeFile)
		if err != nil {
			die(err)
		}
		det := toolchain.New()
		g := grader.New(det)
		v, err := g.Run(*lang, string(src), p.FuncName, p.Tests, p.GraderShape())
		if err != nil {
			v.Err = err.Error()
		}
		failed := 0
		for _, r := range v.Results {
			if !r.Passed {
				failed++
			}
		}
		res := map[string]any{
			"lang": *lang, "id": p.ID, "difficulty": p.Difficulty, "funcName": p.FuncName,
			"passed": v.Passed, "casesTotal": len(v.Results), "casesFailed": failed, "err": v.Err,
		}
		if *expect != "" {
			res["expect"] = *expect
			res["correctlyGraded"] = v.Passed == (*expect == "pass")
		}
		b, _ := json.MarshalIndent(res, "", "  ")
		fmt.Println(string(b))
		return
	}

	if *list {
		type item struct {
			Topic    string `json:"topic"`
			Ex       int    `json:"ex"`
			Check    string `json:"check"`
			Signal   string `json:"signal"`
			FuncName string `json:"funcName,omitempty"`
			Prompt   string `json:"prompt"`
		}
		var out []item
		for _, t := range c.AdvancedTopics {
			if t.Lang != *lang {
				continue
			}
			for i, e := range t.Exercises {
				if e.Check == "" || e.Check == "none" {
					continue
				}
				out = append(out, item{Topic: t.Title, Ex: i + 1, Check: e.Check, Signal: e.Signal, FuncName: e.FuncName, Prompt: e.Prompt})
			}
		}
		b, _ := json.MarshalIndent(out, "", "  ")
		fmt.Println(string(b))
		return
	}

	var found *content.Exercise
	var topicTitle string
	for ti := range c.AdvancedTopics {
		t := c.AdvancedTopics[ti]
		if t.Lang != *lang {
			continue
		}
		if *topic != "" && !strings.Contains(t.Title, *topic) {
			continue
		}
		if *ex >= 1 && *ex <= len(t.Exercises) {
			found = &t.Exercises[*ex-1]
			topicTitle = t.Title
			break
		}
	}
	if found == nil {
		die(fmt.Errorf("no gradeable exercise for lang=%s topic=%q ex=%d", *lang, *topic, *ex))
	}

	src, err := os.ReadFile(*codeFile)
	if err != nil {
		die(err)
	}

	det := toolchain.New()
	g := grader.New(det) // the SAME grader the game's TUI uses

	// Mirror the TUI's solveCheck: the PLAYER's fix for a compile-error/compiles
	// exercise just needs to COMPILE; tests/stdout grade as authored. This is how
	// a real player is graded (not the content-gate's broken-fails semantics).
	var v grader.Verdict
	switch found.Check {
	case "tests":
		v, err = g.Run(*lang, string(src), found.FuncName, found.Tests, found.GraderShape())
	case "compile-error", "compiles":
		v, err = g.Grade(grader.GradeRequest{Lang: *lang, Source: string(src), Check: grader.CheckCompiles})
	case "stdout":
		v, err = g.Grade(grader.GradeRequest{Lang: *lang, Source: string(src), Check: grader.CheckStdout, Signal: found.Signal})
	default:
		die(fmt.Errorf("exercise check %q is not player-gradeable", found.Check))
	}
	if err != nil {
		v.Err = err.Error()
	}

	got := ""
	if len(v.Results) > 0 {
		got = v.Results[0].Got
	}
	res := map[string]any{
		"lang": *lang, "topic": topicTitle, "ex": *ex, "check": found.Check,
		"signal": found.Signal, "passed": v.Passed, "err": v.Err, "got": got,
	}
	if *expect != "" {
		res["expect"] = *expect
		res["correctlyGraded"] = v.Passed == (*expect == "pass")
	}
	b, _ := json.MarshalIndent(res, "", "  ")
	fmt.Println(string(b))
}

func die(err error) { fmt.Fprintln(os.Stderr, "ERROR:", err); os.Exit(1) }
