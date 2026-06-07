package tui

import (
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"devascent/internal/content"
	"devascent/internal/grader"
	"devascent/internal/save"
	"devascent/internal/toolchain"
)

// Correct reference solutions for every CODE diagnostic variant, used by the
// real-grader white-box test. (content.TestAllAuthoredTasksGradePass is the
// canonical grader-correctness guard; this map just drives the TUI flow.)
var tuiSolutions = map[string]string{
	"add":             "def add(a, b):\n    return a + b\n",
	"subtract":        "def subtract(a, b):\n    return a - b\n",
	"multiply":        "def multiply(a, b):\n    return a * b\n",
	"max_two":         "def max_two(a, b):\n    return a if a > b else b\n",
	"min_two":         "def min_two(a, b):\n    return a if a < b else b\n",
	"square":          "def square(n):\n    return n * n\n",
	"triple":          "def triple(n):\n    return n * 3\n",
	"last_element":    "def last_element(items):\n    return items[-1]\n",
	"first_char":      "def first_char(s):\n    return s[0]\n",
	"negate":          "def negate(n):\n    return -n\n",
	"reverse_text":    "def reverse_text(s):\n    return s[::-1]\n",
	"sum_list":        "def sum_list(numbers):\n    return sum(numbers)\n",
	"count_evens":     "def count_evens(nums):\n    return sum(1 for x in nums if x % 2 == 0)\n",
	"count_vowels":    "def count_vowels(s):\n    return sum(1 for c in s.lower() if c in 'aeiou')\n",
	"count_positives": "def count_positives(nums):\n    return sum(1 for n in nums if n > 0)\n",
	"second_largest":  "def second_largest(nums):\n    return sorted(set(nums))[-2]\n",
	"is_palindrome":   "def is_palindrome(s):\n    return s == s[::-1]\n",
	"factorial":       "def factorial(n):\n    r = 1\n    for i in range(2, n + 1):\n        r *= i\n    return r\n",
	"sum_evens":       "def sum_evens(nums):\n    return sum(n for n in nums if n % 2 == 0)\n",
	"count_words":     "def count_words(s):\n    return len(s.split())\n",
}

// ── driver: drives Update directly with a PINNED seed for deterministic runs ──

type driver struct {
	t *testing.T
	m Model
}

func newDriver(t *testing.T) *driver {
	t.Helper()
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	m := New()
	if m.loadErr != nil {
		t.Fatal(m.loadErr)
	}
	m.rng = rand.New(rand.NewSource(7))
	return &driver{t: t, m: m}
}

func (d *driver) step(msg tea.Msg) {
	nm, _ := d.m.Update(msg)
	d.m = nm.(Model)
}
func (d *driver) runes(s string) { d.step(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}) }
func (d *driver) enter()         { d.step(tea.KeyMsg{Type: tea.KeyEnter}) }

func (d *driver) chooseCorrect() {
	idx := -1
	for i, c := range d.m.curDiag.Choices {
		if c.Correct {
			idx = i
			break
		}
	}
	if idx < 0 {
		d.t.Fatal("no correct choice")
	}
	d.runes(strconv.Itoa(idx + 1))
	d.enter()
	d.enter()
}
func (d *driver) chooseWrong() {
	idx := -1
	for i, c := range d.m.curDiag.Choices {
		if !c.Correct {
			idx = i
			break
		}
	}
	if idx < 0 {
		d.t.Fatal("no wrong choice")
	}
	d.runes(strconv.Itoa(idx + 1))
	d.enter()
	d.enter()
}

func (d *driver) passCode() {
	d.step(editorFinishedMsg{code: "x"})
	d.step(gradeMsg{v: grader.Verdict{Passed: true}})
	d.enter()
}
func (d *driver) skipCode() { d.runes("s") }

func (d *driver) solveCodeReal() {
	fn := d.m.task.funcName
	sol, ok := tuiSolutions[fn]
	if !ok {
		d.t.Fatalf("no reference solution for %q", fn)
	}
	d.step(editorFinishedMsg{code: sol})
	v, err := grader.NewNativePython().Run("python", sol, fn, d.m.task.tests, d.m.task.shape)
	if err != nil {
		d.t.Fatalf("grader error for %q: %v", fn, err)
	}
	if !v.Passed {
		d.t.Fatalf("%q reference solution failed: %+v", fn, v.Results)
	}
	d.step(gradeMsg{v: v})
	d.enter()
}

func (d *driver) specAnswerCorrect() {
	var parts []string
	for _, g := range d.m.curDiag.Spec.Required {
		parts = append(parts, strings.TrimSpace(strings.Split(g, "|")[0]))
	}
	d.step(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(strings.Join(parts, " "))})
	d.enter()
	d.enter()
}
func (d *driver) specAnswerWrong() {
	d.step(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("nope nothing relevant here")})
	d.enter()
	d.enter()
}

func (d *driver) devSolveCurrent() {
	t := d.m.curDev
	var ans string
	if len(t.Commands) > 0 {
		ans = t.Commands[0]
		for _, f := range t.Flags {
			ans += " " + f
		}
	} else if len(t.Accept) > 0 {
		ans = t.Accept[0]
	}
	d.step(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(ans)})
	d.enter()
	d.enter()
}

// toIntake drives hook→language→editor→warm-up, picks the self-report option for
// `level`, and lands on the first counted intake item.
func (d *driver) toIntake(level string) {
	d.enter() // hook → language
	d.enter() // language → editor
	d.enter() // editor → intro warm-up
	if d.m.screen != screenIntro {
		d.t.Fatalf("want intro, got screen %d", d.m.screen)
	}
	idx := -1
	for i, c := range d.m.intro.Choices {
		if c.Value == level {
			idx = i
		}
	}
	if idx < 0 {
		d.t.Fatalf("no intro option for level %q", level)
	}
	d.runes(strconv.Itoa(idx + 1))
	d.enter() // confirm
	d.enter() // begin intake
	if d.m.screen != screenDiagnostic {
		d.t.Fatalf("want diagnostic, got screen %d", d.m.screen)
	}
}

// runIntake answers each counted item; ok(i, measures, kind) decides pass/fail.
func (d *driver) runIntake(ok func(i int, measures, kind string) bool) {
	for i := 0; d.m.screen == screenDiagnostic; i++ {
		if i > 50 {
			d.t.Fatal("intake did not terminate")
		}
		pass := ok(i, d.m.curDiag.Measures, d.m.curDiag.Kind)
		switch d.m.curDiag.Kind {
		case "code":
			if pass {
				d.passCode()
			} else {
				d.skipCode()
			}
		case "choice":
			if pass {
				d.chooseCorrect()
			} else {
				d.chooseWrong()
			}
		case "spec":
			if pass {
				d.specAnswerCorrect()
			} else {
				d.specAnswerWrong()
			}
		}
	}
}

func allCorrect(i int, m, k string) bool { return true }

func requirePython(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("python"); err != nil {
		t.Skip("python not found on PATH")
	}
}

// ── difficulty band (the user's explicit request: confirm we follow it) ───────

