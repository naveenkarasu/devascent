package guiapi

// Track A facade: the hint economy (A2), write-up gate (A1), graduation gate
// (A3), and mentor seam (A4) — the same shared core the TUI uses, exposed
// JSON-friendly. Mutex discipline: e.mu is NEVER held across a mentor call
// (up to 45s); token debits commit and persist BEFORE the AI is asked, and a
// failed AI call refunds afterwards.

import (
	"context"
	"fmt"
	"sync"
	"time"

	"devascent/internal/content"
	"devascent/internal/economy"
	"devascent/internal/engine"
	"devascent/internal/mentor"
	"devascent/internal/save"
)

var mentorMu sync.Mutex
var mentorSvc *mentor.Service

// mentorService lazily builds the shared mentor service from mentor.json.
func mentorService() *mentor.Service {
	mentorMu.Lock()
	defer mentorMu.Unlock()
	if mentorSvc == nil {
		cfg, _ := mentor.LoadConfig()
		mentorSvc = mentor.NewService(cfg)
	}
	return mentorSvc
}

// record returns a copy of id's solve record (zero value when absent).
func (sl *slot) record(id string) save.SolveRecord {
	return sl.st.SolveRecords[id]
}

// setRecord stores a solve record (allocating the map on first write).
func (sl *slot) setRecord(id string, rec save.SolveRecord) {
	if sl.st.SolveRecords == nil {
		sl.st.SolveRecords = map[string]save.SolveRecord{}
	}
	sl.st.SolveRecords[id] = rec
}

// trackAFail records a failed grading attempt (pity-rule bookkeeping). The
// failure counter only advances on a DISTINCT attempt — re-running the same
// broken code doesn't count toward the free-strategy unlock.
func (e *Engine) trackAFail(lang, id, code string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	sl := e.getSlot(lang)
	if sl.solved[id] {
		return // re-running an already banked problem
	}
	rec := sl.record(id)
	if h := economy.FailHash(code); h != rec.LastFailHash {
		rec.FailedRuns++
		rec.LastFailHash = h
	}
	if rec.FirstTryAt == "" {
		rec.FirstTryAt = time.Now().UTC().Format(time.RFC3339)
	}
	sl.setRecord(id, rec)
	_ = save.SaveLang(lang, sl.st)
}

// trackABank pays out a clean first bank and reports whether the write-up
// gate is open. Called right after bank() on a passed grade.
func (e *Engine) trackABank(lang string, p content.Problem) (awarded int, writeupPending bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	sl := e.getSlot(lang)
	rec := sl.record(p.ID)
	if rec.WriteupDone {
		return 0, false
	}
	if rec.HintTier == economy.TierNone && !rec.PityUsed {
		w := economy.Load(&sl.st, time.Now())
		awarded = economy.SolveAward(p.Difficulty)
		w.Award(awarded)
		w.Store(&sl.st)
	}
	_ = save.SaveLang(lang, sl.st)
	return awarded, true
}

// ── A2: wallet + hints ───────────────────────────────────────────────────────

// WalletView is the player's hint currency display.
type WalletView struct {
	Tokens          int `json:"tokens"`
	NudgeCharges    int `json:"nudgeCharges"`
	NudgeMax        int `json:"nudgeMax"`
	NextRechargeSec int `json:"nextRechargeSec"` // 0 at cap
}

func walletView(w *economy.Wallet, now time.Time) WalletView {
	return WalletView{
		Tokens: w.Tokens, NudgeCharges: w.NudgeCharges, NudgeMax: economy.NudgeMax,
		NextRechargeSec: int(w.NextRecharge(now).Seconds()),
	}
}

// Wallet returns lang's wallet, granting the starting stash on first call.
func (e *Engine) Wallet(lang string) WalletView {
	e.mu.Lock()
	defer e.mu.Unlock()
	now := time.Now()
	sl := e.getSlot(lang)
	w := economy.Load(&sl.st, now)
	w.Store(&sl.st)
	_ = save.SaveLang(lang, sl.st)
	return walletView(&w, now)
}

// HintInfo is the per-problem hint state the panel needs to render: the wallet
// plus which paid tiers are already owned (free to re-show) and whether the
// earned free Strategy is currently on offer.
type HintInfo struct {
	Wallet           WalletView `json:"wallet"`
	StrategyOwned    bool       `json:"strategyOwned"` // already paid on this problem → re-show is free
	WalkthroughOwned bool       `json:"walkthroughOwned"`
	PityStrategyFree bool       `json:"pityStrategyFree"` // stuck long enough → Strategy offered free (player chooses)
}

