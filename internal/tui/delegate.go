package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pueblomo/lazytest/internal/types"
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(0).Foreground(lipgloss.Color("170"))
)

type TestCaseDelegate struct {
	SpinnerFrame string
}

func (i TestCaseDelegate) Height() int {
	return 1
}

func (i TestCaseDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

func (d TestCaseDelegate) Spacing() int { return 0 }

func (d TestCaseDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	tc, ok := listItem.(*types.TestCase)
	if !ok {
		return
	}

	icon := tc.TestStatusIcon()

	if tc.TestStatus == types.StatusRunning {
		icon = d.SpinnerFrame
	}

	watchIcon := ""
	if tc.Watched.IsWatching {
		watchIcon = types.IconWatching
	}

	checkbox := "[ ]"
	if tc.Selected {
		checkbox = "[x]"
	}

	str := fmt.Sprintf("%s %s %s %s", checkbox, icon, tc.Name, watchIcon)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}
