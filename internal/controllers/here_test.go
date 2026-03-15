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

		if mock.calledName != expectedName {
			t.Fatalf("LaunchSession name = %q, want %q", mock.calledName, expectedName)
		}

		if mock.calledDir != cwd {
			t.Fatalf("LaunchSession dir = %q, want %q", mock.calledDir, cwd)
		}
	})

	t.Run("returns session launcher error", func(t *testing.T) {
		mock := &mockSessionLauncher{returnErr: fmt.Errorf("session failed")}
		controller := NewHereController(mock)

		err := controller.HandleArgs([]string{})
		if err == nil {
			t.Fatalf("HandleArgs() error = nil, want error")
		}
	})
}