// HintInfo returns the per-problem hint state for the hint panel.
func (e *Engine) HintInfo(lang, id string) HintInfo {
	e.mu.Lock()
	defer e.mu.Unlock()
	now := time.Now()
	sl := e.getSlot(lang)
	w := economy.Load(&sl.st, now)
	w.Store(&sl.st)
	_ = save.SaveLang(lang, sl.st)
	rec := sl.record(id)
	return HintInfo{
		Wallet:           walletView(&w, now),
		StrategyOwned:    rec.HintTier >= economy.TierStrategy,
		WalkthroughOwned: rec.HintTier >= economy.TierWalkthrough,
		PityStrategyFree: rec.HintTier < economy.TierStrategy && economy.PityEligible(rec, now),
	}
}

// HintResult is one answered hint request.
type HintResult struct {
	Text     string     `json:"text"`
	Source   string     `json:"source"` // "template" or the backend ID
	Tier     int        `json:"tier"`
	Pity     bool       `json:"pity"`     // served by the one-time pity rule (free)
	Refunded bool       `json:"refunded"` // AI failed; token returned
	Wallet   WalletView `json:"wallet"`
	Err      string     `json:"err"` // "not enough tokens" etc; Text empty
}

// RequestHint serves a tier-1 nudge (charge) or tier-2/3 hint (tokens).
// Paid tiers are recorded on the solve record and discount mastery; asking
// the same tier again on the same problem is free (already paid).
func (e *Engine) RequestHint(lang, id string, tier int, code string) HintResult {
	now := time.Now()
	e.mu.Lock()
	p, ok := e.probByID[id]
	if !ok {
		e.mu.Unlock()
		return HintResult{Err: "unknown problem: " + id}
	}
	sl := e.getSlot(lang)
	w := economy.Load(&sl.st, now)
	rec := sl.record(id)

	if tier == economy.TierNudge {
		if !w.SpendNudge(now) {
			next := int(w.NextRecharge(now).Seconds())
			w.Store(&sl.st)
			e.mu.Unlock()
			return HintResult{Tier: tier, Err: fmt.Sprintf("no nudges left — next recharges in %dm%02ds", next/60, next%60)}
		}
		attempt := sl.nudges[id]
		sl.nudges[id]++
		w.Store(&sl.st)
		_ = save.SaveLang(lang, sl.st)
		res := HintResult{Text: mentor.Nudge(p.Category, attempt), Source: "template", Tier: tier, Wallet: walletView(&w, now)}
		e.mu.Unlock()
		return res
	}

	if tier != economy.TierStrategy && tier != economy.TierWalkthrough {
		e.mu.Unlock()
		return HintResult{Err: fmt.Sprintf("unknown hint tier %d", tier)}
	}

	cost := economy.HintCost(tier)
	pity := false
	switch {
	case rec.HintTier >= tier:
		cost = 0 // already paid for this tier on this problem
	case tier == economy.TierStrategy && economy.PityEligible(rec, now):
		cost = 0
		pity = true
		rec.PityUsed = true
	default:
		if !w.Spend(cost) {
			e.mu.Unlock()
			return HintResult{Tier: tier, Err: fmt.Sprintf("not enough tokens (need %d) — bank a problem cleanly or finish a write-up to earn more", cost)}
		}
	}
	if tier > rec.HintTier {
		rec.HintTier = tier
	}
	sl.setRecord(id, rec)
	w.Store(&sl.st)
	_ = save.SaveLang(lang, sl.st) // debit committed BEFORE the AI call
	req := mentor.Request{
		Kind: mentor.KindStrategy, Lang: lang, Title: p.Title, Prompt: p.Prompt,
		Category: p.Category, Difficulty: p.Difficulty, PlayerCode: code,
		FailedRuns: rec.FailedRuns, FirstFail: firstFailName(rec),
	}
	if tier == economy.TierWalkthrough {
		req.Kind = mentor.KindWalkthrough
	}
	e.mu.Unlock()

	resp := mentorService().Hint(context.Background(), req) // unlocked: up to 45s

	res := HintResult{Text: resp.Text, Source: resp.Source, Tier: tier, Pity: pity}
	e.mu.Lock()
	w2 := economy.Load(&sl.st, time.Now())
	if resp.FellBack && cost > 0 {
		w2.Refund(cost)
		res.Refunded = true
		w2.Store(&sl.st)
		_ = save.SaveLang(lang, sl.st)
	}
	res.Wallet = walletView(&w2, time.Now())
	e.mu.Unlock()
	return res
}

