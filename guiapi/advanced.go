package guiapi

import (
	"devascent/internal/grader"
)

// ── Advanced Topics ───────────────────────────────────────────────────────────
// Per-language deep-dive topics (ownership, concurrency, gotchas, …): a prose
// explainer + primer-style sections + exercises over BROKEN code the player
// fixes. Grading mirrors the TUI's solveCheck mapping: an exercise authored as
// compile-error/compiles passes when the player's fix compiles; stdout compares
// output; tests runs the function-call harness; none = reveal-only reading.

// AdvTopicSummary is one row in the topics list.
type AdvTopicSummary struct {
	Index     int    `json:"index"` // 0-based position in the lang's topic list
	Group     string `json:"group"`
	Title     string `json:"title"`
	Tag       string `json:"tag"` // E | C | P | gotcha
	Exercises int    `json:"exercises"`
	Gradeable int    `json:"gradeable"`
}

// AdvExerciseView is one exercise within a topic. FixedCode/Bug are the reveal
// material — the frontend gates their display, grading never depends on them.
type AdvExerciseView struct {
	Index      int    `json:"index"`
	Kind       string `json:"kind"`
	Prompt     string `json:"prompt"`
	BrokenCode string `json:"brokenCode"`
	FixedCode  string `json:"fixedCode"`
	Bug        string `json:"bug"`
	Check      string `json:"check"` // tests | compile-error | stdout | compiles | none
	Gradeable  bool   `json:"gradeable"`
}

// AdvTopicDetail is a topic opened in the Advanced view.
type AdvTopicDetail struct {
	Found     bool                `json:"found"`
	Index     int                 `json:"index"`
	Lang      string              `json:"lang"`
	Group     string              `json:"group"`
	Title     string              `json:"title"`
	Tag       string              `json:"tag"`
	Summary   string              `json:"summary"`
	Sections  []PrimerSectionView `json:"sections"`
	Exercises []AdvExerciseView   `json:"exercises"`
}

// advSolveCheck mirrors the TUI's solveCheck: how the PLAYER's fix is graded
// for an exercise authored with the given Check. "" = not gradeable.
func advSolveCheck(check string) grader.Check {
	switch check {
	case "tests":
		return grader.CheckTests
	case "compile-error", "compiles":
		return grader.CheckCompiles
	case "stdout":
		return grader.CheckStdout
	default:
		return ""
	}
}

// AdvancedTopics lists lang's advanced topics.
func (e *Engine) AdvancedTopics(lang string) []AdvTopicSummary {
	topics := e.cat.AdvancedTopicsByLang(lang)
	out := make([]AdvTopicSummary, 0, len(topics))
	for i, t := range topics {
		gradeable := 0
		for _, ex := range t.Exercises {
			if advSolveCheck(ex.Check) != "" {
				gradeable++
			}
		}
		out = append(out, AdvTopicSummary{
			Index: i, Group: t.Group, Title: t.Title, Tag: t.Tag,
			Exercises: len(t.Exercises), Gradeable: gradeable,
		})
	}
	return out
}

// AdvancedTopic returns one topic with its sections and exercises.
func (e *Engine) AdvancedTopic(lang string, idx int) AdvTopicDetail {
	topics := e.cat.AdvancedTopicsByLang(lang)
	if idx < 0 || idx >= len(topics) {
		return AdvTopicDetail{Found: false, Index: idx, Lang: lang}
	}
	t := topics[idx]
	d := AdvTopicDetail{
		Found: true, Index: idx, Lang: lang,
		Group: t.Group, Title: t.Title, Tag: t.Tag, Summary: t.Summary,
	}
	for _, s := range t.Sections {
		sv := PrimerSectionView{Title: s.Title}
		for _, op := range s.Ops {
			sv.Ops = append(sv.Ops, PrimerOpView{Label: op.Label, Code: op.Code})
		}
		d.Sections = append(d.Sections, sv)
	}
	for i, ex := range t.Exercises {
		d.Exercises = append(d.Exercises, AdvExerciseView{
			Index: i, Kind: ex.Kind, Prompt: ex.Prompt,
			BrokenCode: ex.BrokenCode, FixedCode: ex.FixedCode, Bug: ex.Bug,
			Check: ex.Check, Gradeable: advSolveCheck(ex.Check) != "",
		})
	}
	return d
}

// GradeAdvanced grades the player's fix for one exercise through the native
// toolchain — the same GradeRequest mapping the TUI uses.
func (e *Engine) GradeAdvanced(lang string, topicIdx, exIdx int, code string) GradeResult {
	topics := e.cat.AdvancedTopicsByLang(lang)
	if topicIdx < 0 || topicIdx >= len(topics) {
		return GradeResult{Err: "no such topic"}
	}
	exs := topics[topicIdx].Exercises
	if exIdx < 0 || exIdx >= len(exs) {
		return GradeResult{Err: "no such exercise"}
	}
	ex := exs[exIdx]
	check := advSolveCheck(ex.Check)
	if check == "" {
		return GradeResult{Err: "this exercise is reveal-only (not graded)"}
	}
	req := grader.GradeRequest{Lang: lang, Check: check, Source: code}
	switch ex.Check {
	case "tests":
		req.FuncName = ex.FuncName
		req.Tests = ex.Tests
		req.Shape = ex.GraderShape()
	case "stdout":
		req.Signal = ex.Signal
	}
	v, err := e.g.Grade(req)
	if err != nil {
		return GradeResult{Err: err.Error()}
	}
	return toGradeResult(v)
}
