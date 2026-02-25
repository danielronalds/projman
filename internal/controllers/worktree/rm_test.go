package worktree

import (
	"errors"
	"strings"
	"testing"
)

type mockWorktreeRemover struct {
	returnErr  error
	calledDir  string
	calledName string
}

func (m *mockWorktreeRemover) RemoveWorktree(dir, name string) error {
	m.calledDir = dir
	m.calledName = name
	return m.returnErr
}

func TestRmControllerHandle(t *testing.T) {
	t.Run("noRemovableWorktrees", func(t *testing.T) {
		lister := &mockWorktreeLister{returnNames: []string{"myapp"}}
		remover := &mockWorktreeRemover{}
		fzf := &mockSelecter{}

		controller := NewRmController(lister, remover, fzf)
		err := controller.Handle("/projects/myapp", "myapp", nil)

		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "no worktrees to remove") {
			t.Fatalf("expected 'no worktrees to remove' error, got: %v", err)
		}
	})

	t.Run("successfulRemove", func(t *testing.T) {
		lister := &mockWorktreeLister{returnNames: []string{"myapp", "feature-auth"}}
		remover := &mockWorktreeRemover{}
		fzf := &mockSelecter{returnSelected: "feature-auth"}

		controller := NewRmController(lister, remover, fzf)
		err := controller.Handle("/projects/myapp", "myapp", nil)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if remover.calledDir != "/projects/myapp" {
			t.Fatalf("RemoveWorktree dir = %q, want %q", remover.calledDir, "/projects/myapp")
		}
		if remover.calledName != "feature-auth" {
			t.Fatalf("RemoveWorktree name = %q, want %q", remover.calledName, "feature-auth")
		}
	})

	t.Run("fuzzySelectCancelled", func(t *testing.T) {
		lister := &mockWorktreeLister{returnNames: []string{"myapp", "feature-auth"}}
		remover := &mockWorktreeRemover{}
		fzf := &mockSelecter{returnErr: errors.New("no selection made")}

		controller := NewRmController(lister, remover, fzf)
		err := controller.Handle("/projects/myapp", "myapp", nil)

		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "no worktree selected") {
			t.Fatalf("expected 'no worktree selected' error, got: %v", err)
		}
	})
}