func TestSelectIntake_BandByLevel(t *testing.T) {
	c, err := content.Load()
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		level string
		ok    map[int]bool
	}{
		{"never", map[int]bool{1: true}},
		{"a-little", map[int]bool{1: true, 2: true}},
		{"regularly", map[int]bool{2: true, 3: true}},
	}
	for _, tc := range cases {
		got := selectIntake(c.Diagnostics, tc.level, nil, rand.New(rand.NewSource(7)))
		if len(got) != 10 {
			t.Errorf("%s: want 10 items, got %d", tc.level, len(got))
		}
		for _, d := range got {
			if !tc.ok[d.Difficulty] {
				t.Errorf("%s: item %q difficulty %d not in band %v", tc.level, d.ID, d.Difficulty, tc.ok)
			}
		}
		if len(got) > 0 && got[0].Measures != "coding" {
			t.Errorf("%s: first item should be coding (the floor), got %q", tc.level, got[0].Measures)
		}
	}
}

// passN passes the first n items, fails the rest (score-based routing only
// cares about the total count).
func passN(n int) func(i int, m, k string) bool {
	return func(i int, m, k string) bool { return i < n }
}

// ── routing branches (percentage bands + self-report clamp) ───────────────────

// Score bands (a-little, no clamp): ≥80% → bench(aced), 40–79% → dev-lit, <40% → tutorial.
func TestRoute_ScoreBands(t *testing.T) {
	cases := []struct {
		pass   int
		place  string
		screen screen
	}{
		{10, "test-out", screenTestOut},
		{8, "test-out", screenTestOut},
		{7, "dev-literacy", screenResults},
		{5, "dev-literacy", screenResults},
		{4, "dev-literacy", screenResults},
		{3, "tutorial-full", screenResults},
		{2, "tutorial-full", screenResults},
	}
	for _, tc := range cases {
		d := newDriver(t)
		d.toIntake("a-little")
		d.runIntake(passN(tc.pass))
		if d.m.placement != tc.place || d.m.screen != tc.screen {
			t.Errorf("pass %d/10: want %q/%d, got %q/%d", tc.pass, tc.place, tc.screen, d.m.placement, d.m.screen)
		}
		if d.m.intakePassed != tc.pass {
			t.Errorf("pass %d/10: intakePassed=%d", tc.pass, d.m.intakePassed)
		}
	}
}

// never: even acing every (easy) item never reaches the bench — clamped to
// tutorial, shown on the results screen (accept → tutorial).
func TestRoute_Never_ClampedToTutorial(t *testing.T) {
	d := newDriver(t)
	d.toIntake("never")
	d.runIntake(allCorrect)
	if d.m.screen != screenResults || d.m.placement != "tutorial-full" {
		t.Fatalf("never+ace should suggest tutorial via results, got screen %d placement %q", d.m.screen, d.m.placement)
	}
	d.enter() // accept
	if d.m.screen != screenLesson {
		t.Fatalf("accept should start tutorial, got %d", d.m.screen)
	}
}

func TestRoute_Regularly_Ace_TestOut(t *testing.T) {
	d := newDriver(t)
	d.toIntake("regularly")
	d.runIntake(allCorrect)
	if d.m.screen != screenTestOut || d.m.placement != "test-out" {
		t.Fatalf("regularly+ace should test-out (aced screen), got screen %d placement %q", d.m.screen, d.m.placement)
	}
}

func TestRoute_ALittle_Ace_TestOut(t *testing.T) {
	d := newDriver(t)
	d.toIntake("a-little")
	d.runIntake(allCorrect)
	if d.m.screen != screenTestOut {
		t.Fatalf("a-little+ace should show the aced screen, got %d", d.m.screen)
	}
}

// regularly + score <40% → base tutorial, clamped UP to dev-literacy (never
// beginner lesson 1 for a self-described regular coder).
func TestRoute_Regularly_LowScore_ClampsToDevLiteracy(t *testing.T) {
	d := newDriver(t)
	d.toIntake("regularly")
	d.runIntake(passN(0)) // fail everything
	if d.m.placement != "dev-literacy" {
		t.Fatalf("regular coder should clamp up to dev-literacy, got %q", d.m.placement)
	}
	if d.m.screen != screenResults || d.m.intakePassed != 0 {
		t.Fatalf("want results screen with 0 passed, got screen %d passed %d", d.m.screen, d.m.intakePassed)
	}
}

// Accepting a dev-literacy suggestion enters the track and completes to handoff.
func TestResults_AcceptDevLiteracy(t *testing.T) {
	d := newDriver(t)
	d.toIntake("a-little")
	d.runIntake(passN(5)) // 50% → dev-literacy
	if d.m.screen != screenResults || d.m.placement != "dev-literacy" {
		t.Fatalf("50%% should suggest dev-literacy, got screen %d placement %q", d.m.screen, d.m.placement)
	}
	d.enter() // accept
	if d.m.screen != screenDevLiteracy {
		t.Fatalf("accept should enter dev-literacy, got %d", d.m.screen)
	}
	for len(d.m.dev) > 0 && d.m.screen == screenDevLiteracy {
		d.devSolveCurrent()
	}
	if d.m.screen != screenHandoff {
		t.Fatalf("after dev-literacy want handoff, got %d", d.m.screen)
	}
}

// Redo re-runs the intake at the same level with (best-effort) fresh questions.
func TestResults_RedoSameLevel(t *testing.T) {
	d := newDriver(t)
	d.toIntake("a-little")
	first := append([]content.Diagnostic(nil), d.m.diag...)
	d.runIntake(passN(2)) // → tutorial suggestion (results screen)
	if d.m.screen != screenResults {
		t.Fatalf("want results screen, got %d", d.m.screen)
	}
	d.runes("r") // redo
	if d.m.screen != screenDiagnostic || d.m.level != "a-little" {
		t.Fatalf("redo should restart the intake at the same level, got screen %d level %q", d.m.screen, d.m.level)
	}
	if d.m.diagIdx != 0 || d.m.intakePassed != 0 {
		t.Fatalf("redo should reset progress; idx %d passed %d", d.m.diagIdx, d.m.intakePassed)
	}
	// fresh-first: the new first item should differ from the old first item
	if len(first) > 0 && d.m.diag[0].ID == first[0].ID {
		t.Logf("redo first item repeated (%q) — acceptable if band is small", first[0].ID)
	}
}

// Aced → [enter] goes straight to the bench (one hop, not via the end screen).
func TestAced_EnterGoesToBench(t *testing.T) {
	d := newDriver(t)
	if len(d.m.cat.Problems) == 0 {
		t.Skip("no bench problems")
	}
	d.toIntake("regularly")
	d.runIntake(allCorrect)
	if d.m.screen != screenTestOut {
		t.Fatalf("want aced screen, got %d", d.m.screen)
	}
	d.enter() // accept → bench browse menu
	if d.m.screen != screenBenchMenu {
		t.Fatalf("aced [enter] should open the bench menu, got %d", d.m.screen)
	}
}

// ── input / dev-literacy behavior ─────────────────────────────────────────────

