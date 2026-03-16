package drivers

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pueblomo/lazytest/internal/types"
)

func TestGoTestDriver_Name(t *testing.T) {
	driver := &GoTestDriver{}
	if driver.Name() != "go test" {
		t.Errorf("Expected name to be 'go test', got %s", driver.Name())
	}
}

func TestGoTestDriver_Detect_NoGoMod(t *testing.T) {
	tmpDir := t.TempDir()

	driver := &GoTestDriver{}
	detected, err := driver.Detect(tmpDir)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if detected {
		t.Error("Expected false when go.mod doesn't exist")
	}
}

func TestGoTestDriver_Detect_GoModButNoTests(t *testing.T) {
	tmpDir := t.TempDir()

	// Create go.mod
	goModPath := filepath.Join(tmpDir, "go.mod")
	goModContent := []byte("module example.com/test\n\ngo 1.21\n")
	if err := os.WriteFile(goModPath, goModContent, 0644); err != nil {
		t.Fatal(err)
	}

	// Create a non-test file
	mainFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(mainFile, []byte("package main\n"), 0644); err != nil {
		t.Fatal(err)
	}

	driver := &GoTestDriver{}
	detected, err := driver.Detect(tmpDir)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if detected {
		t.Error("Expected false when no test files exist")
	}
}

func TestGoTestDriver_Detect_WithTests(t *testing.T) {
	tmpDir := t.TempDir()

	// Create go.mod
	goModPath := filepath.Join(tmpDir, "go.mod")
	goModContent := []byte("module example.com/test\n\ngo 1.21\n")
	if err := os.WriteFile(goModPath, goModContent, 0644); err != nil {
		t.Fatal(err)
	}

	// Create a test file
	testFile := filepath.Join(tmpDir, "main_test.go")
	if err := os.WriteFile(testFile, []byte("package main\n\nimport \"testing\"\n"), 0644); err != nil {
		t.Fatal(err)
	}

	driver := &GoTestDriver{}
	detected, err := driver.Detect(tmpDir)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !detected {
		t.Error("Expected true when test files exist")
	}
}

func TestGoTestDriver_Detect_WithTestsInSubdir(t *testing.T) {
	tmpDir := t.TempDir()

	// Create go.mod
	goModPath := filepath.Join(tmpDir, "go.mod")
	goModContent := []byte("module example.com/test\n\ngo 1.21\n")
	if err := os.WriteFile(goModPath, goModContent, 0644); err != nil {
		t.Fatal(err)
	}

	// Create subdirectory
	subDir := filepath.Join(tmpDir, "pkg")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a test file in subdirectory
	testFile := filepath.Join(subDir, "utils_test.go")
	if err := os.WriteFile(testFile, []byte("package pkg\n"), 0644); err != nil {
		t.Fatal(err)
	}

	driver := &GoTestDriver{}
	detected, err := driver.Detect(tmpDir)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !detected {
		t.Error("Expected true when test files exist in subdirectory")
	}
}

func TestGoTestDriver_BuildTestCommand_Basic(t *testing.T) {
	driver := &GoTestDriver{}
	ctx := context.Background()
	root := "/project"
	filePath := "pkg/utils_test.go"

	cmd := driver.buildTestCommand(ctx, root, filePath)

	if cmd.Args[0] != "go" {
		t.Errorf("Expected command to be 'go', got %s", cmd.Args[0])
	}

	expectedArgs := []string{"go", "test", "-v", "./pkg"}
	if len(cmd.Args) != len(expectedArgs) {
		t.Errorf("Expected %d args, got %d", len(expectedArgs), len(cmd.Args))
	}

	for i, arg := range expectedArgs {
		if cmd.Args[i] != arg {
			t.Errorf("Expected arg[%d] to be %s, got %s", i, arg, cmd.Args[i])
		}
	}

	if cmd.Dir != root {
		t.Errorf("Expected dir to be %s, got %s", root, cmd.Dir)
	}
}

func TestGoTestDriver_BuildTestCommand_RootLevel(t *testing.T) {
	driver := &GoTestDriver{}
	ctx := context.Background()
	root := "/project"
	filePath := "main_test.go"

	cmd := driver.buildTestCommand(ctx, root, filePath)

	expectedArgs := []string{"go", "test", "-v", "./."}
	if len(cmd.Args) != len(expectedArgs) {
		t.Errorf("Expected %d args, got %d", len(expectedArgs), len(cmd.Args))
	}

	for i, arg := range expectedArgs {
		if cmd.Args[i] != arg {
			t.Errorf("Expected arg[%d] to be %s, got %s", i, arg, cmd.Args[i])
		}
	}
}

func TestGoTestDriver_ExecuteTestCommand_PassedStatus(t *testing.T) {
	driver := &GoTestDriver{}
	cmd := exec.Command("echo", "PASS\nok  \texample.com/test\t0.123s")

	status, output, err := driver.executeTestCommand(cmd)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if status != types.StatusPassed {
		t.Errorf("Expected status to be StatusPassed, got %s", status)
	}

	if !strings.Contains(output, "PASS") {
		t.Error("Expected output to contain 'PASS'")
	}
}

