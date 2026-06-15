// Package ticket is the Step-1 ticketing data model: the in-game Jira/Linear-style
// board that is the apprenticeship home and the container for all ladder/git/cloud
// content (every task is a Ticket). It is PURE DATA + a workflow state machine —
// no rendering (that's the board view) and no grading (the workbench calls
// grader.GradeRepo and feeds the boolean outcome into Ticket.Resolve). Types carry
// json tags so a Sprint can persist into the save State slot (#57).
package ticket

// Type is a ticket's work type. Values are the short lowercase tag the board
// renders ("bug"/"story"/…); Label gives the title-case form.
type Type string

const (
	Bug      Type = "bug"
	Feature  Type = "feature"
	Story    Type = "story"
	Task     Type = "task"
	Spike    Type = "spike"
	TechDebt Type = "tech-debt"
	Incident Type = "incident"
)

// Label is the human title-case name of a type.
func (t Type) Label() string {
	switch t {
	case Bug:
		return "Bug"
	case Feature:
		return "Feature"
	case Story:
		return "Story"
	case Task:
		return "Task"
	case Spike:
		return "Spike"
	case TechDebt:
		return "Tech Debt"
	case Incident:
		return "Incident"
	}
	return string(t)
}

// Status is a workflow column. The order BoardColumns renders is To Do → In
// Progress → In Review → Done; Backlog is a separate view and Blocked is an
// orthogonal state a ticket can enter from any active column.
type Status string

const (
	Backlog    Status = "backlog"
	ToDo       Status = "todo"
	InProgress Status = "in-progress"
	InReview   Status = "in-review"
	Done       Status = "done"
	Blocked    Status = "blocked"
)

// BoardColumns is the left-to-right column order the active-sprint board renders.
var BoardColumns = []Status{ToDo, InProgress, InReview, Done}

// Label is the human name of a status, e.g. "In Progress".
func (s Status) Label() string {
	switch s {
	case Backlog:
		return "Backlog"
	case ToDo:
		return "To Do"
	case InProgress:
		return "In Progress"
	case InReview:
		return "In Review"
	case Done:
		return "Done"
	case Blocked:
		return "Blocked"
	}
	return string(s)
}

// Priority is a ticket's urgency. The board palette renders high priorities red
// and the rest faint (see IsHigh).
type Priority string

const (
	PBlocker  Priority = "blocker"
	PCritical Priority = "critical"
	PMajor    Priority = "major"
	PMinor    Priority = "minor"
	PTrivial  Priority = "trivial"
)

// IsHigh reports whether the priority renders as urgent (red on the board).
func (p Priority) IsHigh() bool { return p == PBlocker || p == PCritical || p == PMajor }

// Label is the human title-case name of a priority.
func (p Priority) Label() string {
	switch p {
	case PBlocker:
		return "Blocker"
	case PCritical:
		return "Critical"
	case PMajor:
		return "Major"
	case PMinor:
		return "Minor"
	case PTrivial:
		return "Trivial"
	}
	return string(p)
}

// DueDayFor maps a priority to its SLA due day relative to the day a ticket is
// assigned (Blocker today … Trivial +8) — validated against real-world ITSM SLAs
// (see the Company Sim design doc).
func DueDayFor(p Priority, assignedDay int) int {
	off := 5 // default (Minor)
	switch p {
	case PBlocker:
		off = 0
	case PCritical:
		off = 1
	case PMajor:
		off = 2
	case PMinor:
		off = 5
	case PTrivial:
		off = 8
	}
	return assignedDay + off
}

// Criterion is one acceptance-criteria checkbox.
type Criterion struct {
	Text string `json:"text"`
	Met  bool   `json:"met,omitempty"`
}

// Comment is one entry in a ticket's discussion/review thread. At is an RFC3339
// timestamp set by the caller (the engine stays time-free so logic is testable).
type Comment struct {
	Author string `json:"author"`
	Body   string `json:"body"`
	At     string `json:"at,omitempty"`
}

// LinkKind is the relationship a Link expresses.
type LinkKind string

const (
	LinkBlocks    LinkKind = "blocks"
	LinkBlockedBy LinkKind = "blocked-by"
	LinkRelates   LinkKind = "relates"
)

// Link relates this ticket to another by key.
type Link struct {
	Kind LinkKind `json:"kind"`
	Key  string   `json:"key"`
}

// Grading is a ticket's deterministic grader spec. The trust split the spike (#52)
// flagged lives here: Command and HiddenFiles are TICKET-OWNED (never player-
// editable), StartFiles seed the working copy, and EditPaths are the only paths
// a player may change. Assemble enforces it so a player can't swap in their own
// test or set Command=["true"]. Nil Grading = a non-graded ticket (e.g. a
// ticketing-only scenario graded on structure elsewhere).
type Grading struct {
	Lang        string            `json:"lang"`
	Command     []string          `json:"command"`
	StartFiles  map[string]string `json:"start_files,omitempty"`
	HiddenFiles map[string]string `json:"hidden_files,omitempty"`
	EditPaths   []string          `json:"edit_paths,omitempty"`
}

