package controllers

import (
	"fmt"

	"github.com/danielronalds/projman/internal/services"
	"github.com/danielronalds/projman/internal/ui"
)

type groupedLister interface {
	ListProjectsByDirectory(filter string) ([]services.ProjectGroup, error)
}

type ListController struct {
	projects groupedLister
}

func NewListController(projects groupedLister) ListController {
	return ListController{
		projects: projects,
	}
}

func (c ListController) HandleArgs(args []string) error {
	filter := ""
	if len(args) > 1 {
		filter = args[1]
	}

	groups, err := c.projects.ListProjectsByDirectory(filter)
	if err != nil {
		return err
	}

	fmt.Print(ui.RenderProjectList(groups, filter != ""))

	return nil
}
