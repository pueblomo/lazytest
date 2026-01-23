package tui

import (
	"context"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fsnotify/fsnotify"
	"github.com/pueblomo/lazytest/internal/types"
)

func runAllTestsCmd(m Model) tea.Cmd {
	cases := m.list.VisibleItems()
	driver := m.driver
	root := m.root

	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		var wg sync.WaitGroup
		errCh := make(chan error, len(cases))

		for i := range cases {
			wg.Go(func() {
				err := driver.RunTest(ctx, root, cases[i])
				if err != nil {
					errCh <- err
				}
			})
		}

		wg.Wait()
		close(errCh)

		var err error
		for e := range errCh {
			if err == nil {
				err = e
			}
		}

		return testsFinishedMsg{
			err: err,
		}
	}
}

func watchForFileChanges(m Model, tc *types.TestCase) tea.Cmd {
	return func() tea.Msg {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			tc.Watched.IsWatching = false
			return watcherMsg{
				err:      err,
				testCase: tc,
			}
		}
		defer watcher.Close()
		defer func() {
			tc.Watched.IsWatching = false
		}()

		err = watcher.Add(tc.Filepath)
		if err != nil {
			tc.Watched.IsWatching = false
			return watcherMsg{
				err:      err,
				testCase: tc,
			}
		}

		tc.Watched.IsWatching = true
		tc.Watched.StopWatching = watcher.Close

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					tc.Watched.IsWatching = false
					return watcherMsg{
						err:      nil,
						testCase: tc,
					}
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
					m.driver.RunTest(ctx, m.root, tc)
					cancel()

					return fileChangedMsg{
						testCase: tc,
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					tc.Watched.IsWatching = false
					return watcherMsg{
						err:      nil,
						testCase: tc,
					}
				}
				tc.Watched.IsWatching = false
				return watcherMsg{
					err:      err,
					testCase: tc,
				}
			}
		}
	}
}
