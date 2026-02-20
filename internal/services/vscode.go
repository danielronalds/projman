package services

import (
	"fmt"
	"os"
	"os/exec"
)

type VSCodeProvider struct{}

func NewVSCodeProvider() VSCodeProvider {
	return VSCodeProvider{}
}

func (p VSCodeProvider) Name() string {
	return "vscode"
}

func (p VSCodeProvider) ListActiveSessions() ([]string, error) {
	return nil, fmt.Errorf("vscode provider does not support session management")
}

func (p VSCodeProvider) LaunchSession(name, dir string) error {
	cmd := exec.Command("code", dir)

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func (p VSCodeProvider) OpenActiveSession(name string) error {
	return fmt.Errorf("vscode provider does not support session management")
}
