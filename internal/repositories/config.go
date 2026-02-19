package repositories

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/danielronalds/projman/internal/services"
)

type template struct {
	Name     string   `json:"name"`
	Commands []string `json:"commands"`
}

type tmuxConfig struct {
	Windows        []string `json:"windows"`
	StartingWindow int      `json:"starting_window"`
}

type vsCodeConfig struct{}

type sessionLayout struct {
	Windows        []string `json:"windows"`
	StartingWindow int      `json:"starting_window"`
}

func (t template) GetCommands() []string {
	return t.Commands
}

type config struct {
	Theme           string        `json:"theme"`
	Layout          string        `json:"layout"`
	ProjectDirs     []string      `json:"projectDirs"`
	OpenNewProjects bool          `json:"openNewProjects"`
	Templates       []template    `json:"templates"`
	SessionProvider string        `json:"session_provider"`
	Tmux            tmuxConfig    `json:"tmux"`
	VSCode          vsCodeConfig  `json:"vscode"`
	SessionLayout   sessionLayout `json:"session_layout,omitempty"`
}

type ConfigRepository struct {
	conf config
}

func NewConfigRepository() ConfigRepository {
	conf := config{
		Theme:           "default",
		Layout:          "reverse",
		ProjectDirs:     []string{"Projects/"},
		OpenNewProjects: true,
		Templates:       make([]template, 0),
		SessionProvider: "tmux",
		Tmux: tmuxConfig{
			Windows:        []string{"CLI", "Code", "Server"},
			StartingWindow: 2,
		},
		VSCode: vsCodeConfig{},
	}

	homeDir := getHomeDir()
	configFile := fmt.Sprintf("%v/.config/projman/config.json", homeDir)

	configContents, err := os.ReadFile(configFile)
	if err != nil {
		return ConfigRepository{conf}
	}

	if err = json.Unmarshal(configContents, &conf); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read config file: %v\n", err.Error())
		os.Exit(2)
	}

	if hasLegacySessionLayout(&conf) {
		fmt.Fprintf(os.Stderr, "Error: 'session_layout' is deprecated and no longer supported.\n")
		fmt.Fprintf(os.Stderr, "Please move your config to the 'tmux' block. See README.md for details.\n")
		os.Exit(1)
	}

	if conf.SessionProvider == "" {
		conf.SessionProvider = "tmux"
	}

	return ConfigRepository{conf}
}

func hasLegacySessionLayout(conf *config) bool {
	return len(conf.SessionLayout.Windows) > 0
}

func (r ConfigRepository) Theme() string {
	return r.conf.Theme
}

func (r ConfigRepository) Layout() string {
	return r.conf.Layout
}

func (r ConfigRepository) ProjectDirs() []string {
	homeDir := getHomeDir()

	normedDirs := make([]string, 0)
	for _, dir := range r.conf.ProjectDirs {
		normedDir := strings.TrimSuffix(dir, "/")
		normedDirs = append(normedDirs, fmt.Sprintf("%v/%v/", homeDir, normedDir))
	}

	return normedDirs
}

func (r ConfigRepository) OpenNewProjects() bool {
	return r.conf.OpenNewProjects
}

func (r ConfigRepository) TemplateNames() []string {
	templates := make([]string, 0)
	for _, tmpl := range r.conf.Templates {
		templates = append(templates, tmpl.Name)
	}
	return templates
}

func (r ConfigRepository) GetTemplateCommands(tmpl string) ([]string, error) {
	for _, t := range r.conf.Templates {
		if t.Name == tmpl {
			return t.Commands, nil
		}
	}
	return make([]string, 0), errors.New("no template with that name exists")
}

func (r ConfigRepository) SessionProvider() string {
	return r.conf.SessionProvider
}

func (r ConfigRepository) TmuxConfig() services.TmuxConfig {
	return services.TmuxConfig{
		Windows:        r.conf.Tmux.Windows,
		StartingWindow: r.conf.Tmux.StartingWindow,
	}
}

func (r ConfigRepository) SessionWindows() []string {
	return r.conf.Tmux.Windows
}

func (r ConfigRepository) StartingWindow() int {
	return r.conf.Tmux.StartingWindow
}

func (r ConfigRepository) ConfigFilePath() string {
	homeDir := getHomeDir()
	return fmt.Sprintf("%v/.config/projman/config.json", homeDir)
}

func getHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get users home dir: %v\n", err.Error())
		os.Exit(1)
	}
	return homeDir
}
