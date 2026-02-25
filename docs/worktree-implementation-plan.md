# Worktree Support Implementation Plan

- [ ] TODO: Add copying of files in .gitignore over to new worktrees

projman currently has no worktree management. Users who work with git worktrees need to manually run git commands and manage sessions separately. This adds a `wt` (alias `worktree`) subcommand group that integrates worktree lifecycle management with projman's session launching.

## Design

### Worktree location

Worktrees are created as sibling directories: `<parent>/<project>-<name>/`.

Example for project at `~/Projects/projman/`:

```
~/Projects/projman/                   # main project
~/Projects/projman-feature-auth/      # worktree
~/Projects/projman-bugfix-123/        # worktree
```

Discovery uses `git worktree list --porcelain` (works from any worktree in the repo). The user-facing worktree name strips the project prefix, so `projman-feature-auth` is shown as `feature-auth`.

### Command structure

```
projman wt              -> fuzzy select existing worktree, open session
projman wt <name>       -> open session in named worktree directly
projman wt new <name>   -> create worktree + branch, open session
projman wt checkout     -> fetch remotes, fuzzy select branch, create worktree, open session
projman wt rm           -> fuzzy select worktree, remove it + prune
```

Session names use `<project>-<worktree>` format. All commands error if not inside a git repository.

### Controller architecture

```
controllers/
  worktree.go              # package controllers -- sub-dispatcher + unexported subcommand interface
  worktree/                # package worktree -- individual sub-controllers
    open.go                # Default select + direct open
    new.go                 # Create new worktree
    rm.go                  # Remove worktree
    checkout.go            # Fetch + checkout remote branch
```

The `subcommand` interface is defined (unexported) in `controllers/worktree.go`:

```go
type subcommand interface {
    Handle(projectRoot, projectName string, args []string) error
}
```

Sub-controllers satisfy this implicitly via Go's structural typing.

---

## Stage 1: Foundation + `new` command

- [x] Complete

Sets up the service, dispatcher, and the first subcommand.

### Files to create

**`internal/services/worktree.go`** -- Stateless service. Methods needed for this stage:

| Method | Implementation | Purpose |
|--------|---------------|---------|
| `IsGitRepo(dir)` | `git rev-parse --git-dir` | Check if dir is inside a git repo |
| `MainWorktreePath(dir)` | `git worktree list --porcelain`, parse first entry | Get main worktree path |
| `CreateWorktree(dir, name)` | Resolve context, then `git worktree add ../<project>-<name> -b <name>` | Create sibling worktree with new branch |

Internal `resolveContext(dir)` helper parses `git worktree list --porcelain` to get main worktree path and project name.

**`internal/services/worktree_test.go`** -- Integration tests (guarded by `testing.Short()`) for the above methods.

**`internal/controllers/worktree/new.go`** -- Sub-controller with narrow interfaces:

```go
type worktreeCreator interface {
    CreateWorktree(dir, name string) (string, error)
}
type sessionLauncher interface {
    LaunchSession(name, dir string) error
}
```

Handle: require name in `args[0]` -> `CreateWorktree` -> `LaunchSession` with `<project>-<name>`.

**`internal/controllers/worktree/new_test.go`** -- Mock-based tests: missing name, successful create.

**`internal/controllers/worktree.go`** -- Main dispatcher. Defines `subcommand` interface, `gitRepoChecker` and `mainWorktreePathFinder` interfaces. `HandleArgs` resolves context, dispatches to `new` subcommand. Unknown subcommands error for now (open/direct-open added in Stage 2).

### Files to modify

**`main.go`** -- Create `WorktreeService`, wire `WorktreeController` with `"worktree"` and `"wt"` keys.

**`internal/controllers/help.go`** -- Add `wt` entry, bump version.

### Verification

1. `gofmt -l .` + `go vet ./...` + `go test ./...`
2. `projman wt new test-feature` from inside a git repo -- creates sibling worktree, opens session
3. `projman wt new` with no name -- errors with usage hint
4. `projman wt` outside a git repo -- errors

---

## Stage 2: Listing and opening worktrees

- [x] Complete

Adds default fuzzy select and direct-open by name.

### Files to create

**`internal/controllers/worktree/open.go`** -- Sub-controller with narrow interfaces:

