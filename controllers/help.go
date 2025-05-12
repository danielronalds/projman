package controllers

import "fmt"

const HELP_MENU = `projman v%v

Usage: projman [command | project]

A cli tool for managing projects on your local machine.

To select and open a project, run projman with no arguments

Commands
  new       Create a new project
  local     Open a project currently on your machine
  remote    Open a project from github
  active    Open an existing session
  help      Show this menu
`

const VERSION = "0.2.0"

type HelpController struct{}

func NewHelpController() HelpController {
	return HelpController{}
}

func (c HelpController) HandleArgs(args []string) error {
	fmt.Printf(HELP_MENU, VERSION)
	return nil
}
