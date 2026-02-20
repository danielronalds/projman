package controllers

import (
	"errors"
	"fmt"

	"github.com/danielronalds/projman/internal/ui"
)

type projectPath = string

type projectCloner interface {
	Clone(name, dir string) (projectPath, error)
}

type localProjectsIndex interface {
	IsLocalProject(project string) bool
}

type remoteProjectManager interface {
	projectLister
	projectCloner
}

type remoteConfig interface {
	ProjectDirs() []string
}

type RemoteController struct {
	remote remoteProjectManager
	local  localProjectsIndex
	fzf    selecter
	tmux   sessionLauncher
	config remoteConfig
}

func NewRemoteController(remote remoteProjectManager, local localProjectsIndex, fzf selecter, tmux sessionLauncher, config remoteConfig) RemoteController {
	return RemoteController{remote, local, fzf, tmux, config}
}

func (c RemoteController) HandleArgs(args []string) error {
	remoteProjects, err := ui.WithSpinner("fetching projects...", c.remote.ListProjects)
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

	cloneDir := c.config.ProjectDirs()[0]
	if len(c.config.ProjectDirs()) > 1 {
		cloneDir, err = c.fzf.Select(c.config.ProjectDirs())
		if err != nil {
			return errors.New("no clone dir selected")
		}
	}

	projPath, err := c.remote.Clone(proj, cloneDir)
	if err != nil {
		return fmt.Errorf("unable to get clone project: %v", err.Error())
	}

	return c.tmux.LaunchSession(proj, projPath)
}
