package services

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type WorktreeService struct{}

func NewWorktreeService() WorktreeService {
	return WorktreeService{}
}

func (s WorktreeService) IsGitRepo(dir string) bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = dir
	return cmd.Run() == nil
}

func (s WorktreeService) MainWorktreePath(dir string) (string, error) {
	mainPath, _, err := resolveContext(dir)
	if err != nil {
		return "", err
	}
	return mainPath, nil
}

func (s WorktreeService) CreateWorktree(dir, name string) (string, error) {
	mainPath, projectName, err := resolveContext(dir)
	if err != nil {
		return "", fmt.Errorf("resolving git context: %v", err.Error())
	}

	sanitised := sanitiseForDirectory(name)
	if sanitised == "" {
		return "", fmt.Errorf("invalid worktree name %q: results in empty directory name after sanitisation", name)
	}
	dirName := projectName + "-" + sanitised
	worktreePath := filepath.Join(filepath.Dir(mainPath), dirName)

	cmd := exec.Command("git", "worktree", "add", worktreePath, "-b", name)
	cmd.Dir = mainPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("creating worktree: %v", strings.TrimSpace(string(output)))
	}

	return worktreePath, nil
}

func resolveContext(dir string) (mainPath, projectName string, err error) {
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("listing worktrees: %v", err.Error())
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if path, ok := strings.CutPrefix(line, "worktree "); ok {
			mainPath = path
			break
		}
	}

	if mainPath == "" {
		return "", "", fmt.Errorf("could not determine main worktree path")
	}

	projectName = filepath.Base(mainPath)
	return mainPath, projectName, nil
}

var nonAlphanumericDashDotUnderscore = regexp.MustCompile(`[^a-zA-Z0-9\-._]`)
var multipleDashes = regexp.MustCompile(`-{2,}`)

func sanitiseForDirectory(name string) string {
	result := strings.ReplaceAll(name, "/", "-")
	result = strings.ReplaceAll(result, " ", "-")
	result = nonAlphanumericDashDotUnderscore.ReplaceAllString(result, "")
	result = multipleDashes.ReplaceAllString(result, "-")
	result = strings.Trim(result, "-")
	return result
}