func TestGoTestDriver_ExecuteTestCommand_FailedStatus(t *testing.T) {
	driver := &GoTestDriver{}
	// Use a command that will exit with non-zero but produce output with FAIL
	cmd := exec.Command("sh", "-c", "echo 'FAIL\nexample.com/test\t0.123s' && exit 1")

	status, output, err := driver.executeTestCommand(cmd)

	if err != nil {
		t.Errorf("Expected no error for test failure, got %v", err)
	}

	if status != types.StatusFailed {
		t.Errorf("Expected status to be StatusFailed, got %s", status)
	}

	if !strings.Contains(output, "FAIL") {
		t.Error("Expected output to contain 'FAIL'")
	}
}

func TestGoTestDriver_ExecuteTestCommand_SkippedStatus(t *testing.T) {
	driver := &GoTestDriver{}
	cmd := exec.Command("echo", "SKIP\nok  \texample.com/test\t0.001s")

	status, output, err := driver.executeTestCommand(cmd)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if status != types.StatusSkipped {
		t.Errorf("Expected status to be StatusSkipped, got %s", status)
	}

	if !strings.Contains(output, "SKIP") {
		t.Error("Expected output to contain 'SKIP'")
	}
}

func TestGoTestDriver_ExecuteTestCommand_CommandError(t *testing.T) {
	driver := &GoTestDriver{}
	cmd := exec.Command("sh", "-c", "echo 'build failed' && exit 1")

	status, output, err := driver.executeTestCommand(cmd)

	// ExitError is treated as test failure, not command error
	if err != nil {
		t.Errorf("Expected no error for ExitError (test failure): %v", err)
	}

	if status != types.StatusFailed {
		t.Errorf("Expected status to be StatusFailed, got %s", status)
	}

	if !strings.Contains(output, "build failed") {
		t.Error("Expected output to contain error message")
	}
}

type wrongGoItem struct{}

func (w *wrongGoItem) FilterValue() string { return "" }

func TestGoTestDriver_RunTest_InvalidTestCase(t *testing.T) {
	driver := &GoTestDriver{}
	ctx := context.Background()

	err := driver.RunTest(ctx, "/tmp", &wrongGoItem{})

	if err == nil {
		t.Error("Expected error when item is not a TestCase")
	}

	if !strings.Contains(err.Error(), "can't convert to TestCase") {
		t.Errorf("Expected error message about conversion, got: %v", err)
	}
}

func TestGoTestDriver_RunTest_SetsTestCaseFields(t *testing.T) {
	driver := &GoTestDriver{}
	ctx := context.Background()
	tmpDir := t.TempDir()

	tc := &types.TestCase{
		Name:       "TestExample",
		Filepath:   "main_test.go",
		TestStatus: types.StatusNotStarted,
	}

	// This will fail since there's no actual Go project, but we're testing field assignment
	_ = driver.RunTest(ctx, tmpDir, tc)

	// The test status should have changed from NotStarted
	if tc.TestStatus == types.StatusNotStarted {
		t.Error("Expected test status to change after running")
	}

	// Output should be set (even if empty or contains error)
	// We're not checking the exact value since it depends on the command result
}

func TestGoTestDriver_Name_IsConstant(t *testing.T) {
	driver1 := &GoTestDriver{}
	driver2 := &GoTestDriver{}

	if driver1.Name() != driver2.Name() {
		t.Error("Expected Name() to return consistent value across instances")
	}

	if driver1.Name() != "go test" {
		t.Errorf("Expected name to be 'go test', got %s", driver1.Name())
	}
}

func TestGoTestDriver_DetectTestFiles_NoGoProject(t *testing.T) {
	driver := &GoTestDriver{}
	ctx := context.Background()
	tmpDir := t.TempDir()

	files, err := driver.DetectTestFiles(ctx, tmpDir)

	// Should return an error since it's not a valid Go project
	if err == nil {
		t.Error("Expected error when detecting test files in non-Go project")
	}

	if files != nil {
		t.Error("Expected nil files when error occurs")
	}
}

func TestGoTestDriver_Detect_ReadError(t *testing.T) {
	driver := &GoTestDriver{}

	// Use a path that exists but can't be read (permission denied scenario is hard to test)
	// Instead test with a file (not directory) as root
	tmpFile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// Close and try to use as root directory
	tmpFile.Close()

	detected, err := driver.Detect(tmpFile.Name())

	if detected {
		t.Error("Expected false when path is invalid")
	}

	// Error may or may not occur depending on implementation
	_ = err
}