// Assemble merges the player's edits onto the ticket-owned files and returns the
// (files, command) to hand to grader.GradeRepo. The trust split is enforced:
//   - only paths in EditPaths are taken from playerEdits (default: the StartFiles
//     paths, when EditPaths is empty) — a player can't introduce a new file;
//   - HiddenFiles always overwrite afterward, so the test can't be tampered with;
//   - Command is the ticket's, copied so callers can't mutate the ticket.
func (g Grading) Assemble(playerEdits map[string]string) (files map[string]string, command []string) {
	files = make(map[string]string, len(g.StartFiles)+len(g.HiddenFiles))
	for p, c := range g.StartFiles {
		files[p] = c
	}
	allowed := make(map[string]bool, len(g.EditPaths))
	if len(g.EditPaths) > 0 {
		for _, p := range g.EditPaths {
			allowed[p] = true
		}
	} else {
		for p := range g.StartFiles { // default: the player may edit what they were given
			allowed[p] = true
		}
	}
	for p, c := range playerEdits {
		if allowed[p] {
			files[p] = c
		}
	}
	for p, c := range g.HiddenFiles { // ticket-owned; always wins
		files[p] = c
	}
	return files, append([]string(nil), g.Command...)
}

// Ticket is one unit of work on the board. A graded ticket has a non-nil Grading;
// the workflow guard (MoveTo/Resolve) governs its Status.
type Ticket struct {
	Key         string       `json:"key"` // e.g. "PXF-101"
	Type        Type         `json:"type"`
	Title       string       `json:"title"`
	Desc        string       `json:"desc,omitempty"`
	Status      Status       `json:"status"`
	Priority    Priority     `json:"priority,omitempty"`
	Points      int          `json:"points,omitempty"`
	Assignee    string       `json:"assignee,omitempty"` // "" = unassigned; "you" = the player
	Reporter    string       `json:"reporter,omitempty"` // the NPC who filed it
	Epic        string       `json:"epic,omitempty"`     // epic key
	Labels      []string     `json:"labels,omitempty"`   // domain tags: backend/frontend/cloud/security…
	Components  []string     `json:"components,omitempty"`
	AssignedDay int          `json:"assigned_day,omitempty"` // appears in To Do once the day arrives
	DueDay      int          `json:"due_day,omitempty"`      // SLA deadline (day number); 0 = none
	CreatedDay  int          `json:"created_day,omitempty"`
	ResolvedDay int          `json:"resolved_day,omitempty"` // day it reached Done (for standup recaps)
	Watchers    []string     `json:"watchers,omitempty"`
	Subtasks    []Subtask    `json:"subtasks,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
	Acceptance  []Criterion  `json:"acceptance,omitempty"`
	Links       []Link       `json:"links,omitempty"`
	Comments    []Comment    `json:"comments,omitempty"`
	Learn       string       `json:"learn,omitempty"`   // optional primer/code shown in the detail view
	Grading     *Grading     `json:"grading,omitempty"` // nil = not deterministically graded
}

// Subtask is a checklist item under a ticket.
type Subtask struct {
	Title string `json:"title"`
	Done  bool   `json:"done,omitempty"`
}

// Attachment is repro material shown in the detail (e.g. a log or stack trace).
type Attachment struct {
	Name string `json:"name"`
	Body string `json:"body,omitempty"`
}

// Visible reports whether a ticket has been assigned by the given day (the
// scheduled-reveal gate — you aren't handed the whole sprint at once).
func (t *Ticket) Visible(day int) bool { return t.AssignedDay <= day }

// Overdue reports whether a non-done ticket is past its SLA due day.
func (t *Ticket) Overdue(day int) bool {
	return t.DueDay > 0 && t.Status != Done && day > t.DueDay
}

// NotStarted reports whether work hasn't begun (still Backlog/To Do).
func (t *Ticket) NotStarted() bool { return t.Status == Backlog || t.Status == ToDo }

// Project identifies a board (its key is the ticket-key prefix).
type Project struct {
	Key  string `json:"key"`  // "PXF"
	Name string `json:"name"` // "Pixel Forge"
}

// Sprint is the active set of tickets the board shows.
type Sprint struct {
	Number   int       `json:"number"`
	Goal     string    `json:"goal,omitempty"`
	Day      int       `json:"day,omitempty"`
	Capacity int       `json:"capacity,omitempty"`
	Tickets  []*Ticket `json:"tickets"`
}

// Column returns the sprint's tickets in a given status, in slice order.
func (s *Sprint) Column(st Status) []*Ticket {
	var out []*Ticket
	for _, t := range s.Tickets {
		if t.Status == st {
			out = append(out, t)
		}
	}
	return out
}

// Find returns the ticket with the given key, or nil.
func (s *Sprint) Find(key string) *Ticket {
	for _, t := range s.Tickets {
		if t.Key == key {
			return t
		}
	}
	return nil
}

// Committed is the total story points pulled into the sprint.
func (s *Sprint) Committed() int {
	n := 0
	for _, t := range s.Tickets {
		n += t.Points
	}
	return n
}

// DonePoints is the story points completed (status Done).
func (s *Sprint) DonePoints() int {
	n := 0
	for _, t := range s.Tickets {
		if t.Status == Done {
			n += t.Points
		}
	}
	return n
}

// ColumnVisible returns the tickets in a status that are visible on the given day
// (assigned-day ≤ day). The board renders from this so unassigned-yet tickets stay
// hidden until their day arrives.
func (s *Sprint) ColumnVisible(st Status, day int) []*Ticket {
	var out []*Ticket
	for _, t := range s.Tickets {
		if t.Status == st && t.Visible(day) {
			out = append(out, t)
		}
	}
	return out
}

// Incoming returns tickets not yet visible (assigned on a future day), nearest first.
func (s *Sprint) Incoming(day int) []*Ticket {
	var out []*Ticket
	for _, t := range s.Tickets {
		if !t.Visible(day) {
			out = append(out, t)
		}
	}
	return out
}

// OverdueCount is the number of visible, non-done tickets past their due day.
func (s *Sprint) OverdueCount(day int) int {
	n := 0
	for _, t := range s.Tickets {
		if t.Visible(day) && t.Overdue(day) {
			n++
		}
	}
	return n
}
