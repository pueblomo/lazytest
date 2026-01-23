package tui

import (
	"testing"

	"github.com/charmbracelet/bubbles/key"
)

func TestOutputKeyMap_ShortHelp(t *testing.T) {
	keys := OutputKeys
	shortHelp := keys.ShortHelp()

	if len(shortHelp) != 4 {
		t.Errorf("ShortHelp() length = %v, want 4", len(shortHelp))
	}

	expectedKeys := []key.Binding{keys.ScrollUp, keys.ScrollDown, keys.Focus, keys.Quit}
	for i := range expectedKeys {
		if len(shortHelp) <= i {
			t.Errorf("ShortHelp() missing index %d", i)
		}
	}
}

func TestOutputKeyMap_FullHelp(t *testing.T) {
	keys := OutputKeys
	fullHelp := keys.FullHelp()

	if len(fullHelp) != 1 {
		t.Errorf("FullHelp() length = %v, want 1", len(fullHelp))
	}

	if len(fullHelp[0]) != 4 {
		t.Errorf("FullHelp()[0] length = %v, want 4", len(fullHelp[0]))
	}

	expectedKeys := []key.Binding{keys.ScrollUp, keys.ScrollDown, keys.Focus, keys.Quit}
	for i := range expectedKeys {
		if len(fullHelp[0]) <= i {
			t.Errorf("FullHelp()[0] missing index %d", i)
		}
	}
}

func TestLogsKeyMap_ShortHelp(t *testing.T) {
	keys := LogsKeys
	shortHelp := keys.ShortHelp()

	if len(shortHelp) != 4 {
		t.Errorf("ShortHelp() length = %v, want 4", len(shortHelp))
	}

	expectedKeys := []key.Binding{keys.ScrollUp, keys.ScrollDown, keys.Focus, keys.Quit}
	for i := range expectedKeys {
		if len(shortHelp) <= i {
			t.Errorf("ShortHelp() missing index %d", i)
		}
	}
}

func TestLogsKeyMap_FullHelp(t *testing.T) {
	keys := LogsKeys
	fullHelp := keys.FullHelp()

	if len(fullHelp) != 1 {
		t.Errorf("FullHelp() length = %v, want 1", len(fullHelp))
	}

	if len(fullHelp[0]) != 4 {
		t.Errorf("FullHelp()[0] length = %v, want 4", len(fullHelp[0]))
	}

	expectedKeys := []key.Binding{keys.ScrollUp, keys.ScrollDown, keys.Focus, keys.Quit}
	for i := range expectedKeys {
		if len(fullHelp[0]) <= i {
			t.Errorf("FullHelp()[0] missing index %d", i)
		}
	}
}

func TestListKeyMap_ShortHelp(t *testing.T) {
	keys := ListKeys
	shortHelp := keys.ShortHelp()

	if len(shortHelp) != 8 {
		t.Errorf("ShortHelp() length = %v, want 8", len(shortHelp))
	}

	expectedKeys := []key.Binding{
		keys.Up, keys.Down, keys.Filter, keys.Remove,
		keys.Run, keys.Watch, keys.Focus, keys.Quit,
	}
	for i := range expectedKeys {
		if len(shortHelp) <= i {
			t.Errorf("ShortHelp() missing index %d", i)
		}
	}
}

func TestListKeyMap_FullHelp(t *testing.T) {
	keys := ListKeys
	fullHelp := keys.FullHelp()

	if len(fullHelp) != 1 {
		t.Errorf("FullHelp() length = %v, want 1", len(fullHelp))
	}

	if len(fullHelp[0]) != 8 {
		t.Errorf("FullHelp()[0] length = %v, want 8", len(fullHelp[0]))
	}

	expectedKeys := []key.Binding{
		keys.Up, keys.Down, keys.Filter, keys.Remove,
		keys.Run, keys.Watch, keys.Focus, keys.Quit,
	}
	for i := range expectedKeys {
		if len(fullHelp[0]) <= i {
			t.Errorf("FullHelp()[0] missing index %d", i)
		}
	}
}

