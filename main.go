package main

import (
	"fmt"
	"os"

	"github.com/danielronalds/projman/controllers"
	"github.com/danielronalds/projman/repositories"
	"github.com/danielronalds/projman/services"
)

type controller interface {
	HandleArgs(args []string) error
}

func main() {
	config := repositories.NewConfigRepository()
	
	fzf := services.NewFzfService(&config)
	projects := services.NewProjectsService(&config)
	tmux := services.NewTmuxService()

	c := controllers.NewOpenController(projects, fzf, tmux)
	if err := c.HandleArgs(make([]string, 0)); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
	}
}
