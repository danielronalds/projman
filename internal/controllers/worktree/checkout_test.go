package worktree

import (
	"errors"
	"strings"
	"testing"
)

type mockRemoteBranchLister struct {
	returnBranches []string
	returnErr      error
	calledDir      string
}

func (m *mockRemoteBranchLister) ListRemoteBranches(dir string) ([]string, error) {
	m.calledDir = dir
	return m.returnBranches, m.returnErr
}

type mockWorktreeCheckout struct {
	returnPath     string
	returnErr      error
	calledDir      string
	calledBranch   string
	returnWarnings []string
}

func (m *mockWorktreeCheckout) CheckoutWorktree(dir, remoteBranch string) (string, error) {
	m.calledDir = dir
	m.calledBranch = remoteBranch
	return m.returnPath, m.returnErr
}

func (m *mockWorktreeCheckout) CopyIgnoredFiles(mainPath, worktreePath string) []string {
	return m.returnWarnings
}

func TestCheckoutControllerHandle(t *testing.T) {
	t.Run("noRemoteBranches", func(t *testing.T) {
		branches := &mockRemoteBranchLister{returnBranches: nil}
		checkout := &mockWorktreeCheckout{}
		fzf := &mockSelecter{}
		sessions := &mockSessionLauncher{}

		controller := NewCheckoutController(branches, checkout, fzf, sessions)
		err := controller.Handle("/projects/myapp", "myapp", nil)

		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "no remote branches found") {
			t.Fatalf("expected 'no remote branches found' error, got: %v", err)
		}
	})

	t.Run("fetchError", func(t *testing.T) {
		branches := &mockRemoteBranchLister{returnErr: errors.New("fetch failed")}
		checkout := &mockWorktreeCheckout{}
		fzf := &mockSelecter{}
		sessions := &mockSessionLauncher{}

		controller := NewCheckoutController(branches, checkout, fzf, sessions)
		err := controller.Handle("/projects/myapp", "myapp", nil)

		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "listing remote branches") {
			t.Fatalf("expected 'listing remote branches' error, got: %v", err)
		}
	})

	t.Run("selectionCancelled", func(t *testing.T) {
		branches := &mockRemoteBranchLister{returnBranches: []string{"origin/main", "origin/feature/auth"}}
		checkout := &mockWorktreeCheckout{}
		fzf := &mockSelecter{returnErr: errors.New("cancelled")}
		sessions := &mockSessionLauncher{}

		controller := NewCheckoutController(branches, checkout, fzf, sessions)
		err := controller.Handle("/projects/myapp", "myapp", nil)

		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "no branch selected") {
			t.Fatalf("expected 'no branch selected' error, got: %v", err)
		}
	})

	t.Run("checkoutError", func(t *testing.T) {
		branches := &mockRemoteBranchLister{returnBranches: []string{"origin/main"}}
		checkout := &mockWorktreeCheckout{returnErr: errors.New("worktree already exists")}
		fzf := &mockSelecter{returnSelected: "origin/main"}
		sessions := &mockSessionLauncher{}

		controller := NewCheckoutController(branches, checkout, fzf, sessions)
		err := controller.Handle("/projects/myapp", "myapp", nil)

		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "checking out worktree") {
			t.Fatalf("expected 'checking out worktree' error, got: %v", err)
		}
	})

	t.Run("successfulCheckout", func(t *testing.T) {
		branches := &mockRemoteBranchLister{returnBranches: []string{"origin/main", "origin/feature/auth"}}
		checkout := &mockWorktreeCheckout{returnPath: "/projects/myapp-feature-auth"}
		fzf := &mockSelecter{returnSelected: "origin/feature/auth"}
		sessions := &mockSessionLauncher{}

		controller := NewCheckoutController(branches, checkout, fzf, sessions)
		err := controller.Handle("/projects/myapp", "myapp", nil)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if branches.calledDir != "/projects/myapp" {
			t.Fatalf("ListRemoteBranches dir = %q, want %q", branches.calledDir, "/projects/myapp")
		}
		if checkout.calledDir != "/projects/myapp" {
			t.Fatalf("CheckoutWorktree dir = %q, want %q", checkout.calledDir, "/projects/myapp")
		}
		if checkout.calledBranch != "origin/feature/auth" {
			t.Fatalf("CheckoutWorktree branch = %q, want %q", checkout.calledBranch, "origin/feature/auth")
		}
		if sessions.calledName != "myapp-feature-auth" {
			t.Fatalf("LaunchSession name = %q, want %q", sessions.calledName, "myapp-feature-auth")
		}
		if sessions.calledDir != "/projects/myapp-feature-auth" {
			t.Fatalf("LaunchSession dir = %q, want %q", sessions.calledDir, "/projects/myapp-feature-auth")
		}
	})
}
