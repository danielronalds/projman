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

type healthConfig interface {
	SessionProvider() string
}

type HealthController struct {
	healthService healthChecker
	config        healthConfig
}

func NewHealthController(healthService healthChecker, config healthConfig) HealthController {
	return HealthController{healthService: healthService, config: config}
}

func getProviderDependency(provider string) string {
	switch provider {
	case "vscode":
		return "code"
	case "tmux":
		return "tmux"
	default:
		return "tmux"
	}
}

func (c HealthController) HandleArgs(args []string) error {
	providerDep := getProviderDependency(c.config.SessionProvider())
	deps := []string{"go", providerDep, "gh", "git"}

	results := c.healthService.CheckRequirements(deps)

	hasFailure := false
	for _, dep := range deps {
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
