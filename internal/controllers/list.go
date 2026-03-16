package controllers

import (
	"fmt"
	"strings"

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
		filter = strings.TrimSpace(args[1])
	}

	groups, err := c.projects.ListProjectsByDirectory(filter)
	if err != nil {
		return fmt.Errorf("unable to list projects: %v", err.Error())
	}

	fmt.Print(ui.RenderProjectList(groups, filter != ""))

	return nil
}
