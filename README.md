# projman

A CLI tool for managing projects on your local machine.
projman helps you quickly create, open, and manage your
dev projects with integrated tmux session management.

## Features

- **Project Management**: Create and open projects from local directories or GitHub repositories
- **Tmux Integration**: Automatically manage tmux sessions for your projects
- **Project Templates**: Use customizable templates when creating new projects
- **FZF Integration**: Interactive project selection using fuzzy finding
- **GitHub Integration**: Browse, clone, and open remote repositories
- **Session Management**: Manage active tmux sessions

## Installation

### Prerequisites

- Go 1.24.1 or later
- `fzf` for interactive selection
- `tmux` for session management
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
  active    Open an existing session
  health    Verify required dependencies are installed
  help      Show this menu
```

## Configuration

projman uses a JSON configuration file located at `~/.config/projman/config.json`. If the file doesn't exist, projman will use default settings.

### Default Configuration

```json
{
  "theme": "bw",
  "layout": "reverse",
  "projectDirs": ["Projects/"],
  "openNewProjects": true,
  "templates": [],
  "session_layout": {
    "windows": ["CLI", "Code", "Server"],
    "starting_window": 2
  }
}
```

### Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `theme` | string | `"bw"` | FZF theme |
| `layout` | string | `"reverse"` | FZF layout |
| `projectDirs` | array | `["Projects/"]` | Directories to search for projects (relative to home directory) |
| `openNewProjects` | boolean | `true` | Whether to automatically open new projects in tmux |
| `templates` | array | `[]` | Project templates with commands to run |
| `session_layout.windows` | array | `["CLI", "Code", "Server"]` | Names of tmux windows to create |
| `session_layout.starting_window` | number | `2` | Which window to start in |

### Example Configuration

```json
{
  "theme": "bw",
  "layout": "reverse",
  "projectDirs": ["Projects/", "Work/", "Personal/"],
  "openNewProjects": true,
  "templates": [
    {
      "name": "go-project",
      "commands": [
        "go mod init",
        "touch main.go",
        "git init"
      ]
    },
    {
      "name": "web-app",
      "commands": [
        "npm init -y",
        "npm install express",
        "touch server.js",
        "git init"
      ]
    }
  ],
  "session_layout": {
    "windows": ["Terminal", "Editor", "Server", "Tests"],
    "starting_window": 1
  }
}
```

## Project Templates

Templates allow you to run a series of commands when creating new projects. Define templates in your configuration file with a name and an array of commands to execute.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
