package guiapi

import (
	"devascent/internal/content"
	"devascent/internal/save"
)

// ── Tutorial Island ───────────────────────────────────────────────────────────
// A stepper over per-language lessons. Progress persists into the language's
// save slot using the SAME fields the TUI's linear resume reads (Stage /
// LessonIdx / StageIdx), so a run started in one frontend continues in the
// other. The persisted position is a MONOTONIC FRONTIER: revisiting earlier
// stages never regresses it, and Stage is only moved while the run is at or
// before the tutorial phase (a bench player browsing a lesson keeps their
// Stage).

// LessonStageView is one stage of a lesson (i_do read-only, or we_do/you_do with
// a gradeable task).
type LessonStageView struct {
	Kind     string `json:"kind"` // i_do | we_do | you_do
	Title    string `json:"title"`
	Body     string `json:"body"`
	HasTask  bool   `json:"hasTask"`
	Prompt   string `json:"prompt"`
	FuncName string `json:"funcName"`
	Starter  string `json:"starter"`
}

// LessonView is one lesson with its stages.
type LessonView struct {
	ID     string            `json:"id"`
	Title  string            `json:"title"`
	Index  int               `json:"index"` // 1-based among the lang's lessons
	Total  int               `json:"total"`
	Found  bool              `json:"found"`
	Stages []LessonStageView `json:"stages"`
}

// TutorialPos is the persisted frontier: the furthest (lesson, stage) reached.
// Done means the final lesson's final stage was advanced past.
type TutorialPos struct {
	Lesson int  `json:"lesson"` // 0-based
	Stage  int  `json:"stage"`  // 0-based
	Done   bool `json:"done"`
}

// Tutorial is a live Tutorial-Island session for one language.
type Tutorial struct {
	e       *Engine
	lang    string
	lessons []content.Lesson
}

// StartTutorial loads the lang's ordered lessons.
func (e *Engine) StartTutorial(lang string) *Tutorial {
	return &Tutorial{e: e, lang: lang, lessons: e.cat.LessonsForLang(lang)}
}

// Count is the number of lessons.
func (t *Tutorial) Count() int { return len(t.lessons) }

// Lesson returns the i-th lesson (0-based) as a view.
func (t *Tutorial) Lesson(i int) LessonView {
	if i < 0 || i >= len(t.lessons) {
		return LessonView{Total: len(t.lessons), Found: false}
	}
	l := t.lessons[i]
	v := LessonView{ID: l.ID, Title: l.Title, Index: i + 1, Total: len(t.lessons), Found: true}
	for _, st := range l.Stages {
		sv := LessonStageView{Kind: st.Kind, Title: st.Title, Body: st.Body}
		if st.Task != nil {
			sv.HasTask = true
			sv.Prompt = st.Task.Prompt
			sv.FuncName = st.Task.FuncName
			sv.Starter = st.Task.Starter
		}
		v.Stages = append(v.Stages, sv)
	}
	return v
}

// GradeStage grades a we_do/you_do stage's task with the player's code.
func (t *Tutorial) GradeStage(lessonIdx, stageIdx int, code string) GradeResult {
	if lessonIdx < 0 || lessonIdx >= len(t.lessons) {
		return GradeResult{Err: "no such lesson"}
	}
	l := t.lessons[lessonIdx]
	if stageIdx < 0 || stageIdx >= len(l.Stages) {
		return GradeResult{Err: "no such stage"}
	}
	return t.e.gradeTask(t.lang, code, l.Stages[stageIdx].Task)
}

// stageWritable reports whether tutorial progress may move the run's Stage —
// never regress a run that is already past the tutorial phase.
func stageWritable(s string) bool {
	return s == "" || s == "intake" || s == "tutorial"
}

// Resume returns the persisted frontier for the language's slot.
func (t *Tutorial) Resume() TutorialPos {
	t.e.mu.Lock()
	defer t.e.mu.Unlock()
	return t.frontierLocked(t.e.getSlot(t.lang))
}

func (t *Tutorial) frontierLocked(sl *slot) TutorialPos {
	n := len(t.lessons)
	if n == 0 {
		return TutorialPos{}
	}
	li, si := sl.st.LessonIdx, sl.st.StageIdx
	if li >= n {
		return TutorialPos{Lesson: n, Stage: 0, Done: true}
	}
	if li < 0 {
		li = 0
	}
	if max := len(t.lessons[li].Stages) - 1; si > max {
		si = max
	}
	if si < 0 {
		si = 0
	}
	return TutorialPos{Lesson: li, Stage: si}
}

// Advance moves the persisted frontier forward to (lessonIdx, stageIdx) —
// lessonIdx == Count() marks the tutorial complete. Positions at or behind the
// stored frontier are no-ops (revisiting never regresses). Returns the
// resulting frontier.
func (t *Tutorial) Advance(lessonIdx, stageIdx int) TutorialPos {
	t.e.mu.Lock()
	defer t.e.mu.Unlock()
	sl := t.e.getSlot(t.lang)
	n := len(t.lessons)
	if n == 0 {
		return TutorialPos{}
	}
	done := lessonIdx >= n
	if done {
		lessonIdx, stageIdx = n, 0
	}
	if lessonIdx > sl.st.LessonIdx || (lessonIdx == sl.st.LessonIdx && stageIdx > sl.st.StageIdx) {
		sl.st.LessonIdx, sl.st.StageIdx = lessonIdx, stageIdx
		if stageWritable(sl.st.Stage) {
			if done {
				sl.st.Stage = "done"
			} else {
				sl.st.Stage = "tutorial"
			}
		}
		// Best-effort: the in-memory frontier stands either way.
		_ = save.SaveLang(t.lang, sl.st)
	}
	return t.frontierLocked(sl)
}
