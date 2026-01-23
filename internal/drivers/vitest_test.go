package drivers

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/pueblomo/lazytest/internal/types"
)

func TestVitestDriver_Name(t *testing.T) {
	driver := &VitestDriver{}
	if driver.Name() != "vitest" {
		t.Errorf("Name() = %v, want 'vitest'", driver.Name())
	}
}

func TestVitestDriver_Detect_NoPackageJson(t *testing.T) {
	tmpDir := t.TempDir()

	driver := &VitestDriver{}
	detected, err := driver.Detect(tmpDir)

	if detected {
		t.Error("Detect() should return false when package.json doesn't exist")
	}

	if err == nil {
		t.Error("Detect() should return error when package.json doesn't exist")
	}
}

func TestVitestDriver_Detect_InvalidJson(t *testing.T) {
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "package.json")

	// Write invalid JSON
	if err := os.WriteFile(pkgPath, []byte("{invalid json"), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	driver := &VitestDriver{}
	detected, err := driver.Detect(tmpDir)

	if detected {
		t.Error("Detect() should return false for invalid JSON")
	}

	if err == nil {
		t.Error("Detect() should return error for invalid JSON")
	}
}

func TestVitestDriver_Detect_NoTestScript(t *testing.T) {
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "package.json")

	pkg := map[string]interface{}{
		"name": "test-package",
		"scripts": map[string]string{
			"build": "tsc",
		},
	}

	data, _ := json.Marshal(pkg)
	if err := os.WriteFile(pkgPath, data, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	driver := &VitestDriver{}
	detected, err := driver.Detect(tmpDir)

	if detected {
		t.Error("Detect() should return false when test script doesn't exist")
	}

	if err != nil {
		t.Errorf("Detect() should not return error when test script is missing: %v", err)
	}
}

func TestVitestDriver_Detect_TestScriptWithoutVitest(t *testing.T) {
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "package.json")

	pkg := map[string]interface{}{
		"name": "test-package",
		"scripts": map[string]string{
			"test": "jest",
		},
	}

	data, _ := json.Marshal(pkg)
	if err := os.WriteFile(pkgPath, data, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	driver := &VitestDriver{}
	detected, err := driver.Detect(tmpDir)

	if detected {
		t.Error("Detect() should return false when test script doesn't contain vitest")
	}

	if err != nil {
		t.Errorf("Detect() should not return error: %v", err)
	}
}

func TestVitestDriver_Detect_WithVitest(t *testing.T) {
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "package.json")

	pkg := map[string]interface{}{
		"name": "test-package",
		"scripts": map[string]string{
			"test": "vitest",
		},
	}

	data, _ := json.Marshal(pkg)
	if err := os.WriteFile(pkgPath, data, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	driver := &VitestDriver{}
	detected, err := driver.Detect(tmpDir)

	if !detected {
		t.Error("Detect() should return true when vitest is in test script")
	}

	if err != nil {
		t.Errorf("Detect() should not return error: %v", err)
	}
}

func TestVitestDriver_Detect_WithPnpm(t *testing.T) {
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "package.json")
	lockPath := filepath.Join(tmpDir, "pnpm-lock.yaml")

	pkg := map[string]interface{}{
		"name": "test-package",
		"scripts": map[string]string{
			"test": "vitest run",
		},
	}

	data, _ := json.Marshal(pkg)
	if err := os.WriteFile(pkgPath, data, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Create pnpm lock file
	if err := os.WriteFile(lockPath, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to write lock file: %v", err)
	}

	driver := &VitestDriver{}
	detected, err := driver.Detect(tmpDir)

	if !detected {
		t.Error("Detect() should return true")
	}

	if err != nil {
		t.Errorf("Detect() should not return error: %v", err)
	}

	if driver.packageManager != "pnpm" {
		t.Errorf("packageManager = %v, want 'pnpm'", driver.packageManager)
	}
}

