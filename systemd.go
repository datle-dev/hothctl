package main

import (
	tea "charm.land/bubbletea/v2"
	"os/exec"
	"strings"
)

type statusMsg struct {
	output string
	err    error
}

type activeMsg struct {
	name   string
	active string
	err    error
}

func getServiceStatus(service string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command(
			"systemctl",
			"--no-pager",
			"--full",
			"status",
			service,
		)

		output, err := cmd.CombinedOutput()

		return statusMsg{
			output: string(output),
			err:    err,
		}
	}
}

func getActiveState(name string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command(
			"systemctl",
			"is-active",
			name,
		)

		output, err := cmd.CombinedOutput()

		return activeMsg{
			name:   name,
			active: strings.TrimSpace(string(output)),
			err:    err,
		}
	}
}
