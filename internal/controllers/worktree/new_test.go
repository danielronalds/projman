package worktree

import (
	"strings"
	"testing"
)

type mockWorktreeCreator struct {
	returnPath string
	returnErr  error
	calledDir  string
	calledName string
}

func (m *mockWorktreeCreator) CreateWorktree(dir, name string) (string, error) {
	m.calledDir = dir
	m.calledName = name
	return m.returnPath, m.returnErr
}

type mockIgnoredFileHandler struct {
	hasFiles       bool
	returnWarnings []string
}

func (m *mockIgnoredFileHandler) HasIgnoredFiles(dir string) bool {
	return m.hasFiles
}

func (m *mockIgnoredFileHandler) CopyIgnoredFiles(mainPath, worktreePath string) []string {
	return m.returnWarnings
}

type mockSessionLauncher struct {
	returnErr  error
	calledName string
	calledDir  string
}

func (m *mockSessionLauncher) LaunchSession(name, dir string) error {
	m.calledName = name
	m.calledDir = dir
	return m.returnErr
}

func TestNewControllerHandle(t *testing.T) {
	t.Run("missingName", func(t *testing.T) {
		controller := NewNewController(
			&mockWorktreeCreator{},
			&mockIgnoredFileHandler{},
			&mockSessionLauncher{},
		)

		err := controller.Handle("/projects/myapp", "myapp", []string{})
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "usage") {
			t.Fatalf("expected usage error, got: %v", err)
		}
	})

	t.Run("successfulCreate", func(t *testing.T) {
		worktrees := &mockWorktreeCreator{returnPath: "/projects/myapp-feature-auth"}
		sessions := &mockSessionLauncher{}
		controller := NewNewController(
			worktrees,
			&mockIgnoredFileHandler{},
			sessions,
		)

		err := controller.Handle("/projects/myapp", "myapp", []string{"feature/auth"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if worktrees.calledDir != "/projects/myapp" {
			t.Fatalf("CreateWorktree dir = %q, want %q", worktrees.calledDir, "/projects/myapp")
		}
		if worktrees.calledName != "feature/auth" {
			t.Fatalf("CreateWorktree name = %q, want %q", worktrees.calledName, "feature/auth")
		}
		if sessions.calledName != "myapp-feature-auth" {
			t.Fatalf("LaunchSession name = %q, want %q", sessions.calledName, "myapp-feature-auth")
		}
		if sessions.calledDir != "/projects/myapp-feature-auth" {
			t.Fatalf("LaunchSession dir = %q, want %q", sessions.calledDir, "/projects/myapp-feature-auth")
		}
	})
}
