package tui

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"devascent/internal/ticket"
)

// Step-1 create/edit form (T1): [n] files a new ticket; [E] on a ticket edits it.
// A multi-field form — title, description, type, priority, assignee, points — with
// level-gated assignment (delegate to peers/below, escalate to those above).

var (
	formTypes  = []ticket.Type{ticket.Task, ticket.Bug, ticket.Story, ticket.Feature, ticket.TechDebt, ticket.Spike, ticket.Incident}
	formPris   = []ticket.Priority{ticket.PTrivial, ticket.PMinor, ticket.PMajor, ticket.PCritical, ticket.PBlocker}
	formPoints = []int{1, 2, 3, 5, 8}
)

const (
	fldTitle = iota
	fldDesc
	fldType
	fldPriority
	fldAssignee
	fldPoints
	fldCount
)

func idxOf[T comparable](xs []T, v T) int {
	for i, x := range xs {
		if x == v {
			return i
		}
	}
	return 0
}

func (m Model) openNewTicket() (tea.Model, tea.Cmd) {
	m.ntTitle, m.ntDesc = "", ""
	m.ntType, m.ntPri, m.ntAssignee, m.ntPoints = 0, 1, 0, 0
	m.ntFocus, m.ntEditKey = fldTitle, ""
	m.screen = screenNewTicket
	return m, nil
}

func (m Model) openEditTicket(t *ticket.Ticket) (tea.Model, tea.Cmd) {
	m.ntTitle, m.ntDesc = t.Title, t.Desc
	m.ntType = idxOf(formTypes, t.Type)
	m.ntPri = idxOf(formPris, t.Priority)
	m.ntAssignee = idxOf(assigneeOptions(), t.Assignee)
	m.ntPoints = idxOf(formPoints, t.Points)
	m.ntFocus, m.ntEditKey = fldTitle, t.Key
	m.screen = screenNewTicket
	return m, nil
}

func (m Model) handleFormKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		m.persist()
		m.quitting = true
		return m, tea.Quit
	case tea.KeyEsc:
		return m.closeForm()
	case tea.KeyEnter:
		return m.submitForm()
	case tea.KeyUp, tea.KeyShiftTab:
		if m.ntFocus > 0 {
			m.ntFocus--
		}
	case tea.KeyDown, tea.KeyTab:
		if m.ntFocus < fldCount-1 {
			m.ntFocus++
		}
	case tea.KeyLeft:
		m.cycleField(-1)
	case tea.KeyRight:
		m.cycleField(1)
	case tea.KeyBackspace, tea.KeyDelete:
		switch m.ntFocus {
		case fldTitle:
			m.ntTitle = dropLast(m.ntTitle)
		case fldDesc:
			m.ntDesc = dropLast(m.ntDesc)
		}
	case tea.KeySpace:
		m.typeInto(" ")
	case tea.KeyRunes:
		m.typeInto(string(msg.Runes))
	}
	return m, nil
}

func dropLast(s string) string {
	r := []rune(s)
	if len(r) == 0 {
		return s
	}
	return string(r[:len(r)-1])
}

func (m *Model) typeInto(s string) {
	switch m.ntFocus {
	case fldTitle:
		m.ntTitle += s
	case fldDesc:
		m.ntDesc += s
	}
}

// cycleField changes the focused select-field's value by dir (+1/-1), wrapping.
func (m *Model) cycleField(dir int) {
	wrap := func(i, n int) int { return ((i+dir)%n + n) % n }
	switch m.ntFocus {
	case fldType:
		m.ntType = wrap(m.ntType, len(formTypes))
	case fldPriority:
		m.ntPri = wrap(m.ntPri, len(formPris))
	case fldAssignee:
		m.ntAssignee = wrap(m.ntAssignee, len(assigneeOptions()))
	case fldPoints:
		m.ntPoints = wrap(m.ntPoints, len(formPoints))
	}
}

func (m Model) closeForm() (tea.Model, tea.Cmd) {
	if m.ntEditKey != "" {
		m.screen = screenTicket // editing → back to the detail
	} else {
		m.screen = screenBoard
	}
	return m, nil
}

