package ticket

import (
	"errors"
	"testing"
)

// happy path: Backlog → To Do → In Progress → In Review → (grade passes) → Done.
func TestWorkflow_HappyPath(t *testing.T) {
	tk := &Ticket{Key: "PXF-101", Status: Backlog}
	for _, to := range []Status{ToDo, InProgress, InReview} {
		if err := tk.MoveTo(to); err != nil {
			t.Fatalf("MoveTo(%s): %v", to, err)
		}
	}
	if err := tk.Resolve(true); err != nil {
		t.Fatalf("Resolve(true) from In Review: %v", err)
	}
	if tk.Status != Done {
		t.Fatalf("want Done, got %s", tk.Status)
	}
}

// You cannot skip In Review: there's no In Progress→Done edge, and MoveTo(Done)
// is refused outright with a guiding error.
func TestWorkflow_CannotSkipReview(t *testing.T) {
	tk := &Ticket{Status: InProgress}
	if err := tk.MoveTo(Done); !errors.Is(err, ErrDoneNeedsGrade) {
		t.Fatalf("MoveTo(Done) should be refused with ErrDoneNeedsGrade, got %v", err)
	}
	if tk.Status != InProgress {
		t.Fatalf("status must be unchanged after a refused move, got %s", tk.Status)
	}
	if CanMove(InProgress, Done) {
		t.Fatal("CanMove(InProgress, Done) must be false")
	}
}

// Resolve is gated on BOTH being In Review and the grade passing.
func TestWorkflow_ResolveGating(t *testing.T) {
	notReview := &Ticket{Status: ToDo}
	if err := notReview.Resolve(true); !errors.Is(err, ErrResolveNotInReview) {
		t.Fatalf("Resolve from To Do should error ErrResolveNotInReview, got %v", err)
	}

	failing := &Ticket{Status: InReview}
	if err := failing.Resolve(false); !errors.Is(err, ErrGradeNotPassed) {
		t.Fatalf("Resolve(false) should error ErrGradeNotPassed, got %v", err)
	}
	if failing.Status != InReview {
		t.Fatalf("a failed grade must keep the ticket In Review, got %s", failing.Status)
	}
}

func TestWorkflow_IllegalTransition(t *testing.T) {
	tk := &Ticket{Status: Backlog}
	err := tk.MoveTo(InReview)
	var te *TransitionError
	if !errors.As(err, &te) {
		t.Fatalf("Backlog→In Review should be a *TransitionError, got %v", err)
	}
	if te.From != Backlog || te.To != InReview {
		t.Fatalf("TransitionError fields wrong: %+v", te)
	}
}

func TestWorkflow_BlockedRoundTrip_AndReopen(t *testing.T) {
	tk := &Ticket{Status: InProgress}
	if err := tk.MoveTo(Blocked); err != nil {
		t.Fatalf("InProgress→Blocked: %v", err)
	}
	if err := tk.MoveTo(InProgress); err != nil {
		t.Fatalf("Blocked→InProgress: %v", err)
	}
	// reopen a finished ticket
	done := &Ticket{Status: Done}
	if err := done.MoveTo(InProgress); err != nil {
		t.Fatalf("Done→InProgress (reopen): %v", err)
	}
}

func TestWorkflow_MoveToSelfIsNoop(t *testing.T) {
	tk := &Ticket{Status: InProgress}
	if err := tk.MoveTo(InProgress); err != nil {
		t.Fatalf("MoveTo(self) should be a no-op, got %v", err)
	}
}

// The trust split: only declared paths are taken from the player, and the
// ticket-owned hidden test always wins — a player cannot pass by swapping it out.
func TestGrading_Assemble_TrustSplit(t *testing.T) {
	g := Grading{
		Lang:        "python",
		Command:     []string{"python", "check.py"},
		StartFiles:  map[string]string{"paginate.py": "# start\n"},
		HiddenFiles: map[string]string{"check.py": "ASSERT_REAL\n"},
		EditPaths:   []string{"paginate.py"},
	}
	files, cmd := g.Assemble(map[string]string{
		"paginate.py": "# player fix\n",   // allowed → taken
		"check.py":    "print('cheat')\n", // NOT allowed → ignored (hidden wins anyway)
		"conftest.py": "smuggled\n",       // undeclared new file → ignored
	})

	if files["paginate.py"] != "# player fix\n" {
		t.Errorf("player edit to a declared path should be taken, got %q", files["paginate.py"])
	}
	if files["check.py"] != "ASSERT_REAL\n" {
		t.Errorf("ticket-owned hidden test must win, got %q", files["check.py"])
	}
	if _, ok := files["conftest.py"]; ok {
		t.Error("an undeclared file must not be smuggled in")
	}
	if len(cmd) != 2 || cmd[0] != "python" || cmd[1] != "check.py" {
		t.Errorf("command must be the ticket's, got %v", cmd)
	}
	// returned command is a copy — mutating it must not affect the ticket spec
	cmd[0] = "rm"
	if g.Command[0] != "python" {
		t.Error("Assemble must return a copy of Command, not the underlying slice")
	}
}

