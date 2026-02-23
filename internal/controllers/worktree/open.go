package worktree

import (
	"errors"
	"fmt"
	"strings"
)

type worktreeLister interface {
	ListWorktrees(dir string) ([]string, error)
}

type worktreePathFinder interface {
	WorktreePath(dir, name string) (string, error)
}

type selecter interface {
	Select(options []string) (string, error)
}

const baseSuffix = " (base)"

type OpenController struct {
	worktrees worktreeLister
	paths     worktreePathFinder
	fzf       selecter
	sessions  sessionLauncher
}

func NewOpenController(worktrees worktreeLister, paths worktreePathFinder, fzf selecter, sessions sessionLauncher) OpenController {
	return OpenController{worktrees, paths, fzf, sessions}
}

func (c OpenController) Handle(projectRoot, projectName string, args []string) error {
	name, err := c.resolveWorktreeName(projectRoot, projectName, args)
	if err != nil {
		return err
	}

	path, err := c.paths.WorktreePath(projectRoot, name)
	if err != nil {
		return fmt.Errorf("resolving worktree path: %v", err.Error())
	}

	sessionName := projectName + "-" + name
	if name == projectName {
		sessionName = projectName
	}
	return c.sessions.LaunchSession(sessionName, path)
}

func (c OpenController) resolveWorktreeName(projectRoot, projectName string, args []string) (string, error) {
	if len(args) > 0 {
		return args[0], nil
	}

	worktrees, err := c.worktrees.ListWorktrees(projectRoot)
	if err != nil {
		return "", fmt.Errorf("listing worktrees: %v", err.Error())
	}

	if len(worktrees) == 0 {
		return "", errors.New("no worktrees found, create one with: projman wt new <name>")
	}

	displayNames := make([]string, len(worktrees))
	for i, name := range worktrees {
		if name == projectName {
			displayNames[i] = name + baseSuffix
		} else {
			displayNames[i] = name
		}
	}

	selected, err := c.fzf.Select(displayNames)
	if err != nil {
		return "", errors.New("no worktree selected")
	}

	return strings.TrimSuffix(selected, baseSuffix), nil
}
