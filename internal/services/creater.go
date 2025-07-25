package services

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

type templateGetter interface {
	GetTemplateCommands(tmpl string) ([]string, error)
}

type CreaterService struct {
	templates templateGetter
}

func NewCreaterService(templates templateGetter) CreaterService {
	return CreaterService{templates}
}

func (s CreaterService) CreateProject(name, projectDir string) (projectPath, error) {
	projectPath := fmt.Sprintf("%v%v", projectDir, name)

	err := os.Mkdir(projectPath, 0775) // default permission for folders

	return projectPath, err
}

func (s CreaterService) CreateProjectWithTemplate(name, projectDir, tmpl string) (projectPath, error) {
	commands, err := s.templates.GetTemplateCommands(tmpl)
	if err != nil {
		return "", errors.New("unable to get template")
	}

	projPath, err := s.CreateProject(name, projectDir)
	if err != nil {
		return "", errors.New("project with that name already exists")
	}

	for _, tmplCmd := range commands {
		cmd := exec.Command("bash", "-c", tmplCmd)

		cmd.Dir = projPath

		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, fmt.Sprintf("PROJMAN_PROJECT_NAME=%v", name))

		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		if err := cmd.Run(); err != nil {
			return "", err
		}
	}

	return projPath, nil
}
