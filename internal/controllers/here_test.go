package controllers

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type mockSessionLauncher struct {
	calledName string
	calledDir  string
	returnErr  error
}

func (m *mockSessionLauncher) LaunchSession(name, dir string) error {
	m.calledName = name
	m.calledDir = dir
	return m.returnErr
}

type mockSanitiser struct{}

func (m *mockSanitiser) Sanitise(name string) string {
	return name
}

func TestHereController_HandleArgs(t *testing.T) {
	subdir := "test-project"
	tempDir := filepath.Join(t.TempDir(), subdir)
	err := os.MkdirAll(tempDir, 0755)
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Chdir(tempDir)

	t.Run("launches session with cwd", func(t *testing.T) {
		mock := &mockSessionLauncher{}
		controller := NewHereController(mock, &mockSanitiser{})

		err := controller.HandleArgs([]string{"here"})
		if err != nil {
			t.Fatalf("HandleArgs() error = %v, want nil", err)
		}

		if mock.calledName != subdir {
			t.Fatalf("LaunchSession name = %q, want %q", mock.calledName, subdir)
		}

		if mock.calledDir != tempDir {
			t.Fatalf("LaunchSession dir = %q, want %q", mock.calledDir, tempDir)
		}
	})

	t.Run("returns session launcher error", func(t *testing.T) {
		mock := &mockSessionLauncher{returnErr: fmt.Errorf("session failed")}
		controller := NewHereController(mock, &mockSanitiser{})

		err := controller.HandleArgs([]string{"here"})
		if err == nil {
			t.Fatalf("HandleArgs() error = nil, want error")
		}
	})
}
