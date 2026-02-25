package worktree

import "fmt"

const WORKTREE_HELP_MENU = `Usage: projman wt [subcommand | name]

Manage git worktrees for the current project.

To select and open a worktree, run projman wt with no arguments

Commands
  new <name>  Create a new worktree and branch, then open a session
  checkout    Fetch remote branches, select one, and open a session
  rm          Select and remove a worktree
  help        Show this menu
`

type HelpController struct{}

func NewHelpController() HelpController {
	return HelpController{}
}

func (c HelpController) Handle(projectRoot, projectName string, args []string) error {
	fmt.Print(WORKTREE_HELP_MENU)
	return nil
}
