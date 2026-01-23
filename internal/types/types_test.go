package types

import (
	"testing"
)

func TestTestCase_FilterValue(t *testing.T) {
	tests := []struct {
		name     string
		testCase TestCase
		want     string
	}{
		{
			name: "returns test name",
			testCase: TestCase{
				Name:     "my-test-file",
				Filepath: "/path/to/test.ts",
			},
			want: "my-test-file",
		},
		{
			name: "returns empty string when name is empty",
			testCase: TestCase{
				Name:     "",
				Filepath: "/path/to/test.ts",
			},
			want: "",
		},
		{
			name: "returns name with special characters",
			testCase: TestCase{
				Name:     "test-file_v2.spec",
				Filepath: "/path/to/test.ts",
			},
			want: "test-file_v2.spec",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.testCase.FilterValue()
			if got != tt.want {
				t.Errorf("FilterValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTestCase_TestStatusIcon(t *testing.T) {
	tests := []struct {
		name       string
		testStatus TestStatus
		want       string
	}{
		{
			name:       "running status",
			testStatus: StatusRunning,
			want:       IconRunning,
		},
		{
			name:       "passed status",
			testStatus: StatusPassed,
			want:       IconPassed,
		},
		{
			name:       "failed status",
			testStatus: StatusFailed,
			want:       IconFailed,
		},
		{
			name:       "not started status",
			testStatus: StatusNotStarted,
			want:       IconNotStarted,
		},
		{
			name:       "skipped status",
			testStatus: StatusSkipped,
			want:       IconSkipped,
		},
		{
			name:       "unknown status",
			testStatus: TestStatus("unknown"),
			want:       "❓ ",
		},
		{
			name:       "empty status",
			testStatus: TestStatus(""),
			want:       "❓ ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := &TestCase{
				TestStatus: tt.testStatus,
			}
			got := tc.TestStatusIcon()
			if got != tt.want {
				t.Errorf("TestStatusIcon() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTestStatus_Constants(t *testing.T) {
	// Verify that constants have expected values
	tests := []struct {
		name     string
		status   TestStatus
		expected string
	}{
		{"StatusNotStarted", StatusNotStarted, "not_run"},
		{"StatusPassed", StatusPassed, "passed"},
		{"StatusFailed", StatusFailed, "failed"},
		{"StatusSkipped", StatusSkipped, "skipped"},
		{"StatusRunning", StatusRunning, "running"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.status) != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, tt.status, tt.expected)
			}
		})
	}
}

func TestIconConstants(t *testing.T) {
	// Verify that icon constants are not empty and have expected values
	tests := []struct {
		name     string
		icon     string
		expected string
	}{
		{"IconPassed", IconPassed, "✓"},
		{"IconFailed", IconFailed, "✗"},
		{"IconNotStarted", IconNotStarted, "-"},
		{"IconRunning", IconRunning, "*"},
		{"IconSkipped", IconSkipped, "⚬"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.icon != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, tt.icon, tt.expected)
			}
			if tt.icon == "" {
				t.Errorf("%s should not be empty", tt.name)
			}
		})
	}
}

func TestTestCase_FullStructure(t *testing.T) {
	// Test that all fields can be set and retrieved correctly
	tc := TestCase{
		Name:       "test-name",
		Filepath:   "/path/to/file",
		Output:     "test output\nline 2",
		TestStatus: StatusPassed,
	}

	if tc.Name != "test-name" {
		t.Errorf("Name = %v, want 'test-name'", tc.Name)
	}

	if tc.Filepath != "/path/to/file" {
		t.Errorf("Filepath = %v, want '/path/to/file'", tc.Filepath)
	}

	if tc.Output != "test output\nline 2" {
		t.Errorf("Output = %v, want 'test output\\nline 2'", tc.Output)
	}

	if tc.TestStatus != StatusPassed {
		t.Errorf("TestStatus = %v, want StatusPassed", tc.TestStatus)
	}
}

func TestTestCase_ZeroValue(t *testing.T) {
	// Test zero value behavior
	var tc TestCase

	if tc.Name != "" {
		t.Errorf("Zero value Name = %v, want empty string", tc.Name)
	}

	if tc.Filepath != "" {
		t.Errorf("Zero value Filepath = %v, want empty string", tc.Filepath)
	}

	if tc.Output != "" {
		t.Errorf("Zero value Output = %v, want empty string", tc.Output)
	}

	if tc.TestStatus != "" {
		t.Errorf("Zero value TestStatus = %v, want empty string", tc.TestStatus)
	}

	// Zero value should return unknown icon
	if tc.TestStatusIcon() != "❓ " {
		t.Errorf("Zero value TestStatusIcon() = %v, want '❓ '", tc.TestStatusIcon())
	}
}

func TestTestCase_StatusTransitions(t *testing.T) {
	// Test that status can be changed and icon updates accordingly
	tc := TestCase{
		Name:       "test",
		Filepath:   "/test",
		TestStatus: StatusNotStarted,
	}

	if tc.TestStatusIcon() != IconNotStarted {
		t.Errorf("Initial icon = %v, want %v", tc.TestStatusIcon(), IconNotStarted)
	}

	tc.TestStatus = StatusRunning
	if tc.TestStatusIcon() != IconRunning {
		t.Errorf("Running icon = %v, want %v", tc.TestStatusIcon(), IconRunning)
	}

	tc.TestStatus = StatusPassed
	if tc.TestStatusIcon() != IconPassed {
		t.Errorf("Passed icon = %v, want %v", tc.TestStatusIcon(), IconPassed)
	}

	tc.TestStatus = StatusFailed
	if tc.TestStatusIcon() != IconFailed {
		t.Errorf("Failed icon = %v, want %v", tc.TestStatusIcon(), IconFailed)
	}

	tc.TestStatus = StatusSkipped
	if tc.TestStatusIcon() != IconSkipped {
		t.Errorf("Skipped icon = %v, want %v", tc.TestStatusIcon(), IconSkipped)
	}
}

func TestTestCase_LongOutput(t *testing.T) {
	// Test with long output string
	longOutput := ""
	for i := 0; i < 1000; i++ {
		longOutput += "test line\n"
	}

	tc := TestCase{
		Name:       "test",
		Filepath:   "/test",
		Output:     longOutput,
		TestStatus: StatusPassed,
	}

	if tc.Output != longOutput {
		t.Error("Long output not stored correctly")
	}

	if len(tc.Output) != len(longOutput) {
		t.Errorf("Output length = %v, want %v", len(tc.Output), len(longOutput))
	}
}

func TestTestCase_SpecialCharactersInFields(t *testing.T) {
	// Test with special characters
	tc := TestCase{
		Name:       "test with spaces & special chars: @#$%",
		Filepath:   "/path/with spaces/and-special@chars#test.spec.ts",
		Output:     "Output with\ttabs\nand\nnewlines\rand unicode: 🎉",
		TestStatus: StatusPassed,
	}

	if tc.FilterValue() != "test with spaces & special chars: @#$%" {
		t.Errorf("FilterValue() with special chars = %v", tc.FilterValue())
	}

	if tc.TestStatusIcon() != IconPassed {
		t.Error("TestStatusIcon() should work with special characters in other fields")
	}
}
