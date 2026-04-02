package controllers

import "fmt"

const (
	greenColor = "\033[32m"
	redColor   = "\033[31m"
	resetColor = "\033[0m"
)

type healthChecker interface {
	CheckRequirements(programs []string) map[string]bool
}

type dependency struct {
	name        string
	description string
}

var dependencies = []dependency{
	{name: "tmux", description: "Required for the tmux session provider"},
	{name: "code", description: "Required for the vscode session provider"},
	{name: "gh", description: "Required for the remote command"},
	{name: "git", description: "Required for the clone command"},
}

type HealthController struct {
	healthService healthChecker
}

func NewHealthController(healthService healthChecker) HealthController {
	return HealthController{healthService: healthService}
}

func (c HealthController) HandleArgs(args []string) error {
	depNames := make([]string, len(dependencies))
	for i, dep := range dependencies {
		depNames[i] = dep.name
	}

	results := c.healthService.CheckRequirements(depNames)

	hasFailure := false
	for _, dep := range dependencies {
		statusIcon := fmt.Sprintf("[%s%s%s]", greenColor, "✓", resetColor)
		suffix := ""

		if !results[dep.name] {
			hasFailure = true
			statusIcon = fmt.Sprintf("[%s%s%s]", redColor, "x", resetColor)
			suffix = fmt.Sprintf(" - %s", dep.description)
		}

		fmt.Printf("%s %s%s\n", statusIcon, dep.name, suffix)
	}

	fmt.Println()

	if hasFailure {
		return fmt.Errorf("one or more requirements are missing")
	}

	fmt.Println("All requirements satisfied.")
	return nil
}
