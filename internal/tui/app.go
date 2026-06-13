package tui

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"devascent/internal/content"
	"devascent/internal/economy"
	"devascent/internal/engine"
	"devascent/internal/grader"
	"devascent/internal/mentor"
	"devascent/internal/save"
	"devascent/internal/toolchain"
)

type screen int

const (
	screenHook screen = iota
	screenLanguage
	screenEditor
	screenIntro
	screenDiagnostic
	screenTestOut
	screenResults
	screenDevLiteracy
	screenLesson
	screenHandoff
	screenBenchMenu
	screenBench
	screenStep0Complete
	screenPrimer
	screenAdvancedList
	screenAdvancedTopic
	screenInstallHelp
	screenProfilePick
	screenWriteup // A1: post-solve write-up gate (MCQ + approach note)
	screenGate    // A3: Blind-75 graduation gate progress
	screenMentor  // A4: AI mentor picker
)

// Step 0 completion milestone targets — aliased from internal/engine (the
// authoritative definitions now live there with the bench-math functions).
const (
	step0BankTarget = engine.Step0BankTarget // distinct problems banked
	step0CatTarget  = engine.Step0CatTarget  // distinct categories covered
	step0HardTarget = engine.Step0HardTarget // hard problems solved
)

type taskCtx int

const (
	ctxNone taskCtx = iota
	ctxDiagnostic
	ctxLesson
	ctxBench
)

// codeTask: the shared edit→run→grade unit (diagnostic code items + lesson stages).
type codeTask struct {
	prompt   string
	funcName string
	code     string
	tests    []grader.TestCase
	shape    grader.Shape // zero for plain problems; set for node bench problems
	verdict  *grader.Verdict
}

type editorFinishedMsg struct {
	code string
	err  error
}
type gradeMsg struct{ v grader.Verdict }

// langProbesMsg signals that the background toolchain Presence sweep finished
// (the detector cache is now populated); it just triggers a picker re-render.
type langProbesMsg struct{}

// advGradeMsg carries the result of grading an Advanced-Topics exercise attempt.
type advGradeMsg struct{ v grader.Verdict }

// hintMsg carries an async mentor reply (paid hint tiers); cost is refunded
// when the AI fell back to templates.
type hintMsg struct {
	resp mentor.Response
	tier int
	cost int
}

// mentorSelectMsg carries the result of probing+selecting a mentor backend.
type mentorSelectMsg struct {
	id  string
	err error
}

// langExt is the source-file extension for a language, so the editor opens the
// player's attempt with correct syntax highlighting.
func langExt(lang string) string {
	switch lang {
	case "python":
		return ".py"
	case "javascript":
		return ".js"
	case "typescript":
		return ".ts"
	case "go":
		return ".go"
	case "java":
		return ".java"
	case "csharp":
		return ".cs"
	case "cpp":
		return ".cpp"
	case "rust":
		return ".rs"
	default:
		return ".txt"
	}
}

type Model struct {
	screen   screen
	width    int
	height   int
	quitting bool

	lang    string
	g       grader.Grader
	det     *toolchain.Detector
	langIdx int // cursor on the language picker

	// install-help screen state (ADR-0007: shown when the player needs a toolchain)
	installLang   string // language whose install guide is displayed
	installReason string // detector's reason (e.g. "javac found but compile failed")
	installReturn screen // screen to return to on esc

	// editor selection
	editorChoice     string // command line, e.g. "code -w" ("" = system default)
	editorReturn     screen // screen to return to after the picker (when not first-time)
	editorPickResume bool   // true when the picker should start the diagnostic on pick

	cat     content.Catalog
	loadErr error
	resume  *save.State
	rng     *rand.Rand // run-local RNG (time-seeded; tests pin a fixed seed)

	profiles   []save.Profile // per-language save slots (picker shows when >1)
	profIdx    int            // cursor on the profile picker
	devRevisit bool           // dev-literacy entered as bench-menu practice (returns there)

	ctx    taskCtx
	task   *codeTask
	status string

	// diagnostic state
	level   string               // self-report band: never | a-little | regularly
	intro   content.Diagnostic   // warm-up (not counted, not scored); zero value = none
	diag    []content.Diagnostic // the chosen 10-item intake for THIS run
	diagIdx int
	curDiag content.Diagnostic
	// choice-item UI
	chosen       int
	answered     bool
	itemPassed   bool
	lastFeedback string
	// modal text input (spec items + dev-literacy commands)
	input       string
	inputActive bool
	// 3 signals (codingAbility / machineLiteracy / specReading)
	codingOK, codingTotal   int
	machineOK, machineTotal int
	specOK, specTotal       int
	passedAdd               bool
	placement               string
	intakePassed            int // questions passed this intake (drives the results/aced copy)

	// dev-literacy state
	dev      []content.DevTask // the chosen set for THIS run
	curDev   content.DevTask
	devIdx   int
	devTries int

	// lesson progress
	lessons   []content.Lesson // Tutorial Island, narrowed to the session language
	les       lesson
	lessonIdx int
	stageIdx  int

	// Step 0 bench
	bench       []content.Problem
	benchIdx    int
	benchSolved int
	curProblem  content.Problem
	primer      content.Primer // the primer currently shown ([L] from a bench problem)
	primerPage  int            // section pager cursor within the current primer

	// Stage-2 Advanced Topics (bench browse → language-specific topics)
	advTopics    []content.AdvancedTopic
	advListIdx   int                        // cursor on the topic list
	advTopic     content.AdvancedTopic      // the open topic
	advPage      int                        // pager cursor within the open topic
	advReveal    bool                       // exercise reveal toggle (show bug + fix)
	advAttempt   string                     // the player's current code attempt on this exercise
	advVerdict   *grader.Verdict            // last grade result on this exercise
	advStatus    string                     // status line under the exercise
	benchMenu    []benchOption              // the browse menu built on entry
	benchMenuIdx int                        // cursor on the browse menu
	benchFilter  string                     // label of the active subset (shown on the bench header)
	solvedSet    map[string]bool            // distinct problems banked across the whole bench
	step0Done    bool                       // completion milestone reached
	probByID     map[string]content.Problem // catalog index for stats lookups

	// Track A: hint economy + write-up gate + graduation gate + mentor
	wallet       economy.Wallet
	solveRecords map[string]save.SolveRecord // per-problem hint/write-up state
	milestones   []string                    // gate categories already paid out
	nudgeUsed    map[string]int              // per-problem nudge escalation (session)
	hintMode     bool                        // hint picker open over the bench task
	hintBusy     bool                        // awaiting an async mentor reply
	hintArm      int                         // paid tier armed for confirm (0 = none)
	hintText     string                      // last hint, kept under the task
	hintNote     string                      // source/refund line under the hint
	wuQueue      []string                    // problem IDs queued for write-up
	wuQIdx       int
	wuFromBench  bool // write-up entered from a fresh solve (continue bench after)
	wuMCQ        engine.MCQ
	wuHasMCQ     bool
	wuSel        int
	wuPhase      int // 0 = MCQ, 1 = approach note
	wuErr        string
	wuAward      int // tokens from the bank that opened this write-up (flash)
	mentorRows   []mentor.Status
	mentorIdx    int
	mentorBusy   bool
	mentorNote   string
}

func New() Model {
	det := toolchain.New()
	m := Model{screen: screenHook, det: det, g: grader.New(det), lang: "python",
		rng: rand.New(rand.NewSource(time.Now().UnixNano()))}
	cat, err := content.Load()
	if err != nil {
		m.loadErr = err
		return m
	}
	m.cat = cat
	m.solvedSet = map[string]bool{}
	m.solveRecords = map[string]save.SolveRecord{}
	m.nudgeUsed = map[string]int{}
	m.probByID = map[string]content.Problem{}
	for _, p := range cat.Problems {
		m.probByID[p.ID] = p
	}
	// Resume the most recently played language slot (single-language players get
	// exactly the old behavior; multi-language players resume the last one).
	if s, err := save.LoadLatest(); err == nil && s != nil && s.Stage != "" {
		m.resume = s // offer resume for any prior run (incl. a finished one → its end screen)
	}
	// With more than one slot, [c] opens a profile picker instead of silently
	// resuming the latest (saves are one-per-language, shared with the GUI).
	if ps, err := save.Profiles(); err == nil {
		m.profiles = ps
	}
	return m
}

func (m Model) Init() tea.Cmd {
	// Pre-warm the grader in the background (download + compile) while the player
	// reads the intro screens, so the first real submission isn't a cold stall.
	g := m.g
	warm := func() tea.Msg {
		if w, ok := g.(grader.Warmer); ok {
			w.Warm()
		}
		return nil
	}
	// In parallel, sweep toolchain Presence so the language picker can show
	// ✓/✗ availability without blocking (ADR-0007).
	return tea.Batch(warm, m.probeLangsCmd())
}

