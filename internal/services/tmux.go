package services

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type tmuxConfig interface {
	SessionWindows() []string
	StartingWindow() int
}

type TmuxService struct {
	config tmuxConfig
}

func NewTmuxService(config tmuxConfig) TmuxService {
	return TmuxService{config}
}

func (s TmuxService) ListActiveSessions() ([]string, error) {
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

func (s TmuxService) LaunchSession(name, dir string) error {
	// If an error occurs, its likely the session already existed
	_ = s.createSession(name, dir)

	return s.OpenActiveSession(name)
}

func (s TmuxService) OpenActiveSession(name string) error {
	var cmd *exec.Cmd
	if s.isInTmuxSession() {
		cmd = s.switchToSession(name)
	} else {
		cmd = s.attachToSession(name)
	}

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func (s TmuxService) createSession(name, dir string) error {
	if len(s.config.SessionWindows()) == 0 {
		return errors.New("config needs to define at least one window")
	}

	windows := s.config.SessionWindows()

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

	startingWindow := s.config.StartingWindow()

	cmd = exec.Command("tmux", "select-window", "-t", fmt.Sprintf("%s:%v", name, startingWindow))
	return cmd.Run()
}

func (s TmuxService) isInTmuxSession() bool {
	_, ok := os.LookupEnv("TMUX")
	return ok
}

func (s TmuxService) attachToSession(name string) *exec.Cmd {
	return exec.Command("tmux", "a", "-t", name)
}

func (s TmuxService) switchToSession(name string) *exec.Cmd {
	return exec.Command("tmux", "switch", "-t", name)
}
