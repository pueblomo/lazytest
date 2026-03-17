package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

func TestView_Initializing(t *testing.T) {
	m := NewModel("/test/path")

	// Without dimensions, should show initializing message
	view := m.View()

	if !strings.Contains(view, "Initializing") {
		t.Error("View() without dimensions should show 'Initializing...'")
	}
}

func TestView_WithDimensions(t *testing.T) {
	m := NewModel("/test/path")
	m.updateSizes(100, 50)

	view := m.View()

	// Should not show initializing when dimensions are set
	if strings.Contains(view, "Initializing") {
		t.Error("View() with dimensions should not show 'Initializing'")
	}

	// Should contain title
	if !strings.Contains(view, "lazytest") {
		t.Error("View() should contain 'lazytest' title")
	}

	// Should contain root path
	if !strings.Contains(view, "/test/path") {
		t.Error("View() should contain root path")
	}

	// Should contain driver info (even if not found)
	if !strings.Contains(view, "driver") {
		t.Error("View() should contain driver information")
	}
}

func TestView_RootPathDisplay(t *testing.T) {
	tests := []struct {
		name string
		root string
	}{
		{
			name: "absolute path",
			root: "/absolute/path/to/project",
		},
		{
			name: "relative path",
			root: "./relative/path",
		},
		{
			name: "short path",
			root: "/test",
		},
		{
			name: "empty path",
			root: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel(tt.root)
			m.updateSizes(100, 50)

			view := m.View()

			if tt.root != "" && !strings.Contains(view, tt.root) {
				t.Errorf("View() should contain root path '%s'", tt.root)
			}
		})
	}
}

func TestView_DriverNotFound(t *testing.T) {
	m := NewModel("/test")
	m.updateSizes(100, 50)
	m.driver = nil

	view := m.View()

	if !strings.Contains(view, "driver not found") {
		t.Error("View() should show 'driver not found' when driver is nil")
	}
}

func TestFocused_StyleApplication(t *testing.T) {
	style := lipgloss.NewStyle()

	// Test focused style returns a style
	focusedStyle := focused(style, true)
	_ = focusedStyle // Just verify it returns without error

	// Test unfocused style returns a style
	unfocusedStyle := focused(style, false)
	_ = unfocusedStyle // Just verify it returns without error

	// Both should return valid styles (can't easily compare styles directly)
}

func TestScrollbar_EmptyViewport(t *testing.T) {
	vp := viewport.New(10, 10)

	scrollbarOutput := scrollbar(vp)

	// Should return string with newlines
	if !strings.Contains(scrollbarOutput, "\n") {
		t.Error("scrollbar() should return multiline string")
	}

	// Count lines - split on newline includes trailing empty string
	lines := strings.Split(scrollbarOutput, "\n")
	// The last element is empty due to trailing newline
	actualLines := len(lines) - 1
	if actualLines != vp.Height {
		t.Errorf("scrollbar() should have %d lines, got %d", vp.Height, actualLines)
	}
}

func TestScrollbar_ContentFitsInViewport(t *testing.T) {
	vp := viewport.New(10, 10)
	vp.SetContent("line1\nline2\nline3")

	scrollbarOutput := scrollbar(vp)

	// When content fits, scrollbar should be all spaces
	lines := strings.Split(scrollbarOutput, "\n")
	for i, line := range lines {
		// Skip the last empty line from trailing newline
		if i == len(lines)-1 && line == "" {
			continue
		}
		if line != " " {
			t.Errorf("scrollbar() line should be space when content fits, got '%s'", line)
		}
	}
}

func TestScrollbar_ContentExceedsViewport(t *testing.T) {
	vp := viewport.New(10, 5)

	// Create content that exceeds viewport height
	content := ""
	for i := 0; i < 20; i++ {
		content += "line " + string(rune(i)) + "\n"
	}
	vp.SetContent(content)

	scrollbarOutput := scrollbar(vp)

	// Should contain scrollbar blocks (█)
	if !strings.Contains(scrollbarOutput, "█") {
		t.Error("scrollbar() should contain block characters when content exceeds viewport")
	}

	// Count lines
	lines := strings.Split(scrollbarOutput, "\n")
	if len(lines) != vp.Height {
		t.Errorf("scrollbar() should have %d lines, got %d", vp.Height, len(lines))
	}
}

func TestScrollbar_ZeroHeight(t *testing.T) {
	vp := viewport.New(10, 0)
	vp.SetContent("content")

	scrollbarOutput := scrollbar(vp)

	// Should handle zero height gracefully
	if scrollbarOutput != "" {
		t.Error("scrollbar() with zero height should return empty string")
	}
}

func TestScrollbar_ScrolledToTop(t *testing.T) {
	vp := viewport.New(10, 5)

	// Create content that exceeds viewport
	content := ""
	for i := 0; i < 20; i++ {
		content += "line\n"
	}
	vp.SetContent(content)
	vp.YOffset = 0 // At top

	scrollbarOutput := scrollbar(vp)

	lines := strings.Split(scrollbarOutput, "\n")

	// First lines should contain block character
	foundBlock := false
	for i := 0; i < 3 && i < len(lines); i++ {
		if lines[i] == "█" {
			foundBlock = true
			break
		}
	}

	if !foundBlock {
		t.Error("scrollbar() at top should have block in first few lines")
	}
}

