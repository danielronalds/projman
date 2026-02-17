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

type HealthController struct {
	healthService healthChecker
	dependencies  []string
}

func NewHealthController(healthService healthChecker) HealthController {
	deps := []string{"go", "tmux", "gh", "git"}

	return HealthController{healthService: healthService, dependencies: deps}
}

func (c HealthController) HandleArgs(args []string) error {
	results := c.healthService.CheckRequirements(c.dependencies)

	hasFailure := false
	for _, dep := range c.dependencies {
		statusIcon := fmt.Sprintf("[%s%s%s]", greenColor, "✓", resetColor)

		if !results[dep] {
			hasFailure = true
			statusIcon = fmt.Sprintf("[%s%s%s]", redColor, "x", resetColor)
		}

		fmt.Printf("%s %s\n", statusIcon, dep)
	}

	fmt.Println()

	if hasFailure {
		return fmt.Errorf("one or more requirements are missing")
	}

	fmt.Println("All requirements satisfied.")
	return nil
}
