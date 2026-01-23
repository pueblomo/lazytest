package tui

import "github.com/charmbracelet/bubbles/key"

type outputKeyMap struct {
	ScrollUp   key.Binding
	ScrollDown key.Binding
	Focus      key.Binding
	Quit       key.Binding
}

func (k outputKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.ScrollUp, k.ScrollDown, k.Focus, k.Quit}
}

func (k outputKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.ScrollUp, k.ScrollDown, k.Focus, k.Quit}}
}

type logsKeyMap struct {
	ScrollUp   key.Binding
	ScrollDown key.Binding
	Focus      key.Binding
	Quit       key.Binding
}

func (k logsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.ScrollUp, k.ScrollDown, k.Focus, k.Quit}
}

func (k logsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.ScrollUp, k.ScrollDown, k.Focus, k.Quit}}
}

type listKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Run    key.Binding
	Filter key.Binding
	Focus  key.Binding
	Quit   key.Binding
	Watch  key.Binding
	Remove key.Binding
}

func (k listKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Filter, k.Remove, k.Run, k.Watch, k.Focus, k.Quit}
}

func (k listKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Up, k.Down, k.Filter, k.Remove, k.Run, k.Watch, k.Focus, k.Quit}}
}

var ListKeys = listKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "move down"),
	),
	Run: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "run visible tests"),
	),
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter list"),
	),
	Focus: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch focus"),
	),
	Watch: key.NewBinding(
		key.WithKeys("w"),
		key.WithHelp("w", "toggle watch"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Remove: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "remove filter"),
	),
}

var OutputKeys = outputKeyMap{
	ScrollUp: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "scroll up"),
	),
	ScrollDown: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "scroll down"),
	),
	Focus: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch focus"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

var LogsKeys = logsKeyMap{
	ScrollUp: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "scroll up"),
	),
	ScrollDown: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "scroll down"),
	),
	Focus: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch focus"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
