package controllers

import (
	"fmt"
	"os"
	"path/filepath"
)

type sessionNameSanitiser interface {
	Sanitise(name string) string
}

type HereController struct {
	sessions  sessionLauncher
	sanitiser sessionNameSanitiser
}

func NewHereController(sessions sessionLauncher, sanitiser sessionNameSanitiser) HereController {
	return HereController{sessions, sanitiser}
}

func (c HereController) HandleArgs(args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("unable to get current directory: %v", err.Error())
	}

	name := c.sanitiser.Sanitise(filepath.Base(cwd))

	return c.sessions.LaunchSession(name, cwd)
}
