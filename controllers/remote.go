package controllers

import (
	"errors"
	"fmt"
)

type projectPath = string

type projectCloner interface {
	Clone(name string) (projectPath, error)
}

type remoteProjectManager interface {
	projectLister
	projectCloner
}

type RemoteController struct {
	remote remoteProjectManager
	fzf selecter
	tmux sessionLauncher
}


func NewRemoteController(remote remoteProjectManager, fzf selecter, tmux sessionLauncher) RemoteController { 
	return RemoteController{remote, fzf, tmux}
}

func (c RemoteController) HandleArgs(args []string) error {
	fmt.Println("fetching projects...")
	projects, err := c.remote.ListProjects()
	if err != nil {
		return fmt.Errorf("unable to fetch github projects: %v", err.Error())
	}
	
	proj, err := c.fzf.Select(projects)
	if err != nil {
		// if an error occurs, we assume fzf was Ctrl+c
		return errors.New("no project selected")
	}

	projPath, err := c.remote.Clone(proj)
	if err != nil {
		return fmt.Errorf("unable to get clone project: %v", err.Error())
	}

	return c.tmux.LaunchSession(proj, projPath)
}