func TestGoTestDriver_BuildTestCommand_WithContext(t *testing.T) {
	driver := &GoTestDriver{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	root := "/project"
	filePath := "pkg/test_test.go"

	cmd := driver.buildTestCommand(ctx, root, filePath)

	if cmd.Args[0] != "go" {
		t.Errorf("Expected command to be 'go', got %s", cmd.Args[0])
	}

	// Verify context is properly set (command should respect cancellation)
	if cmd.Cancel == nil {
		t.Error("Expected command to have cancel function from context")
	}
}

func TestGoTestDriver_Detect_MultipleTestFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create go.mod
	goModPath := filepath.Join(tmpDir, "go.mod")
	goModContent := []byte("module example.com/test\n\ngo 1.21\n")
	if err := os.WriteFile(goModPath, goModContent, 0644); err != nil {
		t.Fatal(err)
	}

	// Create multiple test files
	testFiles := []string{"main_test.go", "utils_test.go", "helpers_test.go"}
	for _, testFile := range testFiles {
		path := filepath.Join(tmpDir, testFile)
		if err := os.WriteFile(path, []byte("package main\n"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	driver := &GoTestDriver{}
	detected, err := driver.Detect(tmpDir)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !detected {
		t.Error("Expected true when multiple test files exist")
	}
}

func TestGoTestDriver_BuildTestCommand_NestedPath(t *testing.T) {
	driver := &GoTestDriver{}
	ctx := context.Background()
	root := "/project"
	filePath := "internal/pkg/subpkg/handler_test.go"

	cmd := driver.buildTestCommand(ctx, root, filePath)

	expectedArgs := []string{"go", "test", "-v", "./internal/pkg/subpkg"}
	if len(cmd.Args) != len(expectedArgs) {
		t.Errorf("Expected %d args, got %d", len(expectedArgs), len(cmd.Args))
	}

	for i, arg := range expectedArgs {
		if cmd.Args[i] != arg {
			t.Errorf("Expected arg[%d] to be %s, got %s", i, arg, cmd.Args[i])
		}
	}
}

func TestGoTestDriver_ExecuteTestCommand_EmptyOutput(t *testing.T) {
	driver := &GoTestDriver{}
	cmd := exec.Command("echo", "")

	status, output, err := driver.executeTestCommand(cmd)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// When no specific status markers are found, default to passed
	if status != types.StatusPassed {
		t.Errorf("Expected status to be StatusPassed for empty output, got %s", status)
	}

	if output != "\n" && output != "" {
		t.Errorf("Expected empty or newline output, got %q", output)
	}
}

// Security tests: path containment and go.work support

func TestContainPath_GoTest(t *testing.T) {
	root := "/tmp/project"
	cases := []struct {
		filePath string
		want     string
	}{
		{"pkg/test.go", "/tmp/project/pkg/test.go"},
		{"../outside", "/tmp/project"},
		{"../../etc/passwd", "/tmp/project"},
		{"./test.go", "/tmp/project/test.go"},
		{"pkg/../test.go", "/tmp/project/test.go"},
		{"/tmp/project/pkg/test.go", "/tmp/project/pkg/test.go"},
		{"/tmp/project", "/tmp/project"},
		{"/tmp/project/", "/tmp/project"},
		{"pkg/../../outside", "/tmp/project"},
	}
	for _, c := range cases {
		got := ContainPath(root, c.filePath)
		if got != c.want {
			t.Errorf("ContainPath(%q, %q) = %q, want %q", root, c.filePath, got, c.want)
		}
	}
}

func TestGoTestDriver_Detect_GoWork(t *testing.T) {
	tmpDir := t.TempDir()
	absRoot, _ := filepath.Abs(tmpDir)

	// Create go.work file
	goWorkPath := filepath.Join(absRoot, "go.work")
	goWorkContent := []byte("go 1.21\nuse std\n")
	if err := os.WriteFile(goWorkPath, goWorkContent, 0644); err != nil {
		t.Fatal(err)
	}

	// Create a test file in any subdirectory (workspace)
	testDir := filepath.Join(absRoot, "pkg")
	os.MkdirAll(testDir, 0755)
	testFile := filepath.Join(testDir, "handler_test.go")
	os.WriteFile(testFile, []byte("package pkg\n"), 0644)

	driver := &GoTestDriver{}
	detected, err := driver.Detect(absRoot)
	if err != nil {
		t.Fatalf("Detect error: %v", err)
	}
	if !detected {
		t.Error("Detect should return true when go.work exists and test files present")
	}
}

func TestGoTestDriver_BuildTestCommand_PathEscape(t *testing.T) {
	driver := &GoTestDriver{}
	tmpRoot := t.TempDir()
	absRoot, _ := filepath.Abs(tmpRoot)

	// Build command with an escaping path
	cmd := driver.buildTestCommand(context.Background(), absRoot, "../../../etc/passwd")

	// The command should not include the escaping path in args
	hasEscaping := false
	for _, arg := range cmd.Args {
		if strings.Contains(arg, "..") {
			hasEscaping = true
		}
	}
	if hasEscaping {
		t.Errorf("BuildTestCommand produced escaping argument: %v", cmd.Args)
	}
}