// probeLangsCmd runs a fast Presence sweep over the offered languages, populating
// the (shared) detector cache, then re-renders. Presence is cheap (LookPath +
// --version); the authoritative Capability canary still runs at grade time.
func (m Model) probeLangsCmd() tea.Cmd {
	det := m.det
	langs := m.availableLangs()
	return func() tea.Msg {
		for _, o := range langs {
			det.Presence(o.key)
		}
		return langProbesMsg{}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil
	case editorFinishedMsg:
		if m.screen == screenAdvancedTopic { // editing an Advanced-Topics exercise
			if msg.err == nil {
				m.advAttempt = msg.code
				m.advVerdict = nil
				m.advStatus = "Code updated. Press [r] to grade."
			} else {
				m.advStatus = "Editor error: " + msg.err.Error()
			}
			return m, nil
		}
		if m.task != nil {
			if msg.err == nil {
				m.task.code = msg.code
				m.task.verdict = nil
				m.status = "Code updated. Press [r] to run."
			} else {
				m.status = "Editor error: " + msg.err.Error()
			}
		}
		return m, nil
	case advGradeMsg:
		v := msg.v
		m.advVerdict = &v
		switch {
		case v.Err != "":
			m.advStatus = "✗ " + v.Err
		case v.Passed:
			m.advStatus = "✓ Solved! Your code grades correctly. [b] to compare with the model fix."
		default:
			m.advStatus = "Not yet — edit ([e]) and grade again ([r]). [b] reveals the fix."
		}
		return m, nil
	case gradeMsg:
		if m.task != nil {
			v := msg.v
			m.task.verdict = &v
			switch {
			case v.Err != "":
				m.status = "Error: " + v.Err
			case v.Passed:
				m.status = "✓ All tests passed! Press [enter] to continue."
			default:
				m.status = "Some tests failed — edit ([e]) and run again ([r])."
			}
			// Pity-rule bookkeeping: failed bench attempts accumulate toward the
			// one-time free strategy hint.
			if m.ctx == ctxBench && !v.Passed && !m.solvedSet[m.curProblem.ID] {
				m.recordFail(m.curProblem.ID, m.task.code)
			}
		}
		return m, nil
	case hintMsg:
		m.hintBusy = false
		m.hintText = msg.resp.Text
		m.hintNote = "from " + msg.resp.Source
		if msg.resp.FellBack && msg.cost > 0 {
			m.wallet.Refund(msg.cost)
			m.hintNote = "mentor unavailable — answered from the playbook, token refunded"
			m.persist()
		}
		return m, nil
	case mentorSelectMsg:
		m.mentorBusy = false
		if msg.err != nil {
			m.mentorNote = "✗ " + msg.err.Error()
		} else {
			m.mentorNote = "✓ mentor set: " + msg.id
		}
		m.mentorRows = tuiMentor().Statuses()
		return m, nil
	case langProbesMsg:
		return m, nil // detector cache now populated; re-render the picker
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// MODAL TEXT INPUT (spec + dev-literacy): must intercept BEFORE the global
	// quit/command keys, or a typed answer containing q/e/r/s/x/digits would
	// trigger them. Esc is the explicit exit hatch.
	if m.inputActive {
		return m.handleInputKey(msg)
	}
	if msg.String() == "ctrl+c" || msg.String() == "q" {
		m.persist()
		m.quitting = true
		return m, tea.Quit
	}
	switch m.screen {
	case screenHook:
		switch msg.String() {
		case "c":
			if len(m.profiles) > 1 {
				m.profIdx = 0
				m.screen = screenProfilePick
				return m, nil
			}
			if m.resume != nil {
				return m.applyResume()
			}
		case "enter":
			m.resume = nil // start fresh (saves overwrite as you progress)
			m.langIdx = 0
			m.screen = screenLanguage
			return m, m.probeLangsCmd()
		}
	case screenLanguage:
		langs := m.availableLangs()
		if m.langIdx >= len(langs) {
			m.langIdx = 0
		}
		s := msg.String()
		switch s {
		case "up", "k":
			if m.langIdx > 0 {
				m.langIdx--
			}
			return m, nil
		case "down", "j":
			if m.langIdx < len(langs)-1 {
				m.langIdx++
			}
			return m, nil
		case "i": // install help for the highlighted language (Axis A)
			o := langs[m.langIdx]
			return m.openInstallHelp(o.key, m.det.Get(o.key).Reason, screenLanguage)
		case "enter": // select the highlighted language (defaults to Python at idx 0)
			m.lang = langs[m.langIdx].key
			m.screen = screenEditor
			m.editorPickResume = true
			return m, nil
		}
		if len(s) == 1 && s[0] >= '1' && int(s[0]-'1') < len(langs) {
			m.langIdx = int(s[0] - '1')
			m.lang = langs[m.langIdx].key
			m.screen = screenEditor // pick an editor before the first code item
			m.editorPickResume = true
			return m, nil
		}
	case screenEditor:
		s := msg.String()
		if s == "enter" { // keep current/system default and proceed
			return m.finishEditorPick()
		}
		if len(s) == 1 && s[0] >= '1' && int(s[0]-'1') < len(editorOpts) {
			m.editorChoice = editorOpts[s[0]-'1'].cmd
			return m.finishEditorPick()
		}
	case screenTestOut:
		switch msg.String() {
		case "enter", "y":
			return m.startBench() // aced → straight to the bench (one hop)
		case "t":
			return m.startTutorial()
		}
	case screenResults:
		switch msg.String() {
		case "enter": // accept the suggestion
			return m.acceptPlacement()
		case "r": // redo at the same level with fresh questions
			return m.redoIntake()
		}
	case screenIntro:
		return m.handleIntroKey(msg)
	case screenDevLiteracy:
		return m.handleDevKey(msg)
	case screenDiagnostic:
		return m.handleDiagKey(msg)
	case screenProfilePick:
		switch msg.String() {
		case "up", "k":
			if m.profIdx > 0 {
				m.profIdx--
			}
		case "down", "j":
			if m.profIdx < len(m.profiles)-1 {
				m.profIdx++
			}
		case "esc":
			m.screen = screenHook
		case "enter":
			if m.profIdx >= 0 && m.profIdx < len(m.profiles) {
				if st, err := save.LoadLang(m.profiles[m.profIdx].Lang); err == nil && st != nil {
					m.resume = st
					return m.applyResume()
				}
			}
		}
		return m, nil
	case screenLesson:
		// "test me": skip ahead to this lesson's final hands-on stage — passing
		// it advances to the next lesson exactly like working through the stages.
		if msg.String() == "t" && m.ctx == ctxLesson && !m.inputActive && m.stageIdx < len(m.les.stages)-1 {
			m.stageIdx = len(m.les.stages) - 1
			m.enterStage()
			m.persist()
			return m, nil
		}
		return m.handleTaskKey(msg)
	case screenHandoff:
		switch msg.String() {
		case "b": // enter the Step 0 bench
			return m.startBench()
		case "t": // replay Tutorial Island
			return m.startTutorial()
		case "r": // start a fresh run from the intake
			return m.restart()
		}
	case screenBenchMenu:
		return m.handleBenchMenuKey(msg)
	case screenBench:
		if m.hintMode { // hint picker is modal over the task
			return m.handleHintKey(msg)
		}
		if msg.String() == "m" { // back to the browse menu
			return m.startBench()
		}
		if msg.String() == "l" { // learn: show this pattern's primer
			return m.showPrimer()
		}
		if msg.String() == "h" && m.task != nil { // hints (A2)
			m.ensureWallet()
			m.hintMode = true
			m.hintArm = 0
			return m, nil
		}
		return m.handleTaskKey(msg)
	case screenWriteup:
		return m.handleWriteupKey(msg)
	case screenGate:
		if msg.String() == "esc" || msg.String() == "m" || msg.String() == "enter" {
			return m.startBench()
		}
		return m, nil
	case screenMentor:
		return m.handleMentorKey(msg)
	case screenStep0Complete:
		if msg.String() == "enter" { // keep practicing
			return m.startBench()
		}
	case screenPrimer:
		switch msg.String() {
		case "enter", "esc":
			m.screen = screenBench // back to the problem
		case "right", "l", "n", " ":
			if m.primerPage < m.primerPageCount()-1 {
				m.primerPage++
			}
		case "left", "h", "p":
			if m.primerPage > 0 {
				m.primerPage--
			}
		}
		return m, nil
	case screenAdvancedList:
		switch msg.String() {
		case "up", "k":
			if m.advListIdx > 0 {
				m.advListIdx--
			}
		case "down", "j":
			if m.advListIdx < len(m.advTopics)-1 {
				m.advListIdx++
			}
		case "enter":
			return m.openAdvancedTopic(m.advListIdx)
		case "m", "esc":
			return m.startBench() // back to the bench browse menu
		}
		return m, nil
	case screenAdvancedTopic:
		ex, ok := m.currentAdvExercise()
		switch msg.String() {
		case "right", "n", " ":
			if m.advPage < m.advancedPageCount()-1 {
				m.advPage++
				m.resetAdvExercise()
			}
		case "left", "p":
			if m.advPage > 0 {
				m.advPage--
				m.resetAdvExercise()
			}
		case "b": // reveal the bug + model fix (the spot-the-bug / give-up path)
			if ok {
				m.advReveal = !m.advReveal
			}
		case "e": // edit your attempt (gradeable exercises)
			if ok && solveCheck(ex) != "" {
				start := m.advAttempt
				if start == "" {
					start = ex.BrokenCode
				}
				return m, editorCmd(start, m.editorChoice, ex.Prompt, langExt(m.advTopic.Lang))
			}
		case "c": // change editor, then return here
			if ok {
				m.editorReturn = m.screen
				m.screen = screenEditor
			}
		case "r", "enter": // grade your attempt
			if !ok || solveCheck(ex) == "" {
				if ok {
					m.advStatus = "This one is reveal-only — press [b] to see the bug & fix."
				}
				return m, nil
			}
			if strings.TrimSpace(m.advAttempt) == "" {
				m.advStatus = "Write your fix first ([e] to edit)."
				return m, nil
			}
			// Axis A: need the language's toolchain to grade. Route to install help if absent.
			if p := m.det.Capability(context.Background(), m.advTopic.Lang); p.Status != toolchain.Available {
				return m.openInstallHelp(m.advTopic.Lang, p.Reason, m.screen)
			}
			m.advStatus = "Grading…"
			return m, m.advGradeCmd(ex)
		case "esc":
			m.screen = screenAdvancedList // back to the topic list
		}
		return m, nil
	case screenInstallHelp:
		switch msg.String() {
		case "esc", "enter":
			m.screen = m.installReturn
		case "R": // re-check after installing (Invalidate forces a fresh probe)
			m.det.Invalidate(m.installLang)
			p := m.det.Capability(context.Background(), m.installLang)
			if p.Status == toolchain.Available {
				m.status = langLabel(m.installLang) + " detected — you can grade now."
				m.screen = m.installReturn
			} else {
				m.installReason = p.Reason
				m.status = "Still not detected. Make sure the install finished, then try again."
			}
		}
		return m, nil
	}
	return m, nil
}

// openInstallHelp shows the install guide for a language (ADR-0007), scoped to
// the detected OS, with the detector's reason if the toolchain is broken.
func (m Model) openInstallHelp(lang, reason string, ret screen) (tea.Model, tea.Cmd) {
	m.installLang = lang
	m.installReason = reason
	m.installReturn = ret
	m.screen = screenInstallHelp
	return m, nil
}

// osKey maps the host OS to the install-guide key.
func osKey() string {
	switch runtime.GOOS {
	case "windows":
		return "windows"
	case "darwin":
		return "macos"
	default:
		return "linux"
	}
}

// renderInstallHelp draws the install guide for m.installLang on the host OS
// (ADR-0007 — shown when the player needs a toolchain to grade).
func (m Model) renderInstallHelp() string {
	var b strings.Builder
	label := langLabel(m.installLang)
	b.WriteString(titleStyle.Render("Install "+label) + "\n\n")
	if m.installReason != "" {
		b.WriteString(dimStyle.Render(m.installReason) + "\n\n")
	}
	g, ok := m.cat.InstallGuideForLang(m.installLang)
	if !ok {
		b.WriteString("No install guide is available for " + label + " yet.\n")
		return b.String() + "\n" + dimStyle.Render("[esc] back")
	}
	if g.Notes != "" {
		b.WriteString(g.Notes + "\n\n")
	}
	st, ok := g.OS[osKey()]
	if !ok {
		b.WriteString("No steps for your OS yet — see INSTALL.md.\n")
		return b.String() + "\n" + dimStyle.Render("[esc] back")
	}
	if st.Link != "" {
		b.WriteString(okStyle.Render("Download: ") + st.Link + "\n\n")
	}
	for i, s := range st.Steps {
		b.WriteString(fmt.Sprintf("  %d. %s\n", i+1, s))
	}
	if st.Verify != "" {
		b.WriteString("\n" + dimStyle.Render("Verify with: ") + st.Verify + "\n")
	}
	b.WriteString("\n" + dimStyle.Render("[R] re-check (after installing)   ·   [esc] back"))
	return b.String()
}

// langOption is a selectable session language: storage key + display label.
type langOption struct{ key, label string }

// langCatalog is the canonical language order + display labels. The pick screen
// offers these in this order, but only the ones that actually have an authored
// primer (see availableLangs). Python is first — it's the default and the only
// fully-playable (graded) language today; the rest are reference-only for now.
var langCatalog = []langOption{
	{"python", "Python"},
	{"java", "Java"},
	{"cpp", "C++"},
	{"csharp", "C#"},
	{"javascript", "JavaScript"},
	{"typescript", "TypeScript"},
	{"rust", "Rust"},
	{"go", "Go"},
}

// availableLangs is langCatalog filtered to languages with at least one authored
// primer, so a player can never pick a language with no reference content. Python
// is always present (the original primers normalize to it).
func (m Model) availableLangs() []langOption {
	have := m.cat.PrimerLangs()
	var out []langOption
	for _, o := range langCatalog {
		if have[o.key] {
			out = append(out, o)
		}
	}
	if len(out) == 0 { // defensive: always offer Python
		out = []langOption{langCatalog[0]}
	}
	return out
}

// langLabel maps a language key to its display label (falls back to the key).
func langLabel(key string) string {
	for _, o := range langCatalog {
		if o.key == key {
			return o.label
		}
	}
	return key
}

// showPrimer opens the Learn panel for the current bench problem's category, in
// the session language (m.lang). Reference-only languages reach their primers
// here just like Python does.
func (m Model) showPrimer() (tea.Model, tea.Cmd) {
	if pr, ok := m.cat.PrimerByCategoryAndLang(m.curProblem.Category, m.lang); ok {
		m.primer = pr
		m.primerPage = 0 // always open on the overview page
		m.screen = screenPrimer
	} else {
		m.status = "No " + langLabel(m.lang) + " primer for this category yet."
	}
	return m, nil
}

// primerPageView is one navigable page of the Learn panel (a heading + body).
type primerPageView struct{ heading, body string }

// primerPages splits the current primer into the ordered, navigable pages shown
// one-at-a-time by the section pager: an overview (summary), one page per
// section, then the worked example. Paging by section keeps a rich primer
// readable — no page should overflow a normal terminal.
func (m Model) primerPages() []primerPageView {
	pr := m.primer
	pages := []primerPageView{{heading: "Overview", body: strings.TrimRight(pr.Summary, "\n")}}
	for _, sec := range pr.Sections {
		var sb strings.Builder
		for _, op := range sec.Ops {
			sb.WriteString("  " + dimStyle.Render(op.Label) + "\n")
			sb.WriteString(highlightCode(indentLines(op.Code, "      "), pr.Lang) + "\n\n")
		}
		pages = append(pages, primerPageView{heading: sec.Title, body: strings.TrimRight(sb.String(), "\n")})
	}
	if strings.TrimSpace(pr.Example) != "" {
		pages = append(pages, primerPageView{heading: "Worked example", body: renderExampleBody(pr.Example, pr.Lang)})
	}
	return pages
}

func (m Model) primerPageCount() int { return len(m.primerPages()) }

// ── Stage-2 Advanced Topics ──────────────────────────────────────────────────

// openAdvancedList shows the language's Advanced Topics list (from the bench menu).
func (m Model) openAdvancedList() (tea.Model, tea.Cmd) {
	m.advTopics = m.cat.AdvancedTopicsByLang(m.lang)
	m.advListIdx = 0
	m.screen = screenAdvancedList
	return m, nil
}

// openAdvancedTopic opens one topic in the pager (overview + sections + exercises).
func (m Model) openAdvancedTopic(idx int) (tea.Model, tea.Cmd) {
	if idx < 0 || idx >= len(m.advTopics) {
		return m, nil
	}
	m.advTopic = m.advTopics[idx]
	m.advPage = 0
	m.resetAdvExercise()
	m.screen = screenAdvancedTopic
	return m, nil
}

// advPageView is one navigable page of an open advanced topic: the overview, a
// reference section, or an exercise (ex != nil — rendered with the reveal toggle).
type advPageView struct {
	heading string
	body    string            // overview / section
	ex      *content.Exercise // non-nil for an exercise page
	exNum   int
}

// advancedPages splits the open topic into pages: overview → sections → exercises.
func (m Model) advancedPages() []advPageView {
	at := m.advTopic
	pages := []advPageView{{heading: "Overview", body: strings.TrimRight(at.Summary, "\n")}}
	for _, sec := range at.Sections {
		var sb strings.Builder
		for _, op := range sec.Ops {
			sb.WriteString("  " + dimStyle.Render(op.Label) + "\n")
			sb.WriteString(highlightCode(indentLines(op.Code, "      "), at.Lang) + "\n\n")
		}
		pages = append(pages, advPageView{heading: sec.Title, body: strings.TrimRight(sb.String(), "\n")})
	}
	for i := range at.Exercises {
		pages = append(pages, advPageView{heading: fmt.Sprintf("Exercise %d", i+1), ex: &at.Exercises[i], exNum: i + 1})
	}
	return pages
}

func (m Model) advancedPageCount() int { return len(m.advancedPages()) }

// currentAdvExercise returns the exercise on the current page, if this page is an
// exercise page.
func (m Model) currentAdvExercise() (content.Exercise, bool) {
	pages := m.advancedPages()
	if m.advPage < 0 || m.advPage >= len(pages) {
		return content.Exercise{}, false
	}
	if pages[m.advPage].ex == nil {
		return content.Exercise{}, false
	}
	return *pages[m.advPage].ex, true
}

// resetAdvExercise clears the per-exercise play state on page change.
func (m *Model) resetAdvExercise() {
	m.advReveal = false
	m.advAttempt = ""
	m.advVerdict = nil
	m.advStatus = ""
}

// solveCheck maps an exercise's authored Check (which describes the BROKEN code's
// behavior) to how the PLAYER's fix is graded: a fix-it for a compile error or a
// type error just needs to compile; tests/stdout grade as-is. "" = not gradeable.
func solveCheck(ex content.Exercise) grader.Check {
	switch ex.Check {
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

// advGradeCmd grades the player's attempt for an advanced exercise through the
// native toolchain, mapping the exercise to the right GradeRequest.
func (m Model) advGradeCmd(ex content.Exercise) tea.Cmd {
	lang, g, src := m.advTopic.Lang, m.g, m.advAttempt
	req := grader.GradeRequest{Lang: lang, Check: solveCheck(ex), Source: src}
	switch ex.Check {
	case "tests":
		req.FuncName = ex.FuncName
		req.Tests = ex.Tests
		req.Shape = ex.GraderShape()
	case "stdout":
		req.Signal = ex.Signal
	}
	return func() tea.Msg {
		v, err := g.Grade(req)
		if err != nil {
			v = grader.Verdict{Err: err.Error()}
		}
		return advGradeMsg{v: v}
	}
}

// advTopicName is the distinct short label for the topic list — the title's lead
// (before the " — tagline"). Groups can collide (Python splits Concurrency into
// threading/multiprocessing/asyncio/gil, all group "Concurrency"); titles don't.
func advTopicName(at content.AdvancedTopic) string {
	name := at.Title
	if j := strings.Index(name, " — "); j > 0 {
		name = name[:j]
	}
	if strings.TrimSpace(name) == "" {
		return at.Group
	}
	return name
}

// advGradeLabel describes how the player's fix is graded (via their own toolchain).
func advGradeLabel(ex *content.Exercise, lang string) string {
	switch ex.Check {
	case "", "none":
		return "reveal-only (no auto-grading for this one)"
	case "tests":
		return "graded by tests via your " + langLabel(lang) + " toolchain"
	case "compile-error":
		if ex.Signal != "" {
			return "your fix must compile (the broken version fails with " + ex.Signal + ")"
		}
		return "your fix must compile"
	case "compiles":
		return "your fix must type-check / compile"
	case "stdout":
		if ex.Signal != "" {
			return "your output must match: " + ex.Signal
		}
		return "graded by output (stdout)"
	default:
		return "graded (" + ex.Check + ")"
	}
}

// restart resets run state and returns to language pick for a fresh intake.
func (m Model) restart() (tea.Model, tea.Cmd) {
	m.diag, m.dev = nil, nil
	m.diagIdx, m.devIdx = 0, 0
	m.codingOK, m.codingTotal = 0, 0
	m.machineOK, m.machineTotal = 0, 0
	m.specOK, m.specTotal = 0, 0
	m.passedAdd, m.placement = false, ""
	m.lessonIdx, m.stageIdx = 0, 0
	m.task, m.inputActive = nil, false
	m.ctx = ctxNone
	m.screen = screenLanguage
	return m, nil
}

// ── Modal text input ─────────────────────────────────────────────────────────

func (m Model) handleInputKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		return m.submitInput()
	case tea.KeyEsc, tea.KeyCtrlC:
		if m.screen == screenWriteup && msg.Type == tea.KeyEsc {
			// Leave the note field, back to the question (or [s] to skip).
			m.inputActive = false
			m.wuPhase = 0
			m.wuErr = ""
			return m, nil
		}
		m.persist()
		m.quitting = true
		return m, tea.Quit
	case tea.KeyBackspace, tea.KeyDelete:
		r := []rune(m.input)
		if len(r) > 0 {
			m.input = string(r[:len(r)-1])
		}
	case tea.KeySpace:
		m.input += " "
	case tea.KeyRunes:
		m.input += string(msg.Runes)
	}
	return m, nil
}

func (m Model) submitInput() (tea.Model, tea.Cmd) {
	ans := strings.TrimSpace(m.input)
	if ans == "" {
		m.status = "Type an answer first."
		return m, nil
	}
	if m.screen == screenWriteup { // approach note (A1)
		return m.submitWriteupText(ans)
	}
	if m.screen == screenDiagnostic { // spec item
		passed := specMatch(ans, m.curDiag.Spec)
		m.inputActive = false
		m.answered = true
		m.itemPassed = passed
		if m.curDiag.Spec != nil {
			m.lastFeedback = m.curDiag.Spec.Answer
		}
		m.status = "Press [enter] to continue."
		return m, nil
	}
	// dev-literacy command
	if devMatch(ans, m.curDev) {
		m.inputActive = false
		m.answered = true
		m.status = m.curDev.Success + "  —  Press [enter] to continue."
		return m, nil
	}
	m.devTries++
	m.input = ""
	hint := m.curDev.Hint
	if m.devTries < 2 || hint == "" {
		m.status = "Not quite — try again."
	} else {
		m.status = "Not quite. Hint: " + hint
	}
	return m, nil
}

// specMatch: each Required entry is a synonym GROUP ("list|array|collection");
// the answer must contain at least one synonym from every group. This avoids
// false-negatives on correct paraphrases ("array"/"maximum" vs "list"/"largest").
func specMatch(ans string, sp *content.Spec) bool {
	if sp == nil {
		return true
	}
	return engine.SpecMatch(ans, sp.Required)
}

// devMatch grades a typed command against a DevTask: PASS if the first token is
// one of Commands AND every Flags token appears — OR the whole line exactly
// equals an Accept form. Lenient on args/quoting so correct answers aren't
// false-rejected (this is a checker, not a shell).
func devMatch(ans string, t content.DevTask) bool {
	return engine.DevMatch(ans, t.Commands, t.Flags, t.Accept)
}

// ── Diagnostic engine (7-item ladder, 3 signals, one early-exit) ─────────────

// finishEditorPick proceeds after the editor picker: on first pick it starts the
// run; when reopened via [c] mid-task it returns to the prior screen.
func (m Model) finishEditorPick() (tea.Model, tea.Cmd) {
	if m.editorPickResume {
		m.editorPickResume = false
		if len(m.cat.Diagnostics) == 0 {
			return m.startTutorial()
		}
		return m.startDiagnostic()
	}
	m.screen = m.editorReturn
	m.persist()
	return m, nil
}

func (m Model) startDiagnostic() (tea.Model, tea.Cmd) {
	// reset signals for a fresh intake (selection happens after the warm-up,
	// since the self-report answer sets the question difficulty).
	m.diagIdx = 0
	m.codingOK, m.codingTotal = 0, 0
	m.machineOK, m.machineTotal = 0, 0
	m.specOK, m.specTotal = 0, 0
	m.passedAdd, m.placement = false, ""
	m.ctx = ctxDiagnostic
	if len(m.cat.Intro) > 0 {
		m.intro = m.cat.Intro[m.rng.Intn(len(m.cat.Intro))]
		m.chosen, m.answered, m.lastFeedback, m.level = -1, false, "", ""
		m.status = "Pick a number, then [enter]."
		m.screen = screenIntro
		return m, nil
	}
	m.level = "a-little"
	return m.beginDiagnostic()
}

// beginDiagnostic selects the 10-item intake at the self-report difficulty band
// and enters the first item.
func (m Model) beginDiagnostic() (tea.Model, tea.Cmd) {
	m.diag = selectIntake(m.cat.DiagnosticsForLang(m.lang), m.level, nil, m.rng)
	if len(m.diag) == 0 {
		return m.startTutorial()
	}
	m.screen = screenDiagnostic
	m.ctx = ctxDiagnostic
	m.diagIdx = 0
	m.enterDiagItem()
	m.persist()
	return m, nil
}

// redoIntake re-runs the intake at the SAME self-report level with fresh
// questions (the just-seen items are pushed to the back of the pool). The escape
// from "too hard for me" is [enter] accept, not endless redo.
func (m Model) redoIntake() (tea.Model, tea.Cmd) {
	exclude := map[string]bool{}
	for _, d := range m.diag {
		exclude[d.ID] = true
	}
	m.codingOK, m.codingTotal = 0, 0
	m.machineOK, m.machineTotal = 0, 0
	m.specOK, m.specTotal = 0, 0
	m.passedAdd, m.placement, m.intakePassed = false, "", 0
	m.diag = selectIntake(m.cat.DiagnosticsForLang(m.lang), m.level, exclude, m.rng)
	if len(m.diag) == 0 {
		return m.startTutorial()
	}
	m.screen = screenDiagnostic
	m.ctx = ctxDiagnostic
	m.diagIdx = 0
	m.enterDiagItem()
	m.persist()
	return m, nil
}

// acceptPlacement: the player accepted the results screen's suggestion.
func (m Model) acceptPlacement() (tea.Model, tea.Cmd) {
	if m.placement == "dev-literacy" {
		return m.startDevLiteracy()
	}
	return m.startTutorial()
}

// handleIntroKey: the warm-up choice (re-selectable; sets difficulty, not score).
func (m Model) handleIntroKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	s := msg.String()
	if !m.answered && len(s) == 1 && s[0] >= '1' && int(s[0]-'1') < len(m.intro.Choices) {
		m.chosen = int(s[0] - '1')
		m.status = "Selected " + s + ". Press [enter] to continue, or pick another."
		return m, nil
	}
	if s == "enter" {
		if !m.answered {
			if m.chosen < 0 {
				m.status = "Pick a number first."
				return m, nil
			}
			m.answered = true
			m.level = m.intro.Choices[m.chosen].Value // sets the question difficulty band
			m.lastFeedback = m.intro.Choices[m.chosen].Feedback
			m.status = "Press [enter] to begin the Entrance Test."
			return m, nil
		}
		return m.beginDiagnostic()
	}
	return m, nil
}

