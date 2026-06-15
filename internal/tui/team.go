package tui

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