func TestListKeys_Up(t *testing.T) {
	binding := ListKeys.Up
	if len(binding.Keys()) == 0 {
		t.Error("ListKeys.Up should have keys defined")
	}
}

func TestListKeys_Down(t *testing.T) {
	binding := ListKeys.Down
	if len(binding.Keys()) == 0 {
		t.Error("ListKeys.Down should have keys defined")
	}
}

func TestListKeys_Run(t *testing.T) {
	binding := ListKeys.Run
	if len(binding.Keys()) == 0 {
		t.Error("ListKeys.Run should have keys defined")
	}
}

func TestListKeys_Filter(t *testing.T) {
	binding := ListKeys.Filter
	if len(binding.Keys()) == 0 {
		t.Error("ListKeys.Filter should have keys defined")
	}
}

func TestListKeys_Focus(t *testing.T) {
	binding := ListKeys.Focus
	if len(binding.Keys()) == 0 {
		t.Error("ListKeys.Focus should have keys defined")
	}
}

func TestListKeys_Watch(t *testing.T) {
	binding := ListKeys.Watch
	if len(binding.Keys()) == 0 {
		t.Error("ListKeys.Watch should have keys defined")
	}
}

func TestListKeys_Quit(t *testing.T) {
	binding := ListKeys.Quit
	if len(binding.Keys()) == 0 {
		t.Error("ListKeys.Quit should have keys defined")
	}
}

func TestListKeys_Remove(t *testing.T) {
	binding := ListKeys.Remove
	if len(binding.Keys()) == 0 {
		t.Error("ListKeys.Remove should have keys defined")
	}
}

func TestOutputKeys_ScrollUp(t *testing.T) {
	binding := OutputKeys.ScrollUp
	if len(binding.Keys()) == 0 {
		t.Error("OutputKeys.ScrollUp should have keys defined")
	}
}

func TestOutputKeys_ScrollDown(t *testing.T) {
	binding := OutputKeys.ScrollDown
	if len(binding.Keys()) == 0 {
		t.Error("OutputKeys.ScrollDown should have keys defined")
	}
}

func TestOutputKeys_Focus(t *testing.T) {
	binding := OutputKeys.Focus
	if len(binding.Keys()) == 0 {
		t.Error("OutputKeys.Focus should have keys defined")
	}
}

func TestOutputKeys_Quit(t *testing.T) {
	binding := OutputKeys.Quit
	if len(binding.Keys()) == 0 {
		t.Error("OutputKeys.Quit should have keys defined")
	}
}

func TestLogsKeys_ScrollUp(t *testing.T) {
	binding := LogsKeys.ScrollUp
	if len(binding.Keys()) == 0 {
		t.Error("LogsKeys.ScrollUp should have keys defined")
	}
}

func TestLogsKeys_ScrollDown(t *testing.T) {
	binding := LogsKeys.ScrollDown
	if len(binding.Keys()) == 0 {
		t.Error("LogsKeys.ScrollDown should have keys defined")
	}
}

func TestLogsKeys_Focus(t *testing.T) {
	binding := LogsKeys.Focus
	if len(binding.Keys()) == 0 {
		t.Error("LogsKeys.Focus should have keys defined")
	}
}

func TestLogsKeys_Quit(t *testing.T) {
	binding := LogsKeys.Quit
	if len(binding.Keys()) == 0 {
		t.Error("LogsKeys.Quit should have keys defined")
	}
}

