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

type worktreeManager interface {
	IsGitRepo(dir string) bool
	MainWorktreePath(dir string) (string, error)
	CreateWorktree(dir, name string) (string, error)
	ListWorktrees(dir string) ([]string, error)
	WorktreePath(dir, name string) (string, error)
	CopyIgnoredFiles(mainPath, worktreePath string) []string
	ListRemoteBranches(dir string) ([]string, error)
	CheckoutWorktree(dir, remoteBranch string) (string, error)
	RemoveWorktree(dir, name string) error
}

type WorktreeController struct {
	worktrees      worktreeManager
	subcommands    map[string]subcommand
	defaultHandler subcommand
}

func NewWorktreeController(worktrees worktreeManager, fzf selecter, sessions sessionLauncher) WorktreeController {
	openController := worktree.NewOpenController(worktrees, worktrees, fzf, sessions)

	subcommands := map[string]subcommand{
		"new":      worktree.NewNewController(worktrees, sessions),
		"open":     openController,
		"checkout": worktree.NewCheckoutController(worktrees, worktrees, fzf, sessions),
		"rm":       worktree.NewRmController(worktrees, worktrees, fzf),
	}

	return WorktreeController{worktrees, subcommands, openController}
}

func (c WorktreeController) HandleArgs(args []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %v", err.Error())
	}

	if !c.worktrees.IsGitRepo(dir) {
		return errors.New("not inside a git repository")
	}

	mainPath, err := c.worktrees.MainWorktreePath(dir)
	if err != nil {
		return fmt.Errorf("resolving worktree context: %v", err.Error())
	}

	projectName := filepath.Base(mainPath)
	subArgs := args[1:]

	if len(subArgs) == 0 {
		return c.defaultHandler.Handle(mainPath, projectName, nil)
	}

	if subcmd, ok := c.subcommands[subArgs[0]]; ok {
		return subcmd.Handle(mainPath, projectName, subArgs[1:])
	}

	return c.defaultHandler.Handle(mainPath, projectName, subArgs)
}
