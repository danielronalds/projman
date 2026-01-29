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
	projectName, skipGitCheck := c.parseArgs(args)

	proj, err := c.selectProject(projectName)
	if err != nil {
		return err
	}

	projPath, err := c.projects.GetPath(proj)
	if err != nil {
		return fmt.Errorf("unable to get project path: %v", err.Error())
	}

	if !skipGitCheck && c.git.HasUncommittedChanges(projPath) {
		return errors.New("uncommitted git changes detected, rerun with --without-git-check to proceed")
	}

	if err := os.RemoveAll(projPath); err != nil {
		return err
	}

	fmt.Printf("Removed %v\n", proj)
	return nil
}

func (c RmController) parseArgs(args []string) (projectName string, skipGitCheck bool) {
	for _, arg := range args[1:] {
		if arg == "--without-git-check" {
			skipGitCheck = true
		} else {
			projectName = arg
		}
	}
	return
}

func (c RmController) selectProject(projectName string) (string, error) {
	if projectName != "" {
		return projectName, nil
	}

	projects, err := c.projects.ListProjects()
	if err != nil {
		return "", fmt.Errorf("unable to fetch local projects: %v", err.Error())
	}

	proj, err := c.fzf.Select(projects)
	if err != nil {
		return "", errors.New("no project selected")
	}

	return proj, nil
}
