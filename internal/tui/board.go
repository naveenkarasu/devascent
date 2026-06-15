package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"devascent/internal/ticket"
)

// Step-1 board rendering + navigation (#58 read-only, #59 nav). Colors are the
// LOCKED palette (vault: TUI Board doc): they reuse the bench style vars
// (titleStyle 63, okStyle 42, errStyle 203, codeStyle 117, dimStyle) so the board
// reads as the same app; only the type-tag hues + focus accent are board-local.

const (
	cardW        = 17 // card box content width
	cardInner    = 13 // chars of title that fit on one line
	boardWideMin = 84 // below this terminal width, collapse to one column + tab strip
)

var (
	bdIdle    = lipgloss.Color("238") // idle card border
	bdAccent  = lipgloss.Color("75")  // selected/focused accent (blue)
	cardBox   = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(bdIdle).Width(cardW).Padding(0, 1)
	cardSel   = lipgloss.NewStyle().Border(lipgloss.ThickBorder()).BorderForeground(bdAccent).Width(cardW).Padding(0, 1)
	tagStory  = lipgloss.NewStyle().Foreground(lipgloss.Color("176")) // magenta
	tagTask   = lipgloss.NewStyle().Foreground(lipgloss.Color("75"))  // blue
	headFocus = lipgloss.NewStyle().Foreground(bdAccent).Bold(true)
)

// enterBoard opens the board on the current sprint (seeding a demo sprint until
// real content lands in #62), remembering the screen to return to on esc.
func (m Model) enterBoard() (tea.Model, tea.Cmd) {
	if m.boardSprint == nil {
		m.boardProject, m.boardSprint = seedSprint1()
	}
	m.boardReturn = m.screen
	m.boardCol = defaultFocusCol(m.boardSprint, m.boardSprint.Day)
	m.boardRow = 0
	m.boardHelp = false
	m.screen = screenBoard
	return m, nil
}

// enterStep1 opens the board as the REAL save-backed career home (from the
// graduation screen). It seeds the first sprint if the player doesn't have one
// yet, marks the board as persistent, and saves so the sprint sticks.
func (m Model) enterStep1() (tea.Model, tea.Cmd) {
	if m.boardSprint == nil {
		m.boardProject, m.boardSprint = seedSprint1()
	}
	m.step1Home = true
	m.boardReturn = screenStep0Complete
	m.boardCol = defaultFocusCol(m.boardSprint, m.boardSprint.Day)
	m.boardRow = 0
	m.boardHelp = false
	m.screen = screenBoard
	m.persist()
	return m, nil
}

func (m Model) handleBoardKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.boardHelp { // any of these closes the help overlay
		switch msg.String() {
		case "?", "esc", "enter":
			m.boardHelp = false
		}
		return m, nil
	}
	cols := ticket.BoardColumns
	switch msg.String() {
	case "esc": // layered: close an overlay first, else leave the board
		switch {
		case m.boardAnalytic:
			m.boardAnalytic = false
		case m.boardBacklog:
			m.boardBacklog = false
		case m.boardGroup != 0:
			m.boardGroup = 0
		default:
			m.screen = m.boardReturn
		}
	case "?":
		m.boardHelp = true
	case "f": // quick filter
		m.boardFilter = (m.boardFilter + 1) % len(boardFilters)
		m.boardRow = 0
	case "g": // swimlane grouping
		m.boardGroup = (m.boardGroup + 1) % len(boardGroupings)
	case "b": // backlog view
		m.boardBacklog = !m.boardBacklog
	case "v": // sprint analytics
		m.boardAnalytic = !m.boardAnalytic
	case "left", "a":
		if m.boardCol > 0 {
			m.boardCol--
			m.boardRow = 0
		}
	case "right", "d":
		if m.boardCol < len(cols)-1 {
			m.boardCol++
			m.boardRow = 0
		}
	case "up", "w":
		if m.boardRow > 0 {
			m.boardRow--
		}
	case "down", "s":
		if n := len(m.columnCards(cols[m.boardCol])); m.boardRow < n-1 {
			m.boardRow++
		}
	case "e":
		return m.enterStandup()
	case "n":
		return m.openNewTicket()
	case "enter":
		if t := m.selectedTicket(); t != nil {
			return m.openTicket(t)
		}
	}
	return m, nil
}

