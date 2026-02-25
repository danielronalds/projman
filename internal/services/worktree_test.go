package services

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestSanitiseForDirectory(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "slashToDash", input: "feature/foo", want: "feature-foo"},
		{name: "spaceToDash", input: "feature/foo bar", want: "feature-foo-bar"},
		{name: "underscorePreserved", input: "my_branch", want: "my_branch"},
		{name: "collapsedDashes", input: "feat//double", want: "feat-double"},
		{name: "trimLeadingTrailing", input: "--leading--", want: "leading"},
		{name: "dotPreserved", input: "v1.2.3", want: "v1.2.3"},
		{name: "specialCharsStripped", input: "feat@#$%name", want: "featname"},
		{name: "simple", input: "test-feature", want: "test-feature"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitiseForDirectory(tt.input)
			if got != tt.want {
				t.Fatalf("sanitiseForDirectory(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsGitRepo(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	s := NewWorktreeService()

	t.Run("insideGitRepo", func(t *testing.T) {
		cwd, _ := os.Getwd()
		if !s.IsGitRepo(cwd) {
			t.Fatalf("expected current directory to be a git repo")
		}
	})

	t.Run("outsideGitRepo", func(t *testing.T) {
		if s.IsGitRepo(os.TempDir()) {
			t.Fatalf("expected temp dir to not be a git repo")
		}
	})
}

func TestMainWorktreePath(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	s := NewWorktreeService()

	cwd, _ := os.Getwd()
	got, err := s.MainWorktreePath(cwd)
	if err != nil {
		t.Fatalf("MainWorktreePath() error = %v", err)
	}

	if !strings.HasSuffix(got, "projman") {
		t.Fatalf("MainWorktreePath() = %q, want suffix 'projman'", got)
	}
}

func TestCreateWorktree(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	tmpDir, _ := filepath.EvalSymlinks(t.TempDir())
	repoDir := filepath.Join(tmpDir, "myproject")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatalf("failed to create repo dir: %v", err)
	}

	run := func(dir string, args ...string) {
		t.Helper()
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("command %v failed: %v\n%s", args, err, output)
		}
	}

	run(repoDir, "git", "init")
	run(repoDir, "git", "config", "user.email", "test@test.com")
	run(repoDir, "git", "config", "user.name", "Test")
	run(repoDir, "git", "commit", "--allow-empty", "-m", "initial")

	s := NewWorktreeService()

	t.Run("simpleWorktree", func(t *testing.T) {
		path, err := s.CreateWorktree(repoDir, "test-branch")
		if err != nil {
			t.Fatalf("CreateWorktree() error = %v", err)
		}

		expectedPath := filepath.Join(tmpDir, "myproject-test-branch")
		if path != expectedPath {
			t.Fatalf("CreateWorktree() path = %q, want %q", path, expectedPath)
		}

		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Fatalf("worktree directory does not exist at %q", path)
		}
	})

	t.Run("branchWithSlash", func(t *testing.T) {
		path, err := s.CreateWorktree(repoDir, "feature/my-feature")
		if err != nil {
			t.Fatalf("CreateWorktree() error = %v", err)
		}

		expectedPath := filepath.Join(tmpDir, "myproject-feature-my-feature")
		if path != expectedPath {
			t.Fatalf("CreateWorktree() path = %q, want %q", path, expectedPath)
		}

		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Fatalf("worktree directory does not exist at %q", path)
		}

		branchCmd := exec.Command("git", "branch", "--list", "feature/my-feature")
		branchCmd.Dir = repoDir
		output, err := branchCmd.Output()
		if err != nil {
			t.Fatalf("git branch --list error = %v", err)
		}
		if !strings.Contains(string(output), "feature/my-feature") {
			t.Fatalf("expected branch 'feature/my-feature' to exist, got: %q", string(output))
		}
	})

	t.Run("invalidNameAllSpecialChars", func(t *testing.T) {
		_, err := s.CreateWorktree(repoDir, "@#$%")
		if err == nil {
			t.Fatalf("expected error for invalid name, got nil")
		}
	})
}

