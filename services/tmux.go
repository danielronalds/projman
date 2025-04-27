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
	cmd := exec.Command(
		"tmux", 
		"new",
		"-c",
		dir,
		"-s",
		name,
	)

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}
