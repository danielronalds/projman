package services

import "testing"

func TestSanitiser_Sanitise(t *testing.T) {
	sanitiser := NewSanitiser()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"normal name passes through", "myproject", "myproject"},
		{"colons replaced with hyphens", "my:project", "my-project"},
		{"root directory becomes root", "/", "root"},
		{"leading dot stripped", ".hidden", "hidden"},
		{"empty string becomes session", "", "session"},
		{"multiple bad chars", ".:bad:name", "-bad-name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitiser.Sanitise(tt.input)
			if result != tt.expected {
				t.Fatalf("Sanitise(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