// driveToDevLiteracy: a-little, score 50%, accept the dev-literacy suggestion.
func (d *driver) driveToDevLiteracy() {
	d.toIntake("a-little")
	d.runIntake(passN(5))
	if d.m.screen == screenResults && d.m.placement == "dev-literacy" {
		d.enter() // accept
	}
}

func TestDevLiteracy_WrongCommand_Hints(t *testing.T) {
	d := newDriver(t)
	d.driveToDevLiteracy()
	if d.m.screen != screenDevLiteracy {
		t.Fatalf("want dev-literacy, got %d", d.m.screen)
	}
	d.step(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("zzznotacmd")})
	d.enter()
	d.step(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("stillwrong")})
	d.enter()
	if d.m.devIdx != 0 || d.m.answered {
		t.Fatalf("wrong command should not advance; devIdx=%d answered=%v", d.m.devIdx, d.m.answered)
	}
	if !strings.Contains(d.m.status, "Hint") {
		t.Fatalf("want a hint after 2 tries, got %q", d.m.status)
	}
}

// Modal input: typing command keys (q/e/r/s/x) must not quit; keywords then pass.
func TestSpec_TextInput_DoesNotTriggerCommands(t *testing.T) {
	d := newDriver(t)
	d.toIntake("a-little")
	// advance (answering correctly) until we hit a spec item
	for d.m.screen == screenDiagnostic && d.m.curDiag.Kind != "spec" {
		switch d.m.curDiag.Kind {
		case "code":
			d.passCode()
		case "choice":
			d.chooseCorrect()
		}
	}
	if d.m.screen != screenDiagnostic || d.m.curDiag.Kind != "spec" {
		t.Skip("no spec item in this intake selection")
	}
	if !d.m.inputActive {
		t.Fatal("spec item should be in input mode")
	}
	d.step(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("requires a sex-quirk answer ")})
	if d.m.quitting {
		t.Fatal("typing must never quit in input mode")
	}
	if !strings.Contains(d.m.input, "requires") {
		t.Fatalf("buffer should hold typed text, got %q", d.m.input)
	}
	var kw []string
	for _, g := range d.m.curDiag.Spec.Required {
		kw = append(kw, strings.TrimSpace(strings.Split(g, "|")[0]))
	}
	d.step(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" " + strings.Join(kw, " "))})
	d.enter()
	if !d.m.answered || !d.m.itemPassed {
		t.Fatalf("spec answer with all keywords should pass; answered=%v passed=%v", d.m.answered, d.m.itemPassed)
	}
}

func TestSpecMatch_Synonyms(t *testing.T) {
	d := newDriver(t)
	var sp *content.Spec
	for _, di := range d.m.cat.Diagnostics {
		if di.ID == "spec_largest" {
			sp = di.Spec
		}
	}
	if sp == nil {
		t.Fatal("spec_largest not found")
	}
	for _, s := range []string{
		"input is a list of numbers, output is the largest",
		"an array of integers in, the maximum out",
		"takes a collection, returns the biggest value",
	} {
		if !specMatch(s, sp) {
			t.Errorf("expected PASS: %q", s)
		}
	}
	for _, s := range []string{"i don't know", "it returns something", "a number"} {
		if specMatch(s, sp) {
			t.Errorf("expected FAIL: %q", s)
		}
	}
}

// ── tutorial / resume ─────────────────────────────────────────────────────────

func (d *driver) driveToHandoff() {
	d.toIntake("never")
	d.runIntake(allCorrect) // never → results screen suggesting tutorial
	if d.m.screen != screenResults {
		d.t.Fatalf("want results screen, got %d", d.m.screen)
	}
	d.enter() // accept → tutorial
	if d.m.screen != screenLesson {
		d.t.Fatalf("want tutorial, got %d", d.m.screen)
	}
	d.runes("x") // skip tutorial → handoff
	if d.m.screen != screenHandoff {
		d.t.Fatalf("want handoff, got %d", d.m.screen)
	}
}

func TestExitTutorial_GoesToHandoff(t *testing.T) {
	d := newDriver(t)
	d.driveToHandoff()
}

func TestResume_ContinuesTutorial(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	if err := save.Save(save.State{Language: "python", Stage: "tutorial", LessonIdx: 3, StageIdx: 1}); err != nil {
		t.Fatal(err)
	}
	m := New()
	if m.resume == nil {
		t.Fatal("expected resume")
	}
	nm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("c")})
	m = nm.(Model)
	if m.screen != screenLesson || m.lessonIdx != 3 || m.stageIdx != 1 {
		t.Fatalf("want lesson 3 stage 1, got screen %d %d/%d", m.screen, m.lessonIdx, m.stageIdx)
	}
}

func TestResumeIntake_RebuildsFromIDs(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	c, _ := content.Load()
	intake := selectIntake(c.Diagnostics, "a-little", nil, rand.New(rand.NewSource(11)))
	ids := diagIDs(intake)
	if err := save.Save(save.State{Language: "python", Stage: "intake", Level: "a-little", DiagIdx: 2, DiagIDs: ids, PassedAdd: true}); err != nil {
		t.Fatal(err)
	}
	m := New()
	nm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("c")})
	m = nm.(Model)
	if m.screen != screenDiagnostic {
		t.Fatalf("want diagnostic, got %d", m.screen)
	}
	if strings.Join(diagIDs(m.diag), ",") != strings.Join(ids, ",") {
		t.Fatalf("intake not rebuilt from IDs")
	}
	if m.diagIdx != 2 {
		t.Fatalf("want diagIdx 2, got %d", m.diagIdx)
	}
}

func TestResume_ContinuesDevLiteracy(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	c, _ := content.Load()
	if len(c.DevTasks) < 3 {
		t.Skip("not enough dev tasks")
	}
	set := selectDevSet(c.DevTasks, 5, rand.New(rand.NewSource(2)))
	ids := devTaskIDs(set)
	if err := save.Save(save.State{Language: "python", Stage: "devliteracy", DevIdx: 1, DevIDs: ids, Placement: "dev-literacy"}); err != nil {
		t.Fatal(err)
	}
	m := New()
	nm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("c")})
	m = nm.(Model)
	if m.screen != screenDevLiteracy || m.devIdx != 1 {
		t.Fatalf("want dev-literacy at idx 1, got screen %d idx %d", m.screen, m.devIdx)
	}
}

func TestResume_DoneOffersEndScreen(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	if err := save.Save(save.State{Language: "python", Stage: "done", Placement: "test-out"}); err != nil {
		t.Fatal(err)
	}
	m := New()
	if m.resume == nil {
		t.Fatal("a finished run should still be offered as resume")
	}
	nm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("c")})
	m = nm.(Model)
	if m.screen != screenHandoff {
		t.Fatalf("continue on done should land on end screen, got %d", m.screen)
	}
}

// ── handoff is not a dead-end ──────────────────────────────────────────────────

func TestHandoff_ReplayTutorial(t *testing.T) {
	d := newDriver(t)
	d.driveToHandoff()
	d.runes("t")
	if d.m.screen != screenLesson || d.m.lessonIdx != 0 {
		t.Fatalf("[t] should replay tutorial from lesson 0, got screen %d idx %d", d.m.screen, d.m.lessonIdx)
	}
}

