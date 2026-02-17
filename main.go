package main

import (
	"fmt"
	"os"

	"github.com/danielronalds/projman/internal/controllers"
	"github.com/danielronalds/projman/internal/repositories"
	"github.com/danielronalds/projman/internal/services"
)

type controller interface {
	HandleArgs(args []string) error
}

func run(args []string) {
	config := repositories.NewConfigRepository()

	fzf := services.NewFzfService(config)
	projects := services.NewProjectsService(config)
	github := services.NewGithubService(config)
	tmux := services.NewTmuxService(config)
	creater := services.NewCreaterService(config)
	health := services.NewHealthService()

	cmd := "local"
	if len(args) > 0 {
		cmd = args[0]
	}

	controllerMap := map[string]controller{
		"new":    controllers.NewNewController(fzf, creater, tmux, config),
		"local":  controllers.NewOpenController(projects, fzf, tmux),
		"remote": controllers.NewRemoteController(github, projects, fzf, tmux, config),
		"clone":  controllers.NewCloneController(github, fzf, tmux, config),
		"active": controllers.NewActiveController(projects, fzf, tmux),
		"config": controllers.NewConfigController(config),
		"help":   controllers.NewHelpController(),
		"health": controllers.NewHealthController(health),
	}

	handler, ok := controllerMap[cmd]

	if !ok {
		handler = controllers.NewDirectOpenController(projects, fzf, tmux)
	}

	if err := handler.HandleArgs(args); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
	}
}

func main() {
	args := os.Args[1:]

	run(args)
}
