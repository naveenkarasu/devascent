package guiapi

import (
	"devascent/internal/content"
	"devascent/internal/engine"
)

// ── Orientation (entrance test) ───────────────────────────────────────────────
// A live, server-held session: the GUI renders Step(), submits an answer, and
// gets the graded outcome + the next Step. Scoring + routing reuse internal/engine
// (the same logic the TUI runs), so there is one source of truth.

// OrientationStep is the current diagnostic item to render, or a terminal result.
type OrientationStep struct {
	Done     bool     `json:"done"`
	Index    int      `json:"index"` // 1-based position in the ladder
	Total    int      `json:"total"`
	Kind     string   `json:"kind"` // "code" | "choice" | "spec"
	Slot     string   `json:"slot"`
	Measures string   `json:"measures"`
	Prompt   string   `json:"prompt"`
	Lang     string   `json:"lang"`
	FuncName string   `json:"funcName"` // kind=code
	Starter  string   `json:"starter"`  // kind=code
	Choices  []string `json:"choices"`  // kind=choice (texts only; correctness is hidden)
	// Terminal fields (set when Done):
	Placement    string `json:"placement"` // "test-out" | "dev-literacy" | "tutorial-full"
	Score        int    `json:"score"`
	Level        string `json:"level"`
	CodingOK     int    `json:"codingOK"`
	CodingTotal  int    `json:"codingTotal"`
	MachineOK    int    `json:"machineOK"`
	MachineTotal int    `json:"machineTotal"`
	SpecOK       int    `json:"specOK"`
	SpecTotal    int    `json:"specTotal"`
}

// DiagOutcome is returned after submitting one answer: the grade + the next step.
type DiagOutcome struct {
	Passed   bool            `json:"passed"`
	Feedback string          `json:"feedback"`          // choice/spec feedback or model answer
	Verdict  *GradeResult    `json:"verdict,omitempty"` // kind=code per-case detail
	Next     OrientationStep `json:"next"`
}

// Orientation is a live entrance-test session.
type Orientation struct {
	e                           *Engine
	lang                        string
	level                       string
	diag                        []content.Diagnostic
	idx                         int
	codingOK, machineOK, specOK int
	done                        bool
	place                       string
}

// StartOrientation builds a per-language intake ladder for the given self-report
// level ("never" | "a-little" | "regularly"; default "a-little").
func (e *Engine) StartOrientation(lang, level string) *Orientation {
	if level == "" {
		level = "a-little"
	}
	pool := e.cat.DiagnosticsForLang(lang)
	diag := engine.SelectIntake(pool, level, nil, e.rng)
	o := &Orientation{e: e, lang: lang, level: level, diag: diag}
	if len(diag) == 0 {
		// No intake authored for this language: route by the self-report clamp
		// alone rather than presenting an empty ladder.
		o.done = true
		o.place = engine.Place(0, 0, level)
	}
	return o
}

// Step returns the current item to render, or the terminal result when done.
func (o *Orientation) Step() OrientationStep {
	if o.done || o.idx >= len(o.diag) {
		ct, mt, st := 0, 0, 0
		for _, d := range o.diag {
			switch d.Measures {
			case "coding":
				ct++
			case "machine":
				mt++
			case "spec":
				st++
			}
		}
		return OrientationStep{
			Done: true, Total: len(o.diag), Placement: o.place,
			Score: o.codingOK + o.machineOK + o.specOK, Lang: o.lang, Level: o.level,
			CodingOK: o.codingOK, CodingTotal: ct,
			MachineOK: o.machineOK, MachineTotal: mt,
			SpecOK: o.specOK, SpecTotal: st,
		}
	}
	d := o.diag[o.idx]
	s := OrientationStep{
		Index: o.idx + 1, Total: len(o.diag), Kind: d.Kind, Slot: d.Slot,
		Measures: d.Measures, Prompt: d.Prompt, Lang: o.lang,
	}
	switch d.Kind {
	case "code":
		if d.Task != nil {
			s.FuncName = d.Task.FuncName
			s.Starter = d.Task.Starter
		}
	case "choice":
		for _, c := range d.Choices {
			s.Choices = append(s.Choices, c.Text)
		}
	}
	return s
}

// record scores the current item by its measure, advances, and routes at the end.
func (o *Orientation) record(passed bool) {
	if o.done || o.idx >= len(o.diag) {
		return
	}
	if passed {
		switch o.diag[o.idx].Measures {
		case "coding":
			o.codingOK++
		case "machine":
			o.machineOK++
		case "spec":
			o.specOK++
		}
	}
	o.idx++
	if o.idx >= len(o.diag) {
		o.done = true
		passed := o.codingOK + o.machineOK + o.specOK
		o.place = engine.Place(passed, len(o.diag), o.level)
		o.e.persistPlacement(o.lang, o.place, o.level, passed)
	}
}

func (o *Orientation) outcome(passed bool, feedback string, v *GradeResult) DiagOutcome {
	o.record(passed)
	return DiagOutcome{Passed: passed, Feedback: feedback, Verdict: v, Next: o.Step()}
}

// SubmitCode grades a kind=code item via the real grader.
func (o *Orientation) SubmitCode(code string) DiagOutcome {
	if o.done || o.idx >= len(o.diag) {
		return DiagOutcome{Next: o.Step()}
	}
	d := o.diag[o.idx]
	if d.Kind != "code" || d.Task == nil {
		return o.outcome(false, "this item is not a coding item", nil)
	}
	gr := o.e.gradeTask(o.lang, code, d.Task)
	return o.outcome(gr.Passed, "", &gr)
}

// SubmitChoice grades a kind=choice item by the chosen option index.
func (o *Orientation) SubmitChoice(idx int) DiagOutcome {
	if o.done || o.idx >= len(o.diag) {
		return DiagOutcome{Next: o.Step()}
	}
	d := o.diag[o.idx]
	if d.Kind != "choice" || idx < 0 || idx >= len(d.Choices) {
		return o.outcome(false, "invalid choice", nil)
	}
	c := d.Choices[idx]
	return o.outcome(c.Correct, c.Feedback, nil)
}

// SubmitSpec grades a kind=spec free-text item (synonym-group keyword match).
func (o *Orientation) SubmitSpec(text string) DiagOutcome {
	if o.done || o.idx >= len(o.diag) {
		return DiagOutcome{Next: o.Step()}
	}
	d := o.diag[o.idx]
	if d.Kind != "spec" || d.Spec == nil {
		return o.outcome(false, "", nil)
	}
	ok := engine.SpecMatch(text, d.Spec.Required)
	return o.outcome(ok, d.Spec.Answer, nil)
}

// ── Tutorial Island ───────────────────────────────────────────────────────────
// (moved to tutorial.go — session, stepper positions, persistence)
