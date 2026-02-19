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

## Installation

### Prerequisites

- Go 1.24.1 or later
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
projman v0.3.0

Usage: projman [command | project]

A cli tool for managing projects on your local machine.

To select and open a project, run projman with no arguments

Commands
  new       Create a new project
  local     Open a project currently on your machine
  remote    Open a project from github
  clone     Clone any git URL into a project dir
  active    Open an existing session (tmux only)
  health    Verify required dependencies are installed
  help      Show this menu
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

## Upgrading from older versions

If you're upgrading from an older version that used `session_layout`, you'll see a deprecation warning. Update your config to use the new `tmux` block:

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

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
