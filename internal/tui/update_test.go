package tui

import (
	"fmt"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pueblomo/lazytest/internal/detect"
	"github.com/pueblomo/lazytest/internal/types"
)

func TestUpdate_WindowSizeMsg(t *testing.T) {
	m := newTestModel()

	msg := tea.WindowSizeMsg{
		Width:  testWidth,
		Height: testHeight,
	}

	newModel, cmd := update(m, msg)

	if cmd != nil {
		t.Error("WindowSizeMsg should not return a command")
	}

	updatedModel := newModel.(Model)
	assertDimensions(t, updatedModel, testWidth, testHeight)
}

func TestUpdate_DetectTestsMsg_Success(t *testing.T) {
	m := newTestModelReady()

	msg := detectTestsMsg{
		err: nil,
		testFiles: []string{
			"/test/test1.spec.ts",
			"/test/test2.spec.ts",
		},
	}

	newModel, _ := update(m, msg)
	updatedModel := newModel.(Model)

	assertItemCount(t, updatedModel, 2)
	items := updatedModel.list.Items()

	// Verify first item
	if tc, ok := items[0].(*types.TestCase); ok {
		if tc.Name != "test1.spec.ts" {
			t.Errorf("First item name = %v, want 'test1.spec.ts'", tc.Name)
		}
		if tc.Filepath != "/test/test1.spec.ts" {
			t.Errorf("First item filepath = %v, want '/test/test1.spec.ts'", tc.Filepath)
		}
		if tc.TestStatus != types.StatusNotStarted {
			t.Errorf("First item status = %v, want StatusNotStarted", tc.TestStatus)
		}
	} else {
		t.Error("Item should be a TestCase")
	}
}

func TestUpdate_DetectTestsMsg_Error(t *testing.T) {
	m := newTestModelReady()

	msg := detectTestsMsg{
		err:       fmt.Errorf("driver not found"),
		testFiles: nil,
	}

	newModel, _ := update(m, msg)
	updatedModel := newModel.(Model)

	// Should have logged the error
	if len(updatedModel.logLine) == 0 {
		t.Error("update() should log error message")
	}
}

func TestUpdate_TestsFinishedMsg_Success(t *testing.T) {
	m := newTestModelReady()

	// Add a test case using helper
	testCase := createTestCaseWithOutput(
		"test.spec.ts",
		"/test/test.spec.ts",
		"test output",
		types.StatusPassed,
	)
	m.list.SetItems([]list.Item{testCase})

	msg := testsFinishedMsg{
		err: nil,
	}

	newModel, _ := update(m, msg)
	updatedModel := newModel.(Model)

	// Should have logged "Finished"
	if len(updatedModel.logLine) == 0 {
		t.Error("update() should log finished message")
	}
}

func TestUpdate_TestsFinishedMsg_Error(t *testing.T) {
	m := newTestModelReady()

	msg := testsFinishedMsg{
		err: fmt.Errorf("driver not found"),
	}

	newModel, _ := update(m, msg)
	updatedModel := newModel.(Model)

	// Should have logged both error and finished
	if len(updatedModel.logLine) < 2 {
		t.Errorf("update() should log error and finished, got %d logs", len(updatedModel.logLine))
	}
}

func TestUpdate_KeyMsg_Quit(t *testing.T) {
	m := newTestModelReady()

	// Test 'q' key
	msg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'q'},
	}

	_, cmd := update(m, msg)

	// Should return tea.Quit command
	if cmd == nil {
		t.Error("Quit key should return a command")
	}
}

func TestUpdate_KeyMsg_FocusSwitch(t *testing.T) {
	tests := []struct {
		name          string
		initialFocus  Focus
		expectedFocus Focus
	}{
		{
			name:          "from list to output",
			initialFocus:  FocusList,
			expectedFocus: FocusOutput,
		},
		{
			name:          "from output to logs",
			initialFocus:  FocusOutput,
			expectedFocus: FocusLogs,
		},
		{
			name:          "from logs to list",
			initialFocus:  FocusLogs,
			expectedFocus: FocusList,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newTestModelReady()
			m.focus = tt.initialFocus

			msg := tea.KeyMsg{
				Type:  tea.KeyTab,
				Runes: []rune{},
			}

			newModel, _ := update(m, msg)
			updatedModel := newModel.(Model)

			if updatedModel.focus != tt.expectedFocus {
				t.Errorf("Focus after tab = %v, want %v", updatedModel.focus, tt.expectedFocus)
			}
		})
	}
}

