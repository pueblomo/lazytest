package detect

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// mockDriver implements the drivers.Driver interface for testing
type mockDriver struct {
	name         string
	detectResult bool
	detectErr    error
	wasCalled    bool
}

func (m *mockDriver) Name() string {
	return m.name
}

func (m *mockDriver) Detect(root string) (bool, error) {
	m.wasCalled = true
	return m.detectResult, m.detectErr
}

func (m *mockDriver) DetectTestFiles(ctx context.Context, root string) ([]string, error) {
	return nil, nil
}

func (m *mockDriver) RunTest(ctx context.Context, root string, testCase list.Item) error {
	return nil
}

// Since we can't replace the global drivers.AllDrivers function,
// we'll test with the actual drivers and create integration-style tests
// For true unit tests, we would need to refactor the detect package to accept
// a driver list parameter instead of calling drivers.AllDrivers() directly

func TestDetectDriver_ReturnsTeaCmd(t *testing.T) {
	cmd := DetectDriver("/nonexistent/path")

	if cmd == nil {
		t.Fatal("Expected DetectDriver to return a non-nil tea.Cmd")
	}

	// Verify it implements tea.Cmd by calling it
	msg := cmd()
	if msg == nil {
		t.Error("Expected cmd() to return a message")
	}
}

func TestDetectDriver_NoDriverFound(t *testing.T) {
	// Test with a path that won't match any driver
	cmd := DetectDriver("/nonexistent/path/without/any/test/framework")
	msg := cmd()

	detectMsg, ok := msg.(DriverDetectMsg)
	if !ok {
		t.Fatalf("Expected DriverDetectMsg, got %T", msg)
	}

	if detectMsg.Err == nil {
		t.Error("Expected error when no driver is found")
	}

	if detectMsg.Err.Error() != "no test driver found" {
		t.Errorf("Expected error 'no test driver found', got '%v'", detectMsg.Err)
	}

	if detectMsg.Driver != nil {
		t.Error("Expected driver to be nil when not found")
	}
}

func TestDriverDetectMsg_Structure(t *testing.T) {
	// Test that DriverDetectMsg can hold both driver and error
	msg := DriverDetectMsg{
		Driver: &mockDriver{name: "test"},
		Err:    fmt.Errorf("test error"),
	}

	if msg.Driver == nil {
		t.Error("Expected driver to be set")
	}

	if msg.Err == nil {
		t.Error("Expected error to be set")
	}

	// Verify it can be used as a tea.Msg
	var _ tea.Msg = msg
}

func TestDriverDetectMsg_OnlyError(t *testing.T) {
	// Test that DriverDetectMsg can have just an error
	msg := DriverDetectMsg{
		Driver: nil,
		Err:    fmt.Errorf("test error"),
	}

	if msg.Driver != nil {
		t.Error("Expected driver to be nil")
	}

	if msg.Err == nil {
		t.Error("Expected error to be set")
	}
}

func TestDriverDetectMsg_OnlyDriver(t *testing.T) {
	// Test that DriverDetectMsg can have just a driver
	msg := DriverDetectMsg{
		Driver: &mockDriver{name: "test"},
		Err:    nil,
	}

	if msg.Driver == nil {
		t.Error("Expected driver to be set")
	}

	if msg.Err != nil {
		t.Error("Expected error to be nil")
	}
}

func TestDetectDriver_EmptyRoot(t *testing.T) {
	cmd := DetectDriver("")
	msg := cmd()

	detectMsg, ok := msg.(DriverDetectMsg)
	if !ok {
		t.Fatalf("Expected DriverDetectMsg, got %T", msg)
	}

	// With empty root, detection should fail
	if detectMsg.Err == nil {
		t.Error("Expected error with empty root path")
	}
}

