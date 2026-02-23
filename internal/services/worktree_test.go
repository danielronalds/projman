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