// enterDiagItem loads the current diagnostic item and resets per-item UI state.
func (m *Model) enterDiagItem() {
	d := m.diag[m.diagIdx]
	m.curDiag = d
	m.task = nil
	m.answered = false
	m.itemPassed = false
	m.chosen = -1
	m.lastFeedback = ""
	m.input = ""
	m.inputActive = false
	switch d.Kind {
	case "code":
		t := diagTask(d)
		m.task = &t
		m.applyLangStarter("")
		m.status = "Press [e] to write your code, [r] to run, or [s] to skip. " + editorHint(m.editorChoice)
	case "choice":
		// Shuffle a COPY of the choices so the correct answer isn't always
		// option 1 (and we never mutate the shared catalog slice). Done once on
		// entry, not in View (which re-renders every keystroke).
		shuffled := append([]content.Choice(nil), d.Choices...)
		m.rng.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })
		m.curDiag.Choices = shuffled
		m.status = "Pick a number, then [enter] to confirm."
	case "spec":
		m.inputActive = true
		m.status = "Type your answer, then [enter]."
	}
}

func (m Model) handleDiagKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.curDiag.Kind {
	case "code":
		return m.handleTaskKey(msg)
	case "choice":
		return m.handleChoiceKey(msg)
	case "spec":
		// after submit, [enter] advances (input mode handled modally above)
		if m.answered && msg.String() == "enter" {
			return m.scoreDiag(m.itemPassed)
		}
	}
	return m, nil
}

