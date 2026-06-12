// Package guiapi is the public seam the GUI frontend (a separate Go module that
// cannot import internal/*) calls into the core game engine. It mirrors what the
// TUI does directly via internal/engine + internal/grader, exposed as an
// importable, JSON-friendly surface. The GUI's Wails App methods are thin
// wrappers over an *Engine.
package guiapi

import (
	"math/rand"
	"sync"
	"time"

	"devascent/internal/content"
	"devascent/internal/engine"
	"devascent/internal/grader"
	"devascent/internal/save"
	"devascent/internal/toolchain"
)

// Engine holds the loaded catalog, the live grader (BYO local toolchain), and
// the per-language save slots (the same save files the TUI uses).
type Engine struct {
	cat content.Catalog
	det *toolchain.Detector
	g   grader.Grader
	rng *rand.Rand

	mu       sync.Mutex // guards slots
	slots    map[string]*slot
	probByID map[string]content.Problem
}

// slot is one language's loaded save state + banked index.
type slot struct {
	st     save.State
	solved map[string]bool
	nudges map[string]int // per-problem nudge escalation (session-only)
}

// New loads the content catalog and the grader. Save slots load lazily per
// language on first touch.
func New() (*Engine, error) {
	cat, err := content.Load()
	if err != nil {
		return nil, err
	}
	det := toolchain.New()
	e := &Engine{
		cat:      cat,
		det:      det,
		g:        grader.New(det),
		rng:      rand.New(rand.NewSource(time.Now().UnixNano())),
		slots:    map[string]*slot{},
		probByID: map[string]content.Problem{},
	}
	for _, p := range cat.Problems {
		e.probByID[p.ID] = p
	}
	return e, nil
}

// getSlot returns lang's slot, loading it from disk on first touch. The caller
// must hold e.mu. A corrupt slot must not kill the GUI: it starts fresh in
// memory (the next bank overwrites it).
func (e *Engine) getSlot(lang string) *slot {
	if lang == "" {
		lang = "python"
	}
	if sl, ok := e.slots[lang]; ok {
		return sl
	}
	sl := &slot{solved: map[string]bool{}, nudges: map[string]int{}}
	if st, err := save.LoadLang(lang); err == nil && st != nil {
		sl.st = *st
		for _, id := range st.SolvedIDs {
			sl.solved[id] = true
		}
	}
	e.slots[lang] = sl
	return sl
}

// toGradeResult flattens a grader Verdict into the player-facing shape.
func toGradeResult(v grader.Verdict) GradeResult {
	failed := 0
	rs := make([]CaseResult, len(v.Results))
	for i, r := range v.Results {
		if !r.Passed {
			failed++
		}
		rs[i] = CaseResult{Name: r.Name, Passed: r.Passed, Got: r.Got, Expected: r.Expected, Err: r.Err}
	}
	return GradeResult{Passed: v.Passed, CasesTotal: len(v.Results), CasesFailed: failed, Err: v.Err, Results: rs}
}

// gradeTask runs the player's code for a plain (non-node) Task — the shape used
// by diagnostic code items and lesson you-do stages.
func (e *Engine) gradeTask(lang, code string, t *content.Task) GradeResult {
	if t == nil {
		return GradeResult{Err: "no task"}
	}
	if !engine.GradingAvailable(lang) {
		return GradeResult{Err: lang + " grading isn't available (reference-only)"}
	}
	v, err := e.g.Run(lang, code, t.FuncName, t.Tests, grader.Shape{})
	if err != nil {
		return GradeResult{Err: err.Error()}
	}
	return toGradeResult(v)
}

// GradedLanguages is the set the GUI offers for grading (C++ is reference-only).
func GradedLanguages() []string {
	return []string{"python", "go", "csharp", "javascript", "typescript", "java", "rust"}
}

// ProblemSummary is one row in the bench browse list.
type ProblemSummary struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	Difficulty string   `json:"difficulty"`
	Category   string   `json:"category"`
	Lists      []string `json:"lists"` // curated-list tags (blind75, neetcode150, …)
	Solved     bool     `json:"solved"`
	Writeup    bool     `json:"writeup"` // write-up complete (solved && !writeup = provisional)
}

