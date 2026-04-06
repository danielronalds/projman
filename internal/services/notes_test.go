package services

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

type mockNotesConfig struct {
	projectDirs []string
}

func (m mockNotesConfig) ProjectDirs() []string {
	return m.projectDirs
}

type mockNotesRepository struct {
	notesDir string
	err      error
}

func (m mockNotesRepository) NotesDir() (string, error) {
	return m.notesDir, m.err
}

func TestNoteFilename(t *testing.T) {
	service := NewNotesService(mockNotesConfig{}, mockNotesRepository{})

	tests := []struct {
		name     string
		projPath string
		expected string
	}{
		{
			name:     "standard project path",
			projPath: "/Users/dan/Personal/projman",
			expected: "personal_projman.md",
		},
		{
			name:     "work project path",
			projPath: "/Users/dan/Work/myapp",
			expected: "work_myapp.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.NoteFilename(tt.projPath)
			if result != tt.expected {
				t.Errorf("NoteFilename(%q) = %q, want %q", tt.projPath, result, tt.expected)
			}
		})
	}
}

func TestNotePath(t *testing.T) {
	t.Run("assembles full note path", func(t *testing.T) {
		service := NewNotesService(mockNotesConfig{}, mockNotesRepository{
			notesDir: "/home/dan/.local/share/projman/notes",
		})

		result, err := service.NotePath("/Users/dan/Personal/projman")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := "/home/dan/.local/share/projman/notes/personal_projman.md"
		if result != expected {
			t.Errorf("NotePath() = %q, want %q", result, expected)
		}
	})

	t.Run("returns error when NotesDir fails", func(t *testing.T) {
		service := NewNotesService(mockNotesConfig{}, mockNotesRepository{
			err: errors.New("no data dir"),
		})

		_, err := service.NotePath("/Users/dan/Personal/projman")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestResolveCurrentProject(t *testing.T) {
	tests := []struct {
		name         string
		cwdRelative  string
		expectedName string
		expectedPath string
		expectError  bool
	}{
		{
			name:         "cwd at project root",
			cwdRelative:  "Personal/projman",
			expectedName: "projman",
			expectedPath: "Personal/projman",
			expectError:  false,
		},
		{
			name:         "cwd in subdirectory",
			cwdRelative:  "Personal/projman/internal/services",
			expectedName: "projman",
			expectedPath: "Personal/projman",
			expectError:  false,
		},
		{
			name:        "cwd outside all project dirs",
			cwdRelative: "Other/something",
			expectError: true,
		},
		{
			name:        "cwd is the project dir itself",
			cwdRelative: "Personal",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			cwdPath := filepath.Join(tmpDir, tt.cwdRelative)
			if err := os.MkdirAll(cwdPath, 0755); err != nil {
				t.Fatalf("failed to create cwd: %v", err)
			}
			t.Chdir(cwdPath)

			projectDir := filepath.Join(tmpDir, "Personal") + "/"

			service := NewNotesService(
				mockNotesConfig{projectDirs: []string{projectDir}},
				mockNotesRepository{},
			)

			name, path, err := service.ResolveCurrentProject()

			if tt.expectError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if name != tt.expectedName {
				t.Errorf("name = %q, want %q", name, tt.expectedName)
			}

			expectedFullPath := filepath.Join(tmpDir, tt.expectedPath)
			if path != expectedFullPath {
				t.Errorf("path = %q, want %q", path, expectedFullPath)
			}
		})
	}
}
