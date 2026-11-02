
package main

import tea "charm.land/bubbletea/v2"

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			if m.journal != nil {
				m.journal.cancel()
			}
			return m, tea.Quit
		case "j", "down":
			if m.currentView == serviceListView && len(m.services) > 0 {
				if m.cursor < len(m.services)-1 {
					m.cursor++
				} else {
					m.cursor = 0
				}
			}
		case "k", "up":
			if m.currentView == serviceListView && len(m.services) > 0 {
				if m.cursor > 0 {
					m.cursor--
				} else if len(m.services) > 0 {
					m.cursor = len(m.services) - 1
				}
			}
		case "enter":
			if m.currentView == serviceListView && len(m.services) > 0 {
				m.currentView = serviceMenuView
			}
		case "esc":
			switch m.currentView {
			case serviceStatusView:
				m.currentView = serviceMenuView
			case serviceMenuView:
				m.currentView = serviceListView
			case journalView:
				if m.journal != nil {
					m.journal.cancel()
					m.journal = nil
				}

				m.currentView = serviceMenuView
			}
		case "s":
			if m.currentView == serviceMenuView {
				service := m.services[m.cursor].name

				m.status = "Loading service status..."
				m.statusError = nil
				m.currentView = serviceStatusView

				return m, getServiceStatus(service)
			}
		case "l":
			if m.currentView == serviceMenuView {
				service := m.services[m.cursor].name

				session, err := startJournal(service)
				if err != nil {
					m.status = err.Error()
					m.statusError = err
					m.currentView = serviceStatusView
					return m, nil
				}

				m.journal = session
				m.journalLines = nil
				m.currentView = journalView

				return m, waitForJournalLine(session)
			}
		}
	case statusMsg:
		m.status = msg.output
		m.statusError = msg.err
	case activeMsg:
		for i := range m.services {
			if m.services[i].name == msg.name {
				m.services[i].active = msg.active
				break
			}
		}
	case journalLineMsg:
		if msg.session != m.journal {
			return m, nil
		}

		if !msg.ok {
			m.journal = nil
			return m, nil
		}

		m.journalLines = append(m.journalLines, msg.line)

		const maxJournalLines = 200

		if len(m.journalLines) > maxJournalLines {
			m.journalLines = m.journalLines[len(m.journalLines)-maxJournalLines:]
		}

		return m, waitForJournalLine(msg.session)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}