func (m Model) submitForm() (tea.Model, tea.Cmd) {
	if strings.TrimSpace(m.ntTitle) == "" {
		m.ntFocus = fldTitle // title required
		return m, nil
	}
	assignee := assigneeOptions()[m.ntAssignee]

	if m.ntEditKey != "" { // edit an existing ticket
		t := m.boardSprint.Find(m.ntEditKey)
		if t == nil {
			m.screen = screenBoard
			return m, nil
		}
		t.Title = strings.TrimSpace(m.ntTitle)
		t.Desc = strings.TrimSpace(m.ntDesc)
		t.Type = formTypes[m.ntType]
		t.Priority = formPris[m.ntPri]
		t.Assignee = assignee
		t.Points = formPoints[m.ntPoints]
		t.DueDay = ticket.DueDayFor(t.Priority, t.AssignedDay)
		applyAssignment(t, m.playerLvl, m.boardSprint.Day)
		m.detailTicket = t
		m.persist()
		m.screen = screenTicket
		return m, nil
	}

	day := m.boardSprint.Day
	tk := &ticket.Ticket{
		Key:         nextKey(m.boardSprint, m.boardProject),
		Type:        formTypes[m.ntType],
		Title:       strings.TrimSpace(m.ntTitle),
		Desc:        strings.TrimSpace(m.ntDesc),
		Status:      ticket.ToDo,
		Priority:    formPris[m.ntPri],
		Points:      formPoints[m.ntPoints],
		Assignee:    assignee,
		Reporter:    "you",
		AssignedDay: day,
		CreatedDay:  day,
	}
	tk.DueDay = ticket.DueDayFor(tk.Priority, day)
	applyAssignment(tk, m.playerLvl, day) // delegate → teammate starts; escalate → guidance
	m.boardSprint.Tickets = append(m.boardSprint.Tickets, tk)
	m.screen = screenBoard
	m.focusTicket(tk)
	m.persist()
	return m, nil
}

// focusTicket points the board cursor at t, wherever it currently sits.
func (m *Model) focusTicket(t *ticket.Ticket) {
	for i := range ticket.BoardColumns {
		for r, c := range m.columnCards(ticket.BoardColumns[i]) {
			if c == t {
				m.boardCol, m.boardRow = i, r
				return
			}
		}
	}
	m.boardCol, m.boardRow = 0, 0
}

func (m Model) formView() string {
	heading := "File a ticket"
	if m.ntEditKey != "" {
		heading = "Edit " + m.ntEditKey
	}
	assignee := assigneeOptions()[m.ntAssignee]

	fields := [][2]string{
		{"Title", m.ntTitle},
		{"Description", m.ntDesc},
		{"Type", string(formTypes[m.ntType])},
		{"Priority", formPris[m.ntPri].Label()},
		{"Assignee", roleOf(assignee) + assignTag(assignKind(assignee, m.playerLvl))},
		{"Points", strconv.Itoa(formPoints[m.ntPoints])},
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render(heading) + "\n\n")
	for i, f := range fields {
		cursor, lbl := "  ", dimStyle.Render(fmt.Sprintf("%-12s", f[0]))
		val := f[1]
		if i == m.ntFocus {
			cursor = okStyle.Render("› ")
			lbl = okStyle.Render(fmt.Sprintf("%-12s", f[0]))
			if i == fldTitle || i == fldDesc {
				val += "▏"
			}
		}
		b.WriteString(cursor + lbl + " " + val + "\n")
	}
	b.WriteString("\n" + dimStyle.Render("[↑/↓] field · [←/→] change · type to edit · [enter] save · [esc] cancel"))
	return b.String()
}

func assignTag(kind string) string {
	switch kind {
	case "delegate":
		return tagTask.Render("  → delegate (they'll do it)")
	case "escalate":
		return tagStory.Render("  → escalate (ask for help/guidance)")
	default:
		return dimStyle.Render("  (you'll do it)")
	}
}

// nextKey returns the next sequential ticket key for the project (e.g. PXF-207).
func nextKey(sp *ticket.Sprint, proj *ticket.Project) string {
	prefix := "PXF"
	if proj != nil && proj.Key != "" {
		prefix = proj.Key
	}
	max := 0
	if sp != nil {
		for _, t := range sp.Tickets {
			if i := strings.LastIndexByte(t.Key, '-'); i >= 0 {
				if n, err := strconv.Atoi(t.Key[i+1:]); err == nil && n > max {
					max = n
				}
			}
		}
	}
	return fmt.Sprintf("%s-%d", prefix, max+1)
}
