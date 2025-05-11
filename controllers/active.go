package controllers

import (
	"errors"
	"fmt"
)

type sessionLister interface {
	ListActiveSessions() ([]string, error)
}

type sessionOpener interface {
	OpenActiveSession(name string) error
}

type sessionManager interface {
	sessionLister
	sessionOpener
}

type ActiveController struct {
	projects projectPathFinder
	fzf      selecter
	tmux     sessionManager
}

func NewActiveController(projects projectPathFinder, fzf selecter, tmux sessionManager) ActiveController {
	return ActiveController{projects, fzf, tmux}
}

func (c ActiveController) HandleArgs(args []string) error {
	sessions, err := c.tmux.ListActiveSessions()
	if err != nil {
		return fmt.Errorf("unable to list active sessions: %v", err.Error())
	}

	if len(sessions) == 0 {
		fmt.Println("no active sessions")
		return nil
	}

	session, err := c.fzf.Select(sessions)
	if err != nil {
		// if an error occurs, we assume fzf was Ctrl+c
		return errors.New("no session selected")
	}

	return c.tmux.OpenActiveSession(session)
}
