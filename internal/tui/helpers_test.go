package tui

import (
	"testing"

	"github.com/charmbracelet/bubbles/list"
	"github.com/pueblomo/lazytest/internal/types"
)

// Test fixtures and constants
const (
	testRootPath    = "/test/project"
	testWidth       = 100
	testHeight      = 50
	testSmallWidth  = 20
	testSmallHeight = 10
	testLargeWidth  = 200
	testLargeHeight = 100
)

// newTestModel creates a Model with standard test configuration
func newTestModel() Model {
	return NewModel(testRootPath)
}

// newTestModelWithSize creates a Model with specified dimensions
func newTestModelWithSize(width, height int) Model {
	m := NewModel(testRootPath)
	m.updateSizes(width, height)
	return m
}

// newTestModelReady creates a fully initialized Model ready for testing
func newTestModelReady() Model {
	m := NewModel(testRootPath)
	m.updateSizes(testWidth, testHeight)
	return m
}

// createTestCase creates a TestCase with standard configuration
func createTestCase(name, filepath string, status types.TestStatus) *types.TestCase {
	return &types.TestCase{
		Name:       name,
		Filepath:   filepath,
		TestStatus: status,
	}
}

// createTestCaseWithOutput creates a TestCase with output
func createTestCaseWithOutput(name, filepath, output string, status types.TestStatus) *types.TestCase {
	return &types.TestCase{
		Name:       name,
		Filepath:   filepath,
		Output:     output,
		TestStatus: status,
	}
}

// createTestItems creates a slice of test items for list testing
func createTestItems(count int, status types.TestStatus) []list.Item {
	items := make([]list.Item, count)
	for i := 0; i < count; i++ {
		items[i] = createTestCase(
			testFileName(i),
			testFilePath(i),
			status,
		)
	}
	return items
}

// testFileName returns a standard test file name for index i
func testFileName(i int) string {
	return "test" + string(rune('0'+i)) + ".spec.ts"
}

// testFilePath returns a standard test file path for index i
func testFilePath(i int) string {
	return testRootPath + "/" + testFileName(i)
}

// populateModelWithTests adds test cases to a model
func populateModelWithTests(m *Model, count int, status types.TestStatus) {
	items := createTestItems(count, status)
	m.list.SetItems(items)
}

// assertDimensions checks that model dimensions match expected values
func assertDimensions(t *testing.T, m Model, expectedWidth, expectedHeight int) {
	t.Helper()
	if m.width != expectedWidth {
		t.Errorf("width = %d, want %d", m.width, expectedWidth)
	}
	if m.height != expectedHeight {
		t.Errorf("height = %d, want %d", m.height, expectedHeight)
	}
}

// assertFocus checks that model has expected focus
func assertFocus(t *testing.T, m Model, expectedFocus Focus) {
	t.Helper()
	if m.focus != expectedFocus {
		t.Errorf("focus = %v, want %v", m.focus, expectedFocus)
	}
}

// assertItemCount checks that list has expected number of items
func assertItemCount(t *testing.T, m Model, expectedCount int) {
	t.Helper()
	actualCount := len(m.list.Items())
	if actualCount != expectedCount {
		t.Errorf("item count = %d, want %d", actualCount, expectedCount)
	}
}

// assertLogCount checks that log has expected number of entries
func assertLogCount(t *testing.T, m Model, expectedCount int) {
	t.Helper()
	actualCount := len(m.logLines)
	if actualCount != expectedCount {
		t.Errorf("log count = %d, want %d", actualCount, expectedCount)
	}
}

// assertLogContains checks that log contains expected message
func assertLogContains(t *testing.T, m Model, expectedMsg string) {
	t.Helper()
	for _, log := range m.logLines {
		if log == expectedMsg {
			return
		}
	}
	t.Errorf("log does not contain expected message: %q", expectedMsg)
}

// assertViewNotEmpty checks that view is not empty
func assertViewNotEmpty(t *testing.T, view string) {
	t.Helper()
	if view == "" {
		t.Error("view should not be empty")
	}
}

// assertViewContains checks that view contains expected text
func assertViewContains(t *testing.T, view string, expected string) {
	t.Helper()
	if !contains(view, expected) {
		t.Errorf("view should contain %q", expected)
	}
}

// contains checks if string contains substring (helper to avoid imports)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || indexString(s, substr) >= 0)
}

// indexString finds the index of substr in s
func indexString(s, substr string) int {
	n := len(substr)
	if n == 0 {
		return 0
	}
	for i := 0; i <= len(s)-n; i++ {
		if s[i:i+n] == substr {
			return i
		}
	}
	return -1
}

// testStatuses returns common test status values for table tests
func testStatuses() []types.TestStatus {
	return []types.TestStatus{
		types.StatusNotStarted,
		types.StatusRunning,
		types.StatusPassed,
		types.StatusFailed,
		types.StatusSkipped,
	}
}

// testDimensions returns common dimension pairs for table tests
func testDimensions() []struct {
	width  int
	height int
} {
	return []struct {
		width  int
		height int
	}{
		{testWidth, testHeight},
		{testSmallWidth, testSmallHeight},
		{testLargeWidth, testLargeHeight},
	}
}

// addTestLogs adds multiple log entries to model
func addTestLogs(m *Model, count int) {
	for i := 0; i < count; i++ {
		m.appendToLog("Log entry " + string(rune('0'+i)))
	}
}

// testPaths returns common test paths for table tests
func testPaths() []string {
	return []string{
		"/absolute/path/to/test.spec.ts",
		"relative/path/test.spec.ts",
		"simple.spec.ts",
		"/nested/deeply/nested/path/test.spec.ts",
	}
}
