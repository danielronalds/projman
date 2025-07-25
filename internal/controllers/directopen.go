package controllers

import (
	"errors"
	"fmt"
)

type DirectOpenController struct {
	projects openProjectLister
	fzf      selecter
	tmux     sessionLauncher
}

func NewDirectOpenController(projects openProjectLister, fzf selecter, tmux sessionLauncher) DirectOpenController {
	return DirectOpenController{projects, fzf, tmux}
}

func (c DirectOpenController) HandleArgs(args []string) error {
	if len(args) == 0 {
		return errors.New("no project name passed")
	}

	proj := args[0]

	projPath, err := c.projects.GetPath(proj)
	if err != nil {
		return fmt.Errorf("unable to get project path: %v", err.Error())
	}

	return c.tmux.LaunchSession(proj, projPath)
}
