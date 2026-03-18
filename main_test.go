package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestBuildKnowledgeTree(t *testing.T) {
	roots := BuildKnowledgeTree()
	if len(roots) != 4 {
		t.Fatalf("expected 4 root nodes, got %d", len(roots))
	}

	names := []string{"Algorithms", "Data Structures", "Systems", "Languages"}
	for i, name := range names {
		if roots[i].Name != name {
			t.Errorf("root[%d]: expected %q, got %q", i, name, roots[i].Name)
		}
	}
}

func TestNodeChildren(t *testing.T) {
	roots := BuildKnowledgeTree()

	tests := []struct {
		rootIdx  int
		children []string
	}{
		{0, []string{"Sorting", "Searching", "Graph", "Dynamic Programming"}},
		{1, []string{"Arrays", "Trees", "Graphs", "HashMaps"}},
		{2, []string{"OS", "Networking", "Distributed", "Databases"}},
		{3, []string{"Compiled", "Interpreted", "Functional", "Logic"}},
	}

	for _, tt := range tests {
		root := roots[tt.rootIdx]
		if len(root.Children) != len(tt.children) {
			t.Fatalf("%s: expected %d children, got %d", root.Name, len(tt.children), len(root.Children))
		}
		for i, name := range tt.children {
			if root.Children[i].Name != name {
				t.Errorf("%s child[%d]: expected %q, got %q", root.Name, i, name, root.Children[i].Name)
			}
		}
	}
}

func TestIsLeaf(t *testing.T) {
	roots := BuildKnowledgeTree()
	if roots[0].IsLeaf() {
		t.Error("Algorithms should not be a leaf")
	}
	if !roots[0].Children[0].IsLeaf() {
		t.Error("Sorting should be a leaf")
	}
}

func TestDepth(t *testing.T) {
	roots := BuildKnowledgeTree()
	if d := roots[0].Depth(); d != 0 {
		t.Errorf("root depth: expected 0, got %d", d)
	}
	if d := roots[0].Children[0].Depth(); d != 1 {
		t.Errorf("child depth: expected 1, got %d", d)
	}
}

func TestBreadcrumb(t *testing.T) {
	roots := BuildKnowledgeTree()
	if bc := roots[0].Breadcrumb(); bc != "Algorithms" {
		t.Errorf("root breadcrumb: expected %q, got %q", "Algorithms", bc)
	}
	if bc := roots[0].Children[0].Breadcrumb(); bc != "Algorithms > Sorting" {
		t.Errorf("child breadcrumb: expected %q, got %q", "Algorithms > Sorting", bc)
	}
}

func TestParentPointers(t *testing.T) {
	roots := BuildKnowledgeTree()
	for _, root := range roots {
		if root.Parent != nil {
			t.Errorf("root %q should have nil parent", root.Name)
		}
		for _, child := range root.Children {
			if child.Parent != root {
				t.Errorf("child %q parent should be %q", child.Name, root.Name)
			}
		}
	}
}

func TestFlattenVisibleCollapsed(t *testing.T) {
	roots := BuildKnowledgeTree()
	// All collapsed by default
	vis := FlattenVisible(roots)
	if len(vis) != 4 {
		t.Fatalf("collapsed: expected 4 visible, got %d", len(vis))
	}
}

func TestFlattenVisibleExpanded(t *testing.T) {
	roots := BuildKnowledgeTree()
	roots[0].Expanded = true // Expand Algorithms
	vis := FlattenVisible(roots)
	// 4 roots + 4 children of Algorithms
	if len(vis) != 8 {
		t.Fatalf("one expanded: expected 8 visible, got %d", len(vis))
	}
	if vis[1].Name != "Sorting" {
		t.Errorf("expected Sorting at index 1, got %q", vis[1].Name)
	}
}

func TestFlattenVisibleAllExpanded(t *testing.T) {
	roots := BuildKnowledgeTree()
	for _, r := range roots {
		r.Expanded = true
	}
	vis := FlattenVisible(roots)
	// 4 roots + 4*4 children = 20
	if len(vis) != 20 {
		t.Fatalf("all expanded: expected 20 visible, got %d", len(vis))
	}
}

func TestNavigationDown(t *testing.T) {
	m := initialModel(false)
	if m.cursor != 0 {
		t.Fatalf("initial cursor: expected 0, got %d", m.cursor)
	}

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	m = updated.(model)
	if m.cursor != 1 {
		t.Errorf("after j: expected cursor 1, got %d", m.cursor)
	}
}

