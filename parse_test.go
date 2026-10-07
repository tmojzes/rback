package main

import (
	"os"
	"strings"
	"testing"
)

func TestParseRBAC_ValidData(t *testing.T) {
	file, err := os.Open("testdata/rbac_sample.json")
	if err != nil {
		t.Fatalf("Failed to open test data: %v", err)
	}
	defer file.Close()

	rback := Rback{
		config: Config{
			ignoredPrefixes: []string{"system:"},
		},
	}

	err = rback.parseRBAC(file)
	if err != nil {
		t.Fatalf("parseRBAC failed: %v", err)
	}

	// Check ServiceAccounts
	saMap, exists := rback.permissions.ServiceAccounts["namespace1"]
	if !exists {
		t.Fatalf("expected serviceaccounts in namespace1")
	}
	if _, ok := saMap["my-sa"]; !ok {
		t.Errorf("expected my-sa in namespace1")
	}
	// Check ignored serviceaccount
	if _, ok := rback.permissions.ServiceAccounts["kube-system"]["system:ignored-sa"]; ok {
		t.Errorf("system:ignored-sa should have been filtered out")
	}

	// Check Roles
	roles, exists := rback.permissions.Roles["namespace1"]
	if !exists || len(roles) != 1 {
		t.Errorf("expected 1 role in namespace1, got %v", len(roles))
	}
	if role, ok := roles["pod-reader"]; !ok {
		t.Errorf("expected pod-reader role")
	} else if len(role.rules) != 1 || role.rules[0].verbs[0] != "get" {
		t.Errorf("unexpected rules in pod-reader: %+v", role.rules)
	}

	// Check ClusterRoles
	clusterRoles, exists := rback.permissions.Roles[""]
	if !exists || len(clusterRoles) != 1 {
		t.Errorf("expected 1 clusterrole, got %v", len(clusterRoles))
	}
	if _, ok := clusterRoles["cluster-admin-role"]; !ok {
		t.Errorf("expected cluster-admin-role")
	}

	// Check RoleBindings
	rbs, exists := rback.permissions.RoleBindings["namespace1"]
	if !exists || len(rbs) != 1 {
		t.Errorf("expected 1 rolebinding in namespace1")
	}
	rb := rbs["read-pods-binding"]
	if rb.role.name != "pod-reader" || rb.role.namespace != "namespace1" {
		t.Errorf("unexpected roleRef in binding: %+v", rb.role)
	}
	if len(rb.subjects) != 2 {
		t.Errorf("expected 2 subjects in read-pods-binding, got %d", len(rb.subjects))
	}

	// Check ClusterRoleBindings
	crbs, exists := rback.permissions.RoleBindings[""]
	if !exists || len(crbs) != 1 {
		t.Errorf("expected 1 clusterrolebinding")
	}
}

func TestParseRBAC_InvalidInputs(t *testing.T) {
	rback := Rback{}

	// Invalid JSON
	err := rback.parseRBAC(strings.NewReader("not a json"))
	if err == nil {
		t.Errorf("expected error on invalid JSON, got nil")
	}

	// Kind not List
	err = rback.parseRBAC(strings.NewReader(`{"kind": "Pod"}`))
	if err == nil || !strings.Contains(err.Error(), "expected kind=List") {
		t.Errorf("expected error for non-List kind, got: %v", err)
	}

	// Empty List
	err = rback.parseRBAC(strings.NewReader(`{"kind": "List", "items": []}`))
	if err != nil {
		t.Errorf("unexpected error on empty items list: %v", err)
	}

	// List without items key
	err = rback.parseRBAC(strings.NewReader(`{"kind": "List"}`))
	if err != nil {
		t.Errorf("unexpected error on List with no items key: %v", err)
	}

	// Items containing malformed or missing fields (should not panic)
	malformed := `{
		"kind": "List",
		"items": [
			{},
			{"kind": "Role"},
			{"kind": "Role", "metadata": {}},
			{"kind": "RoleBinding", "metadata": {"name": "test"}, "roleRef": {}},
			{"kind": "ServiceAccount", "metadata": null}
		]
	}`
	err = rback.parseRBAC(strings.NewReader(malformed))
	if err != nil {
		t.Errorf("unexpected error on malformed items: %v", err)
	}
}

func TestShouldIgnore(t *testing.T) {
	r := Rback{
		config: Config{
			ignoredPrefixes: []string{"system:", "kube-"},
		},
	}

	if !r.shouldIgnore("system:admin") {
		t.Errorf("expected system:admin to be ignored")
	}
	if !r.shouldIgnore("kube-proxy") {
		t.Errorf("expected kube-proxy to be ignored")
	}
	if r.shouldIgnore("my-account") {
		t.Errorf("my-account should not be ignored")
	}
}
