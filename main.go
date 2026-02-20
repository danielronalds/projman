package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/danielronalds/projman/internal/controllers"
	"github.com/danielronalds/projman/internal/repositories"
	"github.com/danielronalds/projman/internal/services"
)

type controller interface {
	HandleArgs(args []string) error
}

func parseProviderFlag(args []string) (string, []string) {
	remaining := make([]string, 0, len(args))
	provider := ""
	for i := 0; i < len(args); i++ {
		if args[i] == "--provider" && i+1 < len(args) {
			provider = args[i+1]
			i++
		} else if value, ok := strings.CutPrefix(args[i], "--provider="); ok {
			provider = value
		} else {
			remaining = append(remaining, args[i])
		}
	}
	return provider, remaining
}

func run(args []string) {
	config := repositories.NewConfigRepository()

	providerFlag, args := parseProviderFlag(args)

	var providerCfg services.ProviderConfig = config
	if providerFlag != "" {
		providerCfg = services.NewProviderConfigOverride(config, providerFlag)
	}

	selector := services.NewSelectService(config)
	projects := services.NewProjectsService(config)
	github := services.NewGithubService(config)
	creater := services.NewCreaterService(config)
	health := services.NewHealthService()
	git := services.NewGitService()

	sessionProvider, err := services.NewSessionProvider(providerCfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err.Error())
		os.Exit(1)
	}

	cmd := "local"
	if len(args) > 0 {
		cmd = args[0]
	}

	controllerMap := map[string]controller{
		"new":    controllers.NewNewController(selector, creater, sessionProvider, config),
		"local":  controllers.NewOpenController(projects, selector, sessionProvider),
		"remote": controllers.NewRemoteController(github, projects, selector, sessionProvider, config),
		"clone":  controllers.NewCloneController(github, selector, sessionProvider, config),
		"active": controllers.NewActiveController(projects, selector, sessionProvider),
		"config": controllers.NewConfigController(config),
		"rm":     controllers.NewRmController(projects, selector, git),
		"help":   controllers.NewHelpController(),
		"health": controllers.NewHealthController(health, config),
	}

	handler, ok := controllerMap[cmd]

	if !ok {
		handler = controllers.NewDirectOpenController(projects, selector, sessionProvider)
	}

	if err := handler.HandleArgs(args); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
	}
}

func main() {
	args := os.Args[1:]

	run(args)
}
