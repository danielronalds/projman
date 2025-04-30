package controllers

import (
	"errors"
	"fmt"
)

type projectPath = string

type projectCloner interface {
	Clone(name string) (projectPath, error)
}

type localProjectsIndex interface {
	IsLocalProject(project string) bool
}

type remoteProjectManager interface {
	projectLister
	projectCloner
}

type RemoteController struct {
	remote remoteProjectManager
	local  localProjectsIndex
	fzf    selecter
	tmux   sessionLauncher
}

func NewRemoteController(remote remoteProjectManager, local localProjectsIndex, fzf selecter, tmux sessionLauncher) RemoteController {
	return RemoteController{remote, local, fzf, tmux}
}

func (c RemoteController) HandleArgs(args []string) error {
	fmt.Println("fetching projects...")
	remoteProjects, err := c.remote.ListProjects()
	if err != nil {
		return fmt.Errorf("unable to fetch github projects: %v", err.Error())
	}

	filteredRemoteProjects := make([]string, 0)
	for _, project := range remoteProjects {
		if c.local.IsLocalProject(project) {
			continue
		}
		filteredRemoteProjects = append(filteredRemoteProjects, project)
	}

	proj, err := c.fzf.Select(filteredRemoteProjects)
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