func TestListWorktrees(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	tmpDir, _ := filepath.EvalSymlinks(t.TempDir())
	repoDir := filepath.Join(tmpDir, "myproject")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatalf("failed to create repo dir: %v", err)
	}

	run := func(dir string, args ...string) {
		t.Helper()
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("command %v failed: %v\n%s", args, err, output)
		}
	}

	run(repoDir, "git", "init")
	run(repoDir, "git", "config", "user.email", "test@test.com")
	run(repoDir, "git", "config", "user.name", "Test")
	run(repoDir, "git", "commit", "--allow-empty", "-m", "initial")

	s := NewWorktreeService()

	t.Run("onlyBase", func(t *testing.T) {
		names, err := s.ListWorktrees(repoDir)
		if err != nil {
			t.Fatalf("ListWorktrees() error = %v", err)
		}
		if len(names) != 1 {
			t.Fatalf("ListWorktrees() returned %d names, want 1", len(names))
		}
		if names[0] != "myproject" {
			t.Fatalf("ListWorktrees()[0] = %q, want %q", names[0], "myproject")
		}
	})

	if _, err := s.CreateWorktree(repoDir, "feature-one"); err != nil {
		t.Fatalf("CreateWorktree() error = %v", err)
	}
	if _, err := s.CreateWorktree(repoDir, "feature/two"); err != nil {
		t.Fatalf("CreateWorktree() error = %v", err)
	}

	t.Run("withWorktrees", func(t *testing.T) {
		names, err := s.ListWorktrees(repoDir)
		if err != nil {
			t.Fatalf("ListWorktrees() error = %v", err)
		}
		if len(names) != 3 {
			t.Fatalf("ListWorktrees() returned %d names, want 3", len(names))
		}
		expected := map[string]bool{"myproject": true, "feature-one": true, "feature-two": true}
		for _, name := range names {
			if !expected[name] {
				t.Fatalf("unexpected worktree name %q", name)
			}
		}
	})

	t.Run("fromWorktreeDir", func(t *testing.T) {
		worktreeDir := filepath.Join(tmpDir, "myproject-feature-one")
		names, err := s.ListWorktrees(worktreeDir)
		if err != nil {
			t.Fatalf("ListWorktrees() from worktree dir error = %v", err)
		}
		if len(names) != 3 {
			t.Fatalf("ListWorktrees() from worktree dir returned %d names, want 3", len(names))
		}
	})
}

func TestWorktreePath(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	tmpDir, _ := filepath.EvalSymlinks(t.TempDir())
	repoDir := filepath.Join(tmpDir, "myproject")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatalf("failed to create repo dir: %v", err)
	}

	run := func(dir string, args ...string) {
		t.Helper()
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("command %v failed: %v\n%s", args, err, output)
		}
	}

	run(repoDir, "git", "init")
	run(repoDir, "git", "config", "user.email", "test@test.com")
	run(repoDir, "git", "config", "user.name", "Test")
	run(repoDir, "git", "commit", "--allow-empty", "-m", "initial")

	s := NewWorktreeService()
	if _, err := s.CreateWorktree(repoDir, "feature-one"); err != nil {
		t.Fatalf("CreateWorktree() error = %v", err)
	}

	t.Run("mainWorktree", func(t *testing.T) {
		path, err := s.WorktreePath(repoDir, "myproject")
		if err != nil {
			t.Fatalf("WorktreePath() error = %v", err)
		}
		if path != repoDir {
			t.Fatalf("WorktreePath() = %q, want %q", path, repoDir)
		}
	})

	t.Run("found", func(t *testing.T) {
		path, err := s.WorktreePath(repoDir, "feature-one")
		if err != nil {
			t.Fatalf("WorktreePath() error = %v", err)
		}
		expectedPath := filepath.Join(tmpDir, "myproject-feature-one")
		if path != expectedPath {
			t.Fatalf("WorktreePath() = %q, want %q", path, expectedPath)
		}
	})

	t.Run("notFound", func(t *testing.T) {
		_, err := s.WorktreePath(repoDir, "nonexistent")
		if err == nil {
			t.Fatalf("WorktreePath() expected error for nonexistent worktree")
		}
		if !strings.Contains(err.Error(), "not found") {
			t.Fatalf("WorktreePath() error = %q, want 'not found'", err.Error())
		}
	})
}

