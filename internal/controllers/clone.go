package controllers

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type remoteURLCloner interface {
	CloneURL(url, dir string) (projectPath, error)
}

type CloneController struct {
	remote remoteURLCloner
	fzf    selecter
	tmux   sessionLauncher
	config remoteConfig
}

func NewCloneController(remote remoteURLCloner, fzf selecter, tmux sessionLauncher, config remoteConfig) CloneController {
	return CloneController{remote, fzf, tmux, config}
}

func (c CloneController) HandleArgs(args []string) error {
	if len(args) < 2 {
		return errors.New("git url required")
	}

	repoTarget := args[1]
	repoName, err := parseRepoName(repoTarget)
	if err != nil {
		return err
	}

	cloneDir := c.config.ProjectDirs()[0]
	if len(c.config.ProjectDirs()) > 1 {
		cloneDir, err = c.fzf.Select(c.config.ProjectDirs())
		if err != nil {
			return errors.New("no clone dir selected")
		}
	}

	projPath, err := c.remote.CloneURL(repoTarget, cloneDir)
	if err != nil {
		return fmt.Errorf("unable to clone project: %v", err.Error())
	}

	return c.tmux.LaunchSession(repoName, projPath)
}

var repoNameRegex = regexp.MustCompile(`([\w.-]+)(?:\.git)?$`)

func parseRepoName(target string) (string, error) {
	normalized := strings.TrimSpace(target)

	if idx := strings.LastIndex(normalized, ":"); idx != -1 && strings.Contains(normalized[:idx], "@") {
		normalized = normalized[idx+1:]
	}

	normalized = strings.TrimSuffix(normalized, "/")

	matches := repoNameRegex.FindStringSubmatch(normalized)
	if len(matches) < 2 {
		return "", errors.New("invalid git url")
	}

	name := strings.TrimSuffix(matches[1], ".git")
	name = strings.TrimSpace(name)
	if name == "" {
		return "", errors.New("invalid git url")
	}

	return name, nil
}