// selectedTicket is the card under the cursor (nil if the focused column is empty).
func (m Model) selectedTicket() *ticket.Ticket {
	cards := m.columnCards(ticket.BoardColumns[m.boardCol])
	if m.boardRow >= 0 && m.boardRow < len(cards) {
		return cards[m.boardRow]
	}
	return nil
}

func (m Model) boardView() string {
	if m.boardSprint == nil {
		return boxStyle.Render("No sprint loaded.\n\n" + dimStyle.Render("[esc] back"))
	}
	if m.boardHelp {
		return m.boardHelpView()
	}
	if m.boardAnalytic {
		return m.boardAnalyticsView()
	}
	if m.boardBacklog {
		return m.boardBacklogView()
	}
	if m.boardGroup != 0 {
		return m.boardGroupedView()
	}
	if m.width > 0 && m.width < boardWideMin {
		return m.boardNarrowView()
	}
	return m.boardWideView()
}

// boardWideView: all columns side-by-side (the default).
func (m Model) boardWideView() string {
	cols := make([]string, 0, len(ticket.BoardColumns)*2)
	for i := range ticket.BoardColumns {
		if i > 0 {
			cols = append(cols, "  ") // gutter
		}
		cols = append(cols, m.boardColumn(i))
	}
	board := lipgloss.JoinHorizontal(lipgloss.Top, cols...)
	foot := dimStyle.Render("wasd move · enter open · n new · e end-day · f/g/b/v views · ? help · esc back")
	return m.boardHeader() + "\n\n" + board + "\n\n" + foot
}

// boardNarrowView: one focused column + a tab strip (the gh-dash degrade).
func (m Model) boardNarrowView() string {
	var tabs []string
	for i, c := range ticket.BoardColumns {
		n := len(m.boardSprint.ColumnVisible(c, m.boardSprint.Day))
		label := fmt.Sprintf("%s(%d)", shortCol(c), n)
		if i == m.boardCol {
			tabs = append(tabs, headFocus.Render("‹"+label+"›"))
		} else {
			tabs = append(tabs, dimStyle.Render(label))
		}
	}
	strip := strings.Join(tabs, dimStyle.Render(" · "))
	foot := dimStyle.Render("[a/d][w/s] move · enter open · n new · e end-day · ? help · esc back")
	return m.boardHeader() + "\n\n" + strip + "\n\n" + m.boardColumn(m.boardCol) + "\n\n" + foot
}

func (m Model) boardHeader() string {
	sp := m.boardSprint
	name := "Board"
	if m.boardProject != nil {
		name = m.boardProject.Name
	}
	h := titleStyle.Render(name) + dimStyle.Render(fmt.Sprintf(
		"   Sprint %d · Day %d · cap %d · committed %d · done %dpt",
		sp.Number, sp.Day, sp.Capacity, sp.Committed(), sp.DonePoints()))
	if od := sp.OverdueCount(sp.Day); od > 0 {
		h += "   " + errStyle.Render(fmt.Sprintf("⚠ %d overdue", od))
	}
	if inc := len(sp.Incoming(sp.Day)); inc > 0 {
		h += "\n" + dimStyle.Render(fmt.Sprintf("%d more ticket(s) arrive on later days — end the day to pull them in", inc))
	}
	if m.boardFilter != 0 {
		h += "\n" + okStyle.Render("filter: "+boardFilters[m.boardFilter].name) + dimStyle.Render("  ([f] cycle)")
	}
	return h
}

func (m Model) boardHelpView() string {
	rows := [][2]string{
		{"a / d   ← / →", "move between columns"},
		{"w / s   ↑ / ↓", "move between cards"},
		{"enter", "open the selected ticket"},
		{"n", "file a new ticket"},
		{"e", "end the day (standup → next day)"},
		{"f", "cycle quick filter (mine/overdue/high-pri)"},
		{"g", "cycle swimlane grouping (epic/priority/assignee)"},
		{"b", "backlog view"},
		{"v", "sprint analytics (burndown/throughput)"},
		{"esc", "leave the board"},
		{"?", "close this help"},
	}
	var b strings.Builder
	b.WriteString(titleStyle.Render("Board — keys") + "\n\n")
	for _, r := range rows {
		b.WriteString("  " + okStyle.Render(fmt.Sprintf("%-14s", r[0])) + dimStyle.Render(r[1]) + "\n")
	}
	b.WriteString("\n" + dimStyle.Render("[?] or [esc] close"))
	return b.String()
}