func TestHandoff_StartOver(t *testing.T) {
	d := newDriver(t)
	d.driveToHandoff()
	d.runes("r")
	if d.m.screen != screenLanguage {
		t.Fatalf("[r] should restart at language, got %d", d.m.screen)
	}
	if d.m.placement != "" || len(d.m.diag) != 0 {
		t.Fatalf("restart should clear run state")
	}
}

// ── Step 0 bench ───────────────────────────────────────────────────────────

// enterBenchAll: handoff → [b] → browse menu → pick "All" (first option).
func (d *driver) enterBenchAll() {
	d.driveToHandoff()
	d.runes("b")
	if d.m.screen != screenBenchMenu {
		d.t.Fatalf("want bench menu, got screen %d", d.m.screen)
	}
	d.enter() // first option = All
	if d.m.screen != screenBench || d.m.task == nil {
		d.t.Fatalf("want bench with a problem, got screen %d", d.m.screen)
	}
}

func TestBench_MenuThenSolveAndSkip(t *testing.T) {
	d := newDriver(t)
	if len(d.m.cat.Problems) == 0 {
		t.Skip("no bench problems")
	}
	d.enterBenchAll()
	d.passCode()
	if d.m.benchIdx != 1 || d.m.benchSolved != 1 {
		t.Fatalf("after solve want idx1 solved1, got idx%d solved%d", d.m.benchIdx, d.m.benchSolved)
	}
	d.runes("s")
	if d.m.benchIdx != 2 || d.m.benchSolved != 1 {
		t.Fatalf("after skip want idx2 solved1, got idx%d solved%d", d.m.benchIdx, d.m.benchSolved)
	}
}

func TestSelectBench_TierOrderByPlacement(t *testing.T) {
	c, err := content.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Problems) == 0 {
		t.Skip("no problems")
	}
	easyStart := selectBench(c.Problems, "tutorial-full", rand.New(rand.NewSource(1)))
	if easyStart[0].Difficulty != "easy" {
		t.Fatalf("tutorial-full should start easy, got %q", easyStart[0].Difficulty)
	}
	medStart := selectBench(c.Problems, "test-out", rand.New(rand.NewSource(1)))
	if medStart[0].Difficulty != "medium" {
		t.Fatalf("test-out should start medium, got %q", medStart[0].Difficulty)
	}
	if len(easyStart) != len(c.Problems) {
		t.Fatalf("bench should include all problems")
	}
}

func TestResume_ContinuesBench(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	c, _ := content.Load()
	if len(c.Problems) < 3 {
		t.Skip("not enough problems")
	}
	pool := selectBench(c.Problems, "test-out", rand.New(rand.NewSource(4)))
	ids := benchIDs(pool)
	if err := save.Save(save.State{Language: "python", Stage: "bench", BenchIDs: ids, BenchIdx: 2, BenchSolved: 1, Placement: "test-out"}); err != nil {
		t.Fatal(err)
	}
	m := New()
	nm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("c")})
	m = nm.(Model)
	if m.screen != screenBench || m.benchIdx != 2 || m.benchSolved != 1 {
		t.Fatalf("want bench idx2 solved1, got screen %d idx%d solved%d", m.screen, m.benchIdx, m.benchSolved)
	}
}

// Browse menu: picking a category serves only that category's problems.
func TestBenchMenu_CategoryFilter(t *testing.T) {
	d := newDriver(t)
	if len(d.m.cat.Problems) == 0 {
		t.Skip("no problems")
	}
	d.driveToHandoff()
	d.runes("b")
	if d.m.screen != screenBenchMenu {
		t.Fatalf("want bench menu, got %d", d.m.screen)
	}
	// find the "Dynamic Programming" category option and select it
	target := -1
	for i, o := range d.m.benchMenu {
		if o.kind == "category" && o.value == "Dynamic Programming" {
			target = i
		}
	}
	if target < 0 {
		t.Fatal("no Dynamic Programming category in menu")
	}
	for d.m.benchMenuIdx < target {
		d.runes("j")
	}
	d.enter()
	if d.m.screen != screenBench {
		t.Fatalf("want bench, got %d", d.m.screen)
	}
	for _, p := range d.m.bench {
		if p.Category != "Dynamic Programming" {
			t.Fatalf("category filter leaked: %q", p.Category)
		}
	}
	if !strings.Contains(d.m.benchFilter, "Dynamic Programming") {
		t.Fatalf("bench filter label wrong: %q", d.m.benchFilter)
	}
}

// Gap-fill content is present: Intervals category + Blind 75 tags exist.
func TestBench_GapFillPresent(t *testing.T) {
	c, err := content.Load()
	if err != nil {
		t.Fatal(err)
	}
	cats := map[string]int{}
	blind := 0
	for _, p := range c.Problems {
		cats[p.Category]++
		for _, l := range p.Lists {
			if l == "blind75" {
				blind++
			}
		}
	}
	if cats["Intervals"] < 4 {
		t.Fatalf("expected Intervals gap-fill, got %d", cats["Intervals"])
	}
	if blind < 8 {
		t.Fatalf("expected several blind75-tagged problems, got %d", blind)
	}
	if len(c.Problems) < 128 {
		t.Fatalf("expected 128 problems (117 mined + gap-fill), got %d", len(c.Problems))
	}
}

// The browser exposes a Hard filter so the milestone's hard floor is reachable.
func TestBenchMenu_HasDifficultyFilter(t *testing.T) {
	c, _ := content.Load()
	opts := benchMenuOptions(c.Problems)
	hard := false
	for _, o := range opts {
		if o.kind == "difficulty" && o.value == "hard" {
			hard = true
		}
	}
	if !hard {
		t.Fatal("bench menu should offer a Hard-only filter")
	}
	got := filterProblems(c.Problems, "difficulty", "hard")
	if len(got) < step0HardTarget {
		t.Fatalf("need >= %d hard problems reachable, got %d", step0HardTarget, len(got))
	}
	for _, p := range got {
		if p.Difficulty != "hard" {
			t.Fatalf("difficulty filter leaked %q", p.Difficulty)
		}
	}
}

// Every curated-list tag in the catalog must be reachable via a bench menu
// filter (data tagged but not navigable is a bug we fixed; this prevents regress).
func TestEveryListTagHasMenuFilter(t *testing.T) {
	c, err := content.Load()
	if err != nil {
		t.Fatal(err)
	}
	tags := map[string]bool{}
	for _, p := range c.Problems {
		for _, l := range p.Lists {
			tags[l] = true
		}
	}
	filtered := map[string]bool{}
	for _, o := range benchMenuOptions(c.Problems) {
		if o.kind == "list" {
			filtered[o.value] = true
		}
	}
	for tag := range tags {
		if !filtered[tag] {
			t.Errorf("list tag %q has no bench menu filter", tag)
		}
	}
}

