package controllers

import "fmt"

const HELP_MENU = `projman v%v

A cli tool for managing projects on your local machine.

To open a project, run projman with no arguments

Commands
  new       Create a new project
  local     Open a project currently on your machine
  remote    Open a project from github
  help      Show this menu
`

const VERSION = "0.1.0"

type HelpController struct{}

func NewHelpController() HelpController {
	return HelpController{}
}

func (c HelpController) HandleArgs(args []string) error {
	fmt.Printf(HELP_MENU, VERSION)
	return nil
}
