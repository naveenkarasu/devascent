package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"devascent/internal/grader"
	"devascent/internal/ticket"
	"devascent/internal/toolchain"
)

func workableTicket(status ticket.Status) *ticket.Ticket {
	return &ticket.Ticket{
		Key: "PXF-900", Type: ticket.Bug, Title: "demo", Status: status, Assignee: "you",
		Grading: &ticket.Grading{
			Lang: "python", Command: []string{"python", "check.py"},
			StartFiles: map[string]string{"x.py": "buggy"}, EditPaths: []string{"x.py"},
		},
	}
}

// [s] on a To Do ticket starts work: it moves to In Progress and seeds the buffer.
func TestTicketWork_StartWork(t *testing.T) {
	tk := workableTicket(ticket.ToDo)
	m := Model{screen: screenTicket, detailTicket: tk}
	nm, _ := m.handleTicketKey(mkKey("s"))
	got := nm.(Model)
	if tk.Status != ticket.InProgress {
		t.Fatalf("start work should move To Do → In Progress, got %s", tk.Status)
	}
	if got.workCode != "buggy" {
		t.Fatalf("workCode should seed from StartFiles, got %q", got.workCode)
	}
}

// A passing grade moves In Progress → In Review and posts the reviewer question.
func TestTicketWork_PassMovesToReview(t *testing.T) {
	tk := workableTicket(ticket.InProgress)
	m := Model{screen: screenTicket, detailTicket: tk}
	nm, _ := m.applyTicketGrade(grader.Verdict{Passed: true})
	got := nm.(Model)
	if tk.Status != ticket.InReview {
		t.Fatalf("a passing grade should move In Progress → In Review, got %s", tk.Status)
	}
	if got.reviewQ == "" {
		t.Fatal("a reviewer question should be set")
	}
	if len(tk.Comments) != 1 || tk.Comments[0].Author != "Sam (reviewer)" {
		t.Fatalf("reviewer comment should be posted, got %+v", tk.Comments)
	}
}

// A failing grade keeps the ticket In Progress (no skipping review).
func TestTicketWork_FailStays(t *testing.T) {
	tk := workableTicket(ticket.InProgress)
	m := Model{screen: screenTicket, detailTicket: tk}
	m.applyTicketGrade(grader.Verdict{Passed: false})
	if tk.Status != ticket.InProgress {
		t.Fatalf("a failing grade must keep In Progress, got %s", tk.Status)
	}
	m.applyTicketGrade(grader.Verdict{Err: "traceback"})
	if tk.Status != ticket.InProgress {
		t.Fatalf("a grade error must keep In Progress, got %s", tk.Status)
	}
}

// Answering the reviewer resolves In Review → Done and records the answer.
func TestTicketWork_AnswerResolves(t *testing.T) {
	tk := workableTicket(ticket.InReview)
	m := Model{screen: screenTicket, detailTicket: tk}
	m.answerReview("off-by-one start index; the half-open slice keeps the last page safe")
	if tk.Status != ticket.Done {
		t.Fatalf("answering should resolve In Review → Done, got %s", tk.Status)
	}
	if len(tk.Comments) == 0 || tk.Comments[len(tk.Comments)-1].Author != "you" {
		t.Fatalf("the answer should be recorded as a comment, got %+v", tk.Comments)
	}
}

// Full loop on one ticket: In Progress → (pass) → In Review → (answer) → Done.
func TestTicketWork_FullLoop(t *testing.T) {
	tk := workableTicket(ticket.InProgress)
	m := Model{screen: screenTicket, detailTicket: tk}
	m, _ = mustModel(m.applyTicketGrade(grader.Verdict{Passed: true}))
	if tk.Status != ticket.InReview {
		t.Fatalf("expected In Review after pass, got %s", tk.Status)
	}
	m.answerReview("fixed the index math")
	if tk.Status != ticket.Done {
		t.Fatalf("expected Done after answering, got %s", tk.Status)
	}
}

func mustModel(tm tea.Model, c tea.Cmd) (Model, tea.Cmd) { return tm.(Model), c }

// Real end-to-end: grade the demo PXF-101 through GradeRepo. The buggy start code
// fails; the fix passes and moves the card In Progress → In Review. Skips when
// python isn't installed (the grader shells out to it).
func TestTicketWork_RealGrade_Integration(t *testing.T) {
	det := toolchain.New()
	if det.Presence("python").Status == toolchain.Missing {
		t.Skip("python not installed; skipping the real-grade integration test")
	}
	_, sp := demoBoard()
	tk := sp.Find("PXF-101")
	if tk == nil || tk.Grading == nil {
		t.Fatal("demo PXF-101 should be a graded ticket")
	}
	m := Model{screen: screenTicket, detailTicket: tk, g: grader.NewLocalToolchain(det)}

	// the buggy starting code fails the hidden tests
	m.workCode = startCode(tk)
	if rg := m.gradeTicketCmd(tk)().(repoGradeMsg); rg.v.Passed {
		t.Fatalf("buggy start code must fail, got %+v", rg.v)
	}

	// the fix passes, and the pass moves it to In Review
	m.workCode = "def paginate(items, page, size):\n    start = (page - 1) * size\n    return items[start:start + size]\n"
	rg := m.gradeTicketCmd(tk)().(repoGradeMsg)
	if !rg.v.Passed {
		t.Fatalf("fixed code must pass, got %+v", rg.v)
	}
	if _, _ = m.applyTicketGrade(rg.v); tk.Status != ticket.InReview {
		t.Fatalf("a real pass should move PXF-101 to In Review, got %s", tk.Status)
	}
}
