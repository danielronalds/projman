package services

import (
	"os"
	"os/exec"
)

type TmuxService struct{}

func NewTmuxService() TmuxService {
	return TmuxService{}
}

func (s TmuxService) LaunchSession(name, dir string) error {
	// If an error occurs, its likely the session already existed
	_ = s.createSession(name, dir)

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
	cmd := exec.Command(
		"tmux", 
		"new",
		"-c",
		dir,
		"-s",
		name,
		"-d",
	)

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
