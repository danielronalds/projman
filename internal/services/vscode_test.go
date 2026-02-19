package services

import (
	"testing"
)

func TestVSCodeProvider_Name(t *testing.T) {
	provider := NewVSCodeProvider()

	if provider.Name() != "vscode" {
		t.Errorf("expected Name() to return 'vscode', got '%s'", provider.Name())
	}
}

func TestVSCodeProvider_ListActiveSessions_ReturnsError(t *testing.T) {
	provider := NewVSCodeProvider()

	sessions, err := provider.ListActiveSessions()

	if err == nil {
		t.Error("expected ListActiveSessions() to return error, got nil")
	}

	if sessions != nil {
		t.Errorf("expected sessions to be nil, got %v", sessions)
	}
}

func TestVSCodeProvider_OpenActiveSession_ReturnsError(t *testing.T) {
	provider := NewVSCodeProvider()

	err := provider.OpenActiveSession("test-session")

	if err == nil {
		t.Error("expected OpenActiveSession() to return error, got nil")
	}
}
