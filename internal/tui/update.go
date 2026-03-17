package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	list "github.com/charmbracelet/bubbles/list"
	spinner "github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pueblomo/lazytest/internal/detect"
	"github.com/pueblomo/lazytest/internal/types"
)

func update(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		listCmd    tea.Cmd
		logCmd     tea.Cmd
		spinnerCmd tea.Cmd
		outputCmd  tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.updateSizes(msg.Width, msg.Height)

	case detect.DriverDetectMsg:
		if msg.Err != nil {
			m.appendToLog(msg.Err.Error())
			break
		}
		m.driver = msg.Driver

		return m, func() tea.Msg {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			testFiles, err := m.driver.DetectTestFiles(ctx, m.root)
			if err != nil {
				return detectTestsMsg{err: err}
			}

			return detectTestsMsg{testFiles: testFiles}
		}

	case detectTestsMsg:
		if msg.err != nil {
			m.appendToLog(msg.err.Error())
			break
		}

		items := make([]list.Item, len(msg.testFiles))

		for i, testFile := range msg.testFiles {
			testCase := &types.TestCase{
				Name:       testFile[strings.LastIndex(testFile, "/")+1:],
				Filepath:   testFile,
				TestStatus: types.StatusNotStarted,
			}

			items[i] = testCase
		}

		m.list.SetItems(items)
		// Recompute layout dimensions now that list content has changed
		m.updateSizes(m.width, m.height)

	case testsFinishedMsg:
		if msg.err != nil {
			m.appendToLog(msg.err.Error())
		}
		if item, ok := m.list.SelectedItem().(*types.TestCase); ok {
			m.updateOutputView(item.Output)
		}
		m.appendToLog("Finished")

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, ListKeys.Focus) || key.Matches(msg, OutputKeys.Focus) || key.Matches(msg, LogsKeys.Focus):
			switch m.focus {
			case FocusList:
				m.focus = FocusOutput
			case FocusOutput:
				m.focus = FocusLogs
			case FocusLogs:
				m.focus = FocusList
			}
			return m, nil

		case key.Matches(msg, ListKeys.Quit) || key.Matches(msg, OutputKeys.Quit) || key.Matches(msg, LogsKeys.Quit):
			m.cancel()
			return m, tea.Quit
		}

		if m.focus == FocusList {
			switch {
			case key.Matches(msg, ListKeys.Toggle):
				if item, ok := m.list.SelectedItem().(*types.TestCase); ok {
					item.Selected = !item.Selected
					m.appendToLog(fmt.Sprintf("Toggled selection for %s: %v", item.Name, item.Selected))
				}
				return m, nil
			case key.Matches(msg, ListKeys.RunSelected):
				m.clearLog()
				if m.driver == nil {
					m.appendToLog("No driver")
					break
				}

				// Determine which tests to run
				var testsToRun []*types.TestCase
				for _, item := range m.list.Items() {
					if tc, ok := item.(*types.TestCase); ok && tc.Selected {
						testsToRun = append(testsToRun, tc)
					}
				}

				// If no explicit selection, use focused item
				if len(testsToRun) == 0 {
					if focusedItem := m.list.SelectedItem(); focusedItem != nil {
						if tc, ok := focusedItem.(*types.TestCase); ok {
							testsToRun = append(testsToRun, tc)
						}
					}
				}

				if len(testsToRun) == 0 {
					m.appendToLog("No tests selected")
					break
				}

				m.appendToLog(fmt.Sprintf("Running %d selected test(s)...", len(testsToRun)))

				for _, tc := range testsToRun {
					tc.TestStatus = types.StatusRunning
				}

				return m, runSelectedTestsCmd(m)
			}
		}

		if m.focus == FocusList && m.list.FilterState() != list.Filtering {
			switch {
			case key.Matches(msg, ListKeys.Run):
				m.clearLog()
				if m.driver == nil {
					m.appendToLog("No driver")
					break
				}

				m.appendToLog("Running Tests...")

				for _, item := range m.list.Items() {
					tc := item.(*types.TestCase)
					tc.TestStatus = types.StatusRunning
				}

				return m, runAllTestsCmd(m)
			case key.Matches(msg, ListKeys.Watch):
				if item, ok := m.list.SelectedItem().(*types.TestCase); ok {
					if item.Watched.IsWatching {
						if item.Watched.StopWatching != nil {
							item.Watched.StopWatching()
							item.Watched.IsWatching = false
							item.Watched.StopWatching = nil
						}
					} else {
						m.appendToLog("Watching " + item.Name + "...")
						return m, watchForFileChanges(m, item)
					}
				}
			}
		}

	case spinner.TickMsg:
		m.spinner, spinnerCmd = m.spinner.Update(msg)

	case fileChangedMsg:
		if msg.testCase != nil {
			m.appendToLog("File changed: " + msg.testCase.Name + " - test completed")
			if item, ok := m.list.SelectedItem().(*types.TestCase); ok {
				if item == msg.testCase {
					m.updateOutputView(item.Output)
				}
			}
			return m, watchForFileChanges(m, msg.testCase)
		}

	case watcherMsg:
		if msg.err != nil {
			m.appendToLog("Watch error: " + msg.err.Error())
		}
		if msg.testCase != nil {
			m.appendToLog("Stopped watching " + msg.testCase.Name)
		}
	}

	switch m.focus {
	case FocusList:
		m.list, listCmd = m.list.Update(msg)

		currentIdx := m.list.Index()
		if currentIdx != m.prevSelectedIdx {
			m.prevSelectedIdx = currentIdx

			if item, ok := m.list.SelectedItem().(*types.TestCase); ok {
				m.updateOutputView(item.Output)
			}
		}
	case FocusLogs:
		m.logView, logCmd = m.logView.Update(msg)
	case FocusOutput:
		m.outputView, outputCmd = m.outputView.Update(msg)
	}

	return m, tea.Batch(listCmd, logCmd, spinnerCmd, outputCmd)
}
