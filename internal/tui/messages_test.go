package tui

import (
	"errors"
	"testing"

	"github.com/pueblomo/lazytest/internal/types"
)

func TestDetectTestsMsg(t *testing.T) {
	tests := []struct {
		name      string
		msg       detectTestsMsg
		wantErr   bool
		wantFiles int
	}{
		{
			name: "successful detection with files",
			msg: detectTestsMsg{
				err:       nil,
				testFiles: []string{"test1.spec.ts", "test2.spec.ts"},
			},
			wantErr:   false,
			wantFiles: 2,
		},
		{
			name: "error during detection",
			msg: detectTestsMsg{
				err:       errors.New("detection failed"),
				testFiles: nil,
			},
			wantErr:   true,
			wantFiles: 0,
		},
		{
			name: "no test files found",
			msg: detectTestsMsg{
				err:       nil,
				testFiles: []string{},
			},
			wantErr:   false,
			wantFiles: 0,
		},
		{
			name: "single test file",
			msg: detectTestsMsg{
				err:       nil,
				testFiles: []string{"single.test.ts"},
			},
			wantErr:   false,
			wantFiles: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if (tt.msg.err != nil) != tt.wantErr {
				t.Errorf("detectTestsMsg.err = %v, wantErr %v", tt.msg.err, tt.wantErr)
			}

			if len(tt.msg.testFiles) != tt.wantFiles {
				t.Errorf("detectTestsMsg.testFiles length = %v, want %v", len(tt.msg.testFiles), tt.wantFiles)
			}
		})
	}
}

func TestTestsFinishedMsg(t *testing.T) {
	tests := []struct {
		name    string
		msg     testsFinishedMsg
		wantErr bool
	}{
		{
			name: "successful completion",
			msg: testsFinishedMsg{
				err: nil,
			},
			wantErr: false,
		},
		{
			name: "completion with error",
			msg: testsFinishedMsg{
				err: errors.New("tests failed"),
			},
			wantErr: true,
		},
		{
			name: "completion with wrapped error",
			msg: testsFinishedMsg{
				err: errors.New("execution error: timeout"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if (tt.msg.err != nil) != tt.wantErr {
				t.Errorf("testsFinishedMsg.err = %v, wantErr %v", tt.msg.err, tt.wantErr)
			}
		})
	}
}

func TestWatcherMsg(t *testing.T) {
	tests := []struct {
		name         string
		msg          watcherMsg
		wantErr      bool
		wantTestCase bool
	}{
		{
			name: "watcher stopped successfully",
			msg: watcherMsg{
				err: nil,
				testCase: &types.TestCase{
					Name:     "test.spec.ts",
					Filepath: "/path/to/test.spec.ts",
				},
			},
			wantErr:      false,
			wantTestCase: true,
		},
		{
			name: "watcher error",
			msg: watcherMsg{
				err:      errors.New("watcher failed"),
				testCase: nil,
			},
			wantErr:      true,
			wantTestCase: false,
		},
		{
			name: "watcher stopped without test case",
			msg: watcherMsg{
				err:      nil,
				testCase: nil,
			},
			wantErr:      false,
			wantTestCase: false,
		},
		{
			name: "watcher error with test case",
			msg: watcherMsg{
				err: errors.New("file not found"),
				testCase: &types.TestCase{
					Name:     "missing.spec.ts",
					Filepath: "/path/to/missing.spec.ts",
				},
			},
			wantErr:      true,
			wantTestCase: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if (tt.msg.err != nil) != tt.wantErr {
				t.Errorf("watcherMsg.err = %v, wantErr %v", tt.msg.err, tt.wantErr)
			}

			if (tt.msg.testCase != nil) != tt.wantTestCase {
				t.Errorf("watcherMsg.testCase = %v, wantTestCase %v", tt.msg.testCase, tt.wantTestCase)
			}
		})
	}
}

