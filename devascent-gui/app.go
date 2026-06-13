package main

import (
	"context"

	"devascent/guiapi"
)

// App is the Wails backend. Its exported methods are bound to the JS frontend;
// they delegate to the core engine via the guiapi facade. Orientation/tutorial
// sessions are held here (the GUI drives them via Start/Submit/Grade calls).
type App struct {
	ctx    context.Context
	engine *guiapi.Engine
	orient *guiapi.Orientation
	tut    *guiapi.Tutorial
	devlit *guiapi.DevLiteracy
}

// NewApp creates a new App application struct.
func NewApp() *App {
	return &App{}
}

// startup loads the content catalog + grader once the app context is ready.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	e, err := guiapi.New()
	if err != nil {
		println("guiapi load error:", err.Error())
		return
	}
	a.engine = e
}

// GradedLanguages returns the languages the GUI can grade in.
func (a *App) GradedLanguages() []string { return guiapi.GradedLanguages() }

// GetLanguages lists every offered language for the picker, including the
// reference-only ones (Graded=false ⇒ browse/read works, Run is disabled).
// Catalog-only — no engine state, so no nil-guard is needed.
func (a *App) GetLanguages() []guiapi.LangInfo { return guiapi.Languages() }

// GetVersion returns the build version (ldflags-injected at release time).
func (a *App) GetVersion() string { return version }

// ListProblems returns the bench browse list with lang's banked marks.
func (a *App) ListProblems(lang string) []guiapi.ProblemSummary {
	if a.engine == nil {
		return nil
	}
	return a.engine.Problems(lang)
}

// GetProblem opens one problem with a language-native starter.
func (a *App) GetProblem(id, lang string) guiapi.ProblemDetail {
	if a.engine == nil {
		return guiapi.ProblemDetail{ID: id, Lang: lang}
	}
	return a.engine.Problem(id, lang)
}

// Grade runs the player's code against the problem's hidden tests.
func (a *App) Grade(lang, id, code string) guiapi.GradeResult {
	if a.engine == nil {
		return guiapi.GradeResult{Err: "engine not loaded"}
	}
	return a.engine.Grade(lang, id, code)
}

// ── Capability gating ─────────────────────────────────────────────────────────

// GetLangPresence runs the fast toolchain presence sweep (selector marks).
func (a *App) GetLangPresence() []guiapi.LangStatus {
	if a.engine == nil {
		return nil
	}
	return a.engine.LangPresence()
}

// CheckLang runs the authoritative capability canary for lang (cached).
func (a *App) CheckLang(lang string) guiapi.LangStatus {
	if a.engine == nil {
		return guiapi.LangStatus{Lang: lang, Status: "unknown"}
	}
	return a.engine.CheckLang(lang)
}

// RecheckLang re-probes lang after an install attempt.
func (a *App) RecheckLang(lang string) guiapi.LangStatus {
	if a.engine == nil {
		return guiapi.LangStatus{Lang: lang, Status: "unknown"}
	}
	return a.engine.RecheckLang(lang)
}

// GetInstallGuide returns lang's install guide resolved for this OS.
func (a *App) GetInstallGuide(lang string) guiapi.InstallGuideView {
	if a.engine == nil {
		return guiapi.InstallGuideView{Lang: lang}
	}
	return a.engine.InstallGuideFor(lang)
}

// ── Advanced Topics ───────────────────────────────────────────────────────────

// GetAdvancedTopics lists lang's advanced topics.
func (a *App) GetAdvancedTopics(lang string) []guiapi.AdvTopicSummary {
	if a.engine == nil {
		return nil
	}
	return a.engine.AdvancedTopics(lang)
}

// GetAdvancedTopic returns one topic with sections + exercises.
func (a *App) GetAdvancedTopic(lang string, idx int) guiapi.AdvTopicDetail {
	if a.engine == nil {
		return guiapi.AdvTopicDetail{Index: idx, Lang: lang}
	}
	return a.engine.AdvancedTopic(lang, idx)
}

// GradeAdvanced grades the player's fix for one advanced exercise.
func (a *App) GradeAdvanced(lang string, topicIdx, exIdx int, code string) guiapi.GradeResult {
	if a.engine == nil {
		return guiapi.GradeResult{Err: "engine not loaded"}
	}
	return a.engine.GradeAdvanced(lang, topicIdx, exIdx, code)
}

