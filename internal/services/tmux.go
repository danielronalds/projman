package services

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type TmuxProvider struct {
	config TmuxConfig
}

func NewTmuxProvider(config TmuxConfig) TmuxProvider {
	return TmuxProvider{config}
}

func (p TmuxProvider) Name() string {
	return "tmux"
}

func (p TmuxProvider) ListActiveSessions() ([]string, error) {
	cmd := exec.Command("tmux", "list-sessions", "-F", "#{session_name}")
	out, err := cmd.Output()
	if err != nil {
		return make([]string, 0), err
	}

	sessions := make([]string, 0)

	for session := range strings.SplitSeq(string(out), "\n") {
		if strings.TrimSpace(session) != "" {
			sessions = append(sessions, session)
		}
	}

	return sessions, nil
}

func (p TmuxProvider) LaunchSession(name, dir string) error {
	_ = p.createSession(name, dir)

	return p.OpenActiveSession(name)
}

func (p TmuxProvider) OpenActiveSession(name string) error {
	var cmd *exec.Cmd
	if p.isInTmuxSession() {
		cmd = p.switchToSession(name)
	} else {
		cmd = p.attachToSession(name)
	}

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func (p TmuxProvider) createSession(name, dir string) error {
	if len(p.config.Windows) == 0 {
		return errors.New("config needs to define at least one window")
	}

	windows := p.config.Windows

	cmd := exec.Command("tmux", "new", "-c", dir, "-s", name, "-n", windows[0], "-d")
	if err := cmd.Run(); err != nil {
		return err
	}

	for _, window := range windows[1:] {
		cmd = exec.Command("tmux", "new-window", "-c", dir, "-t", name, "-n", window)
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	startingWindow := p.config.StartingWindow

	cmd = exec.Command("tmux", "select-window", "-t", fmt.Sprintf("%s:%v", name, startingWindow))
	return cmd.Run()
}

func (p TmuxProvider) isInTmuxSession() bool {
	_, ok := os.LookupEnv("TMUX")
	return ok
}

func (p TmuxProvider) attachToSession(name string) *exec.Cmd {
	return exec.Command("tmux", "a", "-t", name)
}

func (p TmuxProvider) switchToSession(name string) *exec.Cmd {
	return exec.Command("tmux", "switch", "-t", name)
}
