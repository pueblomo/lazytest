package tui

import (
	"context"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	spinner "github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	logLines []string

	focus           Focus
	spinner         spinner.Model
	list            list.Model
	logView         viewport.Model
	outputView      viewport.Model
	prevSelectedIdx int
	width           int
	height          int

	help     help.Model
	delegate *TestCaseDelegate

	ctx    context.Context
	cancel context.CancelFunc

	logViewCache    scrollbarCache
	outputViewCache scrollbarCache
}

type scrollbarCache struct {
	lastYOffset     int
	lastTotalLines  int
	lastHeight      int
	cachedScrollbar string
}

func NewModel(root string) Model {
	delegate := &TestCaseDelegate{}
	l := list.New([]list.Item{}, delegate, 0, 0)
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.SetFilteringEnabled(true)
	l.SetShowStatusBar(false)

	logViewport := viewport.New(0, 0)
	outputViewport := viewport.New(0, 0)

	ctx, cancel := context.WithCancel(context.Background())

	return Model{
		root:            root,
		focus:           FocusList,
		spinner:         spinner.New(spinner.WithSpinner(spinner.Jump)),
		list:            l,
		delegate:        delegate,
		logView:         logViewport,
		outputView:      outputViewport,
		prevSelectedIdx: -1,
		width:           0,
		height:          0,
		help:            help.New(),
		ctx:             ctx,
		cancel:          cancel,
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
	m.logView.SetContent(strings.Join(m.logLines, "\n"))
	m.logView.GotoBottom()
}

func (m *Model) updateOutputView(content string) {
	m.outputView.SetContent(content)
	m.outputView.GotoBottom()
}

func (m *Model) appendToLog(log string) {
	m.logLines = append(m.logLines, log)
	m.updateLogView()
}

func (m *Model) clearLog() {
	m.logLines = m.logLines[:0]
	m.updateLogView()
}

func (m *Model) updateSizes(width, height int) {
	m.width = width
	m.height = height

	// Handle uninitialized dimensions
	if width <= 0 || height <= 0 {
		return
	}

	// Compute list width including border and padding.
	// The list is rendered with roundedBorder and focused style.
	// Border adds 1 cell on each side, padding(0,2,0) adds 2 cells each side -> total 4.
	listContentWidth := lipgloss.Width(m.list.View())
	listWidth := listContentWidth + 4

	remainingWidth := width - listWidth - 2
	if remainingWidth < 10 {
		remainingWidth = 10
	}
	outputWidth := remainingWidth - 5

	usableHeight := height - 4
	mainHeight := max(1, usableHeight*80/100)
	logHeight := max(1, usableHeight*10/100)

	m.list.SetHeight(mainHeight + 1)

	m.logView.Width = width - 7
	m.logView.Height = logHeight

	m.outputView.Height = mainHeight
	m.outputView.Width = outputWidth
}

func (m *Model) getLogScrollbar() string {
	cache := &m.logViewCache
	if cache.lastYOffset == m.logView.YOffset &&
		cache.lastTotalLines == m.logView.TotalLineCount() &&
		cache.lastHeight == m.logView.Height {
		return cache.cachedScrollbar
	}
	cache.cachedScrollbar = scrollbar(m.logView)
	cache.lastYOffset = m.logView.YOffset
	cache.lastTotalLines = m.logView.TotalLineCount()
	cache.lastHeight = m.logView.Height
	return cache.cachedScrollbar
}

func (m *Model) getOutputScrollbar() string {
	cache := &m.outputViewCache
	if cache.lastYOffset == m.outputView.YOffset &&
		cache.lastTotalLines == m.outputView.TotalLineCount() &&
		cache.lastHeight == m.outputView.Height {
		return cache.cachedScrollbar
	}
	cache.cachedScrollbar = scrollbar(m.outputView)
	cache.lastYOffset = m.outputView.YOffset
	cache.lastTotalLines = m.outputView.TotalLineCount()
	cache.lastHeight = m.outputView.Height
	return cache.cachedScrollbar
}
