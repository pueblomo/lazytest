package tui

import (
	"context"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pueblomo/lazytest/internal/types"
)

// mockDriver is a simple Driver implementation for tests
type mockDriver struct{}

func (m *mockDriver) Detect(root string) (bool, error) { return true, nil }
func (m *mockDriver) Name() string                     { return "mock" }
func (m *mockDriver) DetectTestFiles(ctx context.Context, root string) ([]string, error) {
	return nil, nil
}
func (m *mockDriver) RunTest(ctx context.Context, root string, testCase list.Item) error {
	return nil
}

// TestToggleSelection ensures that pressing space toggles the Selected field.
func TestToggleSelection(t *testing.T) {
	m := newTestModelReady()
	m.focus = FocusList
	m.driver = &mockDriver{}
	populateModelWithTests(&m, 2, types.StatusNotStarted)
	m.list.Select(0)

	// Toggle selection on first item
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}}
	newModel, _ := update(m, msg)
	updated := newModel.(Model)

	items := updated.list.Items()
	tc0 := items[0].(*types.TestCase)
	if !tc0.Selected {
		t.Errorf("First item should be selected after toggle")
	}

	// Toggle again to deselect
	newModel2, _ := update(updated, msg)
	updated2 := newModel2.(Model)
	items2 := updated2.list.Items()
	tc0 = items2[0].(*types.TestCase)
	if tc0.Selected {
		t.Errorf("First item should be deselected after second toggle")
	}
}

// TestRunSelectedTests ensures that pressing enter runs only selected tests.
func TestRunSelectedTests(t *testing.T) {
	m := newTestModelReady()
	m.focus = FocusList
	m.driver = &mockDriver{}
	populateModelWithTests(&m, 3, types.StatusNotStarted)

	// Select first and third
	items := m.list.Items()
	items[0].(*types.TestCase).Selected = true
	items[2].(*types.TestCase).Selected = true

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, cmd := update(m, msg)
	updated := newModel.(Model)

	if cmd == nil {
		t.Fatal("expected a command when running selected tests")
	}

	// Check that selected items got Running status
	updatedItems := updated.list.Items()
	countRunning := 0
	for _, item := range updatedItems {
		tc := item.(*types.TestCase)
		if tc.Selected {
			if tc.TestStatus != types.StatusRunning {
				t.Errorf("selected test should be Running, got %v", tc.TestStatus)
			}
			countRunning++
		} else {
			if tc.TestStatus != types.StatusNotStarted {
				t.Errorf("unselected test should be NotStarted, got %v", tc.TestStatus)
			}
		}
	}
	if countRunning != 2 {
		t.Errorf("expected 2 tests running, got %d", countRunning)
	}
}

// TestRunSelected_FallbackToFocused ensures that when no explicit selection, the focused item runs.
func TestRunSelected_FallbackToFocused(t *testing.T) {
	m := newTestModelReady()
	m.focus = FocusList
	m.driver = &mockDriver{}
	populateModelWithTests(&m, 3, types.StatusNotStarted)

	// No explicit selections; focus on second item
	items := m.list.Items()
	for _, item := range items {
		item.(*types.TestCase).Selected = false
	}
	m.list.Select(1)

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, cmd := update(m, msg)
	updated := newModel.(Model)

	if cmd == nil {
		t.Fatal("expected a command with fallback to focused")
	}

	updatedItems := updated.list.Items()
	tc := updatedItems[1].(*types.TestCase)
	if tc.TestStatus != types.StatusRunning {
		t.Errorf("focused item should be Running, got %v", tc.TestStatus)
	}
}

// TestSelectionPersistsThroughFiltering ensures that selection survives filter/unfilter.
func TestSelectionPersistsThroughFiltering(t *testing.T) {
	m := newTestModelReady()
	m.focus = FocusList
	m.driver = &mockDriver{}
	populateModelWithTests(&m, 5, types.StatusNotStarted)

	// Select first and third
	items := m.list.Items()
	items[0].(*types.TestCase).Selected = true
	items[2].(*types.TestCase).Selected = true

	// Apply filter (all contain "test")
	m.list.Filter("test", nil)

	// Check that selections persist in filtered view
	filtered := m.list.VisibleItems()
	for _, item := range filtered {
		tc := item.(*types.TestCase)
		if tc.Name == "test0.spec.ts" || tc.Name == "test2.spec.ts" {
			if !tc.Selected {
				t.Errorf("selection for %s should persist after filtering", tc.Name)
			}
		}
	}

	// Clear filter
	m.list.ResetFilter()

	// After reset, selections should be restored
	all := m.list.Items()
	for i, item := range all {
		tc := item.(*types.TestCase)
		expected := (i == 0 || i == 2)
		if tc.Selected != expected {
			t.Errorf("item %d: expected Selected=%v, got %v", i, expected, tc.Selected)
		}
	}
}

// TestRunSelectedWithFilterApplied ensures that filtering doesn’t affect which tests are run by selection.
func TestRunSelectedWithFilterApplied(t *testing.T) {
	m := newTestModelReady()
	m.focus = FocusList
	m.driver = &mockDriver{}
	populateModelWithTests(&m, 5, types.StatusNotStarted)

	items := m.list.Items()
	items[0].(*types.TestCase).Selected = true
	items[2].(*types.TestCase).Selected = true

	// Filter
	m.list.Filter("test", nil)

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, cmd := update(m, msg)
	updated := newModel.(Model)

	if cmd == nil {
		t.Fatal("expected a command while filter applied")
	}

	updatedItems := updated.list.Items()
	running := 0
	for _, item := range updatedItems {
		tc := item.(*types.TestCase)
		if tc.Selected {
			if tc.TestStatus != types.StatusRunning {
				t.Errorf("selected test not set to Running: %v", tc.TestStatus)
			}
			running++
		} else if tc.TestStatus != types.StatusNotStarted {
			t.Errorf("unselected test status changed: %v", tc.TestStatus)
		}
	}
	if running != 2 {
		t.Errorf("expected 2 tests running, got %d", running)
	}
}

// TestDelegateRendering ensures that the delegate renders checkboxes correctly.
// This is more of a visual test but we verify that Selected field is reflected.
func TestDelegateRendering(t *testing.T) {
	m := newTestModelReady()

	selectedTc := &types.TestCase{
		Name:       "selected.spec.ts",
		Filepath:   "/test/selected.spec.ts",
		TestStatus: types.StatusNotStarted,
		Selected:   true,
	}
	unselectedTc := &types.TestCase{
		Name:       "unselected.spec.ts",
		Filepath:   "/test/unselected.spec.ts",
		TestStatus: types.StatusNotStarted,
		Selected:   false,
	}

	m.list.SetItems([]list.Item{selectedTc, unselectedTc})

	// Verify Selected values are as expected
	items := m.list.Items()
	if !items[0].(*types.TestCase).Selected {
		t.Error("first item should be selected")
	}
	if items[1].(*types.TestCase).Selected {
		t.Error("second item should be unselected")
	}
}
