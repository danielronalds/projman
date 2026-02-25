package worktree

import (
	"errors"
	"fmt"
)

type worktreeRemover interface {
	RemoveWorktree(dir, name string) error
}

type RmController struct {
	worktrees worktreeLister
	remover   worktreeRemover
	fzf       selecter
}

func NewRmController(worktrees worktreeLister, remover worktreeRemover, fzf selecter) RmController {
	return RmController{worktrees, remover, fzf}
}

func (c RmController) Handle(projectRoot, projectName string, args []string) error {
	worktrees, err := c.worktrees.ListWorktrees(projectRoot)
	if err != nil {
		return fmt.Errorf("listing worktrees: %v", err.Error())
	}

	var removable []string
	for _, name := range worktrees {
		if name != projectName {
			removable = append(removable, name)
		}
	}

	if len(removable) == 0 {
		return errors.New("no worktrees to remove, create one with: projman wt new <name>")
	}

	selected, err := c.fzf.Select(removable)
	if err != nil {
		return errors.New("no worktree selected")
	}

	if err := c.remover.RemoveWorktree(projectRoot, selected); err != nil {
		return fmt.Errorf("removing worktree: %v", err.Error())
	}

	fmt.Printf("removed worktree %q\n", selected)
	return nil
}
