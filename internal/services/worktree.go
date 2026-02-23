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

func (s WorktreeService) ListWorktrees(dir string) ([]string, error) {
	mainPath, projectName, err := resolveContext(dir)
	if err != nil {
		return nil, fmt.Errorf("resolving git context: %v", err.Error())
	}

	paths, err := listWorktreeEntries(dir)
	if err != nil {
		return nil, err
	}

	var names []string
	prefix := projectName + "-"
	for _, path := range paths {
		if path == mainPath {
			names = append(names, projectName)
			continue
		}
		name := strings.TrimPrefix(filepath.Base(path), prefix)
		names = append(names, name)
	}
	return names, nil
}

func (s WorktreeService) WorktreePath(dir, name string) (string, error) {
	mainPath, projectName, err := resolveContext(dir)
	if err != nil {
		return "", fmt.Errorf("resolving git context: %v", err.Error())
	}

	paths, err := listWorktreeEntries(dir)
	if err != nil {
		return "", err
	}

	if name == projectName {
		return mainPath, nil
	}

	prefix := projectName + "-"
	for _, path := range paths {
		if path == mainPath {
			continue
		}
		if strings.TrimPrefix(filepath.Base(path), prefix) == name {
			return path, nil
		}
	}

	return "", fmt.Errorf("worktree %q not found", name)
}

func listWorktreeEntries(dir string) ([]string, error) {
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("listing worktrees: %v", err.Error())
	}

	var paths []string
	for _, line := range strings.Split(string(output), "\n") {
		if path, ok := strings.CutPrefix(line, "worktree "); ok {
			paths = append(paths, path)
		}
	}
	return paths, nil
}

func resolveContext(dir string) (mainPath, projectName string, err error) {
	paths, err := listWorktreeEntries(dir)
	if err != nil {
		return "", "", err
	}

	if len(paths) == 0 {
		return "", "", fmt.Errorf("could not determine main worktree path")
	}

	mainPath = paths[0]
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