// With no EditPaths, the player may edit exactly the files they were given.
func TestGrading_Assemble_DefaultEditPaths(t *testing.T) {
	g := Grading{
		StartFiles:  map[string]string{"main.go": "package main\n"},
		HiddenFiles: map[string]string{"main_test.go": "owned\n"},
	}
	files, _ := g.Assemble(map[string]string{
		"main.go": "edited\n",  // a start file → editable by default
		"evil.go": "smuggle\n", // undeclared → ignored
	})
	if files["main.go"] != "edited\n" {
		t.Errorf("default-editable start file should be taken, got %q", files["main.go"])
	}
	if _, ok := files["evil.go"]; ok {
		t.Error("undeclared file must be ignored when EditPaths is empty")
	}
	if files["main_test.go"] != "owned\n" {
		t.Errorf("hidden file must be present, got %q", files["main_test.go"])
	}
}

func TestSprint_ColumnsAndMetrics(t *testing.T) {
	s := &Sprint{
		Number:   3,
		Capacity: 13,
		Tickets: []*Ticket{
			{Key: "PXF-101", Status: InProgress, Points: 2},
			{Key: "PXF-104", Status: ToDo, Points: 3},
			{Key: "PXF-097", Status: Done, Points: 3},
			{Key: "PXF-095", Status: Done, Points: 2},
		},
	}
	if got := s.Column(Done); len(got) != 2 {
		t.Fatalf("want 2 Done tickets, got %d", len(got))
	}
	if got := s.Committed(); got != 10 {
		t.Errorf("Committed() = %d, want 10", got)
	}
	if got := s.DonePoints(); got != 5 {
		t.Errorf("DonePoints() = %d, want 5", got)
	}
	if s.Find("PXF-101") == nil {
		t.Error("Find should locate an existing ticket")
	}
	if s.Find("NOPE") != nil {
		t.Error("Find should return nil for a missing key")
	}
}

func TestDueDayFor(t *testing.T) {
	want := map[Priority]int{PBlocker: 3, PCritical: 4, PMajor: 5, PMinor: 8, PTrivial: 11}
	for p, w := range want {
		if got := DueDayFor(p, 3); got != w {
			t.Errorf("DueDayFor(%s, 3) = %d, want %d", p, got, w)
		}
	}
}

func TestVisibilityAndOverdue(t *testing.T) {
	s := &Sprint{Tickets: []*Ticket{
		{Key: "A", Status: ToDo, AssignedDay: 0, DueDay: 2},
		{Key: "B", Status: ToDo, AssignedDay: 3, DueDay: 5}, // assigned later
		{Key: "C", Status: Done, AssignedDay: 0, DueDay: 1}, // done → never overdue
	}}
	if got := s.ColumnVisible(ToDo, 1); len(got) != 1 || got[0].Key != "A" {
		t.Fatalf("day 1: only A visible in To Do, got %v", tkeys(got))
	}
	if got := s.Incoming(1); len(got) != 1 || got[0].Key != "B" {
		t.Fatalf("day 1: B should be incoming, got %v", tkeys(got))
	}
	if !s.Tickets[0].Overdue(3) {
		t.Error("A (due 2) should be overdue on day 3")
	}
	if s.Tickets[2].Overdue(5) {
		t.Error("a Done ticket is never overdue")
	}
	if got := s.OverdueCount(3); got != 1 {
		t.Errorf("OverdueCount(3) = %d, want 1", got)
	}
}

func tkeys(ts []*Ticket) []string {
	out := make([]string, len(ts))
	for i, t := range ts {
		out[i] = t.Key
	}
	return out
}

func TestLabels(t *testing.T) {
	if InProgress.Label() != "In Progress" {
		t.Errorf("Status.Label wrong: %q", InProgress.Label())
	}
	if TechDebt.Label() != "Tech Debt" || Bug.Label() != "Bug" {
		t.Errorf("Type.Label wrong: %q %q", TechDebt.Label(), Bug.Label())
	}
	if !PMajor.IsHigh() || PMinor.IsHigh() {
		t.Error("priority IsHigh classification wrong")
	}
	if PMajor.Label() != "Major" {
		t.Errorf("Priority.Label wrong: %q", PMajor.Label())
	}
}
