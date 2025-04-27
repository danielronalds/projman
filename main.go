package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command("fzf")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	//cmd.Stdout = os.Stdout
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Output: %v\n", string(out))
}