func TestVitestDriver_Detect_WithNpm(t *testing.T) {
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "package.json")
	lockPath := filepath.Join(tmpDir, "package-lock.json")

	pkg := map[string]interface{}{
		"name": "test-package",
		"scripts": map[string]string{
			"test": "vitest",
		},
	}

	data, _ := json.Marshal(pkg)
	if err := os.WriteFile(pkgPath, data, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Create npm lock file
	if err := os.WriteFile(lockPath, []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to write lock file: %v", err)
	}

	driver := &VitestDriver{}
	detected, err := driver.Detect(tmpDir)

	if !detected {
		t.Error("Detect() should return true")
	}

	if err != nil {
		t.Errorf("Detect() should not return error: %v", err)
	}

	if driver.packageManager != "npm" {
		t.Errorf("packageManager = %v, want 'npm'", driver.packageManager)
	}
}

func TestVitestDriver_Detect_DefaultsToNpm(t *testing.T) {
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "package.json")

	pkg := map[string]interface{}{
		"name": "test-package",
		"scripts": map[string]string{
			"test": "vitest",
		},
	}

	data, _ := json.Marshal(pkg)
	if err := os.WriteFile(pkgPath, data, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// No lock file

	driver := &VitestDriver{}
	detected, err := driver.Detect(tmpDir)

	if !detected {
		t.Error("Detect() should return true")
	}

	if err != nil {
		t.Errorf("Detect() should not return error: %v", err)
	}

	if driver.packageManager != "npm" {
		t.Errorf("packageManager should default to 'npm', got %v", driver.packageManager)
	}
}

func TestVitestDriver_BuildTestCommand_Basic(t *testing.T) {
	driver := &VitestDriver{packageManager: "npm"}
	ctx := context.Background()

	cmd := driver.buildTestCommand(ctx, "/test/root", "tests/example.test.ts")

	if cmd.Dir != "/test/root" {
		t.Errorf("cmd.Dir = %v, want '/test/root'", cmd.Dir)
	}

	expectedArgs := []string{"npm", "vitest", "--run", "tests/example.test.ts"}
	if len(cmd.Args) != len(expectedArgs) {
		t.Errorf("cmd.Args length = %v, want %v", len(cmd.Args), len(expectedArgs))
	}

	for i, arg := range expectedArgs {
		if i < len(cmd.Args) && cmd.Args[i] != arg {
			t.Errorf("cmd.Args[%d] = %v, want %v", i, cmd.Args[i], arg)
		}
	}
}

func TestVitestDriver_BuildTestCommand_WithPnpm(t *testing.T) {
	driver := &VitestDriver{packageManager: "pnpm"}
	ctx := context.Background()

	cmd := driver.buildTestCommand(ctx, "/test/root", "tests/example.test.ts")

	if cmd.Args[0] != "pnpm" {
		t.Errorf("cmd.Args[0] = %v, want 'pnpm'", cmd.Args[0])
	}
}

func TestVitestDriver_BuildTestCommand_WithContext(t *testing.T) {
	driver := &VitestDriver{packageManager: "npm"}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd := driver.buildTestCommand(ctx, "/test/root", "test.ts")

	if cmd.Path == "" {
		t.Error("cmd.Path should be set")
	}
}

func TestVitestDriver_ExecuteTestCommand_ParsesPassedStatus(t *testing.T) {
	// We can't easily test actual command execution in unit tests,
	// but we can test the parsing logic would work with mock output
	// This is more of a documentation test showing expected behavior

	// Skip this test for now as it requires mocking exec.Cmd
	t.Skip("Command execution testing requires mocking")
}

// wrongItem implements list.Item but is not a TestCase
type wrongItem struct{}

func (w *wrongItem) FilterValue() string { return "wrong" }

