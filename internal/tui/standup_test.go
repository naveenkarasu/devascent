package tui

import (
	"strings"
	"testing"

	"devascent/internal/ticket"
)

func step1Model() Model {
	proj, sp := seedSprint1() // Day 0
	return Model{screen: screenBoard, lang: "python", step1Home: false,
		boardProject: proj, boardSprint: sp, boardCol: defaultFocusCol(sp, sp.Day)}
}

// [e] → cooldown → [s] skip flips the day (revealing tickets) → [enter] joins the
// morning standup → [enter] starts the working day.
func TestStandup_AdvancesDayAndReveals(t *testing.T) {
	m := step1Model()
	if got := len(m.boardSprint.ColumnVisible(ticket.ToDo, m.boardSprint.Day)); got != 3 {
		t.Fatalf("day 0: 3 To Do visible, got %d", got)
	}
	// [e] ends the working day → the cooldown beat (not the standup yet)
	opened, _ := m.handleBoardKey(mkKey("e"))
	m = opened.(Model)
	if m.screen != screenCooldown || m.boardSprint.Phase != ticket.PhaseCooldown {
		t.Fatalf("e should open the cooldown, got screen %d phase %q", m.screen, m.boardSprint.Phase)
	}
	// [s] skips the wait → the day flips and the standup is pending
	skipped, _ := m.handleCooldownKey(mkKey("s"))
	m = skipped.(Model)
	if m.boardSprint.Day != 1 || m.boardSprint.Phase != ticket.PhaseStandup {
		t.Fatalf("skip should flip to Day 1 standup-pending, got day %d phase %q", m.boardSprint.Day, m.boardSprint.Phase)
	}
	if got := len(m.boardSprint.ColumnVisible(ticket.ToDo, m.boardSprint.Day)); got != 5 {
		t.Fatalf("day 1 should reveal 2 more (5 total) To Do, got %d", got)
	}
	// [enter] joins the morning standup
	joined, _ := m.handleCooldownKey(mkKey("enter"))
	m = joined.(Model)
	if m.screen != screenStandup {
		t.Fatalf("enter should join the standup, got screen %d", m.screen)
	}
	// [enter] starts the working day
	started, _ := m.handleStandupKey(mkKey("enter"))
	m = started.(Model)
	if m.boardSprint.Day != 1 || m.screen != screenBoard || m.boardSprint.Phase != ticket.PhaseWorking {
		t.Fatalf("enter should start Day 1 on the board, got day %d screen %d phase %q", m.boardSprint.Day, m.screen, m.boardSprint.Phase)
	}
}

// The manager flags overdue work at standup.
func TestStandup_ManagerFlagsOverdue(t *testing.T) {
	m := step1Model()
	m.boardSprint.Day = 4 // PXF-201 was due D2 and is untouched → overdue
	note := strings.Join(m.standupNote(), "\n")
	if !strings.Contains(note, "PXF-201") || !strings.Contains(note, "overdue") {
		t.Fatalf("standup note should flag PXF-201 as overdue, got:\n%s", note)
	}
}

// Resolving a ticket records the day it shipped, so the standup can report it.
func TestStandup_DoneTodayCounts(t *testing.T) {
	m := step1Model()
	tk := m.boardSprint.Find("PXF-100") // an ungraded chore
	tk.Status = ticket.Done
	tk.ResolvedDay = m.boardSprint.Day
	if m.doneToday() != 1 {
		t.Fatalf("doneToday should be 1, got %d", m.doneToday())
	}
	m.boardSprint.Day++ // next day → yesterday's done no longer counts as "today"
	if m.doneToday() != 0 {
		t.Errorf("doneToday should reset across days, got %d", m.doneToday())
	}
}
