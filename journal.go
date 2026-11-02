package main

import (
	"bufio"
	tea "charm.land/bubbletea/v2"
	"context"
	"os/exec"
)

type journalSession struct {
	lines  chan string
	cancel context.CancelFunc
}

type journalLineMsg struct {
	session *journalSession
	line    string
	ok      bool
}

func waitForJournalLine(session *journalSession) tea.Cmd {
	return func() tea.Msg {
		line, ok := <-session.lines

		return journalLineMsg{
			session: session,
			line:    line,
			ok:      ok,
		}
	}
}

func startJournal(service string) (*journalSession, error) {
	ctx, cancel := context.WithCancel(context.Background())

	cmd := exec.CommandContext(
		ctx,
		"journalctl",
		"--no-pager",
		"-u",
		service,
		"-n",
		"20",
		"-f",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		cancel()
		return nil, err
	}

	session := &journalSession{
		lines:  make(chan string, 100),
		cancel: cancel,
	}

	go func() {
		defer close(session.lines)
		defer cmd.Wait()

		scanner := bufio.NewScanner(stdout)

		for scanner.Scan() {
			select {
			case session.lines <- scanner.Text():
			case <-ctx.Done():
				return
			}
		}
	}()

	return session, nil
}
