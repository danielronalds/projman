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

func parseProviderFlag(args []string) (string, []string, error) {
	remaining := make([]string, 0, len(args))
	provider := ""
	for i := 0; i < len(args); i++ {
		if args[i] == "--provider" || args[i] == "-p" {
			if i+1 >= len(args) {
				return "", nil, fmt.Errorf("--provider requires a value")
			}
			provider = args[i+1]
			i++
		} else if value, ok := strings.CutPrefix(args[i], "--provider="); ok {
			if value == "" {
				return "", nil, fmt.Errorf("--provider requires a value")
			}
			provider = value
		} else if value, ok := strings.CutPrefix(args[i], "-p="); ok {
			if value == "" {
				return "", nil, fmt.Errorf("--provider requires a value")
			}
			provider = value
		} else {
			remaining = append(remaining, args[i])
		}
	}
	return provider, remaining, nil
}

func run(args []string) {
	config := repositories.NewConfigRepository()
	notesRepo := repositories.NewNotesRepository()

	providerFlag, args, err := parseProviderFlag(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err.Error())
		os.Exit(1)
	}

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
	worktree := services.NewWorktreeService()
	notes := services.NewNotesService(config, notesRepo)

	sanitiser := services.NewSanitiser()

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
		"new":      controllers.NewNewController(selector, creater, sessionProvider, config),
		"local":    controllers.NewOpenController(projects, selector, sessionProvider),
		"remote":   controllers.NewRemoteController(github, projects, selector, sessionProvider, config),
		"clone":    controllers.NewCloneController(github, selector, sessionProvider, config),
		"active":   controllers.NewActiveController(projects, selector, sessionProvider),
		"here":     controllers.NewHereController(sessionProvider, sanitiser),
		"config":   controllers.NewConfigController(config),
		"rm":       controllers.NewRmController(projects, selector, git),
		"notes":    controllers.NewNotesController(notes, notes, projects, selector),
		"list":     controllers.NewListController(projects),
		"help":     controllers.NewHelpController(),
		"health":   controllers.NewHealthController(health),
		"worktree": controllers.NewWorktreeController(worktree, selector, sessionProvider),
		"wt":       controllers.NewWorktreeController(worktree, selector, sessionProvider),
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
