package services

import (
	"errors"
	"fmt"
	"os"
)

type projectName = string
type projectPath = string

type projectsConfig interface {
	ProjectDirs() []string
}

type ProjectsService struct {
	config        projectsConfig
	localProjects map[projectName]projectPath
}

type project struct {
	name string
	path string
}

func NewProjectsService(config projectsConfig) ProjectsService {
	return ProjectsService{
		config:        config,
		localProjects: make(map[projectName]projectPath, 0),
	}
}

func (s ProjectsService) ListProjects() ([]string, error) {
	if len(s.localProjects) != 0 {
		return getProjectNames(s.localProjects), nil
	}

	for _, dir := range s.config.ProjectDirs() {
		if len(dir) == 0 {
			return make([]string, 0), errors.New("invalid project directory")
		}

		contents, err := os.ReadDir(dir)
		if err != nil {
			return nil, err
		}

		for _, entry := range contents {
			if entry.IsDir() {
				name := entry.Name()
				path := fmt.Sprintf("%v%v", dir, name)

				_, alreadyExisting := s.localProjects[name]
				if !alreadyExisting {
					s.localProjects[name] = path
				}
			}
		}
	}

	return getProjectNames(s.localProjects), nil
}

func (s ProjectsService) GetPath(project string) (string, error) {
	if len(s.localProjects) == 0 {
		// Fetching the projects if they haven't already been fetched
		s.ListProjects()
	}

	path, ok := s.localProjects[project]
	if !ok {
		return "", errors.New("project not found")
	}

	return path, nil
}

func (s ProjectsService) IsLocalProject(project string) bool {
	if len(s.localProjects) == 0 {
		// Fetching the projects if they haven't already been fetched
		s.ListProjects()
	}

	_, ok := s.localProjects[project]
	return ok
}

func getProjectNames(m map[projectName]projectPath) []string {
	keys := make([]string, 0)
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}
