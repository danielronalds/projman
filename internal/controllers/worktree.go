package controllers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/danielronalds/projman/internal/controllers/worktree"
)

type subcommand interface {
	Handle(projectRoot, projectName string, args []string) error
}

type gitRepoChecker interface {
	IsGitRepo(dir string) bool
}

type mainWorktreePathFinder interface {
	MainWorktreePath(dir string) (string, error)
}

type worktreeCreator interface {
	CreateWorktree(dir, name string) (string, error)
}

type WorktreeController struct {
	git         gitRepoChecker
	worktrees   mainWorktreePathFinder
	subcommands map[string]subcommand
}

func NewWorktreeController(git gitRepoChecker, worktrees mainWorktreePathFinder, worktreeCreator worktreeCreator, sessions sessionLauncher) WorktreeController {
	subcommands := map[string]subcommand{
		"new": worktree.NewNewController(worktreeCreator, sessions),
	}

	return WorktreeController{git, worktrees, subcommands}
}

func (c WorktreeController) HandleArgs(args []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %v", err.Error())
	}

	if !c.git.IsGitRepo(dir) {
		return errors.New("not inside a git repository")
	}

	mainPath, err := c.worktrees.MainWorktreePath(dir)
	if err != nil {
		return fmt.Errorf("resolving worktree context: %v", err.Error())
	}

	projectName := filepath.Base(mainPath)

	subArgs := args[1:]

	if len(subArgs) == 0 {
		return errors.New("usage: projman wt <new> [args]")
	}

	subcmd, ok := c.subcommands[subArgs[0]]
	if !ok {
		return fmt.Errorf("unknown worktree command %q", subArgs[0])
	}

	return subcmd.Handle(mainPath, projectName, subArgs[1:])
}
