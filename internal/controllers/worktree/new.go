package worktree

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/danielronalds/projman/internal/ui"
)

type worktreeCreator interface {
	CreateWorktree(dir, name string) (string, error)
}

type ignoredFileHandler interface {
	CopyIgnoredFiles(mainPath, worktreePath string) []string
}

type sessionLauncher interface {
	LaunchSession(name, dir string) error
}

type NewController struct {
	worktrees    worktreeCreator
	ignoredFiles ignoredFileHandler
	sessions     sessionLauncher
}

func NewNewController(worktrees worktreeCreator, ignoredFiles ignoredFileHandler, sessions sessionLauncher) NewController {
	return NewController{worktrees, ignoredFiles, sessions}
}

func (c NewController) Handle(projectRoot, projectName string, args []string) error {
	if len(args) == 0 {
		return errors.New("usage: projman wt new <name>")
	}

	branchName := args[0]

	path, err := ui.WithSpinner("creating worktree...", func() (string, error) {
		return c.worktrees.CreateWorktree(projectRoot, branchName)
	})
	if err != nil {
		return fmt.Errorf("creating worktree: %v", err.Error())
	}

	gitignorePath := filepath.Join(projectRoot, ".gitignore")
	if _, statErr := os.Stat(gitignorePath); statErr == nil && ui.Confirm("Copy ignored files to new worktree?") {
		warnings, _ := ui.WithSpinner("copying ignored files...", func() ([]string, error) {
			return c.ignoredFiles.CopyIgnoredFiles(projectRoot, path), nil
		})
		for _, w := range warnings {
			fmt.Fprintf(os.Stderr, "warning: %s\n", w)
		}
	}

	return c.sessions.LaunchSession(filepath.Base(path), path)
}