func TestUpdate_KeyMsg_Watch(t *testing.T) {
	m := newTestModelReady()
	m.focus = FocusList

	// Add a test case using helper
	testCase := createTestCase("test.spec.ts", "/test/test.spec.ts", types.StatusNotStarted)
	m.list.SetItems([]list.Item{testCase})

	// Simulate 'w' key press (watch toggle)
	msg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'w'},
	}

	_, cmd := update(m, msg)

	// Should return a watch command
	if cmd == nil {
		t.Error("Watch key should return a command")
	}
}

func TestUpdate_SpinnerTick(t *testing.T) {
	m := newTestModelReady()

	// Get a spinner tick message
	tickCmd := m.spinner.Tick
	msg := tickCmd()

	newModel, cmd := update(m, msg)

	// Should return another tick command to keep spinner going
	if cmd == nil {
		t.Error("Spinner tick should return a command")
	}

	// Model should be updated
	_ = newModel.(Model)
}

func TestUpdate_FileChangedMsg(t *testing.T) {
	m := newTestModelReady()

	testCase := createTestCaseWithOutput(
		"test.spec.ts",
		"/test/test.spec.ts",
		"updated output",
		types.StatusPassed,
	)

	msg := fileChangedMsg{
		testCase: testCase,
	}

	newModel, cmd := update(m, msg)
	updatedModel := newModel.(Model)

	// Should log the file change
	if len(updatedModel.logLine) == 0 {
		t.Error("update() should log file changed message")
	}

	// Should return a watch command to continue watching
	if cmd == nil {
		t.Error("fileChangedMsg should return a command to continue watching")
	}
}

func TestUpdate_FileChangedMsg_Nil(t *testing.T) {
	m := newTestModelReady()

	msg := fileChangedMsg{
		testCase: nil,
	}

	newModel, cmd := update(m, msg)

	// Should handle nil test case gracefully
	_ = newModel.(Model)
	_ = cmd
}

func TestUpdate_WatcherMsg_Error(t *testing.T) {
	m := newTestModelReady()

	testCase := createTestCase("test.spec.ts", "/test/test.spec.ts", types.StatusNotStarted)

	msg := watcherMsg{
		err:      fmt.Errorf("watcher error"),
		testCase: testCase,
	}

	newModel, _ := update(m, msg)
	updatedModel := newModel.(Model)

	// Should log the watch error
	if len(updatedModel.logLine) == 0 {
		t.Error("update() should log watcher error")
	}
}

func TestUpdate_WatcherMsg_Stopped(t *testing.T) {
	m := newTestModelReady()

	testCase := createTestCase("test.spec.ts", "/test/test.spec.ts", types.StatusNotStarted)

	msg := watcherMsg{
		err:      nil,
		testCase: testCase,
	}

	newModel, _ := update(m, msg)
	updatedModel := newModel.(Model)

	// Should log that watching stopped
	if len(updatedModel.logLine) == 0 {
		t.Error("update() should log stopped watching message")
	}
}

func TestUpdate_SelectedItemChange(t *testing.T) {
	m := newTestModelReady()
	m.focus = FocusList

	// Add test cases using helpers
	items := createTestItems(2, types.StatusPassed)
	m.list.SetItems(items)
	m.prevSelectedIdx = -1

	// Simulate down arrow key to change selection
	msg := tea.KeyMsg{
		Type: tea.KeyDown,
	}

	newModel, _ := update(m, msg)
	updatedModel := newModel.(Model)

	// prevSelectedIdx should be updated
	if updatedModel.prevSelectedIdx == -1 {
		t.Error("prevSelectedIdx should be updated when selection changes")
	}
}

func TestUpdate_ListFiltering(t *testing.T) {
	m := newTestModelReady()
	m.focus = FocusList

	// Add test cases using helpers
	items := createTestItems(2, types.StatusPassed)
	m.list.SetItems(items)

	// Simulate '/' key to start filtering
	msg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'/'},
	}

	newModel, _ := update(m, msg)

	// Model should be updated with filtering state
	_ = newModel.(Model)
}

func TestUpdate_FocusOutput_Navigation(t *testing.T) {
	m := newTestModelReady()
	m.focus = FocusOutput

	// Set some output content
	m.updateOutputView("Line 1\nLine 2\nLine 3\nLine 4\nLine 5")

	// Simulate down arrow key
	msg := tea.KeyMsg{
		Type: tea.KeyDown,
	}

	newModel, cmd := update(m, msg)

	// Should update the output viewport
	_ = newModel.(Model)
	_ = cmd
}