// Problems returns the full bench problem list (browse view) with lang's
// banked marks.
func (e *Engine) Problems(lang string) []ProblemSummary {
	e.mu.Lock()
	defer e.mu.Unlock()
	sl := e.getSlot(lang)
	out := make([]ProblemSummary, 0, len(e.cat.Problems))
	for _, p := range e.cat.Problems {
		out = append(out, ProblemSummary{
			ID: p.ID, Title: p.Title, Difficulty: p.Difficulty, Category: p.Category,
			Lists:   p.Lists,
			Solved:  sl.solved[p.ID],
			Writeup: sl.record(p.ID).WriteupDone,
		})
	}
	return out
}

// ProblemDetail is a problem opened in the workbench, with a starter rendered
// for the chosen language.
type ProblemDetail struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Difficulty string `json:"difficulty"`
	Category   string `json:"category"`
	Prompt     string `json:"prompt"`
	FuncName   string `json:"funcName"`
	Lang       string `json:"lang"`
	Starter    string `json:"starter"` // language-native starter code to seed the editor
	Found      bool   `json:"found"`
}

// Problem returns one problem with a starter for lang (python uses the authored
// starter; other languages get an inferred typed stub via engine.Starter).
func (e *Engine) Problem(id, lang string) ProblemDetail {
	p, ok := e.byID(id)
	if !ok {
		return ProblemDetail{ID: id, Lang: lang, Found: false}
	}
	starter := p.Starter
	if lang != "python" {
		if s, ok := engine.Starter(lang, p.FuncName, p.Solution, p.Tests, p.GraderShape()); ok {
			starter = s
		}
	}
	return ProblemDetail{
		ID: p.ID, Title: p.Title, Difficulty: p.Difficulty, Category: p.Category,
		Prompt: p.Prompt, FuncName: p.FuncName, Lang: lang, Starter: starter, Found: true,
	}
}

// CaseResult is one hidden-test outcome (verdict detail).
type CaseResult struct {
	Name     string `json:"name"`
	Passed   bool   `json:"passed"`
	Got      string `json:"got"`
	Expected string `json:"expected"`
	Err      string `json:"err"`
}

// GradeResult is the player-facing verdict for a submission.
type GradeResult struct {
	Passed      bool         `json:"passed"`
	CasesTotal  int          `json:"casesTotal"`
	CasesFailed int          `json:"casesFailed"`
	Err         string       `json:"err"` // compile/runtime/harness error (empty on normal completion)
	Results     []CaseResult `json:"results"`

	// Bench banking (set only by the bench Grade path on a pass).
	Banked      bool   `json:"banked"`      // this problem is in the banked set
	NewlyBanked bool   `json:"newlyBanked"` // this pass banked it for the first time
	SaveErr     string `json:"saveErr"`     // persistence failure (grade still valid)

	// Track A (set on the bench path): token payout for a clean first bank,
	// and whether the write-up gate is still open for this problem.
	TokensAwarded  int  `json:"tokensAwarded"`
	WriteupPending bool `json:"writeupPending"`
}

// Grade runs the player's code against the problem's hidden tests in lang,
// through the same grader the TUI uses.
func (e *Engine) Grade(lang, id, code string) GradeResult {
	p, ok := e.byID(id)
	if !ok {
		return GradeResult{Err: "unknown problem: " + id}
	}
	if !engine.GradingAvailable(lang) {
		return GradeResult{Err: lang + " grading isn't available (reference-only)"}
	}
	v, err := e.g.Run(lang, code, p.FuncName, p.Tests, p.GraderShape())
	if err != nil {
		return GradeResult{Err: err.Error()}
	}
	res := toGradeResult(v)
	if res.Passed {
		res.Banked = true
		res.NewlyBanked, res.SaveErr = e.bank(lang, id)
		res.TokensAwarded, res.WriteupPending = e.trackABank(lang, p)
	} else {
		e.trackAFail(lang, id)
	}
	return res
}

func (e *Engine) byID(id string) (content.Problem, bool) {
	for i := range e.cat.Problems {
		if e.cat.Problems[i].ID == id {
			return e.cat.Problems[i], true
		}
	}
	return content.Problem{}, false
}
