package tui

import "github.com/pueblomo/lazytest/internal/types"

type detectTestsMsg struct {
	err       error
	testFiles []string
}

type testsFinishedMsg struct {
	err error
}

type watcherMsg struct {
	err      error
	testCase *types.TestCase
}

type fileChangedMsg struct {
	testCase *types.TestCase
}
