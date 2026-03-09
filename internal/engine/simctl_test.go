package engine

import "testing"

func TestParseRuntime(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"com.apple.CoreSimulator.SimRuntime.iOS-18-2", "iOS 18.2"},
		{"com.apple.CoreSimulator.SimRuntime.iOS-17-0", "iOS 17.0"},
		{"com.apple.CoreSimulator.SimRuntime.watchOS-11-1", "watchOS 11.1"},
		{"unknown", "unknown"},
	}

	for _, tt := range tests {
		result := parseRuntime(tt.input)
		if result != tt.expected {
			t.Errorf("parseRuntime(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
