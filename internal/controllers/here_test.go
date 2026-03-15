package controllers

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type mockSessionLauncher struct {
	launchedName string
	launchedDir  string
	err          error
}

func (m *mockSessionLauncher) LaunchSession(name, dir string) error {
	m.launchedName = name
	m.launchedDir = dir
	return m.err
}

func TestHereController_HandleArgs(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	expectedName := filepath.Base(cwd)

	t.Run("launches session with cwd", func(t *testing.T) {
		mock := &mockSessionLauncher{}
		controller := NewHereController(mock)

		err := controller.HandleArgs([]string{})
		if err != nil {
			t.Fatalf("HandleArgs() error = %v, want nil", err)
		}

		if mock.launchedName != expectedName {
			t.Fatalf("LaunchSession name = %q, want %q", mock.launchedName, expectedName)
		}

		if mock.launchedDir != cwd {
			t.Fatalf("LaunchSession dir = %q, want %q", mock.launchedDir, cwd)
		}
	})

	t.Run("returns session launcher error", func(t *testing.T) {
		mock := &mockSessionLauncher{err: fmt.Errorf("session failed")}
		controller := NewHereController(mock)

		err := controller.HandleArgs([]string{})
		if err == nil {
			t.Fatalf("HandleArgs() error = nil, want error")
		}
	})
}
