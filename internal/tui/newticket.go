package tui

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"devascent/internal/ticket"
)

// Step-1 "file a ticket" flow (#64): [n] on the board opens a small form (modal
// title input + Tab-cycled type) that creates a new To Do ticket assigned to you.
// This is the ticketing mechanic; structured-grading bug-report scenarios layer on
// top later.

// newTicketTypes are the types offered in the form ([tab] cycles through them).
var newTicketTypes = []ticket.Type{ticket.Task, ticket.Bug, ticket.Story, ticket.TechDebt}

// openNewTicket opens the file-a-ticket form over the board (title as modal input).
func (m Model) openNewTicket() (tea.Model, tea.Cmd) {
	m.screen = screenNewTicket
	m.inputActive = true
	m.input = ""
	m.ntType = 0
	return m, nil
}

// createTicket files a new To Do ticket (assigned to you) and drops the cursor on it.
func (m Model) createTicket(title string) (tea.Model, tea.Cmd) {
	day := m.boardSprint.Day
	tk := &ticket.Ticket{
		Key:         nextKey(m.boardSprint, m.boardProject),
		Type:        newTicketTypes[m.ntType],
		Title:       title,
		Status:      ticket.ToDo,
		Priority:    ticket.PMinor,
		Assignee:    "you",
		AssignedDay: day, // filed today → visible now
		CreatedDay:  day,
	}
	tk.DueDay = ticket.DueDayFor(tk.Priority, day)
	m.boardSprint.Tickets = append(m.boardSprint.Tickets, tk)
	m.inputActive = false
	m.input = ""
	m.screen = screenBoard
	m.boardCol = 0 // To Do
	m.boardRow = len(m.boardSprint.ColumnVisible(ticket.ToDo, day)) - 1
	m.persist()
	return m, nil
}

func (m Model) newTicketView() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("File a ticket") + "\n\n")
	b.WriteString(dimStyle.Render("Type ") + coloredType(newTicketTypes[m.ntType]) + dimStyle.Render("   ([tab] to change)") + "\n\n")
	b.WriteString(dimStyle.Render("Title") + "\n  " + m.input + "▏\n\n")
	b.WriteString(dimStyle.Render("[enter] create   ·   [tab] type   ·   [esc] cancel"))
	return b.String()
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

// coloredType is a colored short type chip for a bare type (no ticket).
func coloredType(ty ticket.Type) string {
	tag := typeShort(ty)
	switch ty {
	case ticket.Bug, ticket.Incident:
		return errStyle.Render(tag)
	case ticket.Story, ticket.Feature:
		return tagStory.Render(tag)
	default:
		return tagTask.Render(tag)
	}
}