func TestNavigationUp(t *testing.T) {
	m := initialModel(false)
	m.cursor = 2

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	m = updated.(model)
	if m.cursor != 1 {
		t.Errorf("after k: expected cursor 1, got %d", m.cursor)
	}
}

func TestNavigationBounds(t *testing.T) {
	m := initialModel(false)

	// Can't go above 0
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	m = updated.(model)
	if m.cursor != 0 {
		t.Errorf("up at 0: expected 0, got %d", m.cursor)
	}

	// Go to last item
	m.cursor = 3
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	m = updated.(model)
	if m.cursor != 3 {
		t.Errorf("down at end: expected 3, got %d", m.cursor)
	}
}

func TestExpandCollapse(t *testing.T) {
	m := initialModel(false)
	// Cursor is on Algorithms (index 0), a non-leaf
	vis := m.visible()
	if vis[0].Expanded {
		t.Fatal("should start collapsed")
	}

	// Press enter to expand
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(model)
	vis = m.visible()
	if !vis[0].Expanded {
		t.Error("should be expanded after enter")
	}
	if len(vis) != 8 {
		t.Errorf("expanded: expected 8 visible, got %d", len(vis))
	}

	// Press enter again to collapse
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(model)
	vis = m.visible()
	if vis[0].Expanded {
		t.Error("should be collapsed after second enter")
	}
}

func TestLeafDetailToggle(t *testing.T) {
	m := initialModel(false)
	// Expand Algorithms
	m.roots[0].Expanded = true
	m.cursor = 1 // Sorting (leaf)

	if m.showDetail {
		t.Fatal("detail should start hidden")
	}

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(model)
	if !m.showDetail {
		t.Error("enter on leaf should show detail")
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(model)
	if m.showDetail {
		t.Error("second enter on leaf should hide detail")
	}
}

func TestLeftCollapseAndParent(t *testing.T) {
	m := initialModel(false)
	m.roots[0].Expanded = true
	m.cursor = 1 // Sorting (child of Algorithms)

	// Press h to go to parent
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("h")})
	m = updated.(model)
	if m.cursor != 0 {
		t.Errorf("h on leaf: expected cursor 0 (parent), got %d", m.cursor)
	}

	// Press h again to collapse Algorithms
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("h")})
	m = updated.(model)
	if m.roots[0].Expanded {
		t.Error("h on expanded node should collapse it")
	}
}

func TestRightExpand(t *testing.T) {
	m := initialModel(false)
	// Cursor on Algorithms (collapsed)
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")})
	m = updated.(model)
	if !m.roots[0].Expanded {
		t.Error("l on collapsed node should expand it")
	}
}

func TestViewContainsTitle(t *testing.T) {
	m := initialModel(false)
	view := m.View()
	if !containsStr(view, "Knowledge Explorer") {
		t.Error("view should contain title")
	}
}

func TestViewContainsNodes(t *testing.T) {
	m := initialModel(false)
	view := m.View()
	for _, name := range []string{"Algorithms", "Data Structures", "Systems", "Languages"} {
		if !containsStr(view, name) {
			t.Errorf("view should contain %q", name)
		}
	}
}

func TestViewShowsHelp(t *testing.T) {
	m := initialModel(false)
	view := m.View()
	if !containsStr(view, "j/k: navigate") {
		t.Error("view should contain help text")
	}
}

func TestJSONMode(t *testing.T) {
	m := initialModel(true)
	view := m.View()
	if !containsStr(view, `"name": "Algorithms"`) {
		t.Error("json view should contain Algorithms")
	}
	if !containsStr(view, `"name": "Sorting"`) {
		t.Error("json view should contain Sorting")
	}
}

func TestNodeDescriptions(t *testing.T) {
	roots := BuildKnowledgeTree()
	for _, root := range roots {
		if root.Description == "" {
			t.Errorf("root %q should have a description", root.Name)
		}
		for _, child := range root.Children {
			if child.Description == "" {
				t.Errorf("child %q should have a description", child.Name)
			}
		}
	}
}

func TestNodeResources(t *testing.T) {
	roots := BuildKnowledgeTree()
	for _, root := range roots {
		if len(root.Resources) == 0 {
			t.Errorf("root %q should have resources", root.Name)
		}
		for _, child := range root.Children {
			if len(child.Resources) == 0 {
				t.Errorf("child %q should have resources", child.Name)
			}
		}
	}
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && searchStr(s, substr)
}

func searchStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
