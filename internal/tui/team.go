package tui

import (
	"fmt"

	"devascent/internal/ticket"
)

// Pixel Forge's team — the colleagues you can delegate work to (those at/below
// your level, who execute) or escalate to (those above, who give guidance). The
// player's own level rises with the career ladder; for now it defaults to a
// sensible band on the board (see enterBoard) until the ladder is wired to a
// stored career level.

type teammate struct {
	name  string
	role  string
	level int // 1 junior · 2 mid · 3 senior · 4 principal · 5 manager
}

var team = []teammate{
	{"Maya", "Junior Engineer", 1},
	{"Dev", "Mid Engineer", 2},
	{"Priya", "Senior Engineer", 3},
	{"Raj", "Principal Engineer", 4},
	{"Sam", "Engineering Manager", 5},
}

// assigneeOptions are the picker choices: yourself, then the whole team.
func assigneeOptions() []string {
	out := make([]string, 0, len(team)+1)
	out = append(out, "you")
	for _, t := range team {
		out = append(out, t.name)
	}
	return out
}

func teammateByName(name string) (teammate, bool) {
	for _, t := range team {
		if t.name == name {
			return t, true
		}
	}
	return teammate{}, false
}

// assignKind classifies assigning to `name` from the player's level:
//
//	"self"     — you'll do it (real graded work)
//	"delegate" — a teammate at/below your level executes it (on schedule)
//	"escalate" — someone above you: you're asking for help / guidance / requirements
func assignKind(name string, playerLvl int) string {
	if name == "" || name == "you" {
		return "self"
	}
	t, ok := teammateByName(name)
	if !ok {
		return "self"
	}
	if t.level > playerLvl {
		return "escalate"
	}
	return "delegate"
}

// roleOf returns a short "Name · Role" label for the assignee picker.
func roleOf(name string) string {
	if t, ok := teammateByName(name); ok {
		return t.name + " · " + t.role
	}
	return name
}

// applyAssignment reacts to assigning a ticket to a teammate:
//   - delegate → the teammate picks it up now (In Progress) and will deliver by the
//     due day (advanceDay completes it at DueDay — never before).
//   - escalate → an immediate guidance reply from the senior (resolved as advice).
//
// Self/unassigned tickets are untouched (your own graded work).
func applyAssignment(t *ticket.Ticket, playerLvl, day int) {
	switch assignKind(t.Assignee, playerLvl) {
	case "delegate":
		if t.Status == ticket.ToDo || t.Status == ticket.Backlog {
			t.Status = ticket.InProgress
			t.Comments = append(t.Comments, ticket.Comment{
				Author: t.Assignee, Body: fmt.Sprintf("On it — targeting delivery by D%d.", t.DueDay)})
		}
	case "escalate":
		if t.Status != ticket.Done {
			t.Status = ticket.Done
			t.ResolvedDay = day
			t.Comments = append(t.Comments, ticket.Comment{Author: t.Assignee, Body: escalationGuidance(t)})
		}
	}
}

// advanceTeamWork progresses delegated tickets for the given day: not-yet-started
// ones begin, in-progress ones post a daily standup-style progress comment, and
// any that reach their due day are delivered (never before). Your own tickets —
// the real graded work — are left untouched.
func advanceTeamWork(sp *ticket.Sprint, playerLvl, day int) {
	for _, t := range sp.Tickets {
		if t.Assignee == "" || t.Assignee == "you" || assignKind(t.Assignee, playerLvl) != "delegate" {
			continue
		}
		if t.Status == ticket.ToDo && day >= t.AssignedDay {
			t.Status = ticket.InProgress
		}
		if t.Status != ticket.InProgress || t.DueDay <= 0 {
			continue
		}
		if day >= t.DueDay {
			t.Status = ticket.Done
			t.ResolvedDay = day
			t.Comments = append(t.Comments, ticket.Comment{Author: t.Assignee, Body: "Delivered " + t.Key + "."})
		} else {
			t.Comments = append(t.Comments, ticket.Comment{Author: t.Assignee, Body: teammateProgressNote(t, day)})
		}
	}
}

// teammateProgressNote is the templated end-of-day update a teammate posts on a
// ticket they're still working (day k of n, on track). An AI-written update
// replaces this when a mentor backend is wired (S9).
func teammateProgressNote(t *ticket.Ticket, day int) string {
	total := t.DueDay - t.AssignedDay
	if total < 1 {
		total = 1
	}
	dayOf := day - t.AssignedDay
	if dayOf < 1 {
		dayOf = 1
	}
	if dayOf > total {
		dayOf = total
	}
	if t.DueDay-day <= 1 {
		return fmt.Sprintf("Day %d of %d on %s — almost there, delivering by D%d.", dayOf, total, t.Key, t.DueDay)
	}
	return fmt.Sprintf("Day %d of %d on %s — on track for D%d.", dayOf, total, t.Key, t.DueDay)
}

// escalationGuidance is the senior's templated advice when you escalate (an AI
// reply replaces this once a mentor backend is wired — T3).
func escalationGuidance(t *ticket.Ticket) string {
	switch t.Type {
	case ticket.Bug, ticket.Incident:
		return "Reproduce it with a failing test first, then make the smallest fix; watch the edge cases in the acceptance criteria. Ping me on the PR."
	case ticket.Story, ticket.Feature:
		return "Nail the acceptance criteria before you start, slice it thin, and ship behind a flag if it's risky. Happy to pair if you get stuck."
	default:
		return "Keep it small and well-tested, and mirror the patterns already in that module. Ask early if anything's unclear."
	}
}
