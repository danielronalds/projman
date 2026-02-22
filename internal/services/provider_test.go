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

func TestProviderConfigOverride_SessionProvider(t *testing.T) {
	base := mockProviderConfig{
		provider: "tmux",
	}

	override := NewProviderConfigOverride(base, "vscode")

	if override.SessionProvider() != "vscode" {
		t.Errorf("expected SessionProvider() to return 'vscode', got '%s'", override.SessionProvider())
	}
}

func TestProviderConfigOverride_TmuxConfig(t *testing.T) {
	expectedConfig := TmuxConfig{
		Windows:        []string{"CLI", "Code"},
		StartingWindow: 1,
	}

	base := mockProviderConfig{
		provider:   "tmux",
		tmuxConfig: expectedConfig,
	}

	override := NewProviderConfigOverride(base, "vscode")

	result := override.TmuxConfig()

	if len(result.Windows) != len(expectedConfig.Windows) {
		t.Fatalf("expected %d windows, got %d", len(expectedConfig.Windows), len(result.Windows))
	}

	for i, window := range result.Windows {
		if window != expectedConfig.Windows[i] {
			t.Errorf("expected window %d to be '%s', got '%s'", i, expectedConfig.Windows[i], window)
		}
	}

	if result.StartingWindow != expectedConfig.StartingWindow {
		t.Errorf("expected StartingWindow to be %d, got %d", expectedConfig.StartingWindow, result.StartingWindow)
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
