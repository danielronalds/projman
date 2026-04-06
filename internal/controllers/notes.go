package controllers

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

type notesProjectResolver interface {
	ResolveCurrentProject() (string, string, error)
}

type notesPathProvider interface {
	NotePath(projPath string) (string, error)
}

type notesProjectLister interface {
	ListProjects() ([]string, error)
	GetPath(project string) (string, error)
}

type notesSelecter interface {
	Select(options []string) (string, error)
}

type NotesController struct {
	resolver notesProjectResolver
	notes    notesPathProvider
	projects notesProjectLister
	fzf      notesSelecter
}

func NewNotesController(resolver notesProjectResolver, notes notesPathProvider, projects notesProjectLister, fzf notesSelecter) NotesController {
	return NotesController{resolver, notes, projects, fzf}
}

func (c NotesController) HandleArgs(args []string) error {
	projPath, err := c.resolveProjectPath()
	if err != nil {
		return err
	}

	notePath, err := c.notes.NotePath(projPath)
	if err != nil {
		return fmt.Errorf("resolving note path: %v", err.Error())
	}

	editor := getEditor()
	cmd := exec.Command(editor, notePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (c NotesController) resolveProjectPath() (string, error) {
	_, projPath, err := c.resolver.ResolveCurrentProject()
	if err == nil {
		return projPath, nil
	}

	projects, err := c.projects.ListProjects()
	if err != nil {
		return "", fmt.Errorf("unable to fetch local projects: %v", err.Error())
	}

	if len(projects) == 0 {
		return "", errors.New("no local projects found")
	}

	proj, err := c.fzf.Select(projects)
	if err != nil {
		return "", fmt.Errorf("selecting project: %v", err.Error())
	}

	projPath, err = c.projects.GetPath(proj)
	if err != nil {
		return "", fmt.Errorf("unable to get project path: %v", err.Error())
	}

	return projPath, nil
}