// firstFailName is a placeholder until per-case names are tracked on the
// record; the GUI passes live failure info only through the preview for now.
func firstFailName(save.SolveRecord) string { return "" }

// ── A1: write-up gate ────────────────────────────────────────────────────────

// WriteupView is the write-up form for one solved problem. The MCQ's correct
// index never leaves the engine.
type WriteupView struct {
	ProblemID string   `json:"problemId"`
	Title     string   `json:"title"`
	Solved    bool     `json:"solved"`
	Done      bool     `json:"done"`
	Question  string   `json:"question"`
	Options   []string `json:"options"`
	HasMCQ    bool     `json:"hasMcq"`
	MinLen    int      `json:"minLen"`
}

// Writeup returns the write-up form state for one problem.
func (e *Engine) Writeup(lang, id string) WriteupView {
	e.mu.Lock()
	defer e.mu.Unlock()
	p, ok := e.probByID[id]
	if !ok {
		return WriteupView{ProblemID: id}
	}
	sl := e.getSlot(lang)
	v := WriteupView{
		ProblemID: id, Title: p.Title, Solved: sl.solved[id],
		Done: sl.record(id).WriteupDone, MinLen: engine.MinWriteupLen,
	}
	if q, ok := engine.ComplexityMCQ(p); ok {
		v.HasMCQ = true
		v.Question = q.Question
		v.Options = q.Options
	}
	return v
}

// WriteupResult is the outcome of a write-up submission.
type WriteupResult struct {
	Accepted      bool       `json:"accepted"`
	MCQCorrect    bool       `json:"mcqCorrect"`
	TokensAwarded int        `json:"tokensAwarded"` // write-up + any category milestones
	Wallet        WalletView `json:"wallet"`
	Followup      string     `json:"followup"` // one mentor question (AI only; "" otherwise)
	Err           string     `json:"err"`
}

// SubmitWriteup grades the MCQ + text deterministically; acceptance flips the
// solve from provisional to fully banked, pays the write-up token, and pays
// any category milestone the gate newly satisfies.
func (e *Engine) SubmitWriteup(lang, id string, mcqIdx int, text string) WriteupResult {
	now := time.Now()
	e.mu.Lock()
	p, ok := e.probByID[id]
	if !ok {
		e.mu.Unlock()
		return WriteupResult{Err: "unknown problem: " + id}
	}
	sl := e.getSlot(lang)
	if !sl.solved[id] {
		e.mu.Unlock()
		return WriteupResult{Err: "solve the problem first"}
	}
	rec := sl.record(id)
	if rec.WriteupDone {
		e.mu.Unlock()
		return WriteupResult{Err: "write-up already completed"}
	}
	res := WriteupResult{MCQCorrect: true}
	if q, hasMCQ := engine.ComplexityMCQ(p); hasMCQ && mcqIdx != q.Correct {
		res.MCQCorrect = false
	}
	if !engine.WriteupTextOK(text) {
		res.Err = fmt.Sprintf("describe your approach in at least %d characters", engine.MinWriteupLen)
	}
	if !res.MCQCorrect || res.Err != "" {
		w := economy.Load(&sl.st, now)
		res.Wallet = walletView(&w, now)
		e.mu.Unlock()
		return res
	}

	before := e.gateLocked(sl)
	rec.WriteupDone = true
	rec.MCQCorrect = true
	rec.WriteupText = text
	sl.setRecord(id, rec)
	after := e.gateLocked(sl)

	w := economy.Load(&sl.st, now)
	award := economy.WriteupAward
	award += e.payMilestonesLocked(sl, before, after)
	w.Award(award)
	w.Store(&sl.st)
	_ = save.SaveLang(lang, sl.st)
	res.Accepted = true
	res.TokensAwarded = award
	res.Wallet = walletView(&w, now)
	aiOn := mentorService().AIEnabled()
	e.mu.Unlock()

	if aiOn {
		resp := mentorService().Hint(context.Background(), mentor.Request{
			Kind: mentor.KindFollowup, Lang: lang, Title: p.Title, Prompt: p.Prompt,
			Category: p.Category, Difficulty: p.Difficulty, Writeup: text,
		})
		if !resp.FellBack {
			res.Followup = resp.Text
		}
	}
	return res
}

