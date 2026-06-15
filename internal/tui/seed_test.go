package tui

import (
	"strings"
	"testing"

	"devascent/internal/grader"
	"devascent/internal/ticket"
	"devascent/internal/toolchain"
)

// Renders the seed board (and logs it for eyeballing).
func TestSeedSprint1_Renders(t *testing.T) {
	proj, sp := seedSprint1()
	sp.Day = 5 // late enough that every scheduled ticket is revealed
	m := Model{screen: screenBoard, boardProject: proj, boardSprint: sp, boardCol: defaultFocusCol(sp, sp.Day)}
	out := m.boardView()
	for _, want := range []string{"Pixel Forge", "Sprint 1", "committed 10", "PXF-201", "PXF-206", "debt", "@you", "due D"} {
		if !strings.Contains(out, want) {
			t.Errorf("seed board missing %q", want)
		}
	}
	t.Log("\n" + out)
}

// On Day 0 only the day-0 tickets are revealed; later ones are "incoming".
func TestSeedSprint1_ScheduledReveal(t *testing.T) {
	_, sp := seedSprint1()                                      // Day 0
	if got := len(sp.ColumnVisible(ticket.ToDo, 0)); got != 3 { // PXF-100, PXF-110, PXF-201
		t.Fatalf("day 0 should reveal 3 To Do tickets, got %d", got)
	}
	if got := len(sp.Incoming(0)); got != 5 { // 202..206 arrive later
		t.Fatalf("day 0 should have 5 incoming tickets, got %d", got)
	}
	if got := len(sp.Incoming(3)); got != 0 { // by day 3 everything is assigned
		t.Errorf("by day 3 nothing should be incoming, got %d", got)
	}
}

// Structural sanity that needs no toolchain.
func TestSeedSprint1_Structure(t *testing.T) {
	proj, sp := seedSprint1()
	if proj.Name == "" {
		t.Fatal("project should have a name")
	}
	if got := sp.Committed(); got != 10 {
		t.Errorf("committed points = %d, want 10", got)
	}
	if got := sp.DonePoints(); got != 0 {
		t.Errorf("done points = %d, want 0 (fresh start, nothing done)", got)
	}
	if sp.Day != 0 {
		t.Errorf("sprint should start on Day 0, got %d", sp.Day)
	}
	for _, tk := range sp.Tickets {
		if tk.Status == ticket.Backlog {
			continue // backlog items are intentionally off the active board
		}
		if tk.Status != ticket.ToDo {
			t.Errorf("%s should start in To Do, got %s", tk.Key, tk.Status)
		}
	}
	graded := 0
	for _, tk := range sp.Tickets {
		if tk.Grading == nil {
			continue
		}
		graded++
		if tk.Assignee != "you" {
			t.Errorf("%s is graded but not assigned @you", tk.Key)
		}
		if len(tk.Grading.EditPaths) == 0 {
			t.Errorf("%s has no edit path", tk.Key)
		}
		if len(tk.Acceptance) == 0 {
			t.Errorf("%s has no acceptance criteria", tk.Key)
		}
	}
	if graded != 6 {
		t.Fatalf("expected 6 graded tickets, got %d", graded)
	}
}

// Every graded ticket must genuinely grade: the starter FAILS the hidden test and
// a known-good solution PASSES it (the bench-corpus pattern). Needs python.
func TestSeedSprint1_AllTicketsGrade(t *testing.T) {
	det := toolchain.New()
	if det.Presence("python").Status == toolchain.Missing {
		t.Skip("python not installed; skipping the seed grading gate")
	}
	g := grader.NewLocalToolchain(det)

	solutions := map[string]string{
		"PXF-201": "def paginate(items, page, size):\n    start = (page - 1) * size\n    return items[start:start + size]\n",
		"PXF-202": "def is_valid(email):\n    parts = email.split('@')\n    if len(parts) != 2:\n        return False\n    local, domain = parts\n    return bool(local) and '.' in domain\n",
		"PXF-203": `def slugify(title):
    out = []
    prev = False
    for ch in title.lower():
        if ch.isalnum():
            out.append(ch)
            prev = False
        elif not prev:
            out.append('-')
            prev = True
    return ''.join(out).strip('-')
`,
		"PXF-204": "def total(nums):\n    return sum(nums)\n",
		"PXF-205": "def average(nums):\n    if not nums:\n        return 0.0\n    return sum(nums) / len(nums)\n",
		"PXF-206": `def dedupe(items):
    seen = set()
    out = []
    for x in items:
        if x not in seen:
            seen.add(x)
            out.append(x)
    return out
`,
	}

	_, sp := seedSprint1()
	for _, tk := range sp.Tickets {
		if tk.Grading == nil {
			continue
		}
		edit := tk.Grading.EditPaths[0]

		// starter fails
		f0, c0 := tk.Grading.Assemble(nil)
		if v, _ := g.GradeRepo(grader.RepoRequest{Files: f0, Command: c0}); v.Passed {
			t.Errorf("%s: starter code should FAIL the hidden test (the test is too weak)", tk.Key)
		}

		// solution passes
		sol, ok := solutions[tk.Key]
		if !ok {
			t.Errorf("%s: no solution provided in the test", tk.Key)
			continue
		}
		f1, c1 := tk.Grading.Assemble(map[string]string{edit: sol})
		if v, _ := g.GradeRepo(grader.RepoRequest{Files: f1, Command: c1}); !v.Passed {
			t.Errorf("%s: the reference solution should PASS, got err=%q stdout=%q", tk.Key, v.Err, v.Stdout)
		}
	}
}