func (m Model) handleChoiceKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	s := msg.String()
	// Number keys (re)select while not yet confirmed — a mis-press is fixable.
	if !m.answered && len(s) == 1 && s[0] >= '1' && int(s[0]-'1') < len(m.curDiag.Choices) {
		m.chosen = int(s[0] - '1')
		m.status = "Selected " + s + ". Press [enter] to confirm, or pick another number."
		return m, nil
	}
	if s == "enter" {
		if !m.answered { // confirm the highlighted choice
			if m.chosen < 0 {
				m.status = "Pick a number first."
				return m, nil
			}
			m.answered = true
			c := m.curDiag.Choices[m.chosen]
			m.itemPassed = c.Correct
			m.lastFeedback = c.Feedback
			m.status = "Press [enter] to continue."
			return m, nil
		}
		return m.scoreDiag(m.itemPassed) // already confirmed → advance
	}
	return m, nil
}

// scoreDiag records the item's signal, applies the early-exit, then advances.
func (m Model) scoreDiag(passed bool) (tea.Model, tea.Cmd) {
	switch m.curDiag.Measures {
	case "coding":
		if m.codingTotal == 0 { // the first coding item is the "floor" routing reads
			m.passedAdd = passed
		}
		m.codingTotal++
		if passed {
			m.codingOK++
		}
	case "machine":
		m.machineTotal++
		if passed {
			m.machineOK++
		}
	case "spec":
		m.specTotal++
		if passed {
			m.specOK++
		}
	case "self":
		// self-report: informational only, never scored or routed on.
	}
	m.diagIdx++
	if m.diagIdx >= len(m.diag) {
		return m.route()
	}
	m.enterDiagItem()
	m.persist()
	return m, nil
}

// route computes placement from the intake SCORE (percentage), applies the
// self-report band clamp, then shows the aced screen (test-out) or the results
// screen (everyone else — suggestion + accept/redo). Single source of truth.
//
// Bands: ≥80% → bench (aced) · 40–79% → dev-literacy brush-up · <40% → tutorial.
func (m Model) route() (tea.Model, tea.Cmd) {
	passed := m.codingOK + m.machineOK + m.specOK
	m.intakePassed = passed
	m.placement = engine.Place(passed, len(m.diag), m.level)
	m.ctx = ctxNone
	m.task = nil
	if m.placement == "test-out" {
		m.screen = screenTestOut // aced screen (kept)
	} else {
		m.screen = screenResults // suggestion + accept/redo
	}
	m.persist()
	return m, nil
}

// ── Dev-Literacy track (command checker — not a real shell) ───────────────────

// devDone routes dev-literacy completion: a routed (gating) run hands off to
// the bench; a revisit from the bench menu returns to the menu.
func (m Model) devDone() (tea.Model, tea.Cmd) {
	if m.devRevisit {
		m.devRevisit = false
		return m.startBench()
	}
	return m.toHandoff()
}

func (m Model) startDevLiteracy() (tea.Model, tea.Cmd) {
	m.dev = selectDevSet(m.cat.DevTasks, 5, m.rng) // 5 across distinct categories
	if len(m.dev) == 0 {
		return m.devDone() // nothing authored → straight back
	}
	m.screen = screenDevLiteracy
	m.ctx = ctxNone
	m.task = nil
	m.devIdx = 0
	m.enterDevTask()
	m.persist()
	return m, nil
}

