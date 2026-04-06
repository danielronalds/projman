package controllers

import (
	"errors"
	"strings"
	"testing"
)

type mockNotesProjectResolver struct {
	returnName string
	returnPath string
	returnErr  error
}

func (m *mockNotesProjectResolver) ResolveCurrentProject() (string, string, error) {
	return m.returnName, m.returnPath, m.returnErr
}

type mockNotesPathProvider struct {
	returnPath     string
	returnErr      error
	calledProjPath string
}

func (m *mockNotesPathProvider) NotePath(projPath string) (string, error) {
	m.calledProjPath = projPath
	return m.returnPath, m.returnErr
}

type mockNotesProjectLister struct {
	returnProjects []string
	returnListErr  error
	returnPath     string
	returnGetErr   error
	calledProject  string
}

func (m *mockNotesProjectLister) ListProjects() ([]string, error) {
	return m.returnProjects, m.returnListErr
}

func (m *mockNotesProjectLister) GetPath(project string) (string, error) {
	m.calledProject = project
	return m.returnPath, m.returnGetErr
}

type mockNotesSelecter struct {
	returnSelected string
	returnErr      error
	calledOptions  []string
}

func (m *mockNotesSelecter) Select(options []string) (string, error) {
	m.calledOptions = options
	return m.returnSelected, m.returnErr
}

func TestNotesControllerHandleArgs(t *testing.T) {
	t.Run("insideProjectDir", func(t *testing.T) {
		t.Setenv("EDITOR", "true")

		resolver := &mockNotesProjectResolver{
			returnName: "projman",
			returnPath: "/Users/dan/Personal/projman",
		}
		notes := &mockNotesPathProvider{returnPath: "/data/notes/personal_projman.md"}
		projects := &mockNotesProjectLister{}
		fzf := &mockNotesSelecter{}

		controller := NewNotesController(resolver, notes, projects, fzf)
		err := controller.HandleArgs(nil)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if notes.calledProjPath != "/Users/dan/Personal/projman" {
			t.Fatalf("NotePath called with %q, want %q", notes.calledProjPath, "/Users/dan/Personal/projman")
		}
		if fzf.calledOptions != nil {
			t.Fatalf("Select should not be called when inside a project dir")
		}
	})

	t.Run("outsideProjectDir", func(t *testing.T) {
		t.Setenv("EDITOR", "true")

		resolver := &mockNotesProjectResolver{
			returnErr: errors.New("not in a project"),
		}
		notes := &mockNotesPathProvider{returnPath: "/data/notes/personal_projman.md"}
		projects := &mockNotesProjectLister{
			returnProjects: []string{"projman", "myapp"},
			returnPath:     "/Users/dan/Personal/projman",
		}
		fzf := &mockNotesSelecter{returnSelected: "projman"}

		controller := NewNotesController(resolver, notes, projects, fzf)
		err := controller.HandleArgs(nil)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(fzf.calledOptions) != 2 {
			t.Fatalf("Select called with %d options, want 2", len(fzf.calledOptions))
		}
		if projects.calledProject != "projman" {
			t.Fatalf("GetPath called with %q, want %q", projects.calledProject, "projman")
		}
		if notes.calledProjPath != "/Users/dan/Personal/projman" {
			t.Fatalf("NotePath called with %q, want %q", notes.calledProjPath, "/Users/dan/Personal/projman")
		}
	})

	t.Run("fuzzySelectCancelled", func(t *testing.T) {
		resolver := &mockNotesProjectResolver{
			returnErr: errors.New("not in a project"),
		}
		notes := &mockNotesPathProvider{}
		projects := &mockNotesProjectLister{
			returnProjects: []string{"projman"},
		}
		fzf := &mockNotesSelecter{returnErr: errors.New("cancelled")}

		controller := NewNotesController(resolver, notes, projects, fzf)
		err := controller.HandleArgs(nil)

		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "selecting project") {
			t.Fatalf("expected 'selecting project' error, got: %v", err)
		}
	})

	t.Run("noLocalProjects", func(t *testing.T) {
		resolver := &mockNotesProjectResolver{
			returnErr: errors.New("not in a project"),
		}
		notes := &mockNotesPathProvider{}
		projects := &mockNotesProjectLister{
			returnProjects: []string{},
		}
		fzf := &mockNotesSelecter{}

		controller := NewNotesController(resolver, notes, projects, fzf)
		err := controller.HandleArgs(nil)

		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "no local projects found") {
			t.Fatalf("expected 'no local projects found' error, got: %v", err)
		}
		if fzf.calledOptions != nil {
			t.Fatalf("Select should not be called when there are no projects")
		}
	})

	t.Run("notePathError", func(t *testing.T) {
		resolver := &mockNotesProjectResolver{
			returnName: "projman",
			returnPath: "/Users/dan/Personal/projman",
		}
		notes := &mockNotesPathProvider{returnErr: errors.New("notes dir failed")}
		projects := &mockNotesProjectLister{}
		fzf := &mockNotesSelecter{}

		controller := NewNotesController(resolver, notes, projects, fzf)
		err := controller.HandleArgs(nil)

		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "notes dir failed") {
			t.Fatalf("expected error to contain 'notes dir failed', got: %v", err)
		}
	})
}