// Serving a curated list dedupes to one variant per canonical slug (count == distinct slugs).
func TestListServingDedupesBySlug(t *testing.T) {
	c, err := content.Load()
	if err != nil {
		t.Fatal(err)
	}
	tagged := filterProblems(c.Problems, "list", "blind75")
	slugs := map[string]bool{}
	for _, p := range tagged {
		slugs[p.CanonicalSlug()] = true
	}
	served := dedupeBySlug(tagged, rand.New(rand.NewSource(1)))
	if len(served) != len(slugs) {
		t.Fatalf("dedupeBySlug served %d, want %d distinct slugs", len(served), len(slugs))
	}
	if len(slugs) != 75 {
		t.Fatalf("blind75 distinct slugs = %d, want 75", len(slugs))
	}
}

// [L] on a bench problem opens its category's primer; any key returns.
func TestPrimer_LearnKeyOpensAndReturns(t *testing.T) {
	d := newDriver(t)
	if _, ok := d.m.cat.PrimerByCategory("Arrays & Hashing"); !ok {
		t.Skip("no Arrays & Hashing primer")
	}
	d.m.screen = screenBench
	d.m.ctx = ctxBench
	d.m.curProblem = content.Problem{Category: "Arrays & Hashing"}
	d.runes("l")
	if d.m.screen != screenPrimer || d.m.primer.Category != "Arrays & Hashing" {
		t.Fatalf("[l] should open the Arrays & Hashing primer, got screen %d cat %q", d.m.screen, d.m.primer.Category)
	}
	d.enter() // any key returns to the problem
	if d.m.screen != screenBench {
		t.Fatalf("dismiss should return to bench, got %d", d.m.screen)
	}
}

// A non-Python session language (Java) is selectable, renders THAT language's
// primer, and is now GRADEABLE (function-call harness) — so [r] no longer says
// reference-only; with the toolchain absent it routes to install help instead.
func TestPrimer_JavaSession(t *testing.T) {
	d := newDriver(t)
	langs := d.m.availableLangs()
	javaIdx := -1
	for i, o := range langs {
		if o.key == "java" {
			javaIdx = i
		}
	}
	if javaIdx < 0 {
		t.Fatal("Java should be offered on the language pick (it has authored primers)")
	}

	// Pick Java from the language screen → it becomes the session language.
	d.m.screen = screenLanguage
	d.runes(strconv.Itoa(javaIdx + 1))
	if d.m.lang != "java" {
		t.Fatalf("picking option %d should set lang=java, got %q", javaIdx+1, d.m.lang)
	}

	// [l] on a bench problem opens the JAVA primer for that category.
	d.m.screen = screenBench
	d.m.ctx = ctxBench
	d.m.curProblem = content.Problem{Category: "Arrays & Hashing"}
	d.runes("l")
	if d.m.screen != screenPrimer || d.m.primer.Lang != "java" {
		t.Fatalf("[l] in a Java session should open the Java primer, got screen %d lang %q", d.m.screen, d.m.primer.Lang)
	}
	d.enter() // back to the bench problem

	// Java is gradeable now. With the toolchain absent (stub), [r] routes to the
	// install guide — NOT a "reference-only" dead-end.
	d.m.det = toolchain.NewStub(map[string]toolchain.Probe{}) // java missing
	d.m.task = &codeTask{code: "class Solution { public long f(){ return 0L; } }", funcName: "f",
		tests: []grader.TestCase{{Name: "t", Input: []any{}, Expected: 0}}}
	d.runes("r")
	if d.m.screen != screenInstallHelp {
		t.Fatalf("[r] in a Java session with no toolchain should open install help, got screen %d status %q", d.m.screen, d.m.status)
	}

	// The Java primer content must render correctly — the YAML single-quote
	// escaping ('' -> ') should yield  c - 'a'  not the raw  c - ''a''. Cheap
	// guard so a subtle escaping slip can't ship mangled across future languages.
	strp, ok := d.m.cat.PrimerByCategoryAndLang("Strings", "java")
	if !ok {
		t.Fatal("expected a Java Strings primer")
	}
	joined := strp.Summary + strp.Example
	for _, op := range strp.Ops {
		joined += op.Code
	}
	if !strings.Contains(joined, "c - 'a'") || strings.Contains(joined, "''a''") {
		t.Fatalf("Java Strings primer YAML escaping looks wrong (want c - 'a', no doubled quotes)")
	}
}

// ── ADR-0007: BYO runtime — the two-axis grade gate + install help ───────────

// Axis A: a Python session whose toolchain is MISSING routes the grade action
// [r] to the install guide (not a cryptic failure). Uses a stub detector so the
// result doesn't depend on what's installed on the test machine.
func TestActionGate_MissingPythonOpensInstallHelp(t *testing.T) {
	d := newDriver(t)
	d.m.det = toolchain.NewStub(map[string]toolchain.Probe{
		"python": {Status: toolchain.Missing, Reason: "python was not found on your PATH"},
	})
	d.m.lang = "python"
	d.m.screen = screenBench
	d.m.ctx = ctxBench
	d.m.task = &codeTask{code: "def add(a, b):\n    return a + b\n"}
	d.runes("r")
	if d.m.screen != screenInstallHelp {
		t.Fatalf("missing python should route [r] to install help, got screen %d", d.m.screen)
	}
	if d.m.installLang != "python" {
		t.Fatalf("install help should target python, got %q", d.m.installLang)
	}
	if out := d.m.renderInstallHelp(); !strings.Contains(out, "Install Python") {
		t.Fatalf("install help should render the Python guide, got: %q", out)
	}
}

// Axis A: with Python available, [r] proceeds to grading (status flips to Running;
// the returned grade command isn't executed by the driver, so no real Python runs).
func TestActionGate_AvailablePythonGrades(t *testing.T) {
	d := newDriver(t)
	d.m.det = toolchain.NewStub(map[string]toolchain.Probe{"python": {Status: toolchain.Available}})
	d.m.lang = "python"
	d.m.screen = screenBench
	d.m.ctx = ctxBench
	d.m.task = &codeTask{
		code: "def add(a, b):\n    return a + b\n", funcName: "add",
		tests: []grader.TestCase{{Name: "t", Input: []any{2, 3}, Expected: 5}},
	}
	d.runes("r")
	if d.m.screen == screenInstallHelp {
		t.Fatalf("available python should NOT route to install help")
	}
	if !strings.Contains(d.m.status, "Running") {
		t.Fatalf("available python [r] should start grading, got status %q", d.m.status)
	}
}

// Re-check [R] on the install screen returns to the prior screen once the
// toolchain appears (simulating the player completing the install).
func TestInstallHelp_RecheckReturnsWhenAvailable(t *testing.T) {
	d := newDriver(t)
	stub := map[string]toolchain.Probe{"python": {Status: toolchain.Missing}}
	d.m.det = toolchain.NewStub(stub)
	d.m.lang = "python"
	d.m.installLang = "python"
	d.m.installReturn = screenBench
	d.m.screen = screenInstallHelp
	stub["python"] = toolchain.Probe{Status: toolchain.Available} // they installed it
	d.runes("R")
	if d.m.screen != screenBench {
		t.Fatalf("re-check after install should return to the prior screen, got %d", d.m.screen)
	}
}

