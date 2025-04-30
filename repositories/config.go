package repositories

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type template struct {
	Name     string   `json:"name"`
	Commands []string `json:"commands"`
}

func (t template) GetCommands() []string {
	return t.Commands
}

type config struct {
	Theme           string     `json:"theme"`
	Layout          string     `json:"layout"`
	ProjectDirs     []string   `json:"projectDir"`
	OpenNewProjects bool       `json:"openNewProjects"`
	Templates       []template `json:"templates"`
}

type ConfigRepository struct {
	conf config
}

func NewConfigRepository() ConfigRepository {
	conf := config{
		Theme:           "bw",
		Layout:          "reverse",
		ProjectDirs:     []string{"Projects/"},
		OpenNewProjects: true,
		Templates: make([]template, 0),
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

	return ConfigRepository{conf}
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

// Helper function that gets the home dir without an err. If an err occurs, the program exits
func getHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get users home dir: %v\n", err.Error())
		os.Exit(1)
	}
	return homeDir
}
