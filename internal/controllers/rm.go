package controllers

import (
	"errors"
	"fmt"
	"os"
)

type uncommittedChecker interface {
	HasUncommittedChanges(dir string) bool
}

type RmController struct {
	projects openProjectLister
	fzf      selecter
	git      uncommittedChecker
}

func NewRmController(projects openProjectLister, fzf selecter, git uncommittedChecker) RmController {
	return RmController{projects, fzf, git}
}

func (c RmController) HandleArgs(args []string) error {
	projects, err := c.projects.ListProjects()
	if err != nil {
		return fmt.Errorf("unable to fetch local projects: %v", err.Error())
	}

	proj, err := c.fzf.Select(projects)
	if err != nil {
		return errors.New("no project selected")
	}

	projPath, err := c.projects.GetPath(proj)
	if err != nil {
		return fmt.Errorf("unable to get project path: %v", err.Error())
	}

	skipGitCheck := len(args) > 1 && args[1] == "--without-git-check"
	if !skipGitCheck && c.git.HasUncommittedChanges(projPath) {
		return errors.New("uncommitted git changes detected, rerun with --without-git-check to proceed")
	}

	if err := os.RemoveAll(projPath); err != nil {
		return err
	}

	fmt.Printf("Removed %v\n", proj)
	return nil
}
