package guiapi

import (
	"devascent/internal/engine"
	"devascent/internal/save"
)

// ── Scoring surface ───────────────────────────────────────────────────────────
// Progress is computed from the SHARED per-language save slots (the same
// save-<lang>.json files the TUI uses), so banking in either frontend advances
// the same run. The GUI only ever writes the fields it owns (SolvedIDs/
// BenchSolved + the placement block); everything else round-trips untouched
// through Load → mutate → Save.

// Progress is the bench scorecard for one language: banked vs the Step 0
// milestone targets plus the provisional competency profile.
type Progress struct {
	Lang           string `json:"lang"`
	Banked         int    `json:"banked"`
	Cats           int    `json:"cats"`
	Hard           int    `json:"hard"`
	BankTarget     int    `json:"bankTarget"`
	CatTarget      int    `json:"catTarget"`
	HardTarget     int    `json:"hardTarget"`
	Step0Met       bool   `json:"step0Met"`
	ProblemSolving int    `json:"problemSolving"`
	LangProf       int    `json:"langProf"`
	Track          string `json:"track"`
	Placement      string `json:"placement"`
	Level          string `json:"level"` // self-report band from the intake
	TotalProblems  int    `json:"totalProblems"`
}

// ProfileView is one language slot in the profile picker.
type ProfileView struct {
	Lang      string `json:"lang"`
	Stage     string `json:"stage"`
	Placement string `json:"placement"`
	Level     string `json:"level"`
	Banked    int    `json:"banked"`
	UpdatedAt string `json:"updatedAt"`
}

// Progress returns lang's scorecard.
func (e *Engine) Progress(lang string) Progress {
	e.mu.Lock()
	defer e.mu.Unlock()
	sl := e.getSlot(lang)
	b, c, h := engine.BenchStats(sl.solved, e.probByID)
	ps, lp, track := engine.Step0Profile(sl.solved, e.probByID, e.cat.Problems)
	return Progress{
		Lang:   lang,
		Banked: b, Cats: c, Hard: h,
		BankTarget: engine.Step0BankTarget, CatTarget: engine.Step0CatTarget, HardTarget: engine.Step0HardTarget,
		Step0Met:       engine.Step0Met(sl.solved, e.probByID),
		ProblemSolving: ps, LangProf: lp, Track: track,
		Placement:     sl.st.Placement,
		Level:         sl.st.Level,
		TotalProblems: len(e.cat.Problems),
	}
}

// Profiles lists every language slot on disk, most recently played first
// (the Home profile picker).
func (e *Engine) Profiles() []ProfileView {
	ps, err := save.Profiles()
	if err != nil {
		return nil
	}
	out := make([]ProfileView, 0, len(ps))
	for _, p := range ps {
		out = append(out, ProfileView{
			Lang: p.Lang, Stage: p.Stage, Placement: p.Placement,
			Level: p.Level, Banked: p.Banked, UpdatedAt: p.UpdatedAt,
		})
	}
	return out
}

// DeleteProfile removes lang's save slot from disk and memory — the GUI's
// "delete profile". Destructive; the frontend confirms before calling.
// Returns an error message, or "" on success.
func (e *Engine) DeleteProfile(lang string) string {
	if lang == "" {
		lang = "python"
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if err := save.DeleteLang(lang); err != nil {
		return err.Error()
	}
	delete(e.slots, lang)
	return ""
}

// NextProblem returns the first problem unsolved in lang after afterID in
// catalog order (wrapping around), or "" when everything is solved.
func (e *Engine) NextProblem(lang, afterID string) string {
	e.mu.Lock()
	defer e.mu.Unlock()
	sl := e.getSlot(lang)
	ps := e.cat.Problems
	if len(ps) == 0 {
		return ""
	}
	start := 0
	for i := range ps {
		if ps[i].ID == afterID {
			start = i + 1
			break
		}
	}
	for k := 0; k < len(ps); k++ {
		p := ps[(start+k)%len(ps)]
		if !sl.solved[p.ID] {
			return p.ID
		}
	}
	return ""
}

// bank records a newly passed problem into lang's slot. Returns whether this
// bank was new and any persistence error (grading already succeeded, so a
// disk failure is surfaced separately rather than failing the verdict).
func (e *Engine) bank(lang, id string) (newlyBanked bool, saveErr string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	sl := e.getSlot(lang)
	if sl.solved[id] {
		return false, ""
	}
	sl.solved[id] = true
	sl.st.SolvedIDs = append(sl.st.SolvedIDs, id)
	sl.st.BenchSolved++
	if err := save.SaveLang(lang, sl.st); err != nil {
		return true, err.Error()
	}
	return true, ""
}

// persistPlacement records a completed orientation into lang's slot.
func (e *Engine) persistPlacement(lang, place, level string, passed int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	sl := e.getSlot(lang)
	sl.st.Placement = place
	sl.st.Level = level
	sl.st.IntakePassed = passed
	// Best-effort: the placement is still shown to the player either way.
	_ = save.SaveLang(lang, sl.st)
}