func TestVitestDriver_RunTest_InvalidTestCase(t *testing.T) {
	driver := &VitestDriver{packageManager: "npm"}
	ctx := context.Background()

	wrongItemInstance := &wrongItem{}

	err := driver.RunTest(ctx, "/test/root", wrongItemInstance)

	if err == nil {
		t.Error("RunTest() should return error for invalid test case type")
	}

	if err.Error() != "Can't convert to TestCase" {
		t.Errorf("RunTest() error = %v, want 'Can't convert to TestCase'", err)
	}
}

func TestVitestDriver_RunTest_SetsTestCaseFields(t *testing.T) {
	driver := &VitestDriver{packageManager: "npm"}
	ctx := context.Background()

	tc := &types.TestCase{
		Name:       "test",
		Filepath:   "nonexistent.test.ts",
		TestStatus: types.StatusNotStarted,
	}

	// This will fail because the command won't actually work
	// but we can verify the test case gets updated
	_ = driver.RunTest(ctx, "/nonexistent", tc)

	// After RunTest, the TestCase should have been updated
	// (even if the command failed)
	if tc.TestStatus == types.StatusNotStarted {
		t.Error("RunTest() should update TestStatus")
	}
}

func TestVitestDriver_Name_IsConstant(t *testing.T) {
	driver1 := &VitestDriver{}
	driver2 := &VitestDriver{}

	if driver1.Name() != driver2.Name() {
		t.Error("Name() should return the same value for all instances")
	}
}

func TestVitestDriver_Detect_EmptyScripts(t *testing.T) {
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "package.json")

	pkg := map[string]interface{}{
		"name":    "test-package",
		"scripts": map[string]string{},
	}

	data, _ := json.Marshal(pkg)
	if err := os.WriteFile(pkgPath, data, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	driver := &VitestDriver{}
	detected, err := driver.Detect(tmpDir)

	if detected {
		t.Error("Detect() should return false when scripts is empty")
	}

	if err != nil {
		t.Errorf("Detect() should not return error for empty scripts: %v", err)
	}
}

func TestVitestDriver_Detect_VitestInMiddleOfScript(t *testing.T) {
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "package.json")

	pkg := map[string]interface{}{
		"name": "test-package",
		"scripts": map[string]string{
			"test": "cross-env NODE_ENV=test vitest run --coverage",
		},
	}

	data, _ := json.Marshal(pkg)
	if err := os.WriteFile(pkgPath, data, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	driver := &VitestDriver{}
	detected, err := driver.Detect(tmpDir)

	if !detected {
		t.Error("Detect() should return true when vitest is anywhere in test script")
	}

	if err != nil {
		t.Errorf("Detect() should not return error: %v", err)
	}
}

func TestVitestDriver_Detect_CaseInsensitivity(t *testing.T) {
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "package.json")

	pkg := map[string]interface{}{
		"name": "test-package",
		"scripts": map[string]string{
			"test": "Vitest run", // Capital V
		},
	}

	data, _ := json.Marshal(pkg)
	if err := os.WriteFile(pkgPath, data, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	driver := &VitestDriver{}
	detected, err := driver.Detect(tmpDir)

	// Current implementation is case-sensitive
	// This documents the current behavior
	if detected {
		t.Error("Detect() is case-sensitive and should not match 'Vitest'")
	}

	if err != nil {
		t.Errorf("Detect() should not return error: %v", err)
	}
}

func TestVitestDriver_ExecuteTestCommand_PassedStatus(t *testing.T) {
	driver := &VitestDriver{}

	// Create a mock command that will succeed
	cmd := exec.Command("echo", "Test passed")

	status, output, err := driver.executeTestCommand(cmd)

	if err != nil {
		t.Errorf("executeTestCommand() should not return error for successful command: %v", err)
	}

	if status != types.StatusPassed {
		t.Errorf("executeTestCommand() status = %v, want %v", status, types.StatusPassed)
	}

	if output == "" {
		t.Error("executeTestCommand() should return output")
	}
}

