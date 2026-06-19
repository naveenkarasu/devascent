package ticket

import (
	"errors"
	"fmt"
)

// allowed lists the structural workflow transitions. The In Review→Done edge is
// deliberately absent: Done is reachable ONLY through Resolve(passed), which is
// where "the grader gates In Review→Done" and "you can't skip In Review" are
// enforced. Done can be reopened to In Progress.
var allowed = map[Status][]Status{
	Backlog:    {ToDo},
	ToDo:       {Backlog, InProgress, Blocked},
	InProgress: {ToDo, InReview, Blocked},
	InReview:   {InProgress, Blocked}, // → Done only via Resolve
	Blocked:    {ToDo, InProgress, InReview},
	Done:       {InProgress}, // reopen
}

// Workflow errors.
var (
	ErrDoneNeedsGrade     = errors.New("ticket: cannot move directly to Done — submit for review and pass the grade (use Resolve)")
	ErrResolveNotInReview = errors.New("ticket: Resolve requires the ticket to be In Review")
	ErrGradeNotPassed     = errors.New("ticket: cannot resolve — the grade has not passed")
)

// TransitionError is returned for a structurally illegal MoveTo.
type TransitionError struct{ From, To Status }

func (e *TransitionError) Error() string {
	return fmt.Sprintf("ticket: illegal transition %s → %s", e.From, e.To)
}

// CanMove reports whether a structural transition from→to is allowed. It is false
// for any →Done edge (that gate is Resolve, not MoveTo).
func CanMove(from, to Status) bool {
	if to == Done {
		return false
	}
	for _, s := range allowed[from] {
		if s == to {
			return true
		}
	}
	return false
}

// MoveTo applies a structural workflow transition. Moving to the current status
// is a no-op. Moving to Done is refused (use Resolve); any other disallowed edge
// returns a *TransitionError.
func (t *Ticket) MoveTo(to Status) error {
	if to == Done {
		return ErrDoneNeedsGrade
	}
	if t.Status == to {
		return nil
	}
	if !CanMove(t.Status, to) {
		return &TransitionError{From: t.Status, To: to}
	}
	t.Status = to
	return nil
}

// Resolve performs the gated In Review→Done transition: it requires the ticket to
// be In Review and the grade (the boolean outcome of grader.GradeRepo, supplied by
// the workbench) to have passed. This is the single place a GRADED ticket becomes Done.
func (t *Ticket) Resolve(passed bool) error {
	if t.Status != InReview {
		return ErrResolveNotInReview
	}
	if !passed {
		return ErrGradeNotPassed
	}
	t.Status = Done
	return nil
}

// ErrGradedNeedsReview is returned when MarkDone is called on a graded ticket.
var ErrGradedNeedsReview = errors.New("ticket: graded tickets must pass review (use Resolve)")

// MarkDone completes an UNGRADED ticket directly — onboarding chores and similar
// self-attested work that has no deterministic grader. Graded tickets (Grading
// set) must go through the review gate (Resolve) and are refused here.
func (t *Ticket) MarkDone() error {
	if t.Grading != nil {
		return ErrGradedNeedsReview
	}
	t.Status = Done
	return nil
}
