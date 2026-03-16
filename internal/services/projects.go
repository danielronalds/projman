package services

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
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
		if _, err := s.ListProjects(); err != nil {
			return "", err
		}
	}

	path, ok := s.localProjects[project]
	if !ok {
		return "", errors.New("project not found")
	}

	return path, nil
}

func (s ProjectsService) IsLocalProject(project string) bool {
	if len(s.localProjects) == 0 {
		if _, err := s.ListProjects(); err != nil {
			return false
		}
	}

	_, ok := s.localProjects[project]
	return ok
}

type ProjectGroup struct {
	Directory string
	Projects  []string
}

func (s ProjectsService) ListProjectsByDirectory(filter string) ([]ProjectGroup, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("getting home directory: %v", err.Error())
	}

	lowerFilter := strings.ToLower(filter)

	var groups []ProjectGroup
	for _, dir := range s.config.ProjectDirs() {
		if len(dir) == 0 {
			return nil, errors.New("invalid project directory")
		}

		displayName := strings.TrimPrefix(dir, homeDir+"/")

		if filter != "" && !strings.Contains(strings.ToLower(displayName), lowerFilter) {
			continue
		}

		contents, err := os.ReadDir(dir)
		if err != nil {
			return nil, err
		}

		var projects []string
		for _, entry := range contents {
			if entry.IsDir() {
				projects = append(projects, entry.Name())
			}
		}

		if projects == nil {
			projects = []string{}
		}

		sort.Strings(projects)

		groups = append(groups, ProjectGroup{
			Directory: displayName,
			Projects:  projects,
		})
	}

	return groups, nil
}

func getProjectNames(m map[projectName]projectPath) []string {
	keys := make([]string, 0)
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}
