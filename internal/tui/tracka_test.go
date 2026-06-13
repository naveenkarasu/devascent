package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"devascent/internal/economy"
	"devascent/internal/grader"
)

// writeupAnswerCorrect drives the write-up screen to acceptance.
func (d *driver) writeupAnswerCorrect(note string) {
	d.t.Helper()
	if d.m.screen != screenWriteup {
		d.t.Fatalf("want write-up screen, got %d", d.m.screen)
	}
	if d.m.wuHasMCQ {
		for i := 0; i < d.m.wuMCQ.Correct; i++ {
			d.step(tea.KeyMsg{Type: tea.KeyDown})
		}
		d.enter()
		if d.m.wuPhase != 1 {
			d.t.Fatalf("correct MCQ should open the note field (err %q)", d.m.wuErr)
		}
	}
	d.runes(note)
	d.enter()
}

func TestWriteup_AcceptFlow(t *testing.T) {
	d := newDriver(t)
	if len(d.m.cat.Problems) == 0 {
		t.Skip("no bench problems")
	}
	d.enterBenchAll()
	prob := d.m.curProblem
	d.passCode()

	if d.m.screen != screenWriteup {
		t.Fatalf("solve should open the write-up, got screen %d", d.m.screen)
	}
	wantClean := economy.SolveAward(prob.Difficulty)
	if d.m.wallet.Tokens != economy.StartTokens+wantClean {
		t.Fatalf("clean-solve award missing: tokens %d", d.m.wallet.Tokens)
	}

	// Wrong MCQ first: stays on the question with feedback.
	if d.m.wuHasMCQ {
		wrong := (d.m.wuMCQ.Correct + 1) % len(d.m.wuMCQ.Options)
		for i := 0; i < wrong; i++ {
			d.step(tea.KeyMsg{Type: tea.KeyDown})
		}
		d.enter()
		if d.m.wuPhase != 0 || d.m.wuErr == "" {
			t.Fatalf("wrong MCQ should retry, phase %d err %q", d.m.wuPhase, d.m.wuErr)
		}
		for i := 0; i < wrong; i++ { // cursor back to the top
			d.step(tea.KeyMsg{Type: tea.KeyUp})
		}
	}

	before := d.m.wallet.Tokens
	d.writeupAnswerCorrect("Scanned once, kept a hash map of what I'd already seen.")
	if d.m.screen != screenBench {
		t.Fatalf("accepted write-up should resume the bench, got screen %d", d.m.screen)
	}
	if !d.m.solveRecords[prob.ID].WriteupDone {
		t.Fatal("write-up not recorded")
	}
	if d.m.wallet.Tokens < before+economy.WriteupAward {
		t.Fatalf("write-up award missing: %d → %d", before, d.m.wallet.Tokens)
	}
}

func TestWriteup_ShortNoteRejected(t *testing.T) {
	d := newDriver(t)
	if len(d.m.cat.Problems) == 0 {
		t.Skip("no bench problems")
	}
	d.enterBenchAll()
	d.passCode()
	if d.m.wuHasMCQ {
		for i := 0; i < d.m.wuMCQ.Correct; i++ {
			d.step(tea.KeyMsg{Type: tea.KeyDown})
		}
		d.enter()
	}
	d.runes("too short")
	d.enter()
	if d.m.screen != screenWriteup || d.m.wuErr == "" {
		t.Fatalf("short note should be rejected in place: screen %d err %q", d.m.screen, d.m.wuErr)
	}
}

func TestTrackA_StatePersistsThroughResume(t *testing.T) {
	d := newDriver(t)
	if len(d.m.cat.Problems) == 0 {
		t.Skip("no bench problems")
	}
	d.enterBenchAll()
	prob := d.m.curProblem
	d.passCode()
	d.writeupAnswerCorrect("One pass with a frequency map, updating the answer as I went.")

	s := d.m.currentState()
	if !s.WalletInit || s.Tokens != d.m.wallet.Tokens {
		t.Fatalf("wallet not in save: %+v", s)
	}
	if !s.SolveRecords[prob.ID].WriteupDone {
		t.Fatalf("solve record not in save: %+v", s.SolveRecords)
	}

	m2 := Model{cat: d.m.cat, probByID: d.m.probByID, rng: d.m.rng}
	m2.resume = &s
	nm, _ := m2.applyResume()
	m2 = nm.(Model)
	if !m2.wallet.Init || m2.wallet.Tokens != s.Tokens {
		t.Fatalf("wallet not restored: %+v", m2.wallet)
	}
	if !m2.solveRecords[prob.ID].WriteupDone {
		t.Fatal("solve records not restored")
	}
}

func TestRecordFail_CountsDistinctOnly(t *testing.T) {
	d := newDriver(t)
	if len(d.m.cat.Problems) == 0 {
		t.Skip("no bench problems")
	}
	d.enterBenchAll()
	id := d.m.curProblem.ID

	d.m.recordFail(id, "attempt one")
	d.m.recordFail(id, "attempt one") // identical re-run → must NOT count again
	if got := d.m.solveRecords[id].FailedRuns; got != 1 {
		t.Fatalf("re-running the same code counted as a new failure: FailedRuns=%d", got)
	}
	d.m.recordFail(id, "attempt two") // genuinely different → counts
	if got := d.m.solveRecords[id].FailedRuns; got != 2 {
		t.Fatalf("distinct failure not counted: FailedRuns=%d", got)
	}
}

