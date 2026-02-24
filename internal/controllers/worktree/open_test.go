package worktree

import (
	"errors"
	"strings"
	"testing"
)

type mockWorktreeLister struct {
	returnNames []string
	returnErr   error
	calledDir   string
}

func (m *mockWorktreeLister) ListWorktrees(dir string) ([]string, error) {
	m.calledDir = dir
	return m.returnNames, m.returnErr
}

type mockWorktreePathFinder struct {
	returnPath string
	returnErr  error
	calledDir  string
	calledName string
}

func (m *mockWorktreePathFinder) WorktreePath(dir, name string) (string, error) {
	m.calledDir = dir
	m.calledName = name
	return m.returnPath, m.returnErr
}


func TestOpenControllerHandle(t *testing.T) {
	t.Run("noWorktreesExist", func(t *testing.T) {
		lister := &mockWorktreeLister{returnNames: nil}
		paths := &mockWorktreePathFinder{}
		fzf := &mockSelecter{}
		sessions := &mockSessionLauncher{}

		controller := NewOpenController(lister, paths, fzf, sessions)
		err := controller.Handle("/projects/myapp", "myapp", nil)

		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "no worktrees found") {
			t.Fatalf("expected 'no worktrees found' error, got: %v", err)
		}
	})

	t.Run("fuzzySelectSuccess", func(t *testing.T) {
		lister := &mockWorktreeLister{returnNames: []string{"myapp", "feature-auth", "bugfix-login"}}
		paths := &mockWorktreePathFinder{returnPath: "/projects/myapp-feature-auth"}
		fzf := &mockSelecter{returnSelected: "feature-auth"}
		sessions := &mockSessionLauncher{}

		controller := NewOpenController(lister, paths, fzf, sessions)
		err := controller.Handle("/projects/myapp", "myapp", nil)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(fzf.calledOptions) != 3 {
			t.Fatalf("Select called with %d options, want 3", len(fzf.calledOptions))
		}
		if fzf.calledOptions[0] != "myapp (base)" {
			t.Fatalf("Select options[0] = %q, want %q", fzf.calledOptions[0], "myapp (base)")
		}
		if paths.calledName != "feature-auth" {
			t.Fatalf("WorktreePath name = %q, want %q", paths.calledName, "feature-auth")
		}
		if sessions.calledName != "myapp-feature-auth" {
			t.Fatalf("LaunchSession name = %q, want %q", sessions.calledName, "myapp-feature-auth")
		}
		if sessions.calledDir != "/projects/myapp-feature-auth" {
			t.Fatalf("LaunchSession dir = %q, want %q", sessions.calledDir, "/projects/myapp-feature-auth")
		}
	})

	t.Run("fuzzySelectMain", func(t *testing.T) {
		lister := &mockWorktreeLister{returnNames: []string{"myapp", "feature-auth"}}
		paths := &mockWorktreePathFinder{returnPath: "/projects/myapp"}
		fzf := &mockSelecter{returnSelected: "myapp (base)"}
		sessions := &mockSessionLauncher{}

		controller := NewOpenController(lister, paths, fzf, sessions)
		err := controller.Handle("/projects/myapp", "myapp", nil)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if paths.calledName != "myapp" {
			t.Fatalf("WorktreePath name = %q, want %q", paths.calledName, "myapp")
		}
		if sessions.calledName != "myapp" {
			t.Fatalf("LaunchSession name = %q, want %q", sessions.calledName, "myapp")
		}
		if sessions.calledDir != "/projects/myapp" {
			t.Fatalf("LaunchSession dir = %q, want %q", sessions.calledDir, "/projects/myapp")
		}
	})

	t.Run("fuzzySelectCancelled", func(t *testing.T) {
		lister := &mockWorktreeLister{returnNames: []string{"feature-auth"}}
		paths := &mockWorktreePathFinder{}
		fzf := &mockSelecter{returnErr: errors.New("no selection made")}
		sessions := &mockSessionLauncher{}

		controller := NewOpenController(lister, paths, fzf, sessions)
		err := controller.Handle("/projects/myapp", "myapp", nil)

		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "no worktree selected") {
			t.Fatalf("expected 'no worktree selected' error, got: %v", err)
		}
	})

	t.Run("directOpenByName", func(t *testing.T) {
		lister := &mockWorktreeLister{}
		paths := &mockWorktreePathFinder{returnPath: "/projects/myapp-feature-auth"}
		fzf := &mockSelecter{}
		sessions := &mockSessionLauncher{}

		controller := NewOpenController(lister, paths, fzf, sessions)
		err := controller.Handle("/projects/myapp", "myapp", []string{"feature-auth"})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if lister.calledDir != "" {
			t.Fatalf("ListWorktrees should not be called for direct open")
		}
		if paths.calledDir != "/projects/myapp" {
			t.Fatalf("WorktreePath dir = %q, want %q", paths.calledDir, "/projects/myapp")
		}
		if paths.calledName != "feature-auth" {
			t.Fatalf("WorktreePath name = %q, want %q", paths.calledName, "feature-auth")
		}
		if sessions.calledName != "myapp-feature-auth" {
			t.Fatalf("LaunchSession name = %q, want %q", sessions.calledName, "myapp-feature-auth")
		}
	})

	t.Run("directOpenMain", func(t *testing.T) {
		lister := &mockWorktreeLister{}
		paths := &mockWorktreePathFinder{returnPath: "/projects/myapp"}
		fzf := &mockSelecter{}
		sessions := &mockSessionLauncher{}

		controller := NewOpenController(lister, paths, fzf, sessions)
		err := controller.Handle("/projects/myapp", "myapp", []string{"myapp"})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sessions.calledName != "myapp" {
			t.Fatalf("LaunchSession name = %q, want %q", sessions.calledName, "myapp")
		}
		if sessions.calledDir != "/projects/myapp" {
			t.Fatalf("LaunchSession dir = %q, want %q", sessions.calledDir, "/projects/myapp")
		}
	})

	t.Run("directOpenNotFound", func(t *testing.T) {
		lister := &mockWorktreeLister{}
		paths := &mockWorktreePathFinder{returnErr: errors.New("worktree \"nope\" not found")}
		fzf := &mockSelecter{}
		sessions := &mockSessionLauncher{}

		controller := NewOpenController(lister, paths, fzf, sessions)
		err := controller.Handle("/projects/myapp", "myapp", []string{"nope"})

		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "resolving worktree path") {
			t.Fatalf("expected 'resolving worktree path' error, got: %v", err)
		}
	})
}
