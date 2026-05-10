# Git worktrees

`projman wt` (or `projman worktree`) manages git worktrees for the current project. Run it with no arguments to select and open an existing worktree.

```console
Usage: projman wt [subcommand | name]

Manage git worktrees for the current project.

To select and open a worktree, run projman wt with no arguments

Commands
  new <name>  Create a new worktree and branch, then open a session
  checkout    Fetch remote branches, select one, and open a session
  rm          Select and remove a worktree
  help        Show this menu
```

## Creating a worktree

`projman wt new <name>` creates a new git worktree with a new branch and launches a session in it:

```console
projman wt new feature/auth
```

## Checking out a remote branch

`projman wt checkout` fetches all remote branches, presents them in a fuzzy finder, and creates a worktree from the selected branch:

```console
projman wt checkout
```

## Opening a worktree

`projman wt` with no arguments lists existing worktrees in a fuzzy finder. You can also pass a name directly:

```console
projman wt                  # fuzzy select a worktree to open
projman wt my-worktree      # open a named worktree directly
```

## Removing a worktree

`projman wt rm` lists worktrees (excluding the base) in a fuzzy finder and removes the selected one:

```console
projman wt rm
```

## Copying ignored files

When creating or checking out a worktree, projman offers to copy gitignored files (e.g. `.env`, local databases, build outputs) from the main worktree into the new one. This is useful for the small files your tooling needs but git won't track -- without them, the new worktree often won't run.

### Skipping large directories

For ignored directories you'd rather rebuild than copy (`node_modules/`, `target/`, `dist/`), set `worktree_copy_excludes` in your config. Patterns use [doublestar](https://github.com/bmatcuk/doublestar) glob syntax, matched against paths returned by `git ls-files` (relative to the repository root):

```json
{
  "worktree_copy_excludes": ["**/node_modules", "dist", "*.log"]
}
```

- `node_modules` matches only the top-level `node_modules` directory
- `**/node_modules` matches `node_modules` at any depth
- `*.log` matches top-level log files
- `**/*.log` matches log files at any depth

Invalid patterns cause projman to exit at startup with an error.

## Naming conventions

Worktree directories are named `<project>-<sanitised-branch>`, where slashes and spaces in the branch name become dashes. For example, `feature/auth` in a project called `myapp` produces a directory called `myapp-feature-auth`. Sessions follow the same naming pattern.
