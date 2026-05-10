# projman

A CLI tool for managing projects on your local machine. projman helps you quickly create, open, and manage your dev projects with integrated session management via tmux or VS Code.

## Install

```console
go install github.com/danielronalds/projman@latest
```

See the [GitHub repo](https://github.com/danielronalds/projman) for prerequisites and build-from-source instructions.

## Quickstart

```console
projman          # fuzzy-select and open a project
projman new      # create a new project
projman remote   # browse and clone a GitHub repo
projman wt       # manage worktrees for the current project
```

## Commands

```console
Usage: projman [command | project]

Commands
  new       Create a new project
  local     Open a project currently on your machine
  remote    Open a project from github
  clone     Clone any git url to a project dir
  active    Open an existing session
  here      Open a session in the current directory
  list      List all local projects (optionally filter by directory name)
  config    Open the projman config in your editor
  rm        Remove a project from your machine
  worktree  Manage git worktrees -- run 'projman wt help' for details (alias: wt)
  health    Verify required dependencies are installed
  help      Show this menu
```

`--provider, -p <provider>` overrides the configured session provider for a single invocation (`tmux` or `vscode`).

## Next

- [Configuration](configuration.md) -- config file, session providers, templates
- [Worktrees](worktrees.md) -- `projman wt` subcommands
