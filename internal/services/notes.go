package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type notesConfig interface {
	ProjectDirs() []string
}

type notesRepository interface {
	NotesDir() (string, error)
}

type NotesService struct {
	config    notesConfig
	notesRepo notesRepository
}

func NewNotesService(config notesConfig, notesRepo notesRepository) NotesService {
	return NotesService{
		config:    config,
		notesRepo: notesRepo,
	}
}

func (s NotesService) NoteFilename(projPath string) string {
	parent := strings.ToLower(filepath.Base(filepath.Dir(projPath)))
	project := filepath.Base(projPath)
	return fmt.Sprintf("%v_%v.md", parent, project)
}

func (s NotesService) ResolveCurrentProject() (string, string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", "", fmt.Errorf("getting current directory: %v", err.Error())
	}

	for _, dir := range s.config.ProjectDirs() {
		if !strings.HasPrefix(cwd, dir) {
			continue
		}

		relative := strings.TrimPrefix(cwd, dir)
		if relative == "" {
			continue
		}

		projectName := strings.SplitN(relative, string(filepath.Separator), 2)[0]
		projectPath := filepath.Join(dir, projectName)
		return projectName, projectPath, nil
	}

	return "", "", errors.New("current directory is not inside a known project")
}

func (s NotesService) NotePath(projPath string) (string, error) {
	notesDir, err := s.notesRepo.NotesDir()
	if err != nil {
		return "", fmt.Errorf("getting notes directory: %v", err.Error())
	}

	return filepath.Join(notesDir, s.NoteFilename(projPath)), nil
}
