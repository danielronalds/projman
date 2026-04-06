package repositories

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type NotesRepository struct{}

func NewNotesRepository() NotesRepository {
	return NotesRepository{}
}

func (r NotesRepository) NotesDir() (string, error) {
	dataDir, err := userDataDir()
	if err != nil {
		return "", fmt.Errorf("unable to get data directory: %v", err.Error())
	}

	notesDir := filepath.Join(dataDir, "notes")

	if err := os.MkdirAll(notesDir, 0755); err != nil {
		return "", fmt.Errorf("unable to create notes directory: %v", err.Error())
	}

	return notesDir, nil
}

func userDataDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to get user's home dir: %v", err.Error())
	}

	if runtime.GOOS == "darwin" {
		return filepath.Join(homeDir, "Library", "Application Support", "projman"), nil
	}

	return filepath.Join(homeDir, ".local", "share", "projman"), nil
}