func TestHints_NudgeAndStrategyEconomy(t *testing.T) {
	d := newDriver(t)
	if len(d.m.cat.Problems) == 0 {
		t.Skip("no bench problems")
	}
	d.enterBenchAll()
	prob := d.m.curProblem

	d.runes("h")
	if !d.m.hintMode {
		t.Fatal("[h] should open the hint picker")
	}

	// [p] preview: shows EXACTLY the context pack, costs nothing, and never
	// contains the hidden tests / canonical solution / MCQ answer.
	d.runes("p")
	if !strings.Contains(d.m.hintText, prob.Title) {
		t.Fatalf("preview missing the problem: %q", d.m.hintText)
	}
	if prob.Solution != "" && strings.Contains(d.m.hintText, prob.Solution) {
		t.Fatal("preview leaked the canonical solution")
	}
	if d.m.wallet.Tokens != economy.StartTokens || d.m.wallet.NudgeCharges != economy.NudgeMax {
		t.Fatalf("preview was not free: %+v", d.m.wallet)
	}

	d.runes("1")
	if d.m.hintText == "" || d.m.wallet.NudgeCharges != economy.NudgeMax-1 {
		t.Fatalf("nudge: text %q charges %d", d.m.hintText, d.m.wallet.NudgeCharges)
	}
	if d.m.wallet.Tokens != economy.StartTokens {
		t.Fatalf("nudge cost tokens: %d", d.m.wallet.Tokens)
	}
	if d.m.solveRecords[prob.ID].HintTier != economy.TierNone {
		t.Fatal("nudge must not record a paid tier")
	}

	// Strategy: first press arms, second spends (template path, no AI configured).
	d.runes("2")
	if d.m.wallet.Tokens != economy.StartTokens {
		t.Fatalf("arming spent tokens: %d", d.m.wallet.Tokens)
	}
	d.runes("2")
	if d.m.wallet.Tokens != economy.StartTokens-economy.StrategyCost {
		t.Fatalf("strategy not debited: %d", d.m.wallet.Tokens)
	}
	if d.m.solveRecords[prob.ID].HintTier != economy.TierStrategy {
		t.Fatalf("tier not recorded: %+v", d.m.solveRecords[prob.ID])
	}
	if d.m.hintText == "" {
		t.Fatal("no strategy text")
	}

	// A hinted solve banks with NO clean award.
	d.runes("h") // close the picker
	tokens := d.m.wallet.Tokens
	d.passCode()
	if d.m.screen != screenWriteup {
		t.Fatalf("hinted solve still gates on write-up, got %d", d.m.screen)
	}
	if d.m.wallet.Tokens != tokens {
		t.Fatalf("hinted solve was paid the clean award: %d → %d", tokens, d.m.wallet.Tokens)
	}
}

func TestBenchMenu_GateAndMentorEntries(t *testing.T) {
	d := newDriver(t)
	if len(d.m.cat.Problems) == 0 {
		t.Skip("no bench problems")
	}
	d.driveToHandoff()
	d.runes("b")

	gateIdx, mentorIdx := -1, -1
	for i, o := range d.m.benchMenu {
		switch o.kind {
		case "gate":
			gateIdx = i
		case "mentor":
			mentorIdx = i
		}
	}
	if gateIdx < 0 || mentorIdx < 0 {
		t.Fatalf("menu missing gate/mentor entries: %+v", d.m.benchMenu)
	}

	d.m.benchMenuIdx = gateIdx
	d.enter()
	if d.m.screen != screenGate {
		t.Fatalf("gate entry: screen %d", d.m.screen)
	}
	if v := d.m.renderGate(); !strings.Contains(v, "Blind 75") || !strings.Contains(v, "Mandatory") {
		t.Fatalf("gate render: %q", v)
	}
	d.step(tea.KeyMsg{Type: tea.KeyEsc})
	if d.m.screen != screenBenchMenu {
		t.Fatalf("esc from gate: screen %d", d.m.screen)
	}

	d.m.benchMenuIdx = mentorIdx
	d.enter()
	if d.m.screen != screenMentor || len(d.m.mentorRows) < 6 {
		t.Fatalf("mentor entry: screen %d rows %d", d.m.screen, len(d.m.mentorRows))
	}
	if d.m.mentorRows[0].ID != "template" || !d.m.mentorRows[0].Present {
		t.Fatalf("templates row missing: %+v", d.m.mentorRows[0])
	}
}

func TestWriteupQueue_FromBenchMenu(t *testing.T) {
	d := newDriver(t)
	if len(d.m.cat.Problems) == 0 {
		t.Skip("no bench problems")
	}
	d.enterBenchAll()
	d.passCode()
	d.runes("s")                                       // keep provisional
	d.step(gradeMsg{v: grader.Verdict{Passed: false}}) // unrelated fail on next problem

	// Back at the menu the pending write-up shows up and reopens the gate.
	nm, _ := d.m.startBench()
	d.m = nm.(Model)
	wuIdx := -1
	for i, o := range d.m.benchMenu {
		if o.kind == "writeups" {
			wuIdx = i
		}
	}
	if wuIdx < 0 {
		t.Fatalf("no write-ups-pending entry: %+v", d.m.benchMenu)
	}
	d.m.benchMenuIdx = wuIdx
	d.enter()
	if d.m.screen != screenWriteup || len(d.m.wuQueue) != 1 {
		t.Fatalf("write-up queue: screen %d queue %v", d.m.screen, d.m.wuQueue)
	}
	d.writeupAnswerCorrect("Standard pattern for the category; explained the invariant.")
	if d.m.screen != screenBenchMenu {
		t.Fatalf("menu-launched write-up should return to the menu, got %d", d.m.screen)
	}
}
