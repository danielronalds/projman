# Configuration

projman uses a JSON configuration file located at `~/.config/projman/config.json`. If the file doesn't exist, projman uses default settings.

`projman config` opens the configuration file in your editor. It uses the `$EDITOR` environment variable, falling back to `nvim` if unset.

## Default configuration

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

## Options

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
| `worktree_copy_excludes` | array | `[]` | Glob patterns to skip when copying gitignored files into a new worktree (see [Worktrees](worktrees.md#skipping-large-directories)) |

## Session providers

projman supports multiple session providers for opening projects. Override the configured provider for a single invocation with `--provider` / `-p`:

```console
projman --provider=vscode local
projman -p tmux remote
```

### tmux (default)

Creates a new tmux session for each project with configurable windows. Supports:

- Creating sessions with multiple named windows
- Switching between active sessions via `projman active`
- Attaching to sessions from outside tmux or switching from within

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

Opens projects directly in VS Code using the `code` CLI command.

`projman active` is not supported with the vscode provider.

```json
{
  "session_provider": "vscode"
}
```

## Project templates

Templates run a series of commands when creating new projects. Define them in your config with a name and an array of commands.

Template commands have access to the `PROJMAN_PROJECT_NAME` environment variable, which contains the name of the project being created.

```json
{
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

## Upgrading from older versions

If you're upgrading from a version that used `session_layout`, projman will print an error and exit. Update your config to use the new `tmux` block:

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