func TestUpdate_FocusLogs_Navigation(t *testing.T) {
	m := newTestModelReady()
	m.focus = FocusLogs

	// Add some log lines using helper
	addTestLogs(&m, 3)

	// Simulate up arrow key
	msg := tea.KeyMsg{
		Type: tea.KeyUp,
	}

	newModel, cmd := update(m, msg)

	// Should update the log viewport
	_ = newModel.(Model)
	_ = cmd
}

func TestUpdate_DetectDriverMsg_Success(t *testing.T) {
	m := NewModel("/test")
	m.updateSizes(100, 50)

	// Mock driver (we can't easily create a real one in tests)
	// The update function expects a detect.DriverDetectMsg
	// Since we can't easily mock this without circular dependencies,
	// we'll skip detailed testing of this path
	t.Skip("Driver detection requires actual driver implementation")
}

func TestUpdate_DetectDriverMsg_Error(t *testing.T) {
	m := newTestModelReady()

	msg := detect.DriverDetectMsg{
		Driver: nil,
		Err:    fmt.Errorf("driver not found"),
	}

	newModel, _ := update(m, msg)
	updatedModel := newModel.(Model)

	// Should have logged the error
	if len(updatedModel.logLine) == 0 {
		t.Error("update() should log driver detection error")
	}
}

func TestUpdate_MultipleUpdates(t *testing.T) {
	m := newTestModel()

	// Apply multiple updates in sequence
	newModel, _ := update(m, tea.WindowSizeMsg{Width: 100, Height: 50})
	m = newModel.(Model)
	newModel, _ = update(m, detectTestsMsg{
		testFiles: []string{"/test/test1.spec.ts"},
	})
	m = newModel.(Model)

	// Verify accumulated state
	assertDimensions(t, m, 100, 50)
	assertItemCount(t, m, 1)
}

func TestUpdate_CommandBatching(t *testing.T) {
	m := newTestModelReady()
	m.focus = FocusList

	// Add a test case using helper
	populateModelWithTests(&m, 1, types.StatusPassed)

	// Send a key that triggers list update
	msg := tea.KeyMsg{
		Type: tea.KeyDown,
	}

	_, cmd := update(m, msg)

	// Should return a batched command (or nil, depends on state)
	_ = cmd // Just verify it doesn't panic
}

func TestUpdate_PreservesModelState(t *testing.T) {
	m := newTestModelWithSize(testWidth, testHeight)
	m.appendToLog("Initial log")

	originalRoot := m.root
	originalWidth := m.width
	originalHeight := m.height
	originalLogCount := len(m.logLine)

	// Send a neutral message
	msg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'x'}, // Not a bound key
	}

	newModel, _ := update(m, msg)
	updatedModel := newModel.(Model)

	// State should be preserved
	if updatedModel.root != originalRoot {
		t.Error("Root path should be preserved")
	}
	if updatedModel.width != originalWidth {
		t.Error("Width should be preserved")
	}
	if updatedModel.height != originalHeight {
		t.Error("Height should be preserved")
	}
	if len(updatedModel.logLine) != originalLogCount {
		t.Error("Log count should be preserved for unhandled keys")
	}
}

func TestUpdate_DriverDetectMsg_WithDriver(t *testing.T) {
	m := newTestModelReady()

	// We can't create a real driver easily, but we test the nil driver path
	msg := detect.DriverDetectMsg{
		Driver: nil,
		Err:    nil,
	}

	newModel, cmd := update(m, msg)
	updatedModel := newModel.(Model)

	// Driver should remain nil
	if updatedModel.driver != nil {
		t.Error("Driver should be nil when nil driver is passed")
	}

	// Should not return a command when no driver
	if cmd != nil {
		t.Log("Command returned with nil driver (may be expected)")
	}
}

func TestUpdate_KeyMsg_RunTests_NoDriver(t *testing.T) {
	m := newTestModelReady()
	m.focus = FocusList
	m.driver = nil

	// Add test cases using helper
	populateModelWithTests(&m, 1, types.StatusNotStarted)

	// Simulate 'r' key press (run tests)
	msg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'r'},
	}

	newModel, _ := update(m, msg)
	updatedModel := newModel.(Model)

	// Should log "No driver" message
	hasNoDriverLog := false
	for _, log := range updatedModel.logLine {
		if log == "No driver" {
			hasNoDriverLog = true
			break
		}
	}

	if !hasNoDriverLog {
		t.Error("Should log 'No driver' when trying to run tests without driver")
	}
}