// The picker shows per-language toolchain availability (Axis A) from the detector.
func TestPicker_ShowsToolchainStatus(t *testing.T) {
	d := newDriver(t)
	d.m.det = toolchain.NewStub(map[string]toolchain.Probe{
		"python": {Status: toolchain.Available},
		// every other language is unseeded → Missing
	})
	d.m.screen = screenLanguage
	out := d.m.View()
	if !strings.Contains(out, "installed") {
		t.Fatalf("picker should show install status text, got:\n%s", out)
	}
	if !strings.Contains(out, "✓") || !strings.Contains(out, "✗") {
		t.Fatalf("picker should show both ✓ (python) and ✗ (others) icons, got:\n%s", out)
	}
}

// [i] on the picker opens the install guide for the highlighted language and
// returns to the picker on exit (Axis A shortcut from the picker).
func TestPicker_InstallHelpShortcut(t *testing.T) {
	d := newDriver(t)
	d.m.det = toolchain.NewStub(map[string]toolchain.Probe{"python": {Status: toolchain.Missing, Reason: "not found"}})
	d.m.screen = screenLanguage
	d.m.langIdx = 0 // python is first in availableLangs
	d.runes("i")
	if d.m.screen != screenInstallHelp || d.m.installLang != "python" {
		t.Fatalf("[i] should open install help for the highlighted language, got screen %d lang %q", d.m.screen, d.m.installLang)
	}
	if d.m.installReturn != screenLanguage {
		t.Fatalf("install help from the picker should return to the picker")
	}
}

// Cursor nav + enter selects the highlighted language — and selecting a NOT-
// installed language is still allowed (so its primers remain readable; this is
// the whole point of gate-the-action vs blocking selection).
func TestPicker_CursorSelectsUninstalledForReading(t *testing.T) {
	d := newDriver(t)
	d.m.det = toolchain.NewStub(map[string]toolchain.Probe{}) // all Missing
	d.m.screen = screenLanguage
	langs := d.m.availableLangs()
	if len(langs) < 2 {
		t.Skip("need >=2 offered languages")
	}
	d.runes("j") // move highlight to index 1
	if d.m.langIdx != 1 {
		t.Fatalf("[j] should move the cursor to 1, got %d", d.m.langIdx)
	}
	d.enter() // select the highlighted (uninstalled) language
	if d.m.lang != langs[1].key {
		t.Fatalf("enter should select the highlighted language %q, got %q", langs[1].key, d.m.lang)
	}
	if d.m.screen != screenEditor {
		t.Fatalf("selecting an uninstalled language should still proceed (to read primers), got screen %d", d.m.screen)
	}
}

// Advanced Topics are now PLAYABLE: [r] grades the player's attempt via the
// native toolchain, gated on availability; [b] reveals. Uses a Python check:tests
// exercise (the only fully-gradeable advanced content today).
func TestAdvancedTopics_PlayAndGrade(t *testing.T) {
	d := newDriver(t)
	d.m.lang = "python"
	topics := d.m.cat.AdvancedTopicsByLang("python")
	d.m.advTopics = topics
	topicIdx, exIdx := -1, -1
	var theEx content.Exercise
	for ti, tp := range topics {
		for ei, ex := range tp.Exercises {
			if ex.Check == "tests" {
				topicIdx, exIdx, theEx = ti, ei, ex
				break
			}
		}
		if topicIdx >= 0 {
			break
		}
	}
	if topicIdx < 0 {
		t.Skip("no Python check:tests advanced exercise found")
	}
	nm, _ := d.m.openAdvancedTopic(topicIdx)
	d.m = nm.(Model)
	for i, pg := range d.m.advancedPages() {
		if pg.ex != nil && pg.exNum == exIdx+1 {
			d.m.advPage = i
			break
		}
	}

	// Axis A: no Python toolchain → grading routes to install help.
	d.m.det = toolchain.NewStub(map[string]toolchain.Probe{})
	d.m.advAttempt = theEx.FixedCode
	d.runes("r")
	if d.m.screen != screenInstallHelp {
		t.Fatalf("grading without the toolchain should open install help, got screen %d", d.m.screen)
	}

	// With Python available but no attempt → prompts to write.
	d.m.screen = screenAdvancedTopic
	d.m.det = toolchain.NewStub(map[string]toolchain.Probe{"python": {Status: toolchain.Available}})
	d.m.advAttempt = ""
	d.runes("r")
	if !strings.Contains(d.m.advStatus, "Write your fix") {
		t.Fatalf("empty attempt should prompt to write, got %q", d.m.advStatus)
	}

	// Real grade through the player's Python: the model fix PASSES, broken FAILS.
	if toolchain.New().Presence("python").Status == toolchain.Missing {
		return // no system python; the routing assertions above still ran
	}
	d.m.advAttempt = theEx.FixedCode
	if msg, ok := d.m.advGradeCmd(theEx)().(advGradeMsg); !ok || !msg.v.Passed {
		t.Fatalf("model fix should grade PASS, got %#v", msg)
	}
	d.m.advAttempt = theEx.BrokenCode
	if msg, ok := d.m.advGradeCmd(theEx)().(advGradeMsg); !ok || msg.v.Passed {
		t.Fatalf("broken code should grade FAIL, got %#v", msg)
	}
}

// The bench is now playable in Go (Model A function-call grading): a Go session
// gets a generated Go starter stub and [r] grades through the Go toolchain — no
// longer reference-only.
func TestBench_PlayableInGo(t *testing.T) {
	d := newDriver(t)
	d.m.lang = "go"
	d.m.det = toolchain.NewStub(map[string]toolchain.Probe{"go": {Status: toolchain.Available}})

	// A code task in a Go session gets a generated Go stub, not the Python starter.
	d.m.task = &codeTask{funcName: "add", code: "def add(a, b):\n    pass\n",
		tests: []grader.TestCase{{Name: "t", Input: []any{2, 3}, Expected: 5}}}
	(&d.m).applyLangStarter("")
	if !strings.Contains(d.m.task.code, "func add(") {
		t.Fatalf("Go session should get a generated Go stub, got %q", d.m.task.code)
	}

	// [r] must grade (not reference-only) when the Go toolchain is available.
	d.m.ctx = ctxBench
	d.m.screen = screenBench
	d.m.task.code = "func add(a int, b int) int { return a + b }"
	d.runes("r")
	if strings.Contains(d.m.status, "reference-only") {
		t.Fatalf("Go should be gradeable on the bench, got %q", d.m.status)
	}

	// Real end-to-end grade through the Go toolchain (skip if Go absent).
	if toolchain.New().Presence("go").Status == toolchain.Missing {
		return
	}
	if msg, ok := d.m.runTask()().(gradeMsg); !ok || !msg.v.Passed {
		t.Fatalf("a correct Go solution should pass, got %#v", msg)
	}
}