func TestDetectDriver_Logic(t *testing.T) {
	// This is an integration test that verifies the detection logic works
	// with real drivers. It tests the iteration logic and error handling.

	cmd := DetectDriver("/tmp")

	// Verify cmd is callable
	if cmd == nil {
		t.Fatal("Expected DetectDriver to return a non-nil function")
	}

	msg := cmd()

	// Should return a DriverDetectMsg
	detectMsg, ok := msg.(DriverDetectMsg)
	if !ok {
		t.Fatalf("Expected DriverDetectMsg, got %T", msg)
	}

	// Either a driver was found or an error occurred
	if detectMsg.Driver == nil && detectMsg.Err == nil {
		t.Error("Expected either driver or error to be set")
	}
}

func TestDetectDriver_WithRealVitestProject(t *testing.T) {
	// Create a temporary directory with a real Vitest project
	tmpDir := t.TempDir()

	// Create package.json with vitest test script
	pkgJSON := `{
		"name": "test-project",
		"scripts": {
			"test": "vitest"
		}
	}`

	pkgPath := tmpDir + "/package.json"
	if err := os.WriteFile(pkgPath, []byte(pkgJSON), 0644); err != nil {
		t.Fatalf("Failed to create package.json: %v", err)
	}

	// Run detection
	cmd := DetectDriver(tmpDir)
	msg := cmd()

	detectMsg, ok := msg.(DriverDetectMsg)
	if !ok {
		t.Fatalf("Expected DriverDetectMsg, got %T", msg)
	}

	// Should successfully detect Vitest
	if detectMsg.Err != nil {
		t.Errorf("Expected no error, got: %v", detectMsg.Err)
	}

	if detectMsg.Driver == nil {
		t.Fatal("Expected driver to be found")
	}

	if detectMsg.Driver.Name() != "vitest" {
		t.Errorf("Expected vitest driver, got %s", detectMsg.Driver.Name())
	}
}

func TestDetectDriver_IteratesAllDriversOnError(t *testing.T) {
	// Create a directory that will cause detection errors
	tmpDir := t.TempDir()

	// Create an invalid package.json that will cause errors
	pkgPath := tmpDir + "/package.json"
	if err := os.WriteFile(pkgPath, []byte("invalid json {{{"), 0644); err != nil {
		t.Fatalf("Failed to create package.json: %v", err)
	}

	cmd := DetectDriver(tmpDir)
	msg := cmd()

	detectMsg, ok := msg.(DriverDetectMsg)
	if !ok {
		t.Fatalf("Expected DriverDetectMsg, got %T", msg)
	}

	// Should continue past errors and return "no driver found"
	if detectMsg.Err == nil {
		t.Error("Expected error when detection fails")
	}

	if detectMsg.Err.Error() != "no test driver found" {
		t.Errorf("Expected 'no test driver found', got: %v", detectMsg.Err)
	}

	if detectMsg.Driver != nil {
		t.Error("Expected driver to be nil when none found")
	}
}

func TestDetectDriver_StopsAtFirstMatch(t *testing.T) {
	// Create a valid Vitest project
	tmpDir := t.TempDir()

	pkgJSON := `{
		"name": "test-project",
		"scripts": {
			"test": "vitest run"
		}
	}`

	pkgPath := tmpDir + "/package.json"
	if err := os.WriteFile(pkgPath, []byte(pkgJSON), 0644); err != nil {
		t.Fatalf("Failed to create package.json: %v", err)
	}

	cmd := DetectDriver(tmpDir)
	msg := cmd()

	detectMsg, ok := msg.(DriverDetectMsg)
	if !ok {
		t.Fatalf("Expected DriverDetectMsg, got %T", msg)
	}

	// Should find the first matching driver
	if detectMsg.Driver == nil {
		t.Fatal("Expected driver to be found")
	}

	if detectMsg.Err != nil {
		t.Errorf("Expected no error, got: %v", detectMsg.Err)
	}

	// Verify it's the Vitest driver (first and only driver currently)
	if detectMsg.Driver.Name() != "vitest" {
		t.Errorf("Expected vitest driver, got %s", detectMsg.Driver.Name())
	}
}
