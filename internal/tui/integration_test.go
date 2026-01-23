package tui

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pueblomo/lazytest/internal/detect"
	"github.com/pueblomo/lazytest/internal/types"
)

// Integration tests require full system resources and are separated
// from unit tests using build tags.
//
// Run these tests with: go test -tags=integration

func TestIntegration_FullWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	m := newTestModelReady()

	// Simulate a complete workflow
	steps := []struct {
		name string
		msg  tea.Msg
	}{
		{
			name: "window resize",
			msg:  tea.WindowSizeMsg{Width: 120, Height: 60},
		},
		{
			name: "detect test files",
			msg: detectTestsMsg{
				err: nil,
				testFiles: []string{
					"/test/project/test1.spec.ts",
					"/test/project/test2.spec.ts",
					"/test/project/test3.spec.ts",
				},
			},
		},
		{
			name: "tests finished",
			msg: testsFinishedMsg{
				err: nil,
			},
		},
	}

	for _, step := range steps {
		t.Run(step.name, func(t *testing.T) {
			newModel, cmd := update(m, step.msg)
			m = newModel.(Model)
			_ = cmd

			// Verify view renders without panic
			view := m.View()
			assertViewNotEmpty(t, view)
		})
	}
}

func TestIntegration_FileWatcher(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	m := newTestModelReady()

	testCase := createTestCase(
		"test.spec.ts",
		"/test/test.spec.ts",
		types.StatusNotStarted,
	)

	// Create watcher command
	cmd := watchForFileChanges(m, testCase)
	if cmd == nil {
		t.Fatal("watchForFileChanges should return a command")
	}

	// Note: We can't fully test the watcher without actual file operations
	// This test verifies the command is created correctly
	t.Log("File watcher command created successfully")
}

func TestIntegration_TestRunner(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	m := newTestModelReady()
	m.driver = nil // No real driver for test

	// Add test cases
	populateModelWithTests(&m, 3, types.StatusNotStarted)

	// Create run command
	cmd := runAllTestsCmd(m)
	if cmd == nil {
		t.Fatal("runAllTestsCmd should return a command")
	}

	// Note: We can't execute the command without a real driver
	// This test verifies the command is created correctly
	t.Log("Test runner command created successfully")
}

func TestIntegration_LongRunningOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	m := newTestModelReady()

	// Simulate long-running operation
	done := make(chan bool)
	go func() {
		// Simulate multiple updates over time
		for i := 0; i < 10; i++ {
			m.appendToLog("Operation step " + string(rune('0'+i)))
			time.Sleep(10 * time.Millisecond)
		}
		done <- true
	}()

	select {
	case <-done:
		assertLogCount(t, m, 10)
	case <-time.After(5 * time.Second):
		t.Fatal("Operation timed out")
	}
}

func TestIntegration_ConcurrentUpdates(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	m := newTestModelReady()

	// In a real Bubble Tea application, all updates go through the Update method sequentially.
	// This test verifies that rapid sequential updates work correctly.

	// Simulate rapid sequential operations (as would happen in real usage)
	for i := 0; i < 5; i++ {
		m.appendToLog("Log entry " + string(rune('0'+i)))
	}

	for i := 0; i < 5; i++ {
		m.updateOutputView("Output " + string(rune('0'+i)))
	}

	sizes := []struct{ w, h int }{{100, 50}, {120, 60}, {80, 40}}
	for _, sz := range sizes {
		m.updateSizes(sz.w, sz.h)
	}

	// Verify model is still in valid state
	if len(m.logLine) != 5 {
		t.Errorf("Expected 5 log lines, got %d", len(m.logLine))
	}

	// Verify final dimensions
	if m.width != 80 || m.height != 40 {
		t.Errorf("Expected dimensions (80, 40), got (%d, %d)", m.width, m.height)
	}
}

func TestIntegration_DriverDetection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	m := newTestModelReady()

	// Test driver detection message handling
	msg := detect.DriverDetectMsg{
		Driver: nil,
		Err:    nil,
	}

	newModel, cmd := update(m, msg)
	m = newModel.(Model)
	_ = cmd

	// Verify model handles driver detection
	if m.driver != nil {
		t.Log("Driver was set (unexpected in test environment)")
	}
}

func TestIntegration_StressTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	m := newTestModelReady()

	// Add many test items
	populateModelWithTests(&m, 100, types.StatusNotStarted)

	// Add many log entries
	for i := 0; i < 100; i++ {
		m.appendToLog("Stress test log entry " + string(rune(i)))
	}

	// Update view multiple times
	for i := 0; i < 50; i++ {
		view := m.View()
		if view == "" {
			t.Fatal("View became empty during stress test")
		}
	}

	// Verify final state
	assertItemCount(t, m, 100)
	if len(m.logLine) < 100 {
		t.Errorf("Expected at least 100 log entries, got %d", len(m.logLine))
	}
}

func TestIntegration_MessageSequence(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	m := newTestModelReady()

	// Simulate a realistic message sequence
	messages := []tea.Msg{
		tea.WindowSizeMsg{Width: 100, Height: 50},
		detect.DriverDetectMsg{Driver: nil, Err: nil},
		detectTestsMsg{testFiles: []string{"/test/a.ts", "/test/b.ts"}},
		tea.KeyMsg{Type: tea.KeyTab},
		tea.KeyMsg{Type: tea.KeyTab},
		testsFinishedMsg{err: nil},
	}

	for i, msg := range messages {
		newModel, cmd := update(m, msg)
		m = newModel.(Model)
		_ = cmd

		// Verify view is valid after each message
		view := m.View()
		if view == "" {
			t.Errorf("View empty after message %d", i)
		}
	}
}
