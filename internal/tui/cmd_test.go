package tui

import (
	"fmt"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	"github.com/pueblomo/lazytest/internal/types"
)

func TestRunAllTestsCmd_NoDriver(t *testing.T) {
	m := NewModel("/test")
	m.updateSizes(100, 50)
	m.driver = nil

	// Add test cases
	items := []list.Item{
		&types.TestCase{
			Name:       "test1.spec.ts",
			Filepath:   "/test/test1.spec.ts",
			TestStatus: types.StatusNotStarted,
		},
	}
	m.list.SetItems(items)

	cmd := runAllTestsCmd(m)

	// Command should exist even without driver
	if cmd == nil {
		t.Error("runAllTestsCmd should return a command even without driver")
	}
}

func TestRunAllTestsCmd_NoTests(t *testing.T) {
	m := NewModel("/test")
	m.updateSizes(100, 50)

	// No items in the list
	m.list.SetItems([]list.Item{})

	cmd := runAllTestsCmd(m)

	// Command should exist even with no tests
	if cmd == nil {
		t.Error("runAllTestsCmd should return a command even with no tests")
	}
}

func TestRunAllTestsCmd_WithTests(t *testing.T) {
	m := NewModel("/test")
	m.updateSizes(100, 50)

	// Add test cases
	items := []list.Item{
		&types.TestCase{
			Name:       "test1.spec.ts",
			Filepath:   "/test/test1.spec.ts",
			TestStatus: types.StatusNotStarted,
		},
		&types.TestCase{
			Name:       "test2.spec.ts",
			Filepath:   "/test/test2.spec.ts",
			TestStatus: types.StatusNotStarted,
		},
	}
	m.list.SetItems(items)

	cmd := runAllTestsCmd(m)

	// Command should exist
	if cmd == nil {
		t.Error("runAllTestsCmd should return a command")
	}

	// Note: We can't easily test the actual execution without a real driver
	// and test files, but we verify the command is created
}

func TestWatchForFileChanges_FunctionExists(t *testing.T) {
	m := NewModel("/test")
	m.updateSizes(100, 50)

	testCase := &types.TestCase{
		Name:       "test.spec.ts",
		Filepath:   "/test/test.spec.ts",
		TestStatus: types.StatusNotStarted,
	}

	cmd := watchForFileChanges(m, testCase)

	// Command should exist
	if cmd == nil {
		t.Error("watchForFileChanges should return a command")
	}

	// Note: We can't easily test the file watching behavior in unit tests
	// as it requires actual file system operations and would block
}

func TestWatchForFileChanges_NilTestCase(t *testing.T) {
	m := NewModel("/test")
	m.updateSizes(100, 50)

	// This tests defensive programming - function should handle nil gracefully
	// or panic if that's the intended behavior
	defer func() {
		if r := recover(); r != nil {
			t.Logf("watchForFileChanges panics with nil test case (expected): %v", r)
		}
	}()

	cmd := watchForFileChanges(m, nil)

	// If we get here, it didn't panic
	if cmd != nil {
		t.Log("watchForFileChanges handles nil test case")
	}
}

func TestCommandFunctions_ReturnTeaCmd(t *testing.T) {
	// Verify that command functions return the correct type
	m := NewModel("/test")
	m.updateSizes(100, 50)

	testCase := &types.TestCase{
		Name:       "test.spec.ts",
		Filepath:   "/test/test.spec.ts",
		TestStatus: types.StatusNotStarted,
	}

	// Test runAllTestsCmd returns a tea.Cmd
	cmd1 := runAllTestsCmd(m)
	if cmd1 == nil {
		t.Error("runAllTestsCmd should return non-nil tea.Cmd")
	}

	// Test watchForFileChanges returns a tea.Cmd
	cmd2 := watchForFileChanges(m, testCase)
	if cmd2 == nil {
		t.Error("watchForFileChanges should return non-nil tea.Cmd")
	}
}

func TestRunAllTestsCmd_PreservesTestCaseReferences(t *testing.T) {
	m := NewModel("/test")
	m.updateSizes(100, 50)

	// Add test cases
	testCase1 := &types.TestCase{
		Name:       "test1.spec.ts",
		Filepath:   "/test/test1.spec.ts",
		TestStatus: types.StatusNotStarted,
	}
	testCase2 := &types.TestCase{
		Name:       "test2.spec.ts",
		Filepath:   "/test/test2.spec.ts",
		TestStatus: types.StatusNotStarted,
	}

	items := []list.Item{testCase1, testCase2}
	m.list.SetItems(items)

	// Get the command
	cmd := runAllTestsCmd(m)

	// Verify command was created
	if cmd == nil {
		t.Error("runAllTestsCmd should return a command")
	}

	// Verify the test cases still exist and haven't been corrupted
	retrievedItems := m.list.Items()
	if len(retrievedItems) != 2 {
		t.Errorf("Items should still be 2, got %d", len(retrievedItems))
	}
}

