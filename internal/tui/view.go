package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

const (
	titel = "lazytest – test dashboard"
)

var (
	roundedBorder  = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).Padding(0, 2, 0)
	scrollBarStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

func view(m Model) string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	var b strings.Builder
	root := fmt.Sprintf("📁 %s", m.root)
	var driver string
	if m.driver != nil {
		driver = fmt.Sprintf("🔧 %s driver", m.driver.Name())
	} else {
		driver = "🔧 driver not found"
	}

	titlePane := roundedBorder.Render(lipgloss.JoinVertical(lipgloss.Left,
		titel,
		root,
		driver,
	))

	logContent := m.logView.View()
	logBar := scrollBarStyle.Render(scrollbar(m.logView))

	logPane := lipgloss.JoinHorizontal(
		lipgloss.Top,
		logContent,
		logBar,
	)

	listView := focused(roundedBorder, m.focus == FocusList).Render(m.list.View())
	listWidth := lipgloss.Width(listView)

	remainingWidth := m.width - listWidth - 2
	if remainingWidth < 10 {
		remainingWidth = 10
	}

	m.outputView.Width = remainingWidth - 5

	outputContent := m.outputView.View()
	outputBar := scrollBarStyle.Render(scrollbar(m.outputView))
	outputPane := lipgloss.JoinHorizontal(
		lipgloss.Top,
		outputContent,
		outputBar,
	)

	mainPane := lipgloss.JoinHorizontal(
		lipgloss.Top,
		listView,
		focused(roundedBorder, m.focus == FocusOutput).Render(outputPane),
	)

	b.WriteString(lipgloss.JoinVertical(lipgloss.Left,
		titlePane,
		mainPane,
		focused(roundedBorder, m.focus == FocusLogs).Render(logPane)))

	b.WriteString("\n")

	switch m.focus {
	case FocusList:
		b.WriteString(m.help.View(ListKeys))
	case FocusOutput:
		b.WriteString(m.help.View(OutputKeys))
	case FocusLogs:
		b.WriteString(m.help.View(LogsKeys))
	}

	return b.String()
}

func focused(style lipgloss.Style, focused bool) lipgloss.Style {
	if focused {
		return style.BorderForeground(lipgloss.Color("170"))
	}
	return style.BorderForeground(lipgloss.Color("240"))
}

func scrollbar(vp viewport.Model) string {
	if vp.TotalLineCount() <= vp.Height {
		return strings.Repeat(" \n", vp.Height)
	}

	ratio := float64(vp.Height) / float64(vp.TotalLineCount())
	barHeight := int(ratio * float64(vp.Height))
	if barHeight < 1 {
		barHeight = 1
	}

	maxOffset := vp.TotalLineCount() - vp.Height
	barTop := 0
	if maxOffset > 0 {
		barTop = int(
			float64(vp.YOffset) / float64(maxOffset) *
				float64(vp.Height-barHeight),
		)
	}

	lines := make([]string, vp.Height)
	for i := 0; i < vp.Height; i++ {
		if i >= barTop && i < barTop+barHeight {
			lines[i] = "█"
		} else {
			lines[i] = " "
		}
	}

	return strings.Join(lines, "\n")
}
