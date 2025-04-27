package services

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type fzfConfig interface {
	Layout() string
	Theme() string
}

type FzfService struct {
	config fzfConfig
}

func NewFzfService(config fzfConfig) FzfService {
	return FzfService { config }
}

func (s FzfService) layout() string {
	return fmt.Sprintf("--layout=%v", s.config.Layout())
}

func (s FzfService) color() string {
	return fmt.Sprintf("--color=%v", s.config.Theme())
}

func (s FzfService) Select(options []string) (string, error) {
	cmd := exec.Command("fzf", s.layout(), s.color())

	cmd.Stderr = os.Stderr
	inPipe, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}

	for _, option := range options {
		santisedOption := santiseOption(option)
		inPipe.Write([]byte(santisedOption))
	}
	if err = inPipe.Close(); err != nil {
		return "", nil
	}

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	
	return strings.TrimSpace(string(out)), nil
}

func santiseOption(option string) string {
	s := strings.TrimSpace(option)
	return fmt.Sprintf("%v\n", s)
}
