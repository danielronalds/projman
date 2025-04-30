package repositories

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type config struct {
	Theme      string `json:"theme"`
	Layout     string `json:"layout"`
	ProjectDir string `json:"projectDir"`
}

type ConfigRepository struct {
	conf config
}

func NewConfigRepository() ConfigRepository {
	conf := config{
		Theme:      "bw",
		Layout:     "reverse",
		ProjectDir: "Projects/",
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

func (r ConfigRepository) ProjectDir() string {
	homeDir := getHomeDir()

	normedDir := strings.TrimSuffix(r.conf.ProjectDir, "/")

	return fmt.Sprintf("%v/%v/", homeDir, normedDir)
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
