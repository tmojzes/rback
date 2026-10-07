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

func TestStringSliceFlag(t *testing.T) {
	var flag stringSliceFlag
	_ = flag.Set("ns1, ns2")
	_ = flag.Set("ns3")

	expected := []string{"ns1", "ns2", "ns3"}
	if len(flag) != len(expected) {
		t.Fatalf("expected %d elements, got %d", len(expected), len(flag))
	}
	for i, v := range expected {
		if flag[i] != v {
			t.Errorf("expected [%d] = %s, got %s", i, v, flag[i])
		}
	}
}
