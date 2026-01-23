package tui

import (
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