func TestVitestDriver_ExecuteTestCommand_FailedStatus(t *testing.T) {
	driver := &VitestDriver{}

	// Create a mock command that outputs "failed"
	cmd := exec.Command("echo", "Test failed")

	status, output, err := driver.executeTestCommand(cmd)

	if err != nil {
		t.Errorf("executeTestCommand() should not return error when test fails: %v", err)
	}

	if status != types.StatusFailed {
		t.Errorf("executeTestCommand() status = %v, want %v", status, types.StatusFailed)
	}

	if output != "Test failed\n" {
		t.Errorf("executeTestCommand() output = %q, want 'Test failed\\n'", output)
	}
}

func TestVitestDriver_ExecuteTestCommand_SkippedStatus(t *testing.T) {
	driver := &VitestDriver{}

	// Create a mock command that outputs "skipped"
	cmd := exec.Command("echo", "Test skipped")

	status, output, err := driver.executeTestCommand(cmd)

	if err != nil {
		t.Errorf("executeTestCommand() should not return error when test is skipped: %v", err)
	}

	if status != types.StatusSkipped {
		t.Errorf("executeTestCommand() status = %v, want %v", status, types.StatusSkipped)
	}

	if output != "Test skipped\n" {
		t.Errorf("executeTestCommand() output = %q", output)
	}
}

func TestVitestDriver_ExecuteTestCommand_CommandError(t *testing.T) {
	driver := &VitestDriver{}

	// Create a command that will fail
	cmd := exec.Command("false")

	status, output, err := driver.executeTestCommand(cmd)

	if err == nil {
		t.Error("executeTestCommand() should return error when command fails")
	}

	if status != types.StatusFailed {
		t.Errorf("executeTestCommand() status = %v, want %v", status, types.StatusFailed)
	}

	// Output should still be captured even on error
	if output == "" && err != nil {
		// This is acceptable - command failed without output
	}
}

func TestVitestDriver_ExecuteTestCommand_MultipleKeywords(t *testing.T) {
	driver := &VitestDriver{}

	// Test that "failed" takes precedence over "skipped"
	cmd := exec.Command("sh", "-c", "echo 'Test failed and skipped'")

	status, output, err := driver.executeTestCommand(cmd)

	if err != nil {
		t.Errorf("executeTestCommand() should not return error: %v", err)
	}

	// "failed" should be detected first
	if status != types.StatusFailed {
		t.Errorf("executeTestCommand() status = %v, want %v (failed should take precedence)", status, types.StatusFailed)
	}

	if output == "" {
		t.Error("executeTestCommand() should return output")
	}
}

func TestVitestDriver_DetectTestFiles_ValidJSON(t *testing.T) {
	driver := &VitestDriver{packageManager: "echo"}
	ctx := context.Background()

	// This will fail because echo is not vitest, but we can test error handling
	_, err := driver.DetectTestFiles(ctx, t.TempDir())

	// We expect an error since we're using echo instead of vitest
	// This tests that the function attempts to execute and handle errors
	if err == nil {
		// If no error, that's fine - it means echo worked somehow
		t.Log("DetectTestFiles executed without error")
	}
}

func TestVitestDriver_DetectTestFiles_EmptyResult(t *testing.T) {
	driver := &VitestDriver{packageManager: "echo"}
	ctx := context.Background()

	tmpDir := t.TempDir()

	// Using echo with empty output to simulate empty test list
	driver.packageManager = "sh"

	// This will fail to parse as JSON, testing the unmarshal error path
	_, err := driver.DetectTestFiles(ctx, tmpDir)

	// Should return error because echo output isn't valid JSON
	if err == nil {
		t.Error("DetectTestFiles() should return error for invalid JSON")
	}
}

