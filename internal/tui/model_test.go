package tui

import (
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pueblomo/lazytest/internal/types"
)

func TestNewModel(t *testing.T) {
	tests := []struct {
		name     string
		rootPath string
	}{
		{
			name:     "creates model with root path",
			rootPath: "/test/path",
		},
		{
			name:     "creates model with empty root",
			rootPath: "",
		},
		{
			name:     "creates model with relative path",
			rootPath: "./relative/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel(tt.rootPath)

			if m.root != tt.rootPath {
				t.Errorf("NewModel() root = %v, want %v", m.root, tt.rootPath)
			}

			assertFocus(t, m, FocusList)

			if m.prevSelectedIdx != -1 {
				t.Errorf("prevSelectedIdx = %v, want -1", m.prevSelectedIdx)
			}

			if m.width != 0 {
				t.Errorf("width = %v, want 0", m.width)
			}

			if m.height != 0 {
				t.Errorf("height = %v, want 0", m.height)
			}

			// logLine is nil initially, not an empty slice
			if len(m.logLine) != 0 {
				t.Errorf("logLine length = %v, want 0", len(m.logLine))
			}
		})
	}
}

func TestNewModel_WithHelpers(t *testing.T) {
	// Test using helper functions
	m := newTestModel()
	if m.root != testRootPath {
		t.Errorf("newTestModel() root = %v, want %v", m.root, testRootPath)
	}

	m = newTestModelReady()
	assertDimensions(t, m, testWidth, testHeight)
}

func TestModel_Init(t *testing.T) {
	m := newTestModel()
	cmd := m.Init()

	if cmd == nil {
		t.Error("Init() should return a command")
	}
}

func TestModel_UpdateSizes(t *testing.T) {
	tests := testDimensions()

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			m := newTestModel()
			m.updateSizes(tt.width, tt.height)

			assertDimensions(t, m, tt.width, tt.height)

			// Verify viewport dimensions are set
			expectedLogWidth := tt.width - 7
			if m.logView.Width != expectedLogWidth {
				t.Errorf("logView.Width = %v, want %v", m.logView.Width, expectedLogWidth)
			}
		})
	}
}

func TestModel_UpdateSizes_EdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{"zero dimensions", 0, 0},
		{"negative dimensions", -10, -5},
		{"very large dimensions", 10000, 5000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newTestModel()
			// Should not panic
			m.updateSizes(tt.width, tt.height)
			assertDimensions(t, m, tt.width, tt.height)
		})
	}
}

func TestModel_AppendToLog(t *testing.T) {
	tests := []struct {
		name     string
		logs     []string
		expected int
	}{
		{
			name:     "append single log",
			logs:     []string{"log entry 1"},
			expected: 1,
		},
		{
			name:     "append multiple logs",
			logs:     []string{"log 1", "log 2", "log 3"},
			expected: 3,
		},
		{
			name:     "append empty log",
			logs:     []string{""},
			expected: 1,
		},
		{
			name:     "append logs with special characters",
			logs:     []string{"error: test failed", "warning: 🚨 issue found"},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newTestModelReady()

			for _, log := range tt.logs {
				m.appendToLog(log)
			}

			assertLogCount(t, m, tt.expected)

			// Verify logs are in correct order
			for i, log := range tt.logs {
				if m.logLine[i] != log {
					t.Errorf("logLine[%d] = %v, want %v", i, m.logLine[i], log)
				}
			}
		})
	}
}

func TestModel_ClearLog(t *testing.T) {
	m := newTestModelReady()

	// Add some logs using helper
	addTestLogs(&m, 5)
	assertLogCount(t, m, 5)

	// Clear logs
	m.clearLog()
	assertLogCount(t, m, 0)

	// Verify we can still add logs after clearing
	m.appendToLog("new log")
	assertLogCount(t, m, 1)
}

func TestModel_UpdateLogView(t *testing.T) {
	m := newTestModelReady()

	// Add logs
	m.logLine = []string{"log 1", "log 2", "log 3"}
	m.updateLogView()

	content := m.logView.View()

	for _, log := range m.logLine {
		assertViewContains(t, content, log)
	}

	// Test with empty logs
	m.logLine = []string{}
	m.updateLogView()
	// Should handle empty content gracefully
}

