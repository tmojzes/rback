package main

import (
	"testing"
)

func TestNormalizeKind(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"sa", kindServiceAccount},
		{"serviceaccounts", kindServiceAccount},
		{"ServiceAccount", kindServiceAccount},
		{"rb", kindRoleBinding},
		{"rolebindings", kindRoleBinding},
		{"crb", kindClusterRoleBinding},
		{"clusterrolebindings", kindClusterRoleBinding},
		{"r", kindRole},
		{"roles", kindRole},
		{"cr", kindClusterRole},
		{"clusterroles", kindClusterRole},
		{"u", kindUser},
		{"users", kindUser},
		{"g", kindGroup},
		{"groups", kindGroup},
		{"unknown", "unknown"},
		{"CustomResource", "customresource"},
	}

	for _, tt := range tests {
		got := normalizeKind(tt.input)
		if got != tt.expected {
			t.Errorf("normalizeKind(%q) = %q; want %q", tt.input, got, tt.expected)
		}
	}
}

func TestVersionVariables(t *testing.T) {
	if version == "" {
		t.Errorf("version should not be empty")
	}
}