// The Learn panel pages a rich (sectioned) primer one section at a time: opens on
// the overview, →/n advances and clamps at the end, ←/p goes back and clamps at 0,
// enter returns. Uses the Java Strings primer (the first migrated to sections).
func TestPrimer_SectionPager(t *testing.T) {
	d := newDriver(t)
	d.m.lang = "java"
	d.m.screen = screenBench
	d.m.ctx = ctxBench
	d.m.curProblem = content.Problem{Category: "Strings"}
	d.runes("l")
	if d.m.screen != screenPrimer || d.m.primer.Lang != "java" {
		t.Fatalf("[l] should open the Java Strings primer, got screen %d lang %q", d.m.screen, d.m.primer.Lang)
	}
	pages := d.m.primerPages()
	if len(pages) < 4 {
		t.Fatalf("a rich Strings primer should page into several sections, got %d pages", len(pages))
	}
	if d.m.primerPage != 0 {
		t.Fatalf("primer should open on the overview page, got %d", d.m.primerPage)
	}

	// page forward through every page, then confirm it clamps at the end
	for i := 1; i < len(pages); i++ {
		d.runes("n")
		if d.m.primerPage != i {
			t.Fatalf("next should advance to page %d, got %d", i, d.m.primerPage)
		}
	}
	d.runes("n")
	if d.m.primerPage != len(pages)-1 {
		t.Fatalf("next at the last page should clamp at %d, got %d", len(pages)-1, d.m.primerPage)
	}

	// page back to the start, then confirm it clamps at 0
	for i := len(pages) - 2; i >= 0; i-- {
		d.runes("p")
		if d.m.primerPage != i {
			t.Fatalf("prev should go to page %d, got %d", i, d.m.primerPage)
		}
	}
	d.runes("p")
	if d.m.primerPage != 0 {
		t.Fatalf("prev at the first page should clamp at 0, got %d", d.m.primerPage)
	}

	// the built-in-functions section must exist and carry equals(). Assert against
	// the RAW section data — the rendered pg.body is now syntax-highlighted (ANSI
	// escapes interleave the tokens), so substring checks belong on the source.
	found := false
	for _, sec := range d.m.primer.Sections {
		if !strings.Contains(strings.ToLower(sec.Title), "built-in") {
			continue
		}
		for _, op := range sec.Ops {
			if strings.Contains(op.Code, "equals") {
				found = true
			}
		}
	}
	if !found {
		t.Fatal("expected a built-in-functions section containing equals()")
	}

	d.enter() // back to the problem
	if d.m.screen != screenBench {
		t.Fatalf("enter should return to bench, got %d", d.m.screen)
	}
}

// Every primer's worked example must split cleanly into prose + a code block, so
// the syntax highlighter colors the code and leaves the prose alone. This verifies
// the authoring convention (prose at indent ≤2, code at indent ≥4) across ALL
// primers at once — it fails loudly if a future file breaks it.
func TestPrimerExampleSplitsProseAndCode(t *testing.T) {
	c, err := content.Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range c.Primers {
		if strings.TrimSpace(p.Example) == "" {
			continue
		}
		intro, code, _ := splitExample(p.Example)
		if strings.TrimSpace(code) == "" {
			t.Errorf("primer %q (%s): no code block detected in the worked example (indent the code ≥4 spaces)", p.Category, p.Lang)
		}
		if strings.TrimSpace(intro) == "" {
			t.Errorf("primer %q (%s): no intro prose before the code block", p.Category, p.Lang)
		}
	}
}

// The Advanced Topics surface: the bench menu offers it for a language that has
// topics; opening it lists topics; opening a topic pages through the explainer and
// exercises; [r] toggles the reveal on an exercise page.
func TestAdvancedTopicsSurface(t *testing.T) {
	d := newDriver(t)
	d.m.lang = "python" // python has authored advanced topics
	if len(d.m.cat.AdvancedTopicsByLang("python")) == 0 {
		t.Skip("no python advanced topics loaded")
	}

	// the bench browse menu must include the Advanced Topics entry
	nm, _ := d.m.startBench()
	d.m = nm.(Model)
	advIdx := -1
	for i, o := range d.m.benchMenu {
		if o.kind == "advanced" {
			advIdx = i
		}
	}
	if advIdx < 0 {
		t.Fatal("bench menu missing the Advanced Topics entry for a language that has topics")
	}

	// select it → topic list
	d.m.benchMenuIdx = advIdx
	d.enter()
	if d.m.screen != screenAdvancedList {
		t.Fatalf("selecting Advanced Topics should open the list, got screen %d", d.m.screen)
	}
	if len(d.m.advTopics) == 0 {
		t.Fatal("advanced topic list is empty")
	}

	// open the first topic → pager
	d.enter()
	if d.m.screen != screenAdvancedTopic {
		t.Fatalf("enter should open a topic, got screen %d", d.m.screen)
	}
	pages := d.m.advancedPages()
	if len(pages) < 2 { // at least overview + something
		t.Fatalf("topic should have multiple pages, got %d", len(pages))
	}

	// page to the first EXERCISE page and prove [b] toggles the reveal
	exFound := false
	for i, pg := range pages {
		if pg.ex != nil {
			d.m.advPage = i
			d.m.advReveal = false
			d.runes("b")
			if !d.m.advReveal {
				t.Fatal("[b] should reveal the bug/fix on an exercise page")
			}
			d.runes("b")
			if d.m.advReveal {
				t.Fatal("[b] should toggle the reveal back off")
			}
			exFound = true
			break
		}
	}
	if exFound {
		// View must render without panic on a revealed exercise page
		d.m.advReveal = true
		if d.m.View() == "" {
			t.Fatal("empty render on a revealed exercise page")
		}
	}

	// esc returns to the topic list
	d.step(tea.KeyMsg{Type: tea.KeyEsc})
	if d.m.screen != screenAdvancedList {
		t.Fatalf("esc should return to the topic list, got screen %d", d.m.screen)
	}
}

func TestFilterProblems_List(t *testing.T) {
	c, _ := content.Load()
	got := filterProblems(c.Problems, "list", "blind75")
	if len(got) == 0 {
		t.Fatal("no blind75 problems after filter")
	}
	for _, p := range got {
		ok := false
		for _, l := range p.Lists {
			if l == "blind75" {
				ok = true
			}
		}
		if !ok {
			t.Fatalf("list filter leaked %q", p.Title)
		}
	}
}

// Node-harness problems are present + grade-valid is covered by content tests;
// here verify the catalog has the new categories.
func TestBench_NodeCategoriesPresent(t *testing.T) {
	c, err := content.Load()
	if err != nil {
		t.Fatal(err)
	}
	cats := map[string]int{}
	for _, p := range c.Problems {
		cats[p.Category]++
	}
	for _, want := range []string{"Linked List", "Trees & Graphs", "Tries"} {
		if cats[want] == 0 {
			t.Errorf("missing category %q", want)
		}
	}
	if len(c.Problems) < 137 {
		t.Fatalf("expected 137 problems, got %d", len(c.Problems))
	}
}

