package controllers

import (
	"errors"
	"fmt"
	"strings"
)

type newConfig interface {
	ProjectDirs() []string
	OpenNewProjects() bool
	TemplateNames() []string
}

type projectCreator interface {
	CreateProject(name, projectDir string) (projectPath, error)
	CreateProjectWithTemplate(name, projectDir, tmpl string) (projectPath, error)
}

type NewController struct {
	fzf     selecter
	creater projectCreator
	tmux    sessionLauncher
	config  newConfig
}

func NewNewController(fzf selecter, creater projectCreator, tmux sessionLauncher, config newConfig) NewController {
	return NewController{fzf, creater, tmux, config}
}

func (c NewController) HandleArgs(args []string) error {
	if len(args) < 2 {
		return errors.New("no project name supplied")
	}

	projectName := strings.TrimSuffix(strings.TrimSpace(args[1]), "/")

	projectDir := c.config.ProjectDirs()[0]
	if len(c.config.ProjectDirs()) > 1 {
		var err error
		projectDir, err = c.fzf.Select(c.config.ProjectDirs())
		if err != nil {
			return errors.New("no project dir selected")
		}
	}

	projPath, err := c.handleCreateProject(projectName, projectDir)
	if err != nil {
		return fmt.Errorf("unable to create project: %v", err.Error())
	}

	if c.config.OpenNewProjects() {
		return c.tmux.LaunchSession(projectName, projPath)
	}

	return nil
}

func (c NewController) handleCreateProject(projectName, projectDir string) (string, error) {
	if len(c.config.TemplateNames()) == 0 {
		return c.creater.CreateProject(projectName, projectDir)
	}

	noTemplate := "No template"

	templates := append(c.config.TemplateNames(), noTemplate)

	tmpl, err := c.fzf.Select(templates)
	if err != nil {
		return "", err
	}

	if tmpl == noTemplate {
		return c.creater.CreateProject(projectName, projectDir)
	}

	return c.creater.CreateProjectWithTemplate(projectName, projectDir, tmpl)
}
