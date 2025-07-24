package services

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type TmuxService struct{}

func NewTmuxService() TmuxService {
	return TmuxService{}
}

func (s TmuxService) ListActiveSessions() ([]string, error) {
	cmd := exec.Command("tmux", "list-sessions", "-F", "#{session_name}")
	out, err := cmd.Output()
	if err != nil {
		return make([]string, 0), err
	}

	sessions := make([]string, 0)

	for _, session := range strings.Split(string(out), "\n") {
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
	cmd := exec.Command("tmux", "new", "-c", dir, "-s", fmt.Sprintf("%s:1", name), "-n", "CLI", "-d")
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("tmux", "new-window", "-c", dir, "-t", fmt.Sprintf("%s:2", name), "-n", "Code")
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("tmux", "new-window", "-c", dir, "-t", fmt.Sprintf("%s:3", name), "-n", "Server")
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("tmux", "select-window", "-t", fmt.Sprintf("%s:1", name))
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
