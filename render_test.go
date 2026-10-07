package main

import (
	"os"
	"strings"
	"testing"
)

func loadSampleRback(t *testing.T, cfg Config) Rback {
	t.Helper()
	file, err := os.Open("testdata/rbac_sample.json")
	if err != nil {
		t.Fatalf("Failed to open test data: %v", err)
	}
	defer file.Close()

	rback := Rback{config: cfg}
	err = rback.parseRBAC(file)
	if err != nil {
		t.Fatalf("parseRBAC failed: %v", err)
	}
	return rback
}

func TestGenGraph_Default(t *testing.T) {
	cfg := Config{
		showLegend: true,
		showRules:  true,
		namespaces: []string{""}, // all
	}
	rback := loadSampleRback(t, cfg)
	g := rback.genGraph()
	dot := g.String()

	if !strings.Contains(dot, "LEGEND") {
		t.Errorf("graph expected to contain LEGEND")
	}
	if !strings.Contains(dot, "pod-reader") {
		t.Errorf("graph expected to contain pod-reader")
	}
	if !strings.Contains(dot, "cluster-admin-role") {
		t.Errorf("graph expected to contain cluster-admin-role")
	}
}

func TestGenGraph_NoLegend(t *testing.T) {
	cfg := Config{
		showLegend: false,
		showRules:  true,
		namespaces: []string{""},
	}
	rback := loadSampleRback(t, cfg)
	g := rback.genGraph()
	dot := g.String()

	if strings.Contains(dot, "LEGEND") {
		t.Errorf("graph should not contain LEGEND when showLegend=false")
	}
}

func TestGenGraph_NamespaceFilter(t *testing.T) {
	cfg := Config{
		showLegend: false,
		showRules:  true,
		namespaces: []string{"namespace1"},
	}
	rback := loadSampleRback(t, cfg)
	g := rback.genGraph()
	dot := g.String()

	if !strings.Contains(dot, "namespace1") {
		t.Errorf("expected namespace1 subgraph")
	}
	if strings.Contains(dot, "cluster-admin-binding") {
		t.Errorf("cluster-admin-binding should not be rendered when filtering to namespace1")
	}
	if !strings.Contains(dot, "crb-sa-binding") {
		t.Errorf("crb-sa-binding targeting my-sa in namespace1 should be rendered")
	}
}

func TestGenGraph_FocusServiceAccount(t *testing.T) {
	cfg := Config{
		showLegend:    false,
		showRules:     true,
		namespaces:    []string{""},
		resourceKind:  kindServiceAccount,
		resourceNames: []string{"my-sa"},
	}
	rback := loadSampleRback(t, cfg)
	g := rback.genGraph()
	dot := g.String()

	if !strings.Contains(dot, "my-sa") {
		t.Errorf("expected focused my-sa node")
	}
}

func TestWhoCan_Matching(t *testing.T) {
	// 1. Matches pod-reader ("get pods")
	cfg := Config{
		showLegend:   false,
		showRules:    true,
		namespaces:   []string{""},
		resourceKind: kindRule,
		whoCan: WhoCan{
			verb:         "get",
			resourceKind: "pods",
		},
	}
	rback := loadSampleRback(t, cfg)
	g := rback.genGraph()
	dot := g.String()

	if !strings.Contains(dot, "pod-reader") {
		t.Errorf("expected pod-reader in who-can get pods")
	}

	// 2. Matches cluster-admin ("* *")
	if !strings.Contains(dot, "cluster-admin-role") {
		t.Errorf("expected cluster-admin-role in who-can get pods via wildcard")
	}

	// 3. Unmatched query
	cfgUnmatched := Config{
		showLegend:   false,
		showRules:    true,
		namespaces:   []string{"namespace1"},
		resourceKind: kindRule,
		whoCan: WhoCan{
			verb:         "delete",
			resourceKind: "services",
		},
	}
	rbackUnmatched := loadSampleRback(t, cfgUnmatched)
	gUnmatched := rbackUnmatched.genGraph()
	dotUnmatched := gUnmatched.String()

	if strings.Contains(dotUnmatched, "pod-reader") {
		t.Errorf("pod-reader should not match who-can delete services")
	}
}

func TestRuleToHumanReadableString(t *testing.T) {
	rule := Rule{
		verbs:         []string{"get", "list"},
		resources:     []string{"pods"},
		resourceNames: []string{"my-pod"},
		apiGroups:     []string{"apps"},
	}
	s := rule.toHumanReadableString()
	if !strings.Contains(s, "get,list") {
		t.Errorf("expected verbs in string, got %s", s)
	}
	if !strings.Contains(s, "pods") {
		t.Errorf("expected resources in string, got %s", s)
	}
	if !strings.Contains(s, `"my-pod"`) {
		t.Errorf("expected resource name in string, got %s", s)
	}
	if !strings.Contains(s, "(apps)") {
		t.Errorf("expected api group in string, got %s", s)
	}
}