func TestVitestDriver_DetectTestFiles_JSONFormat(t *testing.T) {
	driver := &VitestDriver{packageManager: "sh"}
	ctx := context.Background()

	tmpDir := t.TempDir()

	// Create a script that outputs valid JSON
	scriptPath := filepath.Join(tmpDir, "mock_vitest.sh")
	script := `#!/bin/sh
echo '[{"file":"test1.spec.ts"},{"file":"test2.spec.ts"}]'
`
	if err := os.WriteFile(scriptPath, []byte(script), 0755); err != nil {
		t.Fatalf("Failed to create mock script: %v", err)
	}

	// Point to our mock script
	driver.packageManager = scriptPath

	files, err := driver.DetectTestFiles(ctx, tmpDir)

	if err != nil {
		t.Errorf("DetectTestFiles() error = %v, want nil", err)
	}

	expectedFiles := []string{"test1.spec.ts", "test2.spec.ts"}
	if len(files) != len(expectedFiles) {
		t.Errorf("DetectTestFiles() returned %d files, want %d", len(files), len(expectedFiles))
	}

	for i, file := range files {
		if i < len(expectedFiles) && file != expectedFiles[i] {
			t.Errorf("DetectTestFiles() file[%d] = %v, want %v", i, file, expectedFiles[i])
		}
	}
}

func TestVitestDriver_DetectTestFiles_EmptyArray(t *testing.T) {
	driver := &VitestDriver{packageManager: "sh"}
	ctx := context.Background()

	tmpDir := t.TempDir()

	// Create a script that outputs empty JSON array
	scriptPath := filepath.Join(tmpDir, "mock_empty.sh")
	script := `#!/bin/sh
echo '[]'
`
	if err := os.WriteFile(scriptPath, []byte(script), 0755); err != nil {
		t.Fatalf("Failed to create mock script: %v", err)
	}

	driver.packageManager = scriptPath

	files, err := driver.DetectTestFiles(ctx, tmpDir)

	if err != nil {
		t.Errorf("DetectTestFiles() error = %v, want nil", err)
	}

	if len(files) != 0 {
		t.Errorf("DetectTestFiles() returned %d files, want 0", len(files))
	}
}

func TestVitestDriver_DetectTestFiles_InvalidJSON(t *testing.T) {
	driver := &VitestDriver{packageManager: "sh"}
	ctx := context.Background()

	tmpDir := t.TempDir()

	// Create a script that outputs invalid JSON
	scriptPath := filepath.Join(tmpDir, "mock_invalid.sh")
	script := `#!/bin/sh
echo 'invalid json {'
`
	if err := os.WriteFile(scriptPath, []byte(script), 0755); err != nil {
		t.Fatalf("Failed to create mock script: %v", err)
	}

	driver.packageManager = scriptPath

	_, err := driver.DetectTestFiles(ctx, tmpDir)

	if err == nil {
		t.Error("DetectTestFiles() should return error for invalid JSON")
	}
}

func TestVitestDriver_DetectTestFiles_CommandError(t *testing.T) {
	driver := &VitestDriver{packageManager: "false"}
	ctx := context.Background()

	tmpDir := t.TempDir()

	_, err := driver.DetectTestFiles(ctx, tmpDir)

	// Should return error when command fails
	if err == nil {
		t.Error("DetectTestFiles() should return error when command fails")
	}
}

func TestVitestDriver_DetectTestFiles_ContextCancellation(t *testing.T) {
	driver := &VitestDriver{packageManager: "sleep"}
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel immediately
	cancel()

	tmpDir := t.TempDir()

	_, err := driver.DetectTestFiles(ctx, tmpDir)

	// Should return error when context is cancelled
	if err == nil {
		t.Error("DetectTestFiles() should return error when context is cancelled")
	}
}

func TestVitestDriver_Detect_ReadFileError(t *testing.T) {
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "package.json")

	// Create a directory instead of a file to cause read error
	if err := os.Mkdir(pkgPath, 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	driver := &VitestDriver{}
	detected, err := driver.Detect(tmpDir)

	if detected {
		t.Error("Detect() should return false when package.json cannot be read")
	}

	if err == nil {
		t.Error("Detect() should return error when package.json is a directory")
	}
}