// DeleteProfile removes lang's save slot (destructive; the UI confirms first).
func (a *App) DeleteProfile(lang string) string {
	if a.engine == nil {
		return "engine not loaded"
	}
	return a.engine.DeleteProfile(lang)
}

// GetPrimer returns the category primer (Learn drawer) for lang.
func (a *App) GetPrimer(category, lang string) guiapi.PrimerView {
	if a.engine == nil {
		return guiapi.PrimerView{}
	}
	return a.engine.PrimerFor(category, lang)
}

// GetProgress returns lang's scorecard (banked vs Step 0 targets + profile).
func (a *App) GetProgress(lang string) guiapi.Progress {
	if a.engine == nil {
		return guiapi.Progress{}
	}
	return a.engine.Progress(lang)
}

// GetProfiles lists every language save slot, most recently played first.
func (a *App) GetProfiles() []guiapi.ProfileView {
	if a.engine == nil {
		return nil
	}
	return a.engine.Profiles()
}

// NextProblem returns the first problem unsolved in lang after afterID
// (catalog order, wrapping), or "" when everything is solved.
func (a *App) NextProblem(lang, afterID string) string {
	if a.engine == nil {
		return ""
	}
	return a.engine.NextProblem(lang, afterID)
}

// ── Track A: hint economy, write-up gate, graduation gate, mentor ─────────────

// GetWallet returns lang's hint-currency wallet.
func (a *App) GetWallet(lang string) guiapi.WalletView {
	if a.engine == nil {
		return guiapi.WalletView{}
	}
	return a.engine.Wallet(lang)
}

// RequestHint serves a tier-1 nudge or paid tier-2/3 hint (paid tiers may
// block up to 45s when an AI mentor is configured — the frontend awaits it).
func (a *App) RequestHint(lang, id string, tier int, code string) guiapi.HintResult {
	if a.engine == nil {
		return guiapi.HintResult{Err: "engine not loaded"}
	}
	return a.engine.RequestHint(lang, id, tier, code)
}

// GetWriteup returns the write-up form state for one problem.
func (a *App) GetWriteup(lang, id string) guiapi.WriteupView {
	if a.engine == nil {
		return guiapi.WriteupView{ProblemID: id}
	}
	return a.engine.Writeup(lang, id)
}

// SubmitWriteup grades the MCQ + free text; acceptance banks the solve fully.
func (a *App) SubmitWriteup(lang, id string, mcqIdx int, text string) guiapi.WriteupResult {
	if a.engine == nil {
		return guiapi.WriteupResult{Err: "engine not loaded"}
	}
	return a.engine.SubmitWriteup(lang, id, mcqIdx, text)
}

// GetGate returns lang's graduation-gate (Blind 75) progress.
func (a *App) GetGate(lang string) guiapi.GateView {
	if a.engine == nil {
		return guiapi.GateView{}
	}
	return a.engine.Gate(lang)
}

// MentorStatusView mirrors mentor.Status for the JS binding (this module
// sits outside the devascent tree, so it can't import internal/mentor).
type MentorStatusView struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Present  bool   `json:"present"`
	Info     string `json:"info"`
	Selected bool   `json:"selected"`
	Probed   bool   `json:"probed"`
	ProbeOK  bool   `json:"probeOk"`
	ProbeErr string `json:"probeErr"`
}

// GetMentorBackends lists templates + every detected AI backend.
func (a *App) GetMentorBackends() []MentorStatusView {
	if a.engine == nil {
		return nil
	}
	var out []MentorStatusView
	for _, s := range a.engine.MentorBackends() {
		out = append(out, MentorStatusView{
			ID: s.ID, Name: s.Name, Present: s.Present, Info: s.Info,
			Selected: s.Selected, Probed: s.Probed, ProbeOK: s.ProbeOK, ProbeErr: s.ProbeErr,
		})
	}
	return out
}

// ProbeMentor runs the canary round-trip (may take ~90s); error message or "".
func (a *App) ProbeMentor(id string) string {
	if a.engine == nil {
		return "engine not loaded"
	}
	return a.engine.ProbeMentor(id)
}

// SelectMentor probes then persists a backend choice; error message or "".
func (a *App) SelectMentor(id string) string {
	if a.engine == nil {
		return "engine not loaded"
	}
	return a.engine.SelectMentor(id)
}