func (m *Model) enterDevTask() {
	m.curDev = m.dev[m.devIdx]
	m.answered = false
	m.input = ""
	m.inputActive = true
	m.devTries = 0
	m.status = "Type the command, then [enter]."
}

func (m Model) handleDevKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// non-input keys only reach here once answered (input is modal otherwise)
	if m.answered && msg.String() == "enter" {
		m.devIdx++
		if m.devIdx >= len(m.dev) {
			return m.devDone()
		}
		m.enterDevTask()
		m.persist()
		return m, nil
	}
	return m, nil
}

// ── Code-task / lesson handling (unchanged behavior) ─────────────────────────

func (m Model) handleTaskKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Tutorial Island is optional — [x] leaves it and goes to the bench.
	if msg.String() == "x" && m.ctx == ctxLesson {
		return m.toHandoff()
	}
	// Reading stage (no task): [enter] advances.
	if m.task == nil {
		if msg.String() == "enter" {
			return m.advance()
		}
		return m, nil
	}
	switch msg.String() {
	case "e":
		return m, editorCmd(m.task.code, m.editorChoice, m.task.prompt, langExt(m.lang))
	case "c": // change editor, then come back to this task
		m.editorReturn = m.screen
		m.screen = screenEditor
		return m, nil
	case "r":
		// Two-axis gate (ADR-0007 / Runtime Detection design spec). This [r] case
		// is the single choke point — bench, lessons, and code diagnostics all
		// funnel through handleTaskKey, so gating here covers every graded surface.
		//
		// Axis B (grading maturity): only Python has a grader adapter today, so
		// other languages stay reference-only regardless of whether their toolchain
		// is installed. (Widens as per-language adapters land.)
		if !gradingAvailable(m.lang) {
			status := langLabel(m.lang) + " is reference-only for now — its function-call grader isn't wired yet."
			if m.ctx == ctxBench {
				status += " Use [l] for the " + langLabel(m.lang) + " primer."
			}
			m.status = status
			return m, nil
		}
		// Axis A (runtime availability): grading shells out to the player's OWN
		// installed toolchain (ADR-0007 — DevAscent bundles none). If it isn't
		// installed/working, route to the install guide instead of failing cryptically.
		// Capability caches, so only the first grade in a session pays the canary cost.
		if p := m.det.Capability(context.Background(), m.lang); p.Status != toolchain.Available {
			return m.openInstallHelp(m.lang, p.Reason, m.screen)
		}
		if strings.TrimSpace(m.task.code) == "" {
			m.status = "Write some code first ([e])."
			return m, nil
		}
		m.status = "Running…"
		return m, m.runTask()
	case "s":
		if m.ctx == ctxDiagnostic {
			return m.scoreDiag(false)
		}
		if m.ctx == ctxBench {
			return m.benchNext(false) // skip this problem
		}
	case "enter":
		if m.task.verdict != nil && m.task.verdict.Passed {
			return m.advance()
		}
	}
	return m, nil
}

func (m Model) runTask() tea.Cmd {
	t := m.task
	lang, g := m.lang, m.g
	return func() tea.Msg {
		v, err := g.Run(lang, t.code, t.funcName, t.tests, t.shape)
		if err != nil {
			v = grader.Verdict{Err: err.Error()}
		}
		return gradeMsg{v}
	}
}

// advance: after a reading stage or a passed code task.
func (m Model) advance() (tea.Model, tea.Cmd) {
	switch m.ctx {
	case ctxDiagnostic:
		return m.scoreDiag(true)
	case ctxBench:
		return m.benchNext(true)
	case ctxLesson:
		m.stageIdx++
		if m.stageIdx >= len(m.les.stages) {
			m.lessonIdx++
			if m.lessonIdx >= len(m.lessons) {
				return m.toHandoff()
			}
			m.les = toLesson(m.lessons[m.lessonIdx])
			m.stageIdx = 0
			m.enterStage()
		} else {
			m.enterStage()
		}
	}
	m.persist()
	return m, nil
}

func (m Model) toHandoff() (tea.Model, tea.Cmd) {
	m.screen = screenHandoff
	m.ctx = ctxNone
	m.task = nil
	m.inputActive = false
	m.persist()
	return m, nil
}

// ── Step 0 bench ─────────────────────────────────────────────────────────────

// startBench shows the browse menu (pick All / a category / a list).
func (m Model) startBench() (tea.Model, tea.Cmd) {
	if len(m.cat.Problems) == 0 {
		return m, nil
	}
	m.benchMenu = benchMenuOptions(m.cat.Problems)
	if len(m.cat.AdvancedTopicsByLang(m.lang)) > 0 { // Stage-2 language-specific track
		m.benchMenu = append(m.benchMenu, benchOption{
			label: "⭐ Advanced Topics — " + langLabel(m.lang),
			kind:  "advanced",
		})
	}
	if len(m.cat.DevTasks) > 0 { // revisitable terminal drills (non-gating practice)
		m.benchMenu = append(m.benchMenu, benchOption{
			label: "🖥 Dev-Literacy practice — terminal drills",
			kind:  "devlit",
		})
	}
	// Track A entries: pending write-ups (provisional solves), the graduation
	// gate, and the AI mentor picker.
	m.ensureWallet()
	if n := len(m.provisionalIDs()); n > 0 {
		m.benchMenu = append(m.benchMenu, benchOption{
			label: fmt.Sprintf("✍ Write-ups pending (%d) — explain to bank fully (+1 token each)", n),
			kind:  "writeups",
		})
	}
	g := m.gateProgress()
	m.benchMenu = append(m.benchMenu,
		benchOption{label: fmt.Sprintf("🎓 Graduation gate — Blind 75 (%d/%d)", g.Full, g.Target), kind: "gate"},
		benchOption{label: "🤖 AI mentor — " + m.mentorLabel(), kind: "mentor"},
	)
	m.benchMenuIdx = 0
	m.screen = screenBenchMenu
	m.ctx = ctxNone
	m.task = nil
	return m, nil
}

func (m Model) handleBenchMenuKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.benchMenuIdx > 0 {
			m.benchMenuIdx--
		}
	case "down", "j":
		if m.benchMenuIdx < len(m.benchMenu)-1 {
			m.benchMenuIdx++
		}
	case "enter":
		if m.benchMenuIdx >= 0 && m.benchMenuIdx < len(m.benchMenu) {
			opt := m.benchMenu[m.benchMenuIdx]
			if opt.kind == "advanced" {
				return m.openAdvancedList()
			}
			if opt.kind == "devlit" {
				m.devRevisit = true
				return m.startDevLiteracy()
			}
			if opt.kind == "writeups" {
				return m.openWriteups(m.provisionalIDs(), false)
			}
			if opt.kind == "gate" {
				m.screen = screenGate
				return m, nil
			}
			if opt.kind == "mentor" {
				return m.openMentorPicker()
			}
			return m.startBenchFiltered(opt)
		}
	}
	return m, nil
}

// startBenchFiltered serves the chosen subset, tiered by the player's placement.
func (m Model) startBenchFiltered(opt benchOption) (tea.Model, tea.Cmd) {
	subset := filterProblems(m.cat.Problems, opt.kind, opt.value)
	if opt.kind == "list" { // one random variant per canonical slug (freshness, correct count)
		subset = dedupeBySlug(subset, m.rng)
	}
	m.bench = selectBench(subset, m.placement, m.rng)
	if len(m.bench) == 0 {
		return m, nil
	}
	m.benchFilter = opt.label
	m.benchIdx, m.benchSolved = 0, 0
	m.screen = screenBench
	m.ctx = ctxBench
	m.enterProblem()
	m.persist()
	return m, nil
}

func (m *Model) enterProblem() {
	p := m.bench[m.benchIdx]
	m.curProblem = p
	t := codeTask{prompt: p.Prompt, funcName: p.FuncName, code: p.Starter, tests: p.Tests, shape: p.GraderShape()}
	m.task = &t
	m.applyLangStarter(p.Solution)
	m.hintMode, m.hintText, m.hintNote, m.hintArm = false, "", "", 0
	m.status = "Press [e] to write your code, [r] to run, [h] for a hint, [s] to skip. " + editorHint(m.editorChoice)
}

// gradingAvailable / applyLangStarter delegate the language dispatch to
// internal/engine. applyLangStarter keeps the Model-coupled guard and only
// overwrites the starter when engine generates one (ok), so python and
// reference-only languages keep their authored starter. pySource (may be "")
// gives generated stubs nicer parameter names.
func gradingAvailable(lang string) bool {
	return engine.GradingAvailable(lang)
}

func (m *Model) applyLangStarter(pySource string) {
	if m.task == nil || m.task.funcName == "" {
		return
	}
	if code, ok := engine.Starter(m.lang, m.task.funcName, pySource, m.task.tests, m.task.shape); ok {
		m.task.code = code
	}
}

// benchNext routes a finished problem: solved banks it (provisionally) and
// opens the write-up gate; a skip just advances.
func (m Model) benchNext(solved bool) (tea.Model, tea.Cmd) {
	if solved {
		return m.benchSolvedNow()
	}
	return m.benchAdvance()
}

// benchSolvedNow banks the current problem (provisional until its write-up),
// pays the clean-solve award, and opens the write-up gate (A1) before the
// bench continues.
func (m Model) benchSolvedNow() (tea.Model, tea.Cmd) {
	m.ensureWallet()
	id := m.curProblem.ID
	if !m.solvedSet[id] {
		if m.solvedSet == nil {
			m.solvedSet = map[string]bool{}
		}
		m.solvedSet[id] = true
		m.benchSolved++
		m.wuAward = 0
		rec := m.solveRecords[id]
		if rec.HintTier == economy.TierNone && !rec.PityUsed {
			m.wuAward = economy.SolveAward(m.curProblem.Difficulty)
			m.wallet.Award(m.wuAward)
		}
		m.persist()
		if !rec.WriteupDone {
			return m.openWriteups([]string{id}, true)
		}
	}
	return m.benchContinue()
}

// benchContinue is the post-bank half: milestone screen check, then advance.
func (m Model) benchContinue() (tea.Model, tea.Cmd) {
	if !m.step0Done && m.step0Met() {
		m.step0Done = true
		m.screen = screenStep0Complete
		m.ctx = ctxNone
		m.task = nil
		m.persist()
		return m, nil
	}
	return m.benchAdvance()
}

