package main

import tea "charm.land/bubbletea/v2"

type service struct {
	name   string
	active string
}

type model struct {
	services    []service
	cursor      int
	currentView viewState

	status      string
	statusError error

	journal      *journalSession
	journalLines []string

	width int
	height int
}

type viewState int

const (
	serviceListView viewState = iota
	serviceMenuView
	serviceStatusView
	journalView
)

func (m model) Init() tea.Cmd {
	var cmds []tea.Cmd

	for _, service := range m.services {
		cmds = append(cmds, getActiveState(service.name))
	}

	return tea.Batch(cmds...)
}

func initialModel() model {
	return model{
		services: []service{
			{name: "sshd.service", active: "unknown"},
			{name: "hothctl-test.service", active: "unknown"},
		},
		cursor:      0,
		currentView: serviceListView,
	}
}