```go
type worktreeLister interface {
    ListWorktrees(dir string) ([]string, error)
}
type worktreePathFinder interface {
    WorktreePath(dir, name string) (string, error)
}
type selecter interface {
    Select(options []string) (string, error)
}
type sessionLauncher interface {
    LaunchSession(name, dir string) error
}
```

Handle: if no args, list -> select -> open. If args, treat as direct name -> resolve path -> open. Error with guidance if no worktrees exist.

**`internal/controllers/worktree/open_test.go`** -- Mock-based tests: no worktrees, successful select, direct open, not found.

### Files to modify

**`internal/services/worktree.go`** -- Add methods:

| Method | Implementation | Purpose |
|--------|---------------|---------|
| `ListWorktrees(dir)` | `git worktree list --porcelain`, exclude main, strip `<project>-` prefix | User-facing worktree names |
| `WorktreePath(dir, name)` | Match name against parsed worktree list | Resolve name to full path |

**`internal/services/worktree_test.go`** -- Add tests for the new methods.

**`internal/controllers/worktree.go`** -- Wire `open` sub-controller. Set it as default handler + fallback for unknown subcommands (direct name open).

### Verification

1. `go test ./...`
2. `projman wt` -- shows existing worktrees in fuzzy finder, opens session on select
3. `projman wt feature-auth` -- opens session directly in named worktree
4. `projman wt` with no worktrees -- error with "create one with projman wt new" guidance

---

## Stage 3: `checkout` command

- [x] Complete

Adds fetching remote branches and checking them out into worktrees.

### Files to create

**`internal/controllers/worktree/checkout.go`** -- Sub-controller:

```go
type remoteBranchLister interface {
    ListRemoteBranches(dir string) ([]string, error)
}
type worktreeCheckout interface {
    CheckoutWorktree(dir, remoteBranch string) (string, error)
}
type selecter interface {
    Select(options []string) (string, error)
}
type sessionLauncher interface {
    LaunchSession(name, dir string) error
}
```

Handle: `WithSpinner` around `ListRemoteBranches` -> select -> `CheckoutWorktree` -> `LaunchSession`. Strip `origin/` for session name.

**`internal/controllers/worktree/checkout_test.go`** -- Mock-based tests: no remote branches, successful checkout.

### Files to modify

**`internal/services/worktree.go`** -- Add methods:

| Method | Implementation | Purpose |
|--------|---------------|---------|
| `ListRemoteBranches(dir)` | `git fetch --all --prune` + `git branch -r --format=%(refname:short)` | Fetch and list remote branches |
| `CheckoutWorktree(dir, branch)` | Strip `origin/` prefix, `git worktree add ../<project>-<local> <branch>` | Create worktree from remote branch |

**`internal/services/worktree_test.go`** -- Add integration tests.

**`internal/controllers/worktree.go`** -- Wire `checkout` into subcommands map.

### Verification

1. `go test ./...`
2. `projman wt checkout` -- fetches branches with spinner, shows fuzzy finder, creates worktree and opens session

---

## Stage 4: `rm` command

- [x] Complete

Adds worktree removal with fuzzy selection.

### Files to create

**`internal/controllers/worktree/rm.go`** -- Sub-controller:

```go
type worktreeLister interface {
    ListWorktrees(dir string) ([]string, error)
}
type worktreeRemover interface {
    RemoveWorktree(dir, name string) error
}
type selecter interface {
    Select(options []string) (string, error)
}
```

Handle: list -> select -> remove + prune + print confirmation. Git naturally errors if inside the target worktree.

**`internal/controllers/worktree/rm_test.go`** -- Mock-based tests: no worktrees, successful remove.

### Files to modify

**`internal/services/worktree.go`** -- Add method:

| Method | Implementation | Purpose |
|--------|---------------|---------|
| `RemoveWorktree(dir, name)` | Resolve path, `git worktree remove <path>` + `git worktree prune` | Remove and clean up |

**`internal/services/worktree_test.go`** -- Add integration test.

**`internal/controllers/worktree.go`** -- Wire `rm` into subcommands map.

### Verification

1. `go test ./...`
2. `projman wt rm` -- shows worktrees, select one, removes it and prints confirmation
3. Try removing the worktree you're currently in -- should get a clear error from git