// benchAdvance moves to the next problem in the pool.
func (m Model) benchAdvance() (tea.Model, tea.Cmd) {
	m.screen = screenBench
	m.ctx = ctxBench
	m.hintMode, m.hintText, m.hintNote = false, "", ""
	m.benchIdx++
	if m.benchIdx >= len(m.bench) {
		m.task = nil // exhausted the pool
		m.persist()
		return m, nil
	}
	m.enterProblem()
	m.persist()
	return m, nil
}

// benchStats / step0Met / step0Profile delegate to internal/engine (the bench
// math is UI-neutral and shared with the GUI); these wrappers keep the Model
// receiver so call sites and tests are unchanged.
func (m Model) benchStats() (banked, cats, hard int) {
	return engine.BenchStats(m.solvedSet, m.probByID)
}

func (m Model) step0Met() bool {
	return engine.Step0Met(m.solvedSet, m.probByID)
}

func (m Model) step0Profile() (problemSolving, langProf int, track string) {
	return engine.Step0Profile(m.solvedSet, m.probByID, m.cat.Problems)
}

func (m Model) startTutorial() (tea.Model, tea.Cmd) {
	if len(m.cat.LessonsForLang(m.lang)) == 0 {
		m.loadErr = fmt.Errorf("no lessons loaded")
		return m, nil
	}
	m.screen = screenLesson
	m.ctx = ctxLesson
	m.inputActive = false
	m.lessons = m.cat.LessonsForLang(m.lang)
	m.lessonIdx = 0
	m.les = toLesson(m.lessons[0])
	m.stageIdx = 0
	m.enterStage()
	m.persist()
	return m, nil
}

func (m *Model) enterStage() {
	st := m.les.stages[m.stageIdx]
	if st.task != nil {
		t := *st.task
		m.task = &t
		m.applyLangStarter("")
		m.status = "Press [e] to write your code, then [r] to run."
	} else {
		m.task = nil
		m.status = "Press [enter] to continue."
	}
}

func (m Model) currentState() save.State {
	s := save.State{
		Language:     m.lang,
		Editor:       m.editorChoice,
		Level:        m.level,
		Placement:    m.placement,
		CodingOK:     m.codingOK,
		CodingTotal:  m.codingTotal,
		MachineOK:    m.machineOK,
		MachineTotal: m.machineTotal,
		SpecOK:       m.specOK,
		SpecTotal:    m.specTotal,
		PassedAdd:    m.passedAdd,
		IntakePassed: m.intakePassed,
		DiagIDs:      diagIDs(m.diag),
		DevIDs:       devTaskIDs(m.dev),
		BenchIDs:     benchIDs(m.bench),
		BenchIdx:     m.benchIdx,
		BenchSolved:  m.benchSolved,
		SolvedIDs:    solvedSlice(m.solvedSet),
		Step0Done:    m.step0Done,
	}
	// Track A state rides along on every save (zero values for pre-bench runs).
	m.wallet.Store(&s)
	if len(m.solveRecords) > 0 {
		s.SolveRecords = m.solveRecords
	}
	s.MilestonesAwarded = m.milestones
	switch m.screen {
	case screenLesson:
		s.Stage = "tutorial"
		s.LessonIdx = m.lessonIdx
		s.StageIdx = m.stageIdx
	case screenDevLiteracy:
		s.Stage = "devliteracy"
		s.DevIdx = m.devIdx
	case screenStep0Complete:
		s.Stage = "step0done"
	case screenBench, screenWriteup, screenGate, screenMentor:
		s.Stage = "bench"
	case screenResults:
		s.Stage = "results"
	case screenTestOut:
		s.Stage = "aced"
	case screenHandoff:
		s.Stage = "done"
	default: // diagnostic
		s.Stage = "intake"
		s.DiagIdx = m.diagIdx
	}
	return s
}

// persist best-effort saves progress once the run has actually started.
func (m Model) persist() {
	switch m.screen {
	case screenDiagnostic, screenTestOut, screenResults, screenDevLiteracy, screenLesson, screenHandoff,
		screenBench, screenStep0Complete, screenWriteup, screenGate, screenMentor:
		_ = save.SaveLang(m.lang, m.currentState())
	}
}

// solvedSlice returns the banked problem IDs as a slice (for the save).
func solvedSlice(set map[string]bool) []string {
	if len(set) == 0 {
		return nil
	}
	out := make([]string, 0, len(set))
	for id := range set {
		out = append(out, id)
	}
	return out
}

func (m Model) applyResume() (tea.Model, tea.Cmd) {
	s := m.resume
	m.resume = nil
	if s == nil {
		m.screen = screenLanguage
		return m, nil
	}
	if s.Language != "" {
		m.lang = s.Language
	}
	m.editorChoice = s.Editor
	m.level = s.Level
	// restore signals + placement (degrade gracefully on old saves: zero values)
	m.placement = s.Placement
	m.codingOK, m.codingTotal = s.CodingOK, s.CodingTotal
	m.machineOK, m.machineTotal = s.MachineOK, s.MachineTotal
	m.specOK, m.specTotal = s.SpecOK, s.SpecTotal
	m.passedAdd = s.PassedAdd
	m.intakePassed = s.IntakePassed
	m.step0Done = s.Step0Done
	m.solvedSet = map[string]bool{}
	for _, id := range s.SolvedIDs {
		m.solvedSet[id] = true
	}
	// Track A state (zero values on pre-v4 saves; the wallet grants on load).
	m.solveRecords = map[string]save.SolveRecord{}
	for id, r := range s.SolveRecords {
		m.solveRecords[id] = r
	}
	m.milestones = s.MilestonesAwarded
	m.wallet = economy.Load(s, time.Now())
	if m.nudgeUsed == nil {
		m.nudgeUsed = map[string]int{}
	}
	switch s.Stage {
	case "step0done":
		m.ctx = ctxNone
		m.task = nil
		m.screen = screenStep0Complete
	case "results", "aced":
		// rebuild the intake (for the score total); fall back to a fresh select
		if lad, ok := diagsByIDs(m.cat.Diagnostics, s.DiagIDs); ok {
			m.diag = lad
		}
		m.ctx = ctxNone
		m.task = nil
		if s.Stage == "aced" {
			m.screen = screenTestOut
		} else {
			m.screen = screenResults
		}
	case "tutorial":
		m.ctx = ctxLesson
		m.lessons = m.cat.LessonsForLang(m.lang)
		m.lessonIdx = s.LessonIdx
		if m.lessonIdx < 0 || m.lessonIdx >= len(m.lessons) {
			m.lessonIdx = 0
		}
		m.les = toLesson(m.lessons[m.lessonIdx])
		m.stageIdx = s.StageIdx
		if m.stageIdx < 0 || m.stageIdx >= len(m.les.stages) {
			m.stageIdx = 0
		}
		m.screen = screenLesson
		m.enterStage()
		m.status = "Welcome back — " + m.status
	case "devliteracy":
		// rebuild the exact dev set; if any id is stale, re-select fresh from 0
		if set, ok := devsByIDs(m.cat.DevTasks, s.DevIDs); ok {
			m.dev = set
			m.devIdx = s.DevIdx
			if m.devIdx < 0 || m.devIdx >= len(m.dev) {
				m.devIdx = 0
			}
		} else {
			m.dev = selectDevSet(m.cat.DevTasks, 5, m.rng)
			m.devIdx = 0
		}
		if len(m.dev) == 0 {
			return m.toHandoff()
		}
		m.screen = screenDevLiteracy
		m.enterDevTask()
		m.status = "Welcome back. " + m.status
	case "bench":
		if pool, ok := problemsByIDs(m.cat.Problems, s.BenchIDs); ok {
			m.bench = pool
			m.benchIdx, m.benchSolved = s.BenchIdx, s.BenchSolved
			if m.benchIdx < 0 || m.benchIdx > len(m.bench) {
				m.benchIdx = 0
			}
		} else {
			m.bench = selectBench(m.cat.Problems, m.placement, m.rng)
			m.benchIdx, m.benchSolved = 0, 0
		}
		if len(m.bench) == 0 {
			m.screen = screenHandoff
			return m, nil
		}
		m.ctx = ctxBench
		m.screen = screenBench
		if m.benchIdx < len(m.bench) {
			m.enterProblem()
		} else {
			m.task = nil
		}
		m.status = "Welcome back. " + m.status
	case "intake":
		m.ctx = ctxDiagnostic
		// rebuild the exact ladder; if any id is stale, re-select and restart from 0
		if lad, ok := diagsByIDs(m.cat.Diagnostics, s.DiagIDs); ok {
			m.diag = lad
			m.diagIdx = s.DiagIdx
			if m.diagIdx < 0 || m.diagIdx >= len(m.diag) {
				m.diagIdx = 0
			}
		} else {
			m.diag = selectIntake(m.cat.DiagnosticsForLang(m.lang), m.level, nil, m.rng)
			m.diagIdx = 0
		}
		if len(m.diag) == 0 {
			m.screen = screenLanguage
			return m, nil
		}
		m.screen = screenDiagnostic
		m.enterDiagItem()
		if !m.inputActive {
			m.status = "Welcome back. " + m.status
		}
	default:
		m.screen = screenHandoff
	}
	return m, nil
}

// ── View ───────────────────────────────────────────────────────────────────

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63"))
	dimStyle   = lipgloss.NewStyle().Faint(true)
	okStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	errStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
	// codeStyle renders primer snippets: a bright cyan so code reads as CODE on a
	// dark terminal, not as faint/grey "comment" text (the label is dimmed instead).
	codeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("117"))
	boxStyle  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2)
)

