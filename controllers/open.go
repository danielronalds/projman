package controllers

import (
	"errors"
	"fmt"
)

type projectPathFinder interface {
	GetPath(project string) (string, error)
}

type projectLister interface {
	ListProjects() ([]string, error)
}

type selecter interface {
	Select(options []string) (string, error)
}

type sessionLauncher interface {
	LaunchSession(name, dir string) error
}

type openProjectLister interface {
	projectLister
	projectPathFinder
}

type OpenController struct {
	projects openProjectLister
	fzf selecter
	tmux sessionLauncher
}

func NewOpenController(projects openProjectLister, fzf selecter, tmux sessionLauncher) OpenController { 
	return OpenController{projects, fzf, tmux}
}

func (c OpenController) HandleArgs(args []string) error {
	projects, err := c.projects.ListProjects()
	if err != nil {
		return fmt.Errorf("unable to fetch local projects: %v", err.Error())
	}
	
	proj, err := c.fzf.Select(projects)
	if err != nil {
		// if an error occurs, we assume fzf was Ctrl+c
		return errors.New("no project selected")
	}

	projPath, err := c.projects.GetPath(proj)
	if err != nil {
		return fmt.Errorf("unable to get project path: %v", err.Error())
	}

	return c.tmux.LaunchSession(proj, projPath)
}
