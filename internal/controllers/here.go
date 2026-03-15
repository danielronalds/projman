package controllers

import (
	"fmt"
	"os"
	"path/filepath"
)

type HereController struct {
	sessions sessionLauncher
}

func NewHereController(sessions sessionLauncher) HereController {
	return HereController{sessions}
}

func (c HereController) HandleArgs(args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("unable to get current directory: %v", err.Error())
	}

	name := filepath.Base(cwd)

	return c.sessions.LaunchSession(name, cwd)
}