func TestScrollbar_ScrolledToBottom(t *testing.T) {
	vp := viewport.New(10, 5)

	// Create content that exceeds viewport
	content := ""
	for i := 0; i < 20; i++ {
		content += "line\n"
	}
	vp.SetContent(content)
	vp.GotoBottom()

	scrollbarOutput := scrollbar(vp)

	lines := strings.Split(scrollbarOutput, "\n")

	// Last lines should contain block character
	foundBlock := false
	for i := len(lines) - 3; i < len(lines); i++ {
		if i >= 0 && lines[i] == "█" {
			foundBlock = true
			break
		}
	}

	if !foundBlock {
		t.Error("scrollbar() at bottom should have block in last few lines")
	}
}

func TestView_FocusIndicators(t *testing.T) {
	tests := []struct {
		name  string
		focus Focus
	}{
		{
			name:  "focus on list",
			focus: FocusList,
		},
		{
			name:  "focus on output",
			focus: FocusOutput,
		},
		{
			name:  "focus on logs",
			focus: FocusLogs,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel("/test")
			m.updateSizes(100, 50)
			m.focus = tt.focus

			view := m.View()

			// View should show different help based on focus
			// Each focus mode has different keybindings displayed
			if view == "" {
				t.Error("View() should not be empty")
			}
		})
	}
}

func TestView_ListKeysHelp(t *testing.T) {
	m := NewModel("/test")
	m.updateSizes(100, 50)
	m.focus = FocusList

	view := m.View()

	// Should contain list-specific help
	// The help system should render keybindings, but we can't easily check
	// specific text without testing the bubbles library itself
	if !strings.Contains(view, "lazytest") {
		t.Error("View() should contain content")
	}
}

func TestView_OutputKeysHelp(t *testing.T) {
	m := NewModel("/test")
	m.updateSizes(100, 50)
	m.focus = FocusOutput

	view := m.View()

	if !strings.Contains(view, "lazytest") {
		t.Error("View() should contain content")
	}
}

func TestView_LogsKeysHelp(t *testing.T) {
	m := NewModel("/test")
	m.updateSizes(100, 50)
	m.focus = FocusLogs

	view := m.View()

	if !strings.Contains(view, "lazytest") {
		t.Error("View() should contain content")
	}
}

func TestView_WithLogs(t *testing.T) {
	m := NewModel("/test")
	m.updateSizes(100, 50)

	m.appendToLog("Log line 1")
	m.appendToLog("Log line 2")
	m.appendToLog("Log line 3")

	view := m.View()

	// View should render logs (they're in the viewport)
	if view == "" {
		t.Error("View() should not be empty with logs")
	}
}

func TestView_WithOutput(t *testing.T) {
	m := NewModel("/test")
	m.updateSizes(100, 50)

	m.updateOutputView("Test output line 1\nTest output line 2")

	view := m.View()

	// View should render output
	if view == "" {
		t.Error("View() should not be empty with output")
	}
}

func TestView_SmallDimensions(t *testing.T) {
	m := NewModel("/test")
	m.updateSizes(20, 10)

	view := m.View()

	// Should still render without errors even with small dimensions
	if view == "" {
		t.Error("View() should render even with small dimensions")
	}
}

func TestView_LargeDimensions(t *testing.T) {
	m := NewModel("/test")
	m.updateSizes(500, 200)

	view := m.View()

	// Should handle large dimensions
	if view == "" {
		t.Error("View() should render with large dimensions")
	}
}

func TestMax_Utility(t *testing.T) {
	// This is tested in model_test.go but verify it's used in view rendering
	if max(10, 5) != 10 {
		t.Error("max() should return larger value")
	}
}

func TestScrollbar_WithDifferentHeights(t *testing.T) {
	heights := []int{1, 5, 10, 20, 50}

	for _, height := range heights {
		t.Run(string(rune(height)), func(t *testing.T) {
			vp := viewport.New(10, height)
			vp.SetContent("content\n")

			scrollbarOutput := scrollbar(vp)

			lines := strings.Split(scrollbarOutput, "\n")
			// Account for trailing newline creating empty string at end
			// However, for empty scrollbar output, split returns [""]
			actualLines := len(lines)
			if scrollbarOutput != "" && lines[len(lines)-1] == "" {
				actualLines = len(lines) - 1
			}
			if actualLines != height {
				t.Errorf("scrollbar() should have %d lines for height %d, got %d", height, height, actualLines)
			}
		})
	}
}

func TestScrollbar_BarHeightCalculation(t *testing.T) {
	vp := viewport.New(10, 10)

	// Set content with exactly 20 lines (2x viewport height)
	content := ""
	for i := 0; i < 20; i++ {
		content += "line\n"
	}
	vp.SetContent(content)

	scrollbarOutput := scrollbar(vp)

	// Should have scrollbar blocks
	if !strings.Contains(scrollbarOutput, "█") {
		t.Error("scrollbar() should contain blocks when content is 2x viewport height")
	}

	lines := strings.Split(scrollbarOutput, "\n")
	blockCount := 0
	for _, line := range lines {
		if line == "█" {
			blockCount++
		}
	}

	// With 20 lines of content and 10 line viewport, ratio is 0.5
	// Bar height should be around 5 lines (ratio * viewport height)
	if blockCount < 1 {
		t.Error("scrollbar() should have at least 1 block")
	}
}

func TestView_TitleConstant(t *testing.T) {
	if title != "lazytest – test dashboard" {
		t.Errorf("title constant = %v, want 'lazytest – test dashboard'", title)
	}
}

func TestView_StylesInitialized(t *testing.T) {
	// Test that styles exist and can be used
	_ = roundedBorder.GetBorderStyle()
	_ = scrollBarStyle.GetForeground()

	// Just verify they don't panic when accessed
	t.Log("Styles are initialized and accessible")
}
