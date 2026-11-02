package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"strings"
)

var titleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#7AA2F7"))

var selectedStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#9ECE6A")).
	Bold(true)

var mutedStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#565F89"))

var panelStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	Padding(1, 2)

func (m model) viewServiceList() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("HOTH CONTROL"))
	s.WriteString("\n\n")

	for i, service := range m.services {
		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}

		name := service.name

		if i == m.cursor {
			name = selectedStyle.Render(name)
		}

		state := mutedStyle.Render("[" + service.active + "]")

		s.WriteString(cursor + name + " " + state + "\n")
	}
	s.WriteString("\n")
	s.WriteString(mutedStyle.Render("[j/k] navigate • [enter] select • [q] quit"))

	return panelStyle.
		Width(m.width).
		Height(m.height).
		Render(s.String())
}

func (m model) viewServiceMenu() string {
	var s strings.Builder

	service := m.services[m.cursor].name

	s.WriteString("HOTH CONTROL\n\n")
	s.WriteString("Service: " + service + "\n\n")

	s.WriteString("  [s] Show status\n")
	s.WriteString("  [r] Restart service\n")
	s.WriteString("  [x] Stop service\n")
	s.WriteString("  [l] View journal\n")

	s.WriteString("\n[esc] Back  [q] Quit\n")

	return s.String()
}

func (m model) viewServiceStatus() string {
	var s strings.Builder

	s.WriteString("HOTH CONTROL\n\n")
	s.WriteString("Service: " + m.services[m.cursor].name + "\n\n")

	s.WriteString(m.status)

	if m.statusError != nil {
		s.WriteString("\nCommand returned: ")
		s.WriteString(m.statusError.Error())
		s.WriteString("\n")
	}

	s.WriteString("\n[esc] Back  [q] Quit\n")

	return s.String()
}

func (m model) viewJournal() string {
	var s strings.Builder

	s.WriteString("HOTH CONTROL - JOURNAL\n\n")
	s.WriteString("Service: " + m.services[m.cursor].name + "\n\n")

	start := len(m.journalLines) - 20
	if start < 0 {
		start = 0
	}

	for _, line := range m.journalLines[start:] {
		s.WriteString(line + "\n")
	}

	s.WriteString("\n[esc] Back  [q] Quit\n")

	return s.String()
}

func (m model) View() tea.View {
	var content string

	switch m.currentView {
	case serviceListView:
		content = m.viewServiceList()
	case serviceMenuView:
		content = m.viewServiceMenu()
	case serviceStatusView:
		content = m.viewServiceStatus()
	case journalView:
		content = m.viewJournal()
	}

	v := tea.NewView(content)
	v.AltScreen = true

	return v
}
