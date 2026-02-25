package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
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
	mainPath, _, _, err := resolveContext(dir)
	if err != nil {
		return "", err
	}
	return mainPath, nil
}

func (s WorktreeService) CreateWorktree(dir, name string) (string, error) {
	mainPath, projectName, _, err := resolveContext(dir)
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
	mainPath, projectName, paths, err := resolveContext(dir)
	if err != nil {
		return nil, fmt.Errorf("resolving git context: %v", err.Error())
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
	mainPath, projectName, paths, err := resolveContext(dir)
	if err != nil {
		return "", fmt.Errorf("resolving git context: %v", err.Error())
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

func resolveContext(dir string) (mainPath, projectName string, paths []string, err error) {
	paths, err = listWorktreeEntries(dir)
	if err != nil {
		return "", "", nil, err
	}

	if len(paths) == 0 {
		return "", "", nil, fmt.Errorf("could not determine main worktree path")
	}

	mainPath = paths[0]
	projectName = filepath.Base(mainPath)
	return mainPath, projectName, paths, nil
}

func listIgnoredPaths(dir string) ([]string, error) {
	cmd := exec.Command("git", "ls-files", "--ignored", "--exclude-standard", "--others", "--directory")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git ls-files: %v", err)
	}

	var paths []string
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		if line == "" {
			continue
		}
		paths = append(paths, line)
	}
	return paths, nil
}

// NOTE: cp -a is not available on Windows
func copyPath(src, dst string) error {
	src = strings.TrimSuffix(src, string(filepath.Separator))
	dst = strings.TrimSuffix(dst, string(filepath.Separator))

	if _, err := os.Stat(src); os.IsNotExist(err) {
		return nil
	}

	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("creating parent directory: %v", err)
	}

	cmd := exec.Command("cp", "-a", src, dst)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v", strings.TrimSpace(string(output)))
	}
	return nil
}

func (s WorktreeService) CopyIgnoredFiles(mainPath, worktreePath string) []string {
	paths, err := listIgnoredPaths(mainPath)
	if err != nil {
		return []string{fmt.Sprintf("listing ignored files: %v", err)}
	}

	var warnings []string
	var mu sync.Mutex
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 8)

	for _, relPath := range paths {
		wg.Add(1)
		go func(relPath string) {
			defer wg.Done()
			// Limit concurrent cp processes to avoid exhausting OS resources
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			src := filepath.Join(mainPath, relPath)
			dst := filepath.Join(worktreePath, relPath)

			if err := copyPath(src, dst); err != nil {
				mu.Lock()
				warnings = append(warnings, fmt.Sprintf("copying %s: %v", relPath, err))
				mu.Unlock()
			}
		}(relPath)
	}

	wg.Wait()
	return warnings
}

func (s WorktreeService) ListRemoteBranches(dir string) ([]string, error) {
	fetchCmd := exec.Command("git", "fetch", "--all", "--prune")
	fetchCmd.Dir = dir
	if output, err := fetchCmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("fetching remotes: %v", strings.TrimSpace(string(output)))
	}

	branchCmd := exec.Command("git", "branch", "-r", "--format=%(refname:short)")
	branchCmd.Dir = dir
	output, err := branchCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("listing remote branches: %v", err.Error())
	}

	var branches []string
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		if line == "" || strings.Contains(line, "HEAD") {
			continue
		}
		branches = append(branches, line)
	}
	return branches, nil
}

func (s WorktreeService) CheckoutWorktree(dir, remoteBranch string) (string, error) {
	mainPath, projectName, _, err := resolveContext(dir)
	if err != nil {
		return "", fmt.Errorf("resolving git context: %v", err.Error())
	}

	parts := strings.SplitN(remoteBranch, "/", 2)
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid remote branch %q: expected format <remote>/<branch>", remoteBranch)
	}
	localBranch := parts[1]

	sanitised := sanitiseForDirectory(localBranch)
	if sanitised == "" {
		return "", fmt.Errorf("invalid branch name %q: results in empty directory name after sanitisation", remoteBranch)
	}

	dirName := projectName + "-" + sanitised
	worktreePath := filepath.Join(filepath.Dir(mainPath), dirName)

	cmd := exec.Command("git", "worktree", "add", worktreePath, localBranch)
	cmd.Dir = mainPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		if _, statErr := os.Stat(worktreePath); statErr == nil {
			return worktreePath, nil
		}
		return "", fmt.Errorf("checking out worktree: %v", strings.TrimSpace(string(output)))
	}

	return worktreePath, nil
}

func (s WorktreeService) RemoveWorktree(dir, name string) error {
	path, err := s.WorktreePath(dir, name)
	if err != nil {
		return fmt.Errorf("resolving worktree path: %v", err.Error())
	}

	removeCmd := exec.Command("git", "worktree", "remove", path)
	removeCmd.Dir = dir
	if output, err := removeCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("removing worktree: %v", strings.TrimSpace(string(output)))
	}

	pruneCmd := exec.Command("git", "worktree", "prune")
	pruneCmd.Dir = dir
	if output, err := pruneCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("pruning worktrees: %v", strings.TrimSpace(string(output)))
	}

	return nil
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
