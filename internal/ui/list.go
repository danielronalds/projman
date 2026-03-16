package ui

import (
	"fmt"
	"strings"

	"github.com/danielronalds/projman/internal/services"
)

func RenderProjectList(groups []services.ProjectGroup, filtered bool) string {
	var buf strings.Builder

	totalProjects := 0
	for i, group := range groups {
		if i > 0 {
			buf.WriteString("\n")
		}
		fmt.Fprintf(&buf, "\033[1m%s\033[0m\n", group.Directory)
		for _, project := range group.Projects {
			fmt.Fprintf(&buf, " \u2022 %s\n", project)
		}
		totalProjects += len(group.Projects)
	}

	if len(groups) > 0 {
		buf.WriteString("\n")
	}

	label := "projects"
	if totalProjects == 1 {
		label = "project"
	}

	if filtered {
		fmt.Fprintf(&buf, "%d %s matched\n", totalProjects, label)
	} else {
		fmt.Fprintf(&buf, "%d %s total on your system\n", totalProjects, label)
	}

	return buf.String()
}
