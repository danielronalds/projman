package worktree

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/danielronalds/projman/internal/ui"
)

type remoteBranchLister interface {
	ListRemoteBranches(dir string) ([]string, error)
}

type worktreeCheckout interface {
	CheckoutWorktree(dir, remoteBranch string) (string, error)
	CopyIgnoredFiles(mainPath, worktreePath string) []string
}

type CheckoutController struct {
	branches remoteBranchLister
	checkout worktreeCheckout
	fzf      selecter
	sessions sessionLauncher
}

func NewCheckoutController(branches remoteBranchLister, checkout worktreeCheckout, fzf selecter, sessions sessionLauncher) CheckoutController {
	return CheckoutController{branches, checkout, fzf, sessions}
}

func (c CheckoutController) Handle(projectRoot, projectName string, args []string) error {
	branches, err := ui.WithSpinner("fetching remote branches...", func() ([]string, error) {
		return c.branches.ListRemoteBranches(projectRoot)
	})
	if err != nil {
		return fmt.Errorf("listing remote branches: %v", err)
	}

	if len(branches) == 0 {
		return errors.New("no remote branches found")
	}

	selected, err := c.fzf.Select(branches)
	if err != nil {
		return errors.New("no branch selected")
	}

	path, err := ui.WithSpinner("checking out worktree...", func() (string, error) {
		return c.checkout.CheckoutWorktree(projectRoot, selected)
	})
	if err != nil {
		return fmt.Errorf("checking out worktree: %v", err)
	}

	gitignorePath := filepath.Join(projectRoot, ".gitignore")
	if _, statErr := os.Stat(gitignorePath); statErr == nil && ui.Confirm("Copy ignored files to new worktree?") {
		warnings, _ := ui.WithSpinner("copying ignored files...", func() ([]string, error) {
			return c.checkout.CopyIgnoredFiles(projectRoot, path), nil
		})
		for _, w := range warnings {
			fmt.Fprintf(os.Stderr, "warning: %s\n", w)
		}
	}

	return c.sessions.LaunchSession(filepath.Base(path), path)
}
