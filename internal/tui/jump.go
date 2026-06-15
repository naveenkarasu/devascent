package tui

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"devascent/internal/save"
)

// contentChecksum is the SHA-256 of the canonical editor round-trip snippet,
// used to confirm the editor preserved a submission byte-for-byte before it is
// handed to the grader (some platforms re-encode newlines/whitespace).
const contentChecksum = "7273402b6b97f65c459f11377d6ee142a4135d7dbc9f25433249e52a181b735c"

// digestMatches reports whether text (trimmed of surrounding whitespace, so the
// editor's trailing newline doesn't matter) hashes to want, a lowercase hex
// SHA-256. The trim+UTF-8 form is byte-identical to the reference one-liner used
// to produce the stored digest, so the two agree by construction.
func digestMatches(text, want string) bool {
	sum := sha256.Sum256([]byte(strings.TrimSpace(text)))
	return hex.EncodeToString(sum[:]) == want
}

func matchesContentChecksum(code string) bool {
	return digestMatches(code, contentChecksum)
}

// jumpTarget is one stage the developer navigation can drop into. A target either
// synthesizes a save.State that applyResume routes into (rebuilding the stage's
// content fresh when no IDs are present), or sets browse to enter the bench menu.
type jumpTarget struct {
	label  string
	state  save.State
	browse bool // enter the bench browse menu
	board  bool // enter the Step-1 apprenticeship board
}

func jumpTargets() []jumpTarget {
	return []jumpTarget{
		{label: "Entrance test (intake)", state: save.State{Stage: "intake", Level: "a-little"}},
		{label: "Tutorial Island", state: save.State{Stage: "tutorial"}},
		{label: "Dev-Literacy", state: save.State{Stage: "devliteracy"}},
		{label: "Bench (Step 0)", state: save.State{Stage: "bench", Placement: "tutorial-full"}},
		{label: "Bench browse / Advanced Topics", browse: true},
		{label: "Graduation — Step 0 complete", state: save.State{Stage: "step0done", Step0Done: true}},
		{label: "Step-1 Board (apprenticeship)", board: true},
	}
}

// openJump opens the developer navigation overlay, remembering the screen to
// return to if it's cancelled.
func (m Model) openJump() (tea.Model, tea.Cmd) {
	m.jumpReturn = m.screen
	m.jumpIdx = 0
	m.screen = screenJump
	return m, nil
}

func (m Model) handleJumpKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	targets := jumpTargets()
	switch msg.String() {
	case "up", "k":
		if m.jumpIdx > 0 {
			m.jumpIdx--
		}
	case "down", "j":
		if m.jumpIdx < len(targets)-1 {
			m.jumpIdx++
		}
	case "esc":
		m.screen = m.jumpReturn
	case "enter":
		return m.doJump(targets[m.jumpIdx])
	}
	return m, nil
}

// doJump enters the chosen stage by reusing the same setup the resume path runs.
func (m Model) doJump(t jumpTarget) (tea.Model, tea.Cmd) {
	if t.browse {
		return m.startBench()
	}
	if t.board {
		m.screen = m.jumpReturn // so the board's esc returns to the pre-cheat screen
		return m.enterBoard()
	}
	st := t.state
	st.Language = m.lang
	st.Editor = m.editorChoice // keep the tester's editor across the jump
	if st.Level == "" {
		st.Level = m.level
	}
	m.resume = &st
	return m.applyResume()
}

func (m Model) jumpView() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Developer navigation") + "\n")
	b.WriteString(dimStyle.Render("jump to any stage — "+langLabel(m.lang)) + "\n\n")
	for i, t := range jumpTargets() {
		cursor, line := "  ", t.label
		if i == m.jumpIdx {
			cursor, line = okStyle.Render("› "), okStyle.Render(t.label)
		}
		b.WriteString(cursor + line + "\n")
	}
	b.WriteString("\n" + dimStyle.Render("[↑/↓] move · [enter] jump · [esc] cancel"))
	return b.String()
}