// SetMentorEndpoint stores openai-compat connection details (+model override).
func (a *App) SetMentorEndpoint(endpoint, model, apiKey string) string {
	if a.engine == nil {
		return "engine not loaded"
	}
	return a.engine.SetMentorEndpoint(endpoint, model, apiKey)
}

// GetMentorPreview returns the exact prompt a hint request would send.
func (a *App) GetMentorPreview(lang, id string, tier int, code string) string {
	if a.engine == nil {
		return ""
	}
	return a.engine.MentorPreview(lang, id, tier, code)
}

// ── Orientation (entrance test) ───────────────────────────────────────────────

// StartOrientation begins a new entrance test for lang at the given self-report
// level ("never" | "a-little" | "regularly") and returns the first item.
func (a *App) StartOrientation(lang, level string) guiapi.OrientationStep {
	if a.engine == nil {
		return guiapi.OrientationStep{Done: true}
	}
	a.orient = a.engine.StartOrientation(lang, level)
	return a.orient.Step()
}

// SubmitOrientationCode grades a code item and returns the outcome + next step.
func (a *App) SubmitOrientationCode(code string) guiapi.DiagOutcome {
	if a.orient == nil {
		return guiapi.DiagOutcome{}
	}
	return a.orient.SubmitCode(code)
}

// AdvanceOrientation commits the current code item's latest grade and moves to
// the next step (Continue after a pass, or Skip on a failing item).
func (a *App) AdvanceOrientation() guiapi.DiagOutcome {
	if a.orient == nil {
		return guiapi.DiagOutcome{}
	}
	return a.orient.AdvanceOrientation()
}

// SubmitOrientationChoice grades a multiple-choice item by option index.
func (a *App) SubmitOrientationChoice(idx int) guiapi.DiagOutcome {
	if a.orient == nil {
		return guiapi.DiagOutcome{}
	}
	return a.orient.SubmitChoice(idx)
}

// SubmitOrientationSpec grades a free-text spec item.
func (a *App) SubmitOrientationSpec(text string) guiapi.DiagOutcome {
	if a.orient == nil {
		return guiapi.DiagOutcome{}
	}
	return a.orient.SubmitSpec(text)
}

// ── Dev-Literacy track ────────────────────────────────────────────────────────

// StartDevLiteracy begins a 5-task dev-literacy session and returns the first task.
func (a *App) StartDevLiteracy() guiapi.DevLitStep {
	if a.engine == nil {
		return guiapi.DevLitStep{Done: true}
	}
	a.devlit = a.engine.StartDevLiteracy()
	return a.devlit.Step()
}

// SubmitDevCommand grades a typed command and returns the outcome + next task.
func (a *App) SubmitDevCommand(ans string) guiapi.DevLitOutcome {
	if a.devlit == nil {
		return guiapi.DevLitOutcome{}
	}
	return a.devlit.Submit(ans)
}

// ── Tutorial Island ───────────────────────────────────────────────────────────

// StartTutorial loads lang's lessons and returns how many there are.
func (a *App) StartTutorial(lang string) int {
	if a.engine == nil {
		return 0
	}
	a.tut = a.engine.StartTutorial(lang)
	return a.tut.Count()
}

// GetLesson returns the i-th lesson (0-based) with its stages.
func (a *App) GetLesson(i int) guiapi.LessonView {
	if a.tut == nil {
		return guiapi.LessonView{}
	}
	return a.tut.Lesson(i)
}

// ResumeTutorial returns the persisted frontier for the active session's lang.
func (a *App) ResumeTutorial() guiapi.TutorialPos {
	if a.tut == nil {
		return guiapi.TutorialPos{}
	}
	return a.tut.Resume()
}

// AdvanceTutorial moves the frontier forward (lessonIdx == count → complete);
// revisits never regress. Returns the resulting frontier.
func (a *App) AdvanceTutorial(lessonIdx, stageIdx int) guiapi.TutorialPos {
	if a.tut == nil {
		return guiapi.TutorialPos{}
	}
	return a.tut.Advance(lessonIdx, stageIdx)
}

// GradeLessonStage grades a we_do/you_do stage's task with the player's code.
func (a *App) GradeLessonStage(lessonIdx, stageIdx int, code string) guiapi.GradeResult {
	if a.tut == nil {
		return guiapi.GradeResult{Err: "no tutorial session"}
	}
	return a.tut.GradeStage(lessonIdx, stageIdx, code)
}