func (m Model) View() string {
	if m.quitting {
		return "See you at the Studio.\n"
	}
	if m.loadErr != nil {
		return boxStyle.Render(errStyle.Render("Failed to load content:") + "\n\n" + m.loadErr.Error() + "\n\n" + dimStyle.Render("[q] quit"))
	}
	if m.width > 0 && (m.width < 62 || m.height < 18) {
		return fmt.Sprintf("\n  Terminal too small (%dx%d).\n  Please resize to at least 62 x 18.\n\n  [q] quit\n", m.width, m.height)
	}
	var b strings.Builder
	switch m.screen {
	case screenHook:
		b.WriteString(titleStyle.Render("DevAscent — Day One") + "\n\n")
		b.WriteString("You step into the Studio: a small dev shop that's agreed to take\n")
		b.WriteString("you on as an apprentice. Before you touch real client code, the\n")
		b.WriteString("lead wants to see where you're at.\n\n")
		if m.resume != nil {
			if len(m.profiles) > 1 {
				b.WriteString(okStyle.Render(fmt.Sprintf("Saved progress found (%d language profiles).", len(m.profiles))) + "\n")
			} else {
				b.WriteString(okStyle.Render("Saved progress found.") + "\n")
			}
			b.WriteString(dimStyle.Render("[c] continue   ·   [enter] start over   ·   [q] quit"))
		} else {
			b.WriteString(dimStyle.Render("[enter] begin   ·   [q] quit"))
		}
	case screenProfilePick:
		b.WriteString(titleStyle.Render("Choose a profile") + "\n\n")
		b.WriteString("One save slot per language — progress is shared with the desktop app.\n\n")
		for i, p := range m.profiles {
			sel := "  "
			line := fmt.Sprintf("%-12s %d banked", langLabel(p.Lang), p.Banked)
			if p.Placement != "" {
				line += " · " + p.Placement
			}
			if len(p.UpdatedAt) >= 10 {
				line += "   (" + p.UpdatedAt[:10] + ")"
			}
			if i == m.profIdx {
				sel = "› "
				line = okStyle.Render(line)
			}
			b.WriteString(sel + line + "\n")
		}
		b.WriteString("\n" + dimStyle.Render("[↑/↓] move   ·   [enter] continue   ·   [esc] back"))
	case screenLanguage:
		b.WriteString(titleStyle.Render("Pick your language") + "\n\n")
		b.WriteString("This is your language for the whole session — it sets which\n")
		b.WriteString("language the Learn primers ([l] on a problem) are written in.\n\n")
		for i, o := range m.availableLangs() {
			// Axis A — toolchain availability (from the background Presence sweep).
			iconText, iconStyle := "· checking…", dimStyle
			switch m.det.Get(o.key).Status {
			case toolchain.Available:
				iconText, iconStyle = "✓ installed", okStyle
			case toolchain.Missing:
				iconText, iconStyle = "✗ not installed", errStyle
			case toolchain.Broken:
				iconText, iconStyle = "⚠ not working", errStyle
			}
			icon := iconStyle.Render(fmt.Sprintf("%-15s", iconText))
			// Only flag the exception: a language DevAscent can't grade yet (no
			// harness — currently just C++). Every other installed language plays the
			// full game, so no per-row "graded" label is needed (the ✓ says it all).
			tag := ""
			if !gradingAvailable(o.key) {
				tag = dimStyle.Render("reference-only")
			}
			cursor := "  "
			if i == m.langIdx {
				cursor = "› "
			}
			b.WriteString(fmt.Sprintf("  %s%d. %-11s %s %s\n", cursor, i+1, o.label, icon, tag))
		}
		b.WriteString("\n" + dimStyle.Render("Pick the language you want to play in. ✓ = installed & ready. You write and\nrun your own code, so its toolchain must be installed — [i] shows how.") + "\n\n")
		b.WriteString(dimStyle.Render("[↑/↓] move   ·   [1-9] pick   ·   [i] install help   ·   [enter] select   ·   [q] quit"))
	case screenEditor:
		b.WriteString(titleStyle.Render("Pick your code editor") + "\n\n")
		b.WriteString("You'll write code in a real editor. Pick one you have installed —\n")
		b.WriteString("you can change it anytime later with [c] while coding.\n\n")
		for i, o := range editorOpts {
			avail := dimStyle.Render("(not found)")
			if editorAvailable(o) {
				avail = okStyle.Render("(found)")
			}
			b.WriteString(fmt.Sprintf("  %d. %-8s %s  %s\n", i+1, o.label, dimStyle.Render(o.cmd), avail))
		}
		b.WriteString("\n" + dimStyle.Render("notepad is fine but a poor Python editor (no indent help)."))
		b.WriteString("\n" + dimStyle.Render("[1-4] pick   ·   [enter] use system default   ·   [q] quit"))
	case screenIntro:
		b.WriteString(titleStyle.Render("Before we start") + "\n\n")
		b.WriteString(m.renderIntro())
	case screenDiagnostic:
		b.WriteString(titleStyle.Render(fmt.Sprintf("Entrance Test — question %d of %d", m.diagIdx+1, len(m.diag))) + "\n\n")
		b.WriteString(m.renderDiag())
	case screenTestOut:
		b.WriteString(titleStyle.Render(fmt.Sprintf("You aced the intake — %d of %d!", m.intakePassed, m.intakeTotal())) + "\n\n")
		b.WriteString("You can clearly already code, and you know your way around the\n")
		b.WriteString("terminal. You may skip straight to the bench, or run through\n")
		b.WriteString("Tutorial Island anyway for a warm-up.\n\n")
		b.WriteString(dimStyle.Render("[enter] skip to the bench   ·   [t] do Tutorial Island   ·   [q] quit"))
	case screenResults:
		b.WriteString(titleStyle.Render("Entrance Test — results") + "\n\n")
		b.WriteString(m.renderResults())
	case screenDevLiteracy:
		b.WriteString(titleStyle.Render(fmt.Sprintf("Dev-Literacy — %d of %d · %s", m.devIdx+1, len(m.dev), m.curDev.Title)) + "\n\n")
		b.WriteString(m.renderDev())
	case screenLesson:
		st := m.les.stages[m.stageIdx]
		b.WriteString(titleStyle.Render(fmt.Sprintf("Tutorial Island · Lesson %d/%d · %s — %s", m.lessonIdx+1, len(m.lessons), m.les.title, st.title)) + "\n\n")
		b.WriteString(st.body + "\n\n")
		if st.task != nil {
			b.WriteString(m.renderTask())
		} else {
			b.WriteString(dimStyle.Render("[enter] continue   ·   [t] test me (jump to the task)   ·   [x] skip tutorial   ·   [q] quit"))
		}
	case screenHandoff:
		b.WriteString(titleStyle.Render("Orientation complete — nice work") + "\n\n")
		b.WriteString("You can write and run real functions, and you know your way\n")
		b.WriteString("around the basics. That's everything Orientation — the Entrance\n")
		b.WriteString("Test and Tutorial Island — has for now.\n\n")
		b.WriteString("Next up is The Apprenticeship: the bench, where you solve real\n")
		b.WriteString("problems and build toward your competency profile.\n\n")
		b.WriteString(dimStyle.Render("[b] enter The Apprenticeship (the bench)   ·   [t] replay Tutorial Island\n[r] start over   ·   [q] quit"))
	case screenBenchMenu:
		b.WriteString(titleStyle.Render("The Bench — pick your practice") + "\n")
		bk, ca, hd := m.benchStats()
		prog := fmt.Sprintf("Apprenticeship progress: banked %d/%d · categories %d/%d · hard %d/%d",
			bk, step0BankTarget, ca, step0CatTarget, hd, step0HardTarget)
		if m.step0Done {
			prog = okStyle.Render("Apprenticeship milestone reached ✓ — keep practicing or [q] quit")
		}
		b.WriteString(dimStyle.Render(prog) + "\n\n")
		for i, o := range m.benchMenu {
			sel := "  "
			line := o.label
			if i == m.benchMenuIdx {
				sel = "› "
				line = okStyle.Render(o.label)
			}
			b.WriteString(sel + line + "\n")
		}
		b.WriteString("\n" + dimStyle.Render("[↑/↓] move   ·   [enter] choose   ·   [q] quit"))
	case screenBench:
		if m.task == nil { // worked through the whole subset
			b.WriteString(titleStyle.Render("Bench — set complete") + "\n\n")
			b.WriteString(fmt.Sprintf("You worked through all %d problems and solved %d.\n\n", len(m.bench), m.benchSolved))
			b.WriteString(dimStyle.Render("[m] back to the menu   ·   [q] quit"))
			break
		}
		p := m.curProblem
		hdr := fmt.Sprintf("Bench · %d/%d · %s · %s", m.benchIdx+1, len(m.bench), p.Difficulty, p.Pattern)
		if len(p.Lists) > 0 {
			hdr += " · " + strings.ToUpper(p.Lists[0])
		}
		b.WriteString(titleStyle.Render(hdr) + "\n")
		sub := fmt.Sprintf("solved %d", m.benchSolved)
		if m.benchFilter != "" {
			sub += "  ·  " + m.benchFilter
		}
		sub += "  ·  " + m.walletLine()
		b.WriteString(dimStyle.Render(sub) + "\n\n")
		b.WriteString(m.renderTask())
		b.WriteString(m.renderHintPanel())
	case screenWriteup:
		b.WriteString(m.renderWriteup())
	case screenGate:
		b.WriteString(m.renderGate())
	case screenMentor:
		b.WriteString(m.renderMentor())
	case screenPrimer:
		pr := m.primer
		pages := m.primerPages()
		idx := m.primerPage
		if idx < 0 {
			idx = 0
		}
		if idx >= len(pages) {
			idx = len(pages) - 1
		}
		pg := pages[idx]
		b.WriteString(titleStyle.Render(pr.Title) + "  " + dimStyle.Render("["+langLabel(pr.Lang)+"]") + "\n")
		b.WriteString(dimStyle.Render(pr.Category) + "\n\n")
		b.WriteString(titleStyle.Render(pg.heading) + "\n\n")
		b.WriteString(pg.body + "\n\n")
		b.WriteString(dimStyle.Render(fmt.Sprintf("page %d/%d  ·  [→/space] next   [←] prev   [enter] back", idx+1, len(pages))))
	case screenAdvancedList:
		b.WriteString(titleStyle.Render("Advanced Topics — "+langLabel(m.lang)) + "\n")
		b.WriteString(dimStyle.Render("Language-specific topics: explainers + fix-it / spot-the-bug drills.") + "\n\n")
		for i, at := range m.advTopics {
			name := advTopicName(at) // distinct short name (groups can collide, e.g. Concurrency)
			var tag string
			switch at.Tag {
			case "E":
				tag = dimStyle.Render(fmt.Sprintf("  (%d exercises)", len(at.Exercises)))
			case "gotcha":
				tag = dimStyle.Render("  (spot the bug)")
			case "C":
				tag = dimStyle.Render("  (reading)")
			case "P":
				tag = dimStyle.Render("  (reference · project — Internship later)")
			}
			sel, shown := "  ", name
			if i == m.advListIdx {
				sel, shown = "› ", okStyle.Render(name)
			}
			b.WriteString(sel + shown + tag + "\n")
		}
		b.WriteString("\n" + dimStyle.Render("[↑/↓] move   ·   [enter] open   ·   [m] back to bench   ·   [q] quit"))
	case screenAdvancedTopic:
		at := m.advTopic
		pages := m.advancedPages()
		idx := m.advPage
		if idx < 0 {
			idx = 0
		}
		if idx >= len(pages) {
			idx = len(pages) - 1
		}
		pg := pages[idx]
		b.WriteString(titleStyle.Render(at.Title) + "  " + dimStyle.Render("["+langLabel(at.Lang)+" · "+at.Group+"]") + "\n\n")
		if pg.ex == nil {
			b.WriteString(titleStyle.Render(pg.heading) + "\n\n")
			b.WriteString(pg.body + "\n\n")
		} else {
			ex := pg.ex
			label := "fix-it"
			if ex.Kind != "" {
				label = ex.Kind
			}
			b.WriteString(titleStyle.Render(fmt.Sprintf("Exercise %d — %s", pg.exNum, label)) + "\n\n")
			b.WriteString(strings.TrimRight(ex.Prompt, "\n") + "\n\n")
			gradeable := solveCheck(*ex) != ""
			// Show the player's current attempt (if any), else the broken starter.
			codeLabel, code := "broken:", ex.BrokenCode
			if m.advAttempt != "" {
				codeLabel, code = "your code:", m.advAttempt
			}
			b.WriteString(dimStyle.Render("  "+codeLabel) + "\n")
			b.WriteString(highlightCode(indentLines(code, "      "), at.Lang) + "\n\n")
			if m.advStatus != "" {
				line := "  " + m.advStatus
				if m.advVerdict != nil && m.advVerdict.Passed {
					line = okStyle.Render(line)
				} else if m.advVerdict != nil {
					line = errStyle.Render(line)
				}
				b.WriteString(line + "\n\n")
			}
			if m.advReveal {
				b.WriteString(okStyle.Render("  ▸ the bug:") + "\n")
				b.WriteString("    " + strings.ReplaceAll(strings.TrimRight(ex.Bug, "\n"), "\n", "\n    ") + "\n\n")
				b.WriteString(okStyle.Render("  ▸ the model fix:") + "\n")
				b.WriteString(highlightCode(indentLines(ex.FixedCode, "      "), at.Lang) + "\n\n")
			}
			b.WriteString(dimStyle.Render("  grading: "+advGradeLabel(ex, at.Lang)) + "\n")
			if gradeable {
				b.WriteString(dimStyle.Render("  [e] edit your fix   ·   [r] grade   ·   [b] reveal bug + fix") + "\n")
			} else {
				b.WriteString(dimStyle.Render("  [b] reveal the bug + fix (reveal-only)") + "\n")
			}
		}
		nav := "[→/space] next   [←] prev   ·   [esc] topics"
		b.WriteString("\n" + dimStyle.Render(fmt.Sprintf("page %d/%d  ·  %s", idx+1, len(pages), nav)))
	case screenStep0Complete:
		bk, ca, hd := m.benchStats()
		ps, lp, track := m.step0Profile()
		b.WriteString(titleStyle.Render("The Apprenticeship — milestone reached") + "\n\n")
		b.WriteString(fmt.Sprintf("You've banked %d problems across %d categories (%d hard).\n", bk, ca, hd))
		b.WriteString("That clears the Apprenticeship practice milestone. Your competency so far:\n\n")
		b.WriteString(fmt.Sprintf("  Problem-Solving        %d / 100\n", ps))
		b.WriteString(fmt.Sprintf("  Language Proficiency   %d / 100\n", lp))
		b.WriteString(dimStyle.Render("  Code Quality           — (needs write-up review)\n"))
		b.WriteString(dimStyle.Render("  Speed / Fluency        — (needs timing)\n"))
		b.WriteString(fmt.Sprintf("  Suggested track        %s\n\n", track))
		b.WriteString(dimStyle.Render("Next is The Job — a real role simulator (not built yet); it will\nconsume this profile. For now you can keep practicing the bench.\n\n"))
		b.WriteString(dimStyle.Render("[enter] keep practicing   ·   [q] quit"))
	case screenInstallHelp:
		b.WriteString(m.renderInstallHelp())
	}
	box := boxStyle
	if m.width > 0 { // constrain to the terminal so long lines wrap instead of clipping
		w := m.width - 6
		if w > 96 {
			w = 96
		}
		if w < 40 {
			w = 40
		}
		box = box.Width(w)
	}
	return box.Render(b.String())
}