// Step 0 completion milestone: banking enough distinct problems across enough
// categories with enough hards reaches the milestone.
func TestStep0_Milestone(t *testing.T) {
	d := newDriver(t)
	// pick a set spanning >=6 categories incl >=2 hard, >=15 problems
	bySeen := map[string]bool{}
	hard := 0
	for _, p := range d.m.cat.Problems {
		take := false
		if !bySeen[p.Category] && len(bySeen) < 8 {
			take = true
		} else if p.Difficulty == "hard" && hard < 3 {
			take = true
		} else if len(d.m.solvedSet) < 15 {
			take = true
		}
		if take {
			d.m.solvedSet[p.ID] = true
			bySeen[p.Category] = true
			if p.Difficulty == "hard" {
				hard++
			}
		}
	}
	bk, ca, hd := d.m.benchStats()
	if bk < step0BankTarget || ca < step0CatTarget || hd < step0HardTarget {
		t.Skipf("seed set didn't reach targets (banked %d cats %d hard %d)", bk, ca, hd)
	}
	if !d.m.step0Met() {
		t.Fatalf("step0Met should be true with banked %d cats %d hard %d", bk, ca, hd)
	}
	ps, lp, track := d.m.step0Profile()
	if ps <= 0 || lp <= 0 || track == "" {
		t.Fatalf("profile should be populated: ps=%d lp=%d track=%q", ps, lp, track)
	}
}

func TestStep0_NotMetEarly(t *testing.T) {
	d := newDriver(t)
	// only a couple solved → not met
	count := 0
	for _, p := range d.m.cat.Problems {
		d.m.solvedSet[p.ID] = true
		count++
		if count >= 3 {
			break
		}
	}
	if d.m.step0Met() {
		t.Fatal("step0 should NOT be met after only 3 solves")
	}
}

func TestResume_Step0Done(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	if err := save.Save(save.State{Language: "python", Stage: "step0done", Step0Done: true, SolvedIDs: []string{"reverse-linked-list", "merge-intervals"}}); err != nil {
		t.Fatal(err)
	}
	m := New()
	nm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("c")})
	m = nm.(Model)
	if m.screen != screenStep0Complete {
		t.Fatalf("want step0-complete screen on resume, got %d", m.screen)
	}
	if !m.step0Done || len(m.solvedSet) != 2 {
		t.Fatalf("solved set / done flag not restored: done=%v set=%d", m.step0Done, len(m.solvedSet))
	}
}

// ── real grader, full intake (selection-agnostic) ─────────────────────────────

func TestIntake_RealGrader_TestOut(t *testing.T) {
	requirePython(t)
	d := newDriver(t)
	d.toIntake("regularly")
	for d.m.screen == screenDiagnostic {
		switch d.m.curDiag.Kind {
		case "code":
			d.solveCodeReal()
		case "choice":
			d.chooseCorrect()
		case "spec":
			d.specAnswerCorrect()
		}
	}
	if d.m.screen != screenTestOut {
		t.Fatalf("regularly + all-correct (real grader) should test-out, got %d", d.m.screen)
	}
}

func TestWrongSolution_RealGrader_Fails(t *testing.T) {
	requirePython(t)
	d := newDriver(t)
	d.toIntake("a-little") // floor is a code item
	fn := d.m.task.funcName
	if fn == "" {
		t.Fatal("expected a code item at the floor")
	}
	starter := d.m.task.code
	d.step(editorFinishedMsg{code: starter})
	v, err := grader.NewNativePython().Run("python", starter, fn, d.m.task.tests, d.m.task.shape)
	if err != nil {
		t.Fatalf("grader error: %v", err)
	}
	d.step(gradeMsg{v: v})
	if v.Passed {
		t.Fatalf("starter stub for %q should NOT pass", fn)
	}
	if d.m.screen != screenDiagnostic {
		t.Fatalf("failing solution should stay on the diagnostic, got %d", d.m.screen)
	}
}

// ── live wiring: New() → grader.New() → wazero, Init() pre-warm ───────────────

// TestLiveWazeroWiring exercises the integrated path the wazero default actually
// creates — which every other test bypasses by injecting gradeMsg directly:
// tui.New() routes through grader.New() (forced to wazero here), Init() returns
// the background Warm() command, and a real grade then runs through m.g.Run.
// Opt-in (needs a local python.wasm; slow) so CI stays hermetic. This is the
// automatable half; the keyboard-driven submit remains human-playtest territory.
func TestLiveWazeroWiring(t *testing.T) {
	if testing.Short() {
		t.Skip("live wazero wiring is slow; skipped in -short")
	}
	wasm := os.Getenv("DEVASCENT_PYTHON_WASM")
	if _, err := os.Stat(wasm); err != nil {
		t.Skip("python.wasm not found; set DEVASCENT_PYTHON_WASM to run the live wazero wiring test")
	}
	t.Setenv("DEVASCENT_GRADER", "wazero")
	t.Setenv("DEVASCENT_PYTHON_WASM", wasm)

	m := New() // routes through grader.New() → WazeroPython with WasmPath=wasm
	if _, ok := m.g.(grader.Warmer); !ok {
		t.Fatalf("default grader is not the warmable wazero backend: %T", m.g)
	}

	// Run Init()'s background command (Warm) synchronously; must not panic/hang.
	if cmd := m.Init(); cmd != nil {
		_ = cmd()
	}

	// Now grade a trivial problem through the LIVE grader the TUI constructed.
	v, err := m.g.Run("python", "def add(a, b):\n    return a + b\n", "add",
		[]grader.TestCase{{Name: "t", Input: []any{2, 3}, Expected: 5}}, grader.Shape{})
	if err != nil {
		t.Fatalf("live wazero grade errored: %v", err)
	}
	if !v.Passed {
		t.Fatalf("live wazero grade did not pass: %+v", v)
	}
}

// TestStartTutorial_RendersSessionLanguage is the end-to-end assertion of the
// reported bug: Tutorial Island must teach the player's chosen language, not
// Python. Starting the tutorial as a Go player yields Go lessons whose first
// teaching body shows `func` (not `def`); Python still shows `def`.
func TestStartTutorial_RendersSessionLanguage(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	for _, tc := range []struct{ lang, want, notWant string }{
		{"go", "func ", "def "},
		{"rust", "fn ", "def "},
		{"python", "def ", "func "},
	} {
		d := newDriver(t)
		d.m.lang = tc.lang
		nm, _ := d.m.startTutorial()
		m := nm.(Model)
		if m.screen != screenLesson {
			t.Fatalf("%s: startTutorial should reach the lesson screen, got %d", tc.lang, m.screen)
		}
		if len(m.lessons) != 10 {
			t.Fatalf("%s: expected 10 lessons in the session language, got %d", tc.lang, len(m.lessons))
		}
		body := m.lessons[0].Stages[0].Body
		if !strings.Contains(body, tc.want) {
			t.Errorf("%s: first lesson body should teach %q:\n%s", tc.lang, tc.want, body)
		}
		if strings.Contains(body, tc.notWant) {
			t.Errorf("%s: first lesson body should NOT contain %q (wrong language)", tc.lang, tc.notWant)
		}
	}
}