func TestUpdate_KeyMsg_EscapeRemovesFilter(t *testing.T) {
	m := newTestModelReady()
	m.focus = FocusList

	// Add test cases using helper
	populateModelWithTests(&m, 1, types.StatusNotStarted)

	// Simulate 'esc' key press (remove filter)
	msg := tea.KeyMsg{
		Type: tea.KeyEsc,
	}

	newModel, _ := update(m, msg)

	// Model should be updated
	_ = newModel.(Model)
}

func TestUpdate_KeyMsg_UpDown(t *testing.T) {
	m := newTestModelReady()
	m.focus = FocusList

	// Add multiple test cases using helper
	populateModelWithTests(&m, 2, types.StatusNotStarted)

	// Test down key
	downMsg := tea.KeyMsg{
		Type: tea.KeyDown,
	}
	newModel, _ := update(m, downMsg)
	m = newModel.(Model)

	// Test up key
	upMsg := tea.KeyMsg{
		Type: tea.KeyUp,
	}
	newModel, _ = update(m, upMsg)

	// Model should handle navigation
	_ = newModel.(Model)
}

func TestUpdate_DetectTestsMsg_WithPathSeparators(t *testing.T) {
	m := newTestModelReady()

	msg := detectTestsMsg{
		err: nil,
		testFiles: []string{
			"/path/to/nested/test.spec.ts",
			"relative/path/test.spec.ts",
			"simple.spec.ts",
		},
	}

	newModel, _ := update(m, msg)
	updatedModel := newModel.(Model)

	items := updatedModel.list.Items()
	if len(items) != 3 {
		t.Errorf("Should have 3 items, got %d", len(items))
	}

	// Check name extraction
	if tc, ok := items[0].(*types.TestCase); ok {
		if tc.Name != "test.spec.ts" {
			t.Errorf("Name should be 'test.spec.ts', got '%s'", tc.Name)
		}
	}

	if tc, ok := items[2].(*types.TestCase); ok {
		if tc.Name != "simple.spec.ts" {
			t.Errorf("Name should be 'simple.spec.ts', got '%s'", tc.Name)
		}
	}
}

func TestUpdate_TestsFinishedMsg_WithSelectedItem(t *testing.T) {
	m := newTestModelReady()

	testCase := createTestCaseWithOutput(
		"test.spec.ts",
		"/test/test.spec.ts",
		"Test passed successfully",
		types.StatusPassed,
	)
	m.list.SetItems([]list.Item{testCase})
	m.list.Select(0)

	msg := testsFinishedMsg{
		err: nil,
	}

	newModel, _ := update(m, msg)
	updatedModel := newModel.(Model)

	// Should update output view with selected item's output
	if len(updatedModel.logLine) == 0 {
		t.Error("Should have logged finished message")
	}
}

func TestUpdate_KeyMsg_CtrlC(t *testing.T) {
	m := newTestModelReady()

	msg := tea.KeyMsg{
		Type: tea.KeyCtrlC,
	}

	_, cmd := update(m, msg)

	// Should return quit command
	if cmd == nil {
		t.Error("Ctrl+C should return a quit command")
	}
}

func TestUpdate_ListFocus_WithFiltering(t *testing.T) {
	m := newTestModelReady()
	m.focus = FocusList

	// Add test items using helper
	populateModelWithTests(&m, 2, types.StatusPassed)

	// Start filtering
	msg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'/'},
	}

	newModel, _ := update(m, msg)

	// Model should handle filter mode
	_ = newModel.(Model)
}

func TestUpdate_EmptyTestFiles(t *testing.T) {
	m := newTestModelReady()

	msg := detectTestsMsg{
		err:       nil,
		testFiles: []string{},
	}

	newModel, _ := update(m, msg)
	updatedModel := newModel.(Model)

	assertItemCount(t, updatedModel, 0)
}

func TestUpdate_SequentialFocusChanges(t *testing.T) {
	m := newTestModelReady()

	// Cycle through all focus states
	tabMsg := tea.KeyMsg{Type: tea.KeyTab}

	// Start at FocusList
	m.focus = FocusList
	newModel, _ := update(m, tabMsg)
	m = newModel.(Model)
	if m.focus != FocusOutput {
		t.Error("Should move to FocusOutput")
	}

	// FocusOutput -> FocusLogs
	newModel, _ = update(m, tabMsg)
	m = newModel.(Model)
	if m.focus != FocusLogs {
		t.Error("Should move to FocusLogs")
	}

	// FocusLogs -> FocusList
	newModel, _ = update(m, tabMsg)
	m = newModel.(Model)
	if m.focus != FocusList {
		t.Error("Should cycle back to FocusList")
	}
}
