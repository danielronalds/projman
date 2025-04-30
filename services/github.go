package services

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type GithubService struct{
	config projectsConfig
}

func NewGithubService(config projectsConfig) GithubService {
	return GithubService{config}
}

type githubJsonOutput struct {
	Name string `json:"name"`
}

func (s GithubService) ListProjects() ([]projectName, error) {
	projects := make([]string, 0)

	cmd := exec.Command(
		"gh",
		"repo",
		"list",
		"--json",
		"name",
		"--limit", // Have to override default limit, otherwise not all repos returned
		"5000",
	)

	output, err := cmd.Output()
	if err != nil {
		return projects, err
	}

	var outputProjects []githubJsonOutput
	if err = json.Unmarshal(output, &outputProjects); err != nil {
		return projects, err
	}

	for _, project := range outputProjects {
		projects = append(projects, project.Name)
	}

	return projects, nil
}

func (s GithubService) Clone(name, dir string) (projectPath, error) {
	cmd := exec.Command("gh", "repo", "clone", name)

	cmd.Dir = dir

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	projPath := fmt.Sprintf("%v%v", dir, name)

	return projPath, cmd.Run()
}