// renderDiag dispatches by the current item's kind.
func (m Model) renderDiag() string {
	switch m.curDiag.Kind {
	case "choice":
		return m.renderChoice()
	case "spec":
		return m.renderSpec()
	default:
		return m.renderTask()
	}
}

// intakeTotal is the number of counted intake questions (10 in practice; falls
// back to the constant if m.diag wasn't rebuilt, e.g. an old save).
func (m Model) intakeTotal() int {
	if len(m.diag) > 0 {
		return len(m.diag)
	}
	return intakeSize
}

// renderResults: friendly one-liner with the score + the suggested next step.
func (m Model) renderResults() string {
	total, passed := m.intakeTotal(), m.intakePassed
	var line, accept string
	switch {
	case passed == 0:
		line = fmt.Sprintf("Couldn't really gauge you from that (%d of %d). Here's a starting point — or take it again.", passed, total)
	case m.placement == "dev-literacy":
		line = fmt.Sprintf("Solid — %d of %d. A quick terminal & Git brush-up and you're set.", passed, total)
	default: // tutorial-full
		line = fmt.Sprintf("You got %d of %d. A run through Tutorial Island will build the foundation.", passed, total)
	}
	if m.placement == "dev-literacy" {
		accept = "accept — Dev-Literacy brush-up"
	} else {
		accept = "accept — start Tutorial Island"
	}
	var b strings.Builder
	b.WriteString(line + "\n\n")
	b.WriteString(dimStyle.Render("[enter] " + accept + "   ·   [r] redo intake   ·   [q] quit"))
	return b.String()
}

func (m Model) renderIntro() string {
	var b strings.Builder
	b.WriteString(m.intro.Prompt + "\n\n")
	for i, c := range m.intro.Choices {
		sel := "  "
		if i == m.chosen {
			sel = "› "
		}
		b.WriteString(fmt.Sprintf("%s%d. %s\n", sel, i+1, c.Text))
	}
	b.WriteString("\n")
	if m.answered {
		if m.lastFeedback != "" {
			b.WriteString(m.lastFeedback + "\n\n")
		}
		b.WriteString(dimStyle.Render("[enter] begin the Entrance Test   ·   [q] quit"))
	} else {
		b.WriteString(dimStyle.Render(m.status + "   ·   [q] quit"))
	}
	return b.String()
}

func (m Model) renderChoice() string {
	var b strings.Builder
	b.WriteString(m.curDiag.Prompt + "\n\n")
	for i, c := range m.curDiag.Choices {
		sel := "  "
		if i == m.chosen {
			sel = "› " // current selection (highlight, before confirm)
		}
		line := fmt.Sprintf("%s%d. %s", sel, i+1, c.Text)
		if m.answered && i == m.chosen {
			if m.itemPassed {
				line = okStyle.Render(line)
			} else {
				line = errStyle.Render(line)
			}
		}
		b.WriteString(line + "\n")
	}
	b.WriteString("\n")
	if m.answered {
		if m.lastFeedback != "" {
			b.WriteString(m.lastFeedback + "\n\n")
		}
		b.WriteString(dimStyle.Render("[enter] continue   ·   [q] quit"))
	} else {
		b.WriteString(dimStyle.Render(m.status + "   ·   [q] quit"))
	}
	return b.String()
}

func (m Model) renderSpec() string {
	var b strings.Builder
	b.WriteString(m.curDiag.Prompt + "\n\n")
	if m.answered {
		mark := okStyle.Render("✓ Looks good.")
		if !m.itemPassed {
			mark = dimStyle.Render("That's okay — here's how we'd read it:")
		}
		b.WriteString("Your answer: " + dimStyle.Render(strings.TrimSpace(m.input)) + "\n\n")
		b.WriteString(mark + "\n" + m.lastFeedback + "\n\n")
		b.WriteString(dimStyle.Render("[enter] continue   ·   [q] quit"))
	} else {
		b.WriteString("> " + m.input + "▏\n\n")
		b.WriteString(dimStyle.Render("type your answer · [enter] submit · [esc] quit"))
	}
	return b.String()
}

func (m Model) renderDev() string {
	var b strings.Builder
	b.WriteString(m.curDev.Prompt + "\n\n")
	if m.answered {
		b.WriteString(okStyle.Render("$ "+strings.TrimSpace(m.input)) + "\n\n")
		b.WriteString(okStyle.Render("✓ ") + m.curDev.Success + "\n\n")
		b.WriteString(dimStyle.Render("[enter] continue   ·   [q] quit"))
	} else {
		b.WriteString("$ " + m.input + "▏\n\n")
		if m.status != "" {
			b.WriteString(m.status + "\n\n")
		}
		b.WriteString(dimStyle.Render("type a command · [enter] submit · [esc] quit"))
	}
	return b.String()
}

func (m Model) renderTask() string {
	t := m.task
	var b strings.Builder
	b.WriteString(t.prompt + "\n\n")
	if t.verdict != nil {
		for _, r := range t.verdict.Results {
			mark := errStyle.Render("✗")
			if r.Passed {
				mark = okStyle.Render("✓")
			}
			line := fmt.Sprintf("  %s %s", mark, r.Name)
			if !r.Passed {
				if r.Err != "" {
					line += dimStyle.Render("  (error: " + oneline(r.Err) + ")")
				} else {
					line += dimStyle.Render(fmt.Sprintf("  (got %s, want %s)", r.Got, r.Expected))
				}
			}
			b.WriteString(line + "\n")
		}
		b.WriteString("\n")
	}
	if m.status != "" {
		st := m.status
		if t.verdict != nil && t.verdict.Passed {
			st = okStyle.Render(st)
		}
		b.WriteString(st + "\n\n")
	}
	hints := "[e] edit code   ·   [r] run   ·   [c] change editor"
	if m.ctx == ctxDiagnostic || m.ctx == ctxBench {
		hints += "   ·   [s] skip"
	}
	if m.ctx == ctxBench {
		hints += "   ·   [l] learn   ·   [m] menu"
	}
	if m.ctx == ctxLesson {
		hints += "   ·   [t] test me   ·   [x] skip tutorial"
	}
	hints += "   ·   [q] quit"
	b.WriteString(dimStyle.Render(hints))
	return b.String()
}

func oneline(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) > 80 {
		s = s[:80] + "…"
	}
	return s
}
