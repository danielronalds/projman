package controllers

import (
	"testing"
)

type mockHealthChecker struct {
	results map[string]bool
}

func (m mockHealthChecker) CheckRequirements(programs []string) map[string]bool {
	return m.results
}

func TestHealthController_ChecksAllProviderDependencies(t *testing.T) {
	checker := mockHealthChecker{
		results: map[string]bool{
			"tmux": true,
			"code": true,
			"gh":   true,
			"git":  true,
		},
	}

	controller := NewHealthController(checker)
	err := controller.HandleArgs([]string{})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestHealthController_ReturnsErrorWhenDependencyMissing(t *testing.T) {
	checker := mockHealthChecker{
		results: map[string]bool{
			"tmux": true,
			"code": false,
			"gh":   true,
			"git":  true,
		},
	}

	controller := NewHealthController(checker)
	err := controller.HandleArgs([]string{})

	if err == nil {
		t.Fatal("expected error for missing dependency, got nil")
	}
}