func TestKeyMaps_AreNotNil(t *testing.T) {
	if ListKeys.Up.Keys() == nil {
		t.Error("ListKeys.Up should not be nil")
	}
	if ListKeys.Down.Keys() == nil {
		t.Error("ListKeys.Down should not be nil")
	}
	if ListKeys.Run.Keys() == nil {
		t.Error("ListKeys.Run should not be nil")
	}
	if ListKeys.Filter.Keys() == nil {
		t.Error("ListKeys.Filter should not be nil")
	}
	if ListKeys.Focus.Keys() == nil {
		t.Error("ListKeys.Focus should not be nil")
	}
	if ListKeys.Watch.Keys() == nil {
		t.Error("ListKeys.Watch should not be nil")
	}
	if ListKeys.Quit.Keys() == nil {
		t.Error("ListKeys.Quit should not be nil")
	}
	if ListKeys.Remove.Keys() == nil {
		t.Error("ListKeys.Remove should not be nil")
	}

	if OutputKeys.ScrollUp.Keys() == nil {
		t.Error("OutputKeys.ScrollUp should not be nil")
	}
	if OutputKeys.ScrollDown.Keys() == nil {
		t.Error("OutputKeys.ScrollDown should not be nil")
	}
	if OutputKeys.Focus.Keys() == nil {
		t.Error("OutputKeys.Focus should not be nil")
	}
	if OutputKeys.Quit.Keys() == nil {
		t.Error("OutputKeys.Quit should not be nil")
	}

	if LogsKeys.ScrollUp.Keys() == nil {
		t.Error("LogsKeys.ScrollUp should not be nil")
	}
	if LogsKeys.ScrollDown.Keys() == nil {
		t.Error("LogsKeys.ScrollDown should not be nil")
	}
	if LogsKeys.Focus.Keys() == nil {
		t.Error("LogsKeys.Focus should not be nil")
	}
	if LogsKeys.Quit.Keys() == nil {
		t.Error("LogsKeys.Quit should not be nil")
	}
}

func TestKeyMaps_HaveHelpText(t *testing.T) {
	tests := []struct {
		name    string
		binding key.Binding
	}{
		{"ListKeys.Up", ListKeys.Up},
		{"ListKeys.Down", ListKeys.Down},
		{"ListKeys.Run", ListKeys.Run},
		{"ListKeys.Filter", ListKeys.Filter},
		{"ListKeys.Focus", ListKeys.Focus},
		{"ListKeys.Watch", ListKeys.Watch},
		{"ListKeys.Quit", ListKeys.Quit},
		{"ListKeys.Remove", ListKeys.Remove},
		{"OutputKeys.ScrollUp", OutputKeys.ScrollUp},
		{"OutputKeys.ScrollDown", OutputKeys.ScrollDown},
		{"OutputKeys.Focus", OutputKeys.Focus},
		{"OutputKeys.Quit", OutputKeys.Quit},
		{"LogsKeys.ScrollUp", LogsKeys.ScrollUp},
		{"LogsKeys.ScrollDown", LogsKeys.ScrollDown},
		{"LogsKeys.Focus", LogsKeys.Focus},
		{"LogsKeys.Quit", LogsKeys.Quit},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			help := tt.binding.Help()
			if help.Key == "" {
				t.Errorf("%s should have help key", tt.name)
			}
			if help.Desc == "" {
				t.Errorf("%s should have help description", tt.name)
			}
		})
	}
}

func TestKeyBindings_Consistency(t *testing.T) {
	// Test that Focus key is defined across all key maps
	if len(ListKeys.Focus.Keys()) == 0 {
		t.Error("ListKeys.Focus should have keys defined")
	}
	if len(OutputKeys.Focus.Keys()) == 0 {
		t.Error("OutputKeys.Focus should have keys defined")
	}
	if len(LogsKeys.Focus.Keys()) == 0 {
		t.Error("LogsKeys.Focus should have keys defined")
	}

	// Test that Quit key is defined across all key maps
	quitKeys := []key.Binding{ListKeys.Quit, OutputKeys.Quit, LogsKeys.Quit}
	for _, binding := range quitKeys {
		if len(binding.Keys()) == 0 {
			t.Error("Quit should have keys defined")
		}
	}
}
