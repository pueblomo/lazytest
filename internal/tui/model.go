package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	spinner "github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pueblomo/lazytest/internal/detect"
	"github.com/pueblomo/lazytest/internal/drivers"
)

type Focus int

const (
	FocusList Focus = iota
	FocusLogs
	FocusOutput
)

type Model struct {
	driver  drivers.Driver
	root    string
	logLine []string

	focus           Focus
	spinner         spinner.Model
	list            list.Model
	logView         viewport.Model
	outputView      viewport.Model
	prevSelectedIdx int
	width           int
	height          int

	help help.Model
}

func NewModel(root string) Model {
	l := list.New([]list.Item{}, TestCaseDelegate{}, 0, 0)
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.SetFilteringEnabled(true)
	l.SetShowStatusBar(false)

	logViewport := viewport.New(0, 0)
	outputViewport := viewport.New(0, 0)
	return Model{
		root:            root,
		focus:           FocusList,
		spinner:         spinner.New(spinner.WithSpinner(spinner.Jump)),
		list:            l,
		logView:         logViewport,
		outputView:      outputViewport,
		prevSelectedIdx: -1,
		width:           0,
		height:          0,
		help:            help.New(),
	}
}

var _ tea.Model = (*Model)(nil)

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, detect.DetectDriver(m.root))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return update(m, msg)
}

func (m Model) View() string {
	return view(m)
}

func (m *Model) updateLogView() {
	m.logView.SetContent(strings.Join(m.logLine, "\n"))
	m.logView.GotoBottom()
}

func (m *Model) updateOutputView(content string) {
	m.outputView.SetContent(content)
	m.outputView.GotoBottom()
}

func (m *Model) appendToLog(log string) {
	m.logLine = append(m.logLine, log)
	m.updateLogView()
}

func (m *Model) clearLog() {
	m.logLine = m.logLine[:0]
	m.updateLogView()
}

func (m *Model) updateSizes(width, height int) {
	m.width = width
	m.height = height

	usableHeight := height - 4
	mainHeight := max(1, usableHeight*80/100)
	logHeight := max(1, usableHeight*10/100)

	m.list.SetHeight(mainHeight + 1)

	m.logView.Width = width - 7
	m.logView.Height = logHeight

	m.outputView.Height = mainHeight
}


