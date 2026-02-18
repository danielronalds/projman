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
  clone     Clone any git url to a project dir
  active    Open an existing session
  config    Open the projman config in your editor
  rm        Remove a project from your machine
  health    Verify required dependencies are installed
  help      Show this menu
`

const VERSION = "0.6.0"

type HelpController struct{}

func NewHelpController() HelpController {
	return HelpController{}
}

func (c HelpController) HandleArgs(args []string) error {
	fmt.Printf(HELP_MENU, VERSION)
	return nil
}