func setupRemoteAndClone(t *testing.T) (repoDir, tmpDir string) {
	t.Helper()
	tmpDir, _ = filepath.EvalSymlinks(t.TempDir())
	remoteDir := filepath.Join(tmpDir, "remote")
	repoDir = filepath.Join(tmpDir, "myproject")

	run := func(dir string, args ...string) {
		t.Helper()
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("command %v failed: %v\n%s", args, err, output)
		}
	}

	run(tmpDir, "git", "init", remoteDir)
	run(remoteDir, "git", "config", "user.email", "test@test.com")
	run(remoteDir, "git", "config", "user.name", "Test")
	run(remoteDir, "git", "commit", "--allow-empty", "-m", "initial")
	run(remoteDir, "git", "branch", "feature/test-branch")
	run(remoteDir, "git", "branch", "feature/second-branch")
	run(tmpDir, "git", "clone", remoteDir, repoDir)
	return repoDir, tmpDir
}

func TestListRemoteBranches(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repoDir, _ := setupRemoteAndClone(t)
	s := NewWorktreeService()

	t.Run("includesRemoteBranches", func(t *testing.T) {
		branches, err := s.ListRemoteBranches(repoDir)
		if err != nil {
			t.Fatalf("ListRemoteBranches() error = %v", err)
		}
		found := false
		for _, b := range branches {
			if b == "origin/feature/test-branch" {
				found = true
			}
		}
		if !found {
			t.Fatalf("expected 'origin/feature/test-branch' in branches, got %v", branches)
		}
	})

	t.Run("excludesHead", func(t *testing.T) {
		branches, err := s.ListRemoteBranches(repoDir)
		if err != nil {
			t.Fatalf("ListRemoteBranches() error = %v", err)
		}
		for _, b := range branches {
			if strings.Contains(b, "HEAD") {
				t.Fatalf("expected HEAD to be filtered out, got %q in branches", b)
			}
		}
	})
}

func TestCheckoutWorktree(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repoDir, tmpDir := setupRemoteAndClone(t)
	s := NewWorktreeService()

	t.Run("checkoutRemoteBranch", func(t *testing.T) {
		path, err := s.CheckoutWorktree(repoDir, "origin/feature/test-branch")
		if err != nil {
			t.Fatalf("CheckoutWorktree() error = %v", err)
		}

		expectedPath := filepath.Join(tmpDir, "myproject-feature-test-branch")
		if path != expectedPath {
			t.Fatalf("CheckoutWorktree() path = %q, want %q", path, expectedPath)
		}

		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Fatalf("worktree directory does not exist at %q", path)
		}
	})

	t.Run("alreadyCheckedOut", func(t *testing.T) {
		path1, err := s.CheckoutWorktree(repoDir, "origin/feature/second-branch")
		if err != nil {
			t.Fatalf("first CheckoutWorktree() error = %v", err)
		}
		path2, err := s.CheckoutWorktree(repoDir, "origin/feature/second-branch")
		if err != nil {
			t.Fatalf("second CheckoutWorktree() error = %v", err)
		}
		if path1 != path2 {
			t.Fatalf("second CheckoutWorktree() path = %q, want %q", path2, path1)
		}
	})

	t.Run("invalidRemoteBranch", func(t *testing.T) {
		_, err := s.CheckoutWorktree(repoDir, "noslash")
		if err == nil {
			t.Fatalf("expected error for branch with no remote prefix, got nil")
		}
	})
}

