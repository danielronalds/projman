package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	output := runFzf()
	fmt.Printf("Output: %v\n", output)
	runTmux()
}

func runTmux() {
	cmd := exec.Command("tmux")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func runFzf() string {
	cmd := exec.Command("fzf")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(out))
}
