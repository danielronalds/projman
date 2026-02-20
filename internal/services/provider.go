package services

import "fmt"

type SessionProvider interface {
	LaunchSession(name, dir string) error
	ListActiveSessions() ([]string, error)
	OpenActiveSession(name string) error
	Name() string
}

type ProviderConfig interface {
	SessionProvider() string
	TmuxConfig() TmuxConfig
}

type TmuxConfig struct {
	Windows        []string
	StartingWindow int
}

type ProviderConfigOverride struct {
	ProviderConfig
	provider string
}

func NewProviderConfigOverride(base ProviderConfig, provider string) ProviderConfig {
	return ProviderConfigOverride{base, provider}
}

func (p ProviderConfigOverride) SessionProvider() string {
	return p.provider
}

func NewSessionProvider(cfg ProviderConfig) (SessionProvider, error) {
	switch cfg.SessionProvider() {
	case "tmux":
		return NewTmuxProvider(cfg.TmuxConfig()), nil
	case "vscode":
		return NewVSCodeProvider(), nil
	default:
		return nil, fmt.Errorf("unknown session provider: %s", cfg.SessionProvider())
	}
}