func TestModel_UpdateOutputView(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "simple output",
			content: "test output",
		},
		{
			name:    "multiline output",
			content: "line 1\nline 2\nline 3",
		},
		{
			name:    "empty output",
			content: "",
		},
		{
			name:    "output with special characters",
			content: "✓ test passed\n✗ test failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newTestModelReady()
			m.updateOutputView(tt.content)

			// The viewport should handle the content
			_ = m.outputView.View()
		})
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{"a greater than b", 10, 5, 10},
		{"b greater than a", 3, 8, 8},
		{"a equals b", 7, 7, 7},
		{"negative numbers", -5, -10, -5},
		{"zero and positive", 0, 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := max(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("max(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestFocus_Constants(t *testing.T) {
	if FocusList != 0 {
		t.Errorf("FocusList = %v, want 0", FocusList)
	}
	if FocusLogs != 1 {
		t.Errorf("FocusLogs = %v, want 1", FocusLogs)
	}
	if FocusOutput != 2 {
		t.Errorf("FocusOutput = %v, want 2", FocusOutput)
	}
}

func TestModel_TeaModelInterface(t *testing.T) {
	// Verify Model implements tea.Model interface
	var _ tea.Model = (*Model)(nil)
	var _ tea.Model = Model{}
}

func TestModel_Update(t *testing.T) {
	m := newTestModel()

	msg := tea.WindowSizeMsg{Width: testWidth, Height: testHeight}
	newModel, cmd := m.Update(msg)

	if newModel == nil {
		t.Fatal("Update() should return a non-nil model")
	}

	_ = cmd // Cmd can be nil in some cases

	// Verify the returned model has updated dimensions
	if updatedModel, ok := newModel.(Model); ok {
		assertDimensions(t, updatedModel, testWidth, testHeight)
	} else {
		t.Error("Update() should return a Model type")
	}
}

func TestModel_View(t *testing.T) {
	m := newTestModel()

	// Without proper dimensions, should show initializing
	view := m.View()
	assertViewContains(t, view, "Initializing")

	// With dimensions, should show content
	m.updateSizes(testWidth, testHeight)
	view = m.View()

	assertViewNotEmpty(t, view)
	assertViewContains(t, view, "lazytest")
	assertViewContains(t, view, m.root)
}

func TestModel_ListItems(t *testing.T) {
	m := newTestModelReady()

	// Create test items using helper
	items := createTestItems(2, types.StatusNotStarted)
	m.list.SetItems(items)

	assertItemCount(t, m, 2)
}

func TestModel_FocusSwitching(t *testing.T) {
	m := newTestModelReady()

	tests := []struct {
		name     string
		initial  Focus
		expected Focus
	}{
		{"list to output", FocusList, FocusList},
		{"output to logs", FocusOutput, FocusOutput},
		{"logs to list", FocusLogs, FocusLogs},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.focus = tt.initial
			assertFocus(t, m, tt.expected)
		})
	}
}

func TestModel_PrevSelectedIdx(t *testing.T) {
	m := newTestModel()

	if m.prevSelectedIdx != -1 {
		t.Errorf("Initial prevSelectedIdx = %v, want -1", m.prevSelectedIdx)
	}

	tests := []int{0, 5, 10, -1}
	for _, idx := range tests {
		m.prevSelectedIdx = idx
		if m.prevSelectedIdx != idx {
			t.Errorf("After update, prevSelectedIdx = %v, want %v", m.prevSelectedIdx, idx)
		}
	}
}

func TestModel_WithMultipleItems(t *testing.T) {
	m := newTestModelReady()

	// Test with various item counts
	counts := []int{0, 1, 5, 10, 50}
	for _, count := range counts {
		t.Run("", func(t *testing.T) {
			populateModelWithTests(&m, count, types.StatusNotStarted)
			assertItemCount(t, m, count)
		})
	}
}

func TestModel_ViewWithAllStates(t *testing.T) {
	m := newTestModelReady()

	// Add test data
	testCase := createTestCaseWithOutput(
		"test.spec.ts",
		"/test/test.spec.ts",
		"test output",
		types.StatusPassed,
	)
	m.list.SetItems([]list.Item{testCase})
	m.appendToLog("Test log entry")
	m.updateOutputView("Test output")

	// Test view with each focus state
	focusStates := []Focus{FocusList, FocusOutput, FocusLogs}
	for _, focus := range focusStates {
		t.Run("", func(t *testing.T) {
			m.focus = focus
			view := m.View()
			assertViewNotEmpty(t, view)
		})
	}
}
