package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/danielronalds/projman/repositories"
	"github.com/danielronalds/projman/services"
)

func main() {
	config := repositories.NewConfigRepository()
	
	fzf := services.NewFzfService(config)

	output, _ := fzf.Select([]string{ "option one", "option two"})

	fmt.Printf("Output: %v\n", output)
}

func runTmux() {
	cmd := exec.Command("tmux")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

