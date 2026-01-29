package services

import "os/exec"

type GitService struct{}

func NewGitService() GitService {
	return GitService{}
}

func (s GitService) HasUncommittedChanges(dir string) bool {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return len(output) > 0
}
