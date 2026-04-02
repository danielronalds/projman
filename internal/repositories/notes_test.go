package repositories

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestUserDataDir(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)

	dataDir, err := userDataDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(dataDir, "projman") {
		t.Errorf("expected path to contain 'projman', got %s", dataDir)
	}

	if runtime.GOOS == "darwin" {
		if !strings.Contains(dataDir, filepath.Join("Library", "Application Support")) {
			t.Errorf("expected macOS path to contain 'Library/Application Support', got %s", dataDir)
		}
	}
}

func TestNotesDir(t *testing.T) {
	t.Run("creates the directory on disk", func(t *testing.T) {
		tempHome := t.TempDir()
		t.Setenv("HOME", tempHome)

		repo := NewNotesRepository()
		notesDir, err := repo.NotesDir()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		info, err := os.Stat(notesDir)
		if err != nil {
			t.Fatalf("expected notes directory to exist: %v", err)
		}

		if !info.IsDir() {
			t.Error("expected notes path to be a directory")
		}

		if !strings.HasSuffix(notesDir, "notes") {
			t.Errorf("expected path to end with 'notes', got %s", notesDir)
		}
	})

	t.Run("is idempotent", func(t *testing.T) {
		tempHome := t.TempDir()
		t.Setenv("HOME", tempHome)

		repo := NewNotesRepository()

		first, err := repo.NotesDir()
		if err != nil {
			t.Fatalf("first call: unexpected error: %v", err)
		}

		second, err := repo.NotesDir()
		if err != nil {
			t.Fatalf("second call: unexpected error: %v", err)
		}

		if first != second {
			t.Errorf("expected same path, got %s and %s", first, second)
		}
	})
}