// payMilestonesLocked awards MilestoneAward for every gate category that just
// crossed its minimum (idempotent via MilestonesAwarded). Caller holds e.mu.
func (e *Engine) payMilestonesLocked(sl *slot, before, after engine.GateProgress) int {
	awarded := map[string]bool{}
	for _, m := range sl.st.MilestonesAwarded {
		awarded[m] = true
	}
	wasMet := map[string]bool{}
	for _, c := range before.Categories {
		wasMet[c.Category] = c.Done >= c.Required
	}
	total := 0
	for _, c := range after.Categories {
		if c.Done >= c.Required && !wasMet[c.Category] && !awarded[c.Category] {
			total += economy.MilestoneAward
			sl.st.MilestonesAwarded = append(sl.st.MilestonesAwarded, c.Category)
		}
	}
	return total
}

// ── A3: graduation gate ──────────────────────────────────────────────────────

// GateCategoryView is one category row of the gate progress view.
type GateCategoryView struct {
	Category  string `json:"category"`
	Done      int    `json:"done"`
	Required  int    `json:"required"`
	Available int    `json:"available"`
}

// GateItemView is one mandatory problem row.
type GateItemView struct {
	Slug  string `json:"slug"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// GateView is the full graduation-gate progress for one language.
type GateView struct {
	Full        int                `json:"full"`
	Provisional int                `json:"provisional"`
	Target      int                `json:"target"`
	Categories  []GateCategoryView `json:"categories"`
	Mandatory   []GateItemView     `json:"mandatory"`
	CountMet    bool               `json:"countMet"`
	CatsMet     bool               `json:"catsMet"`
	MandatoryOK bool               `json:"mandatoryOk"`
	Met         bool               `json:"met"`
}

func (e *Engine) gateLocked(sl *slot) engine.GateProgress {
	return engine.Blind75Progress(e.cat.Problems, sl.solved, func(id string) bool {
		return sl.record(id).WriteupDone
	})
}

// Gate returns lang's graduation-gate progress.
func (e *Engine) Gate(lang string) GateView {
	e.mu.Lock()
	defer e.mu.Unlock()
	g := e.gateLocked(e.getSlot(lang))
	v := GateView{
		Full: g.Full, Provisional: g.Provisional, Target: g.Target,
		CountMet: g.CountMet, CatsMet: g.CatsMet, MandatoryOK: g.MandatoryOK, Met: g.Met,
	}
	for _, c := range g.Categories {
		v.Categories = append(v.Categories, GateCategoryView(c))
	}
	for _, m := range g.Mandatory {
		v.Mandatory = append(v.Mandatory, GateItemView(m))
	}
	return v
}

// ── A4: mentor management ────────────────────────────────────────────────────

// MentorBackends lists templates + every detected backend for the picker.
func (e *Engine) MentorBackends() []mentor.Status {
	return mentorService().Statuses()
}

// ProbeMentor runs the canary probe; returns an error message or "".
func (e *Engine) ProbeMentor(id string) string {
	if err := mentorService().Probe(context.Background(), id); err != nil {
		return err.Error()
	}
	return ""
}

// SelectMentor probes and persists a backend choice ("template" = offline).
func (e *Engine) SelectMentor(id string) string {
	if err := mentorService().Select(context.Background(), id); err != nil {
		return err.Error()
	}
	return ""
}

// SetMentorEndpoint stores openai-compat connection details (and the model
// override shared by Ollama/CLI backends).
func (e *Engine) SetMentorEndpoint(endpoint, model, apiKey string) string {
	if err := mentorService().SetEndpoint(endpoint, model, apiKey); err != nil {
		return err.Error()
	}
	return ""
}

// MentorPreview shows exactly what a hint request would send (transparency).
func (e *Engine) MentorPreview(lang, id string, tier int, code string) string {
	e.mu.Lock()
	p, ok := e.probByID[id]
	rec := e.getSlot(lang).record(id)
	e.mu.Unlock()
	if !ok {
		return ""
	}
	kind := mentor.KindStrategy
	if tier == economy.TierWalkthrough {
		kind = mentor.KindWalkthrough
	}
	return mentorService().Preview(mentor.Request{
		Kind: kind, Lang: lang, Title: p.Title, Prompt: p.Prompt,
		Category: p.Category, Difficulty: p.Difficulty, PlayerCode: code,
		FailedRuns: rec.FailedRuns,
	})
}
