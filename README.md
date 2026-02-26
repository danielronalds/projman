# projman

A CLI tool for managing projects on your local machine.
projman helps you quickly create, open, and manage your
dev projects with integrated session management.

## Features

- **Project Management**: Create and open projects from local directories or GitHub repositories
- **Session Providers**: Choose between tmux or VS Code for opening projects
- **Project Templates**: Use customizable templates when creating new projects
- **Fuzzy Finder**: Built-in interactive fuzzy finding for project selection
- **GitHub Integration**: Browse, clone, and open remote repositories
- **Session Management**: Manage active sessions (tmux provider only)
- **Project Removal**: Remove projects with uncommitted change protection
- **Git Worktrees**: Create, open, checkout, and remove git worktrees with automatic session launching
- **Config Management**: Open your config file directly in your editor

## Installation

### Prerequisites

- Go 1.24.2 or later
- `tmux` for tmux session provider (default)
- `code` for VS Code session provider
- `gh` for remote repository management
- `git` for repository operations

These can be checked using the `health` command mentioned below

### Install from source

```console
go install github.com/danielronalds/projman@latest
```

### Build locally

```console
git clone https://github.com/danielronalds/projman.git
cd projman
go build -o projman .
```

## Usage

```console
projman v0.8.0

Usage: projman [command | project]

A cli tool for managing projects on your local machine.

To select and open a project, run projman with no arguments

Commands
  new       Create a new project
  local     Open a project currently on your machine
  remote    Open a project from github
  clone     Clone any git url to a project dir
  active    Open an existing session
  config    Open the projman config in your editor
  rm        Remove a project from your machine
  worktree  Manage git worktrees -- run 'projman wt help' for details (alias: wt)
  health    Verify required dependencies are installed
  help      Show this menu
```

### Global Flags

`--provider, -p <provider>` overrides the configured session provider for a single invocation:

```console
projman --provider=vscode local
projman --provider tmux remote
projman -p tmux local
```

## Configuration

projman uses a JSON configuration file located at `~/.config/projman/config.json`. If the file doesn't exist, projman will use default settings.

### Default Configuration

```json
{
  "theme": "default",
  "layout": "reverse",
  "projectDirs": ["Projects/"],
  "openNewProjects": true,
  "templates": [],
  "session_provider": "tmux",
  "tmux": {
    "windows": ["CLI", "Code", "Server"],
    "starting_window": 2
  }
}
```

### Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `theme` | string | `"default"` | Selector theme preset (`default`, `bw`, `minimal`) |
| `layout` | string | `"reverse"` | Selector input position (`reverse` = top, `default` = bottom) |
| `projectDirs` | array | `["Projects/"]` | Directories to search for projects (relative to home directory) |
| `openNewProjects` | boolean | `true` | Whether to automatically open new projects |
| `templates` | array | `[]` | Project templates with commands to run |
| `session_provider` | string | `"tmux"` | Session provider to use (`tmux` or `vscode`) |
| `tmux.windows` | array | `["CLI", "Code", "Server"]` | Names of tmux windows to create |
| `tmux.starting_window` | number | `2` | Which window to start in |

## Session Providers

projman supports multiple session providers for opening projects:

### tmux (default)

The tmux provider creates a new tmux session for each project with configurable windows. It supports:
- Creating sessions with multiple named windows
- Switching between active sessions via `projman active`
- Attaching to sessions from outside tmux or switching from within

Configure via the `tmux` block:

```json
{
  "session_provider": "tmux",
  "tmux": {
    "windows": ["Terminal", "Editor", "Server", "Tests"],
    "starting_window": 1
  }
}
```

### vscode

The VS Code provider opens projects directly in VS Code using the `code` CLI command. 

**Note**: The `projman active` command is not supported with the vscode provider.

```json
{
  "session_provider": "vscode"
}
```

### Example Configurations

#### tmux with custom windows

```json
{
  "theme": "bw",
  "layout": "reverse",
  "projectDirs": ["Projects/", "Work/", "Personal/"],
  "openNewProjects": true,
  "session_provider": "tmux",
  "tmux": {
    "windows": ["Terminal", "Editor", "Server", "Tests"],
    "starting_window": 1
  },
  "templates": [
    {
      "name": "go-project",
      "commands": [
        "go mod init",
        "touch main.go",
        "git init"
      ]
    }
  ]
}
```

#### VS Code

```json
{
  "theme": "minimal",
  "layout": "reverse",
  "projectDirs": ["Projects/"],
  "openNewProjects": true,
  "session_provider": "vscode"
}
```

## Project Templates

Templates allow you to run a series of commands when creating new projects. Define templates in your configuration file with a name and an array of commands to execute.

Template commands have access to the `PROJMAN_PROJECT_NAME` environment variable, which contains the name of the project being created.

## Editing Configuration

`projman config` opens the configuration file in your editor. It uses the `$EDITOR` environment variable, falling back to `nvim` if unset.

## Removing Projects

`projman rm` removes a project directory from your machine. It checks for uncommitted git changes before deleting:

```console
projman rm              # fuzzy select a project to remove
projman rm my-project   # remove a named project directly
```

If the project has uncommitted changes, projman will refuse to delete it. Use `--without-git-check` to bypass this:

```console
projman rm my-project --without-git-check
```

## Git Worktrees

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

### Creating a worktree

`projman wt new <name>` creates a new git worktree with a new branch and launches a session in it:

```console
projman wt new feature/auth
```

### Checking out a remote branch

`projman wt checkout` fetches all remote branches, presents them in a fuzzy finder, and creates a worktree from the selected branch:

```console
projman wt checkout
```

### Opening a worktree

`projman wt` with no arguments lists existing worktrees in a fuzzy finder. You can also pass a name directly:

```console
projman wt                  # fuzzy select a worktree to open
projman wt my-worktree      # open a named worktree directly
```

### Removing a worktree

`projman wt rm` lists worktrees (excluding the base) in a fuzzy finder and removes the selected one:

```console
projman wt rm
```

### Ignored files

When creating or checking out a worktree, projman offers to copy gitignored files (e.g. `.env`, `node_modules/`) from the main worktree into the new one.

### Naming conventions

Worktree directories are named `<project>-<sanitised-branch>`, where slashes and spaces in the branch name become dashes. For example, `feature/auth` in a project called `myapp` produces a directory called `myapp-feature-auth`. Sessions follow the same naming pattern.

## Upgrading from older versions

If you're upgrading from an older version that used `session_layout`, projman will print an error and exit. Update your config to use the new `tmux` block:

**Old format (deprecated):**
```json
{
  "session_layout": {
    "windows": ["CLI", "Code", "Server"],
    "starting_window": 2
  }
}
```

**New format:**
```json
{
  "session_provider": "tmux",
  "tmux": {
    "windows": ["CLI", "Code", "Server"],
    "starting_window": 2
  }
}
```

## Development

This project uses [mise](https://mise.jdx.dev) for tool management and task running.

### Setup

```console
mise install
```

### Tasks

| Task | Description |
|------|-------------|
| `mise dev` | Run the project |
| `mise build` | Build the binary |
| `mise build:install` | Build and install to `~/.local/bin/` |
| `mise test` | Run all tests |
| `mise lint:fmt` | Check for unformatted code |
| `mise lint:vet` | Run go vet |
| `mise lint:static` | Run golangci-lint |
| `mise lint:deps` | Check for unused or missing dependencies |
| `mise clean` | Remove built binary |

### CI

A GitHub Actions pipeline runs on pushes to `main` and on pull requests. It runs formatting, vet, static analysis, tests, and dependency checks.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
