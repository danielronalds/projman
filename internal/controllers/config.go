package controllers

import (
	"os"
	"os/exec"
)

type configPathProvider interface {
	ConfigFilePath() string
}

type ConfigController struct {
	configRepo configPathProvider
}

func NewConfigController(configRepo configPathProvider) ConfigController {
	return ConfigController{configRepo: configRepo}
}

func (c ConfigController) HandleArgs(args []string) error {
	editor := getEditor()
	configPath := c.configRepo.ConfigFilePath()

	cmd := exec.Command(editor, configPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func getEditor() string {
	if editor, ok := os.LookupEnv("EDITOR"); ok {
		return editor
	}
	return "nvim"
}
