package controllers

import "fmt"

type projectLister interface {
	ListLocal() ([]string, error)
	GetPath(project string) (string, error)
}

type selecter interface {
	Select(options []string) (string, error)
}

type sessionLauncher interface {
	LaunchSession(name, dir string) error
}

type OpenController struct {
	projects projectLister
	fzf selecter
	tmux sessionLauncher
}

func NewOpenController(projects projectLister, fzf selecter, tmux sessionLauncher) OpenController { 
	return OpenController{projects, fzf, tmux}
}

func (c OpenController) HandleArgs(args []string) error {
	projects, err := c.projects.ListLocal()
	if err != nil {
		return fmt.Errorf("unable to fetch local projects: %v", err.Error())
	}
	
	proj, err := c.fzf.Select(projects)
	if err != nil {
		// if an error occurs, we assuke fzf was Ctrl+c
		return fmt.Errorf("no project selected", err.Error())
	}

	projPath, err := c.projects.GetPath(proj)
	if err != nil {
		return fmt.Errorf("unable to get project path: %v", err.Error())
	}

	return c.tmux.LaunchSession(proj, projPath)
}