func TestFileChangedMsg(t *testing.T) {
	tests := []struct {
		name         string
		msg          fileChangedMsg
		wantTestCase bool
	}{
		{
			name: "file changed with test case",
			msg: fileChangedMsg{
				testCase: &types.TestCase{
					Name:       "test.spec.ts",
					Filepath:   "/path/to/test.spec.ts",
					TestStatus: types.StatusPassed,
				},
			},
			wantTestCase: true,
		},
		{
			name: "file changed without test case",
			msg: fileChangedMsg{
				testCase: nil,
			},
			wantTestCase: false,
		},
		{
			name: "file changed with running test",
			msg: fileChangedMsg{
				testCase: &types.TestCase{
					Name:       "running.spec.ts",
					Filepath:   "/path/to/running.spec.ts",
					TestStatus: types.StatusRunning,
				},
			},
			wantTestCase: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if (tt.msg.testCase != nil) != tt.wantTestCase {
				t.Errorf("fileChangedMsg.testCase = %v, wantTestCase %v", tt.msg.testCase, tt.wantTestCase)
			}

			if tt.msg.testCase != nil {
				if tt.msg.testCase.Name == "" {
					t.Error("fileChangedMsg.testCase.Name should not be empty")
				}
				if tt.msg.testCase.Filepath == "" {
					t.Error("fileChangedMsg.testCase.Filepath should not be empty")
				}
			}
		})
	}
}

func TestMessages_ErrorHandling(t *testing.T) {
	// Test various error scenarios
	testErr := errors.New("test error")

	detectMsg := detectTestsMsg{err: testErr}
	if detectMsg.err == nil {
		t.Error("detectTestsMsg should preserve error")
	}
	if detectMsg.err.Error() != "test error" {
		t.Errorf("detectTestsMsg error = %v, want 'test error'", detectMsg.err.Error())
	}

	finishedMsg := testsFinishedMsg{err: testErr}
	if finishedMsg.err == nil {
		t.Error("testsFinishedMsg should preserve error")
	}
	if finishedMsg.err.Error() != "test error" {
		t.Errorf("testsFinishedMsg error = %v, want 'test error'", finishedMsg.err.Error())
	}

	watchMsg := watcherMsg{err: testErr}
	if watchMsg.err == nil {
		t.Error("watcherMsg should preserve error")
	}
	if watchMsg.err.Error() != "test error" {
		t.Errorf("watcherMsg error = %v, want 'test error'", watchMsg.err.Error())
	}
}

func TestMessages_TestCaseHandling(t *testing.T) {
	testCase := &types.TestCase{
		Name:       "example.spec.ts",
		Filepath:   "/path/to/example.spec.ts",
		Output:     "test output",
		TestStatus: types.StatusPassed,
	}

	watchMsg := watcherMsg{testCase: testCase}
	if watchMsg.testCase == nil {
		t.Error("watcherMsg should preserve test case")
	}
	if watchMsg.testCase.Name != "example.spec.ts" {
		t.Errorf("watcherMsg.testCase.Name = %v, want 'example.spec.ts'", watchMsg.testCase.Name)
	}

	fileMsg := fileChangedMsg{testCase: testCase}
	if fileMsg.testCase == nil {
		t.Error("fileChangedMsg should preserve test case")
	}
	if fileMsg.testCase.Name != "example.spec.ts" {
		t.Errorf("fileChangedMsg.testCase.Name = %v, want 'example.spec.ts'", fileMsg.testCase.Name)
	}
}

func TestMessages_TestFilesHandling(t *testing.T) {
	testFiles := []string{
		"test1.spec.ts",
		"test2.spec.ts",
		"test3.spec.ts",
	}

	msg := detectTestsMsg{testFiles: testFiles}

	if len(msg.testFiles) != 3 {
		t.Errorf("detectTestsMsg.testFiles length = %v, want 3", len(msg.testFiles))
	}

	for i, file := range testFiles {
		if msg.testFiles[i] != file {
			t.Errorf("detectTestsMsg.testFiles[%d] = %v, want %v", i, msg.testFiles[i], file)
		}
	}
}

func TestMessages_NilValues(t *testing.T) {
	// Test that messages can have nil values
	detectMsg := detectTestsMsg{err: nil, testFiles: nil}
	if detectMsg.err != nil {
		t.Error("detectTestsMsg with nil err should be nil")
	}
	if detectMsg.testFiles != nil {
		t.Error("detectTestsMsg with nil testFiles should be nil")
	}

	finishedMsg := testsFinishedMsg{err: nil}
	if finishedMsg.err != nil {
		t.Error("testsFinishedMsg with nil err should be nil")
	}

	watchMsg := watcherMsg{err: nil, testCase: nil}
	if watchMsg.err != nil {
		t.Error("watcherMsg with nil err should be nil")
	}
	if watchMsg.testCase != nil {
		t.Error("watcherMsg with nil testCase should be nil")
	}

	fileMsg := fileChangedMsg{testCase: nil}
	if fileMsg.testCase != nil {
		t.Error("fileChangedMsg with nil testCase should be nil")
	}
}
