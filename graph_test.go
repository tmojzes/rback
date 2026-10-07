package main

import (
	"strings"
	"testing"

	"github.com/emicklei/dot"
)

func TestEscapeHTML(t *testing.T) {
	input := "<script>foo bar\nbaz</script>"
	got := escapeHTML(input)
	if strings.Contains(got, "<script>") {
		t.Errorf("escapeHTML failed to sanitize: %s", got)
	}
	if !strings.Contains(got, "&lt;script&gt;") {
		t.Errorf("expected &lt;script&gt; in %s", got)
	}
	if !strings.Contains(got, "&nbsp;") {
		t.Errorf("expected &nbsp; for space in %s", got)
	}
	if !strings.Contains(got, "<br/>") {
		t.Errorf("expected <br/> for newline in %s", got)
	}
}

func TestFormatLabel(t *testing.T) {
	normal := formatLabel("mylabel", false)
	if normal != "mylabel" {
		t.Errorf("expected normal string label, got %v", normal)
	}

	highlighted := formatLabel("mylabel", true)
	if html, ok := highlighted.(dot.HTML); !ok || !strings.Contains(string(html), "<b>mylabel</b>") {
		t.Errorf("expected bold dot.HTML label for highlight, got %v", highlighted)
	}
}

func TestGraphCreationAndNodes(t *testing.T) {
	g := newGraph()
	ns := newNamespaceSubgraph(g, "test-ns")

	subNode := newSubjectNode0(ns, "ServiceAccount", "app-sa", true, false)
	rbNode := newRoleBindingNode(ns, "app-binding", false)
	rNode := newRoleNode(ns, "test-ns", "app-role", true, false)
	rulesNode := newRulesNode0(ns, "test-ns", "app-role", "get pods<br align=\"left\"/>", false)

	newSubjectToBindingEdge(subNode, rbNode)
	newBindingToRoleEdge(rbNode, rNode)
	newRoleToRulesEdge(rNode, rulesNode)

	dotOutput := g.String()
	if !strings.Contains(dotOutput, "app-sa") {
		t.Errorf("graph missing subject node: %s", dotOutput)
	}
	if !strings.Contains(dotOutput, "app-binding") {
		t.Errorf("graph missing binding node: %s", dotOutput)
	}
	if !strings.Contains(dotOutput, "app-role") {
		t.Errorf("graph missing role node: %s", dotOutput)
	}

	// Test duplicate edge avoidance
	_ = edge(subNode, rbNode)
	_ = edge(subNode, rbNode)
	if edges := subNode.EdgesTo(rbNode); len(edges) != 1 {
		t.Errorf("expected exactly 1 edge between nodes, got %d", len(edges))
	}
}