func TestRunSelectedTestsCmd_NoTests(t *testing.T) {
	m := NewModel("/test")
	m.updateSizes(100, 50)

	// No items in the list
	m.list.SetItems([]list.Item{})

	cmd := runSelectedTestsCmd(m)

	// Should not return a command if no tests to run
	if cmd != nil {
		t.Error("runSelectedTestsCmd should return nil when no tests available")
	}
}

func TestRunSelectedTestsCmd_WithSelectedTests(t *testing.T) {
	m := NewModel("/test")
	m.updateSizes(100, 50)

	// Add test cases
	testCase1 := &types.TestCase{
		Name:       "test1.spec.ts",
		Filepath:   "/test/test1.spec.ts",
		TestStatus: types.StatusNotStarted,
		Selected:   true,
	}
	testCase2 := &types.TestCase{
		Name:       "test2.spec.ts",
		Filepath:   "/test/test2.spec.ts",
		TestStatus: types.StatusNotStarted,
		Selected:   false,
	}
	testCase3 := &types.TestCase{
		Name:       "test3.spec.ts",
		Filepath:   "/test/test3.spec.ts",
		TestStatus: types.StatusNotStarted,
		Selected:   true,
	}

	m.list.SetItems([]list.Item{testCase1, testCase2, testCase3})

	cmd := runSelectedTestsCmd(m)

	// Should return a command
	if cmd == nil {
		t.Error("runSelectedTestsCmd should return a command when selected tests exist")
	}
}

func TestRunSelectedTestsCmd_FallbackToFocused(t *testing.T) {
	m := NewModel("/test")
	m.updateSizes(100, 50)

	// Add test cases - none selected
	testCase1 := &types.TestCase{
		Name:       "test1.spec.ts",
		Filepath:   "/test/test1.spec.ts",
		TestStatus: types.StatusNotStarted,
		Selected:   false,
	}
	testCase2 := &types.TestCase{
		Name:       "test2.spec.ts",
		Filepath:   "/test/test2.spec.ts",
		TestStatus: types.StatusNotStarted,
		Selected:   false,
	}

	m.list.SetItems([]list.Item{testCase1, testCase2})

	// Select second item (focus, not selection)
	m.list.Select(1)

	cmd := runSelectedTestsCmd(m)

	// Should return a command (fallback to focused item)
	if cmd == nil {
		t.Error("runSelectedTestsCmd should return a command when using fallback")
	}
}

func TestRunSelectedTestsCmd_AllSelected(t *testing.T) {
	m := NewModel("/test")
	m.updateSizes(100, 50)

	// Add multiple test cases, all selected
	items := make([]list.Item, 5)
	for i := 0; i < 5; i++ {
		items[i] = &types.TestCase{
			Name:       fmt.Sprintf("test%d.spec.ts", i),
			Filepath:   fmt.Sprintf("/test/test%d.spec.ts", i),
			TestStatus: types.StatusNotStarted,
			Selected:   true,
		}
	}
	m.list.SetItems(items)

	cmd := runSelectedTestsCmd(m)

	if cmd == nil {
		t.Error("runSelectedTestsCmd should return a command when all tests selected")
	}
}

func TestRunSelectedTestsCmd_PreservesSelection(t *testing.T) {
	m := NewModel("/test")
	m.updateSizes(100, 50)

	testCase1 := &types.TestCase{
		Name:       "test1.spec.ts",
		Filepath:   "/test/test1.spec.ts",
		TestStatus: types.StatusNotStarted,
		Selected:   true,
	}
	testCase2 := &types.TestCase{
		Name:       "test2.spec.ts",
		Filepath:   "/test/test2.spec.ts",
		TestStatus: types.StatusNotStarted,
		Selected:   false,
	}

	m.list.SetItems([]list.Item{testCase1, testCase2})

	// Get the command (but don't execute it - we just want to ensure the model state is preserved)
	cmd := runSelectedTestsCmd(m)

	if cmd == nil {
		t.Error("Command should be created")
	}

	// Verify selection states are unchanged
	if !testCase1.Selected {
		t.Error("testCase1 should still be selected")
	}
	if testCase2.Selected {
		t.Error("testCase2 should still be unselected")
	}
}