func (m Model) boardColumn(idx int) string {
	col := ticket.BoardColumns[idx]
	tickets := m.columnCards(col)
	focused := idx == m.boardCol

	count := fmt.Sprintf("%d", len(tickets))
	over := false
	if lim := wipLimit(col); lim > 0 {
		count = fmt.Sprintf("%d/%d", len(tickets), lim) // B3: WIP limit
		over = len(tickets) > lim
	}
	label := fmt.Sprintf(" %s (%s)", col.Label(), count)
	if focused {
		label += "  ◄"
	}
	if over {
		label += " ⚠"
	}
	var head string
	switch {
	case focused:
		head = headFocus.Render(label)
	case over:
		head = errStyle.Render(label)
	default:
		head = dimStyle.Render(label)
	}

	parts := []string{head, ""}
	for r, t := range tickets {
		parts = append(parts, boardCard(t, focused && r == m.boardRow, m.boardSprint.Day))
	}
	if len(tickets) == 0 {
		parts = append(parts, dimStyle.Render("  —"))
	}
	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

// boardCard renders one ticket card. The selected card (the cursor) gets the
// thick accent border; every other card is a plain border.
func boardCard(t *ticket.Ticket, selected bool, day int) string {
	box := cardBox
	if selected {
		box = cardSel
	}
	body := codeStyle.Render(t.Key) + "  " + typeTag(t) + "\n" +
		trunc(t.Title, cardInner) + "\n" +
		cardMeta(t)
	if d := dueLine(t, day); d != "" {
		body += "\n" + d
	}
	return box.Render(body)
}

// dueLine is the card's SLA line: faint "due Dn", red when due within a day or
// already overdue. Empty for Done or no-deadline tickets.
func dueLine(t *ticket.Ticket, day int) string {
	if t.Status == ticket.Done || t.DueDay == 0 {
		return ""
	}
	switch {
	case t.Overdue(day):
		return errStyle.Render("⚠ OVERDUE")
	case t.DueDay-day <= 1:
		return errStyle.Render(fmt.Sprintf("due D%d", t.DueDay))
	default:
		return dimStyle.Render(fmt.Sprintf("due D%d", t.DueDay))
	}
}

// typeTag is the colored type chip; Done cards show a green ✓ instead. The tag is
// a short form so it fits the card width (e.g. tech-debt → "debt").
func typeTag(t *ticket.Ticket) string {
	if t.Status == ticket.Done {
		return okStyle.Render("✓")
	}
	tag := typeShort(t.Type)
	switch t.Type {
	case ticket.Bug, ticket.Incident:
		return errStyle.Render(tag)
	case ticket.Story, ticket.Feature:
		return tagStory.Render(tag)
	case ticket.Task, ticket.TechDebt, ticket.Spike:
		return tagTask.Render(tag)
	}
	return dimStyle.Render(tag)
}

// typeShort is the abbreviated type label used on cards (keeps them one line).
func typeShort(t ticket.Type) string {
	switch t {
	case ticket.TechDebt:
		return "debt"
	case ticket.Incident:
		return "inc"
	case ticket.Feature:
		return "feat"
	default: // bug, story, task, spike
		return string(t)
	}
}

// cardMeta is the bottom line: priority (Major=red, else faint) · points · @you.
func cardMeta(t *ticket.Ticket) string {
	if t.Status == ticket.Done {
		if t.Points > 0 {
			return dimStyle.Render(fmt.Sprintf("%dpt", t.Points))
		}
		return ""
	}
	var parts []string
	if t.Priority != "" {
		if t.Priority.IsHigh() {
			parts = append(parts, errStyle.Render(priShort(t.Priority)))
		} else {
			parts = append(parts, dimStyle.Render(priShort(t.Priority)))
		}
	}
	if t.Points > 0 {
		parts = append(parts, fmt.Sprintf("%dpt", t.Points))
	}
	meta := strings.Join(parts, " ")
	if t.Assignee == "you" {
		if meta != "" {
			meta += " "
		}
		meta += codeStyle.Render("@you")
	}
	return meta
}

// priShort keeps the 7-char priorities from overflowing the card width (the
// detail view still shows the full label).
func priShort(p ticket.Priority) string {
	switch p {
	case ticket.PBlocker:
		return "Block"
	case ticket.PCritical:
		return "Crit"
	case ticket.PTrivial:
		return "Triv"
	default: // Major, Minor
		return p.Label()
	}
}

// shortCol is the abbreviated column name for the narrow tab strip.
func shortCol(s ticket.Status) string {
	switch s {
	case ticket.ToDo:
		return "To Do"
	case ticket.InProgress:
		return "In Prog"
	case ticket.InReview:
		return "Review"
	case ticket.Done:
		return "Done"
	}
	return s.Label()
}

// defaultFocusCol focuses In Progress if it has cards, else the first non-empty
// column, else the first column.
func defaultFocusCol(sp *ticket.Sprint, day int) int {
	for i, c := range ticket.BoardColumns {
		if c == ticket.InProgress && len(sp.ColumnVisible(c, day)) > 0 {
			return i
		}
	}
	for i, c := range ticket.BoardColumns {
		if len(sp.ColumnVisible(c, day)) > 0 {
			return i
		}
	}
	return 0
}

// trunc shortens s to at most n runes, adding an ellipsis when cut.
func trunc(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	if n <= 1 {
		return string(r[:n])
	}
	return string(r[:n-1]) + "…"
}

// demoBoard is a placeholder sprint so the board renders before real content
// exists (#62 seeds the real Sprint 1; #63 wires the bench→board entry).
func demoBoard() (*ticket.Project, *ticket.Sprint) {
	proj := &ticket.Project{Key: "PXF", Name: "Pixel Forge"}
	sp := &ticket.Sprint{
		Number: 3, Day: 4, Capacity: 13, Goal: "Checkout hardening",
		Tickets: []*ticket.Ticket{
			{Key: "PXF-104", Type: ticket.Bug, Title: "Email validator rejects + tags", Status: ticket.ToDo, Priority: ticket.PMajor, Points: 3, Assignee: "you"},
			{Key: "PXF-105", Type: ticket.Task, Title: "Bump lodash", Status: ticket.ToDo, Priority: ticket.PMinor, Points: 1},
			{Key: "PXF-101", Type: ticket.Bug, Title: "Off-by-one in the paginator", Status: ticket.InProgress, Priority: ticket.PMajor, Points: 2, Assignee: "you",
				Desc: "The results list drops the first row of every page. The page math is off by one — page 1 should start at index 0, not at `size`. Fix the slice without breaking the last (partial) page.",
				Acceptance: []ticket.Criterion{
					{Text: "page 1 returns items 1..size"},
					{Text: "no rows overlap or go missing between pages"},
					{Text: "the existing tests stay green"},
				},
				Comments: []ticket.Comment{
					{Author: "Sam (manager)", Body: "Grab it once you're set up — ping me if you're stuck, but try first."},
				},
				Learn: "def paginate(items, page, size):\n    start = (page - 1) * size   # 1-indexed pages\n    return items[start:start + size]",
				Grading: &ticket.Grading{
					Lang:        "python",
					Command:     []string{"python", "check.py"},
					EditPaths:   []string{"paginate.py"},
					StartFiles:  map[string]string{"paginate.py": "def paginate(items, page, size):\n    start = page * size  # BUG: page 1 should start at index 0\n    return items[start:start + size]\n"},
					HiddenFiles: map[string]string{"check.py": "from paginate import paginate\n\nitems = list(range(1, 11))\nassert paginate(items, 1, 3) == [1, 2, 3], paginate(items, 1, 3)\nassert paginate(items, 2, 3) == [4, 5, 6], paginate(items, 2, 3)\nassert paginate(items, 4, 3) == [10], paginate(items, 4, 3)\nprint('OK')\n"},
				},
			},
			{Key: "PXF-099", Type: ticket.Story, Title: "Export data to CSV", Status: ticket.InReview, Priority: ticket.PMinor, Points: 2},
			{Key: "PXF-097", Type: ticket.Story, Title: "First commit", Status: ticket.Done, Points: 2},
			{Key: "PXF-095", Type: ticket.Task, Title: "README fix", Status: ticket.Done, Points: 2},
			{Key: "PXF-093", Type: ticket.Task, Title: "Unit test scaffold", Status: ticket.Done, Points: 1},
		},
	}
	return proj, sp
}
