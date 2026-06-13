package guiapi

import (
	"devascent/internal/content"
	"devascent/internal/engine"
)

// ── Dev-Literacy track ────────────────────────────────────────────────────────
// A command CHECKER, not a shell: the player types the command they'd run; a
// wrong answer stays on the same task with the hint, a right one advances.
// Matching reuses engine.DevMatch (the same logic the TUI track runs).

// DevLitStep is the current task to render, or the terminal result.
type DevLitStep struct {
	Done     bool   `json:"done"`
	Index    int    `json:"index"` // 1-based
	Total    int    `json:"total"`
	Category string `json:"category"`
	Title    string `json:"title"`
	Prompt   string `json:"prompt"`
	Passed   int    `json:"passed"` // terminal: how many were solved
}

// DevLitOutcome is returned after submitting one command.
type DevLitOutcome struct {
	Passed  bool       `json:"passed"`
	Hint    string     `json:"hint"`    // on a wrong answer
	Success string     `json:"success"` // on a right answer
	Next    DevLitStep `json:"next"`
}

// DevLiteracy is a live dev-literacy session.
type DevLiteracy struct {
	tasks []content.DevTask
	idx   int
	ok    int
}

// StartDevLiteracy selects 5 tasks across distinct categories.
func (e *Engine) StartDevLiteracy() *DevLiteracy {
	return &DevLiteracy{tasks: engine.SelectDevSet(e.cat.DevTasks, 5, e.rng)}
}

// Step returns the current task, or the terminal result when done.
func (d *DevLiteracy) Step() DevLitStep {
	if d.idx >= len(d.tasks) {
		return DevLitStep{Done: true, Total: len(d.tasks), Passed: d.ok}
	}
	t := d.tasks[d.idx]
	return DevLitStep{
		Index: d.idx + 1, Total: len(d.tasks),
		Category: t.Category, Title: t.Title, Prompt: t.Prompt,
	}
}

// Submit grades a typed command. Wrong answers do NOT advance.
func (d *DevLiteracy) Submit(ans string) DevLitOutcome {
	if d.idx >= len(d.tasks) {
		return DevLitOutcome{Next: d.Step()}
	}
	t := d.tasks[d.idx]
	if engine.DevMatch(ans, t.Commands, t.Flags, t.Accept) {
		d.ok++
		d.idx++
		return DevLitOutcome{Passed: true, Success: t.Success, Next: d.Step()}
	}
	return DevLitOutcome{Passed: false, Hint: t.Hint, Next: d.Step()}
}
