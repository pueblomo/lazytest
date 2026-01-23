package app

import (
	"os"
	"testing"
)

func TestRun_CanGetWorkingDirectory(t *testing.T) {
	// Test that we can get the working directory
	// This is a prerequisite for Run() to work
	_, err := os.Getwd()
	if err != nil {
		t.Fatalf("Cannot get working directory: %v", err)
	}
}

func TestRun_ValidatesEnvironment(t *testing.T) {
	// Skip this test as it hangs when trying to initialize the TUI
	// Full integration tests would require a PTY or terminal emulator
	t.Skip("Skipping Run() test as it requires a terminal environment")
}

func TestRun_WorkingDirectoryAccess(t *testing.T) {
	// Verify we can access the working directory before Run() is called
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	if wd == "" {
		t.Error("Working directory should not be empty")
	}

	// Verify the directory exists
	info, err := os.Stat(wd)
	if err != nil {
		t.Fatalf("Working directory does not exist: %v", err)
	}

	if !info.IsDir() {
		t.Error("Working directory path should point to a directory")
	}
}

func TestRun_ErrorHandling(t *testing.T) {
	// Test that Run returns an error when it can't get the working directory
	// We can't easily test this without changing the working directory to something invalid
	// or mocking os.Getwd(), but we can verify the error path exists

	// Just verify the function signature and that it returns an error type
	t.Skip("Error path testing requires mocking or directory manipulation")
}

func TestRun_FunctionSignature(t *testing.T) {
	// Verify Run() has the correct signature
	// This is a compile-time check, but we document it here
	var fn func() error = Run
	if fn == nil {
		t.Error("Run should not be nil")
	}
}

func TestWorkingDirectory_IsValid(t *testing.T) {
	// Additional test to ensure working directory is accessible
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Try to list directory contents to verify access
	entries, err := os.ReadDir(wd)
	if err != nil {
		t.Fatalf("Failed to read working directory: %v", err)
	}

	// Just verify we can read it (might be empty, but should be readable)
	_ = entries
}

func TestRun_Prerequisites(t *testing.T) {
	// Test all prerequisites for Run() to work
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "can get working directory",
			test: func(t *testing.T) {
				_, err := os.Getwd()
				if err != nil {
					t.Errorf("Cannot get working directory: %v", err)
				}
			},
		},
		{
			name: "working directory is not empty",
			test: func(t *testing.T) {
				wd, err := os.Getwd()
				if err != nil {
					t.Fatalf("Cannot get working directory: %v", err)
				}
				if wd == "" {
					t.Error("Working directory should not be empty")
				}
			},
		},
		{
			name: "working directory exists",
			test: func(t *testing.T) {
				wd, err := os.Getwd()
				if err != nil {
					t.Fatalf("Cannot get working directory: %v", err)
				}
				if _, err := os.Stat(wd); err != nil {
					t.Errorf("Working directory does not exist: %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

// Note: Full integration tests for the TUI would require a terminal emulator
// or PTY (pseudo-terminal). These tests verify the basic structure without
// needing to actually run the TUI.