func TestCopyIgnoredFiles(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	s := NewWorktreeService()

	type testHelpers struct {
		run   func(dir string, args ...string)
		mkdir func(path string)
		write func(path, content string)
	}

	setupRepo := func(t *testing.T) (repoDir string, tmpDir string, h testHelpers) {
		t.Helper()
		tmpDir, _ = filepath.EvalSymlinks(t.TempDir())
		repoDir = filepath.Join(tmpDir, "myproject")
		if err := os.MkdirAll(repoDir, 0755); err != nil {
			t.Fatalf("failed to create repo dir: %v", err)
		}

		h.run = func(dir string, args ...string) {
			t.Helper()
			cmd := exec.Command(args[0], args[1:]...)
			cmd.Dir = dir
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("command %v failed: %v\n%s", args, err, output)
			}
		}

		h.mkdir = func(path string) {
			t.Helper()
			if err := os.MkdirAll(path, 0755); err != nil {
				t.Fatalf("failed to create directory %s: %v", path, err)
			}
		}

		h.write = func(path, content string) {
			t.Helper()
			if err := os.WriteFile(path, []byte(content), 0644); err != nil {
				t.Fatalf("failed to write file %s: %v", path, err)
			}
		}

		h.run(repoDir, "git", "init")
		h.run(repoDir, "git", "config", "user.email", "test@test.com")
		h.run(repoDir, "git", "config", "user.name", "Test")
		h.run(repoDir, "git", "commit", "--allow-empty", "-m", "initial")
		return repoDir, tmpDir, h
	}

	t.Run("noGitignore", func(t *testing.T) {
		repoDir, tmpDir, h := setupRepo(t)
		worktreePath := filepath.Join(tmpDir, "myproject-test")

		h.mkdir(worktreePath)
		warnings := s.CopyIgnoredFiles(repoDir, worktreePath)
		if len(warnings) != 0 {
			t.Fatalf("expected no warnings, got %v", warnings)
		}
	})

	t.Run("ignoredFileCopied", func(t *testing.T) {
		repoDir, tmpDir, h := setupRepo(t)

		h.write(filepath.Join(repoDir, ".gitignore"), "*.log\n")
		h.write(filepath.Join(repoDir, "debug.log"), "log content")
		h.run(repoDir, "git", "add", ".gitignore")
		h.run(repoDir, "git", "commit", "-m", "add gitignore")

		worktreePath := filepath.Join(tmpDir, "myproject-test")
		h.mkdir(worktreePath)

		warnings := s.CopyIgnoredFiles(repoDir, worktreePath)
		if len(warnings) != 0 {
			t.Fatalf("expected no warnings, got %v", warnings)
		}

		content, err := os.ReadFile(filepath.Join(worktreePath, "debug.log"))
		if err != nil {
			t.Fatalf("expected debug.log to be copied, got error: %v", err)
		}
		if string(content) != "log content" {
			t.Fatalf("debug.log content = %q, want %q", string(content), "log content")
		}
	})

	t.Run("ignoredDirectoryCopied", func(t *testing.T) {
		repoDir, tmpDir, h := setupRepo(t)

		h.write(filepath.Join(repoDir, ".gitignore"), "node_modules/\n")
		h.mkdir(filepath.Join(repoDir, "node_modules", "pkg"))
		h.write(filepath.Join(repoDir, "node_modules", "pkg", "index.js"), "module.exports = {}")
		h.run(repoDir, "git", "add", ".gitignore")
		h.run(repoDir, "git", "commit", "-m", "add gitignore")

		worktreePath := filepath.Join(tmpDir, "myproject-test")
		h.mkdir(worktreePath)

		warnings := s.CopyIgnoredFiles(repoDir, worktreePath)
		if len(warnings) != 0 {
			t.Fatalf("expected no warnings, got %v", warnings)
		}

		content, err := os.ReadFile(filepath.Join(worktreePath, "node_modules", "pkg", "index.js"))
		if err != nil {
			t.Fatalf("expected node_modules/pkg/index.js to be copied, got error: %v", err)
		}
		if string(content) != "module.exports = {}" {
			t.Fatalf("index.js content = %q, want %q", string(content), "module.exports = {}")
		}
	})

	t.Run("nestedGitignore", func(t *testing.T) {
		repoDir, tmpDir, h := setupRepo(t)

		h.mkdir(filepath.Join(repoDir, "subdir"))
		h.write(filepath.Join(repoDir, "subdir", ".gitignore"), "*.tmp\n")
		h.write(filepath.Join(repoDir, "subdir", "cache.tmp"), "cached")
		h.run(repoDir, "git", "add", "subdir/.gitignore")
		h.run(repoDir, "git", "commit", "-m", "add nested gitignore")

		worktreePath := filepath.Join(tmpDir, "myproject-test")
		h.mkdir(worktreePath)

		warnings := s.CopyIgnoredFiles(repoDir, worktreePath)
		if len(warnings) != 0 {
			t.Fatalf("expected no warnings, got %v", warnings)
		}

		content, err := os.ReadFile(filepath.Join(worktreePath, "subdir", "cache.tmp"))
		if err != nil {
			t.Fatalf("expected subdir/cache.tmp to be copied, got error: %v", err)
		}
		if string(content) != "cached" {
			t.Fatalf("cache.tmp content = %q, want %q", string(content), "cached")
		}
	})
}
