package services

import (
	"testing"
)

type mockProviderConfig struct {
	provider   string
	tmuxConfig TmuxConfig
}

func (m mockProviderConfig) SessionProvider() string {
	return m.provider
}

func (m mockProviderConfig) TmuxConfig() TmuxConfig {
	return m.tmuxConfig
}

func TestNewSessionProvider_Tmux(t *testing.T) {
	cfg := mockProviderConfig{
		provider: "tmux",
		tmuxConfig: TmuxConfig{
			Windows:        []string{"CLI", "Code"},
			StartingWindow: 1,
		},
	}

	provider, err := NewSessionProvider(cfg)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if provider == nil {
		t.Fatal("expected provider to be non-nil")
	}

	if provider.Name() != "tmux" {
		t.Errorf("expected provider name 'tmux', got '%s'", provider.Name())
	}
}

func TestNewSessionProvider_VSCode(t *testing.T) {
	cfg := mockProviderConfig{
		provider: "vscode",
	}

	provider, err := NewSessionProvider(cfg)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if provider == nil {
		t.Fatal("expected provider to be non-nil")
	}

	if provider.Name() != "vscode" {
		t.Errorf("expected provider name 'vscode', got '%s'", provider.Name())
	}
}

func TestNewSessionProvider_UnknownProvider(t *testing.T) {
	cfg := mockProviderConfig{
		provider: "unknown",
	}

	provider, err := NewSessionProvider(cfg)

	if err == nil {
		t.Error("expected error for unknown provider, got nil")
	}

	if provider != nil {
		t.Errorf("expected provider to be nil, got %v", provider)
	}
}

func TestTmuxProvider_Name(t *testing.T) {
	provider := NewTmuxProvider(TmuxConfig{
		Windows:        []string{"CLI"},
		StartingWindow: 1,
	})

	if provider.Name() != "tmux" {
		t.Errorf("expected Name() to return 'tmux', got '%s'", provider.Name())
	}
}
