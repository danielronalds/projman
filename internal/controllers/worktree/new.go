package worktree

import (
	"errors"
	"fmt"
	"path/filepath"
)

type worktreeCreator interface {
	CreateWorktree(dir, name string) (string, error)
}

type sessionLauncher interface {
	LaunchSession(name, dir string) error
}

type NewController struct {
	worktrees worktreeCreator
	sessions  sessionLauncher
}

func NewNewController(worktrees worktreeCreator, sessions sessionLauncher) NewController {
	return NewController{worktrees, sessions}
}

func (c NewController) Handle(projectRoot, projectName string, args []string) error {
	if len(args) == 0 {
		return errors.New("usage: projman wt new <name>")
	}

	branchName := args[0]

	path, err := c.worktrees.CreateWorktree(projectRoot, branchName)
	if err != nil {
		return fmt.Errorf("creating worktree: %v", err.Error())
	}

	return c.sessions.LaunchSession(filepath.Base(path), path)
}
