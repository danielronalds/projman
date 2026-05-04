package controllers

import (
	"fmt"
	"strings"
)

type llmConfigProvider interface {
	ProjectDirs() []string
	SessionProvider() string
	TemplateNames() []string
}

type LlmController struct {
	config llmConfigProvider
}

func NewLlmController(config llmConfigProvider) LlmController {
	return LlmController{config: config}
}

func (c LlmController) HandleArgs(args []string) error {
	var b strings.Builder

	b.WriteString("# projman\n\n")
	b.WriteString("projman is a CLI tool for managing development projects on a local machine. ")
	b.WriteString("It discovers projects from configured directories, and can create, open, clone, ")
	b.WriteString("and remove them. Projects are opened in sessions via a configurable session provider ")
	b.WriteString("(tmux or VS Code).\n\n")

	writeCurrentSetup(&b, c.config)
	writeCommandReference(&b)
	writeAgentGuidance(&b)

	fmt.Print(b.String())
	return nil
}

func writeCurrentSetup(b *strings.Builder, config llmConfigProvider) {
	b.WriteString("## Current setup\n\n")

	b.WriteString("### Project directories\n\n")
	b.WriteString("Projects are discovered from the following directories:\n\n")
	for _, dir := range config.ProjectDirs() {
		fmt.Fprintf(b, "- `%s`\n", dir)
	}

	fmt.Fprintf(b, "\n### Session provider\n\n")
	fmt.Fprintf(b, "The active session provider is `%s`.\n\n", config.SessionProvider())

	templateNames := config.TemplateNames()
	if len(templateNames) > 0 {
		b.WriteString("### Templates\n\n")
		b.WriteString("Available project templates:\n\n")
		for _, name := range templateNames {
			fmt.Fprintf(b, "- `%s`\n", name)
		}
		b.WriteString("\n")
	}
}

func writeCommandReference(b *strings.Builder) {
	b.WriteString("## Commands\n\n")

	commands := []struct {
		name        string
		usage       string
		description string
	}{
		{
			name:        "local",
			usage:       "projman local",
			description: "Opens an interactive fuzzy finder to select and open a local project. This is also the default when projman is run with no arguments.",
		},
		{
			name:        "<project>",
			usage:       "projman my-project",
			description: "Opens a project by name directly, without the fuzzy finder. Use this when you know the exact project name.",
		},
		{
			name:        "new",
			usage:       "projman new",
			description: "Creates a new project. Prompts for a name, directory, and optional template. The template runs a series of configured shell commands in the new project directory. The PROJMAN_PROJECT_NAME env var is available to template commands.",
		},
		{
			name:        "remote",
			usage:       "projman remote",
			description: "Lists GitHub repositories (via the gh CLI), lets the user select one, clones it locally, and opens a session. Repositories already cloned locally are filtered out.",
		},
		{
			name:        "clone",
			usage:       "projman clone <git-url>",
			description: "Clones any git URL into a project directory and opens a session. Accepts HTTPS and SSH URLs.",
		},
		{
			name:        "active",
			usage:       "projman active",
			description: "Lists existing tmux sessions and lets the user switch to one. Only available with the tmux session provider.",
		},
		{
			name:        "here",
			usage:       "projman here",
			description: "Opens a session in the current working directory. Useful when you are already in a project directory and want to launch a session without navigating the project list.",
		},
		{
			name:        "list",
			usage:       "projman list [filter]",
			description: "Lists all local projects grouped by directory. Accepts an optional filter string to narrow results by directory name.",
		},
		{
			name:        "notes",
			usage:       "projman notes [project]",
			description: "Opens a markdown notes file for a project in your editor. If inside an active session, it resolves the current project automatically. Otherwise prompts for selection.",
		},
		{
			name:        "config",
			usage:       "projman config",
			description: "Opens the projman config file (~/.config/projman/config.json) in the user's $EDITOR.",
		},
		{
			name:        "rm",
			usage:       "projman rm [project] [--without-git-check]",
			description: "Removes a project directory. By default it checks for uncommitted git changes and refuses to delete if any exist. Pass --without-git-check to bypass this safety check.",
		},
		{
			name:        "worktree (wt)",
			usage:       "projman wt [subcommand | name]",
			description: "Manages git worktrees for the current project. Subcommands: `new <name>` creates a worktree and branch, `checkout` fetches remote branches to create a worktree from, `rm` removes a worktree. With no arguments, lists and opens existing worktrees.",
		},
		{
			name:        "health",
			usage:       "projman health",
			description: "Checks that required dependencies (tmux, code, gh, git) are installed and reports their status.",
		},
		{
			name:        "help",
			usage:       "projman help",
			description: "Displays the help menu with a summary of all commands.",
		},
	}

	for _, cmd := range commands {
		fmt.Fprintf(b, "### %s\n\n", cmd.name)
		fmt.Fprintf(b, "%s\n\n", cmd.description)
		fmt.Fprintf(b, "```\n%s\n```\n\n", cmd.usage)
	}
}

func writeAgentGuidance(b *strings.Builder) {
	b.WriteString("## Agent usage guide\n\n")
	b.WriteString("When using projman as an AI agent, prefer non-interactive commands. ")
	b.WriteString("Several commands launch a fuzzy finder which requires user interaction and is not suitable for automated use.\n\n")

	b.WriteString("### Non-interactive commands (safe for agents)\n\n")
	b.WriteString("- `projman list` -- discover available projects\n")
	b.WriteString("- `projman <project>` -- open a known project by name\n")
	b.WriteString("- `projman here` -- open a session in the current directory\n")
	b.WriteString("- `projman clone <url>` -- clone a repository\n")
	b.WriteString("- `projman health` -- check dependencies\n")
	b.WriteString("- `projman wt new <name>` -- create a worktree by name\n\n")

	b.WriteString("### Interactive commands (require user input)\n\n")
	b.WriteString("- `projman local` / `projman` -- fuzzy finder selection\n")
	b.WriteString("- `projman remote` -- fuzzy finder selection\n")
	b.WriteString("- `projman active` -- fuzzy finder selection\n")
	b.WriteString("- `projman new` -- prompts for name, directory, and template\n")
	b.WriteString("- `projman rm` -- fuzzy finder selection (without a project name argument)\n")
	b.WriteString("- `projman notes` -- fuzzy finder selection (without a project name argument)\n")
	b.WriteString("- `projman wt` -- fuzzy finder selection (without a subcommand)\n")
	b.WriteString("- `projman wt checkout` -- fuzzy finder selection\n")
	b.WriteString("- `projman wt rm` -- fuzzy finder selection\n\n")

	b.WriteString("### Recommended workflow\n\n")
	b.WriteString("1. Run `projman list` to see what projects exist and where they are located\n")
	b.WriteString("2. Open a project by name with `projman <project>` if a session is needed\n")
	b.WriteString("3. Use `projman here` when already in the correct directory\n")
	b.WriteString("4. Use `projman health` to verify the environment is set up correctly\n")
	b.WriteString("5. Use `projman clone <url>` to bring in external repositories\n")

	b.WriteString("\n### Global flags\n\n")
	b.WriteString("- `--provider <tmux|vscode>` (or `-p`) -- override the session provider for a single invocation\n")
}
