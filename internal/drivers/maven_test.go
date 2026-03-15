package drivers

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pueblomo/lazytest/internal/types"
)

func TestMavenDriver_Name(t *testing.T) {
	driver := &MavenDriver{}
	if driver.Name() != "maven" {
		t.Errorf("Name() = %v, want 'maven'", driver.Name())
	}
}

func TestMavenDriver_Name_IsConstant(t *testing.T) {
	driver1 := &MavenDriver{}
	driver2 := &MavenDriver{}

	if driver1.Name() != driver2.Name() {
		t.Error("Name() should return the same value for all instances")
	}
}

func TestMavenDriver_Detect_NoPomXml(t *testing.T) {
	tmpDir := t.TempDir()

	driver := &MavenDriver{}
	detected, err := driver.Detect(tmpDir)

	if err != nil {
		t.Errorf("Detect() should not return error when pom.xml doesn't exist, got %v", err)
	}

	if detected {
		t.Error("Detect() should return false when pom.xml doesn't exist")
	}
}

func TestMavenDriver_Detect_PomXmlNoTestDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Create pom.xml
	pomPath := filepath.Join(tmpDir, "pom.xml")
	pomContent := `<project><modelVersion>4.0.0</modelVersion></project>`
	if err := os.WriteFile(pomPath, []byte(pomContent), 0644); err != nil {
		t.Fatal(err)
	}

	driver := &MavenDriver{}
	detected, err := driver.Detect(tmpDir)

	if err != nil {
		t.Errorf("Detect() should not return error, got %v", err)
	}

	if detected {
		t.Error("Detect() should return false when src/test/java doesn't exist")
	}
}

func TestMavenDriver_Detect_PomXmlNoTestFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create pom.xml
	pomPath := filepath.Join(tmpDir, "pom.xml")
	pomContent := `<project><modelVersion>4.0.0</modelVersion></project>`
	if err := os.WriteFile(pomPath, []byte(pomContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create src/test/java but no test files
	testDir := filepath.Join(tmpDir, "src", "test", "java")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatal(err)
	}

	driver := &MavenDriver{}
	detected, err := driver.Detect(tmpDir)

	if err != nil {
		t.Errorf("Detect() should not return error, got %v", err)
	}

	if detected {
		t.Error("Detect() should return false when no *Test.java files exist")
	}
}

func TestMavenDriver_Detect_PomXmlWithTestFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create pom.xml
	pomPath := filepath.Join(tmpDir, "pom.xml")
	pomContent := `<project><modelVersion>4.0.0</modelVersion></project>`
	if err := os.WriteFile(pomPath, []byte(pomContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create src/test/java with a test file
	pkgDir := filepath.Join(tmpDir, "src", "test", "java", "com", "example")
	if err := os.MkdirAll(pkgDir, 0755); err != nil {
		t.Fatal(err)
	}

	testFile := filepath.Join(pkgDir, "MyTest.java")
	if err := os.WriteFile(testFile, []byte("package com.example;"), 0644); err != nil {
		t.Fatal(err)
	}

	driver := &MavenDriver{}
	detected, err := driver.Detect(tmpDir)

	if err != nil {
		t.Errorf("Detect() should not return error, got %v", err)
	}

	if !detected {
		t.Error("Detect() should return true when *Test.java files exist")
	}
}

func TestMavenDriver_Detect_PomXmlButOnlyNonTestJavaFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create pom.xml
	pomPath := filepath.Join(tmpDir, "pom.xml")
	pomContent := `<project><modelVersion>4.0.0</modelVersion></project>`
	if err := os.WriteFile(pomPath, []byte(pomContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create src/test/java with a non-test file
	pkgDir := filepath.Join(tmpDir, "src", "test", "java", "com", "example")
	if err := os.MkdirAll(pkgDir, 0755); err != nil {
		t.Fatal(err)
	}

	nonTestFile := filepath.Join(pkgDir, "MyClass.java")
	if err := os.WriteFile(nonTestFile, []byte("package com.example;"), 0644); err != nil {
		t.Fatal(err)
	}

	driver := &MavenDriver{}
	detected, err := driver.Detect(tmpDir)

	if err != nil {
		t.Errorf("Detect() should not return error, got %v", err)
	}

	if detected {
		t.Error("Detect() should return false when no *Test.java files exist")
	}
}

func TestMavenDriver_DetectTestFiles_FindsTestFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create src/test/java with test files in packages
	pkgDir := filepath.Join(tmpDir, "src", "test", "java", "com", "example")
	if err := os.MkdirAll(pkgDir, 0755); err != nil {
		t.Fatal(err)
	}

	testFiles := []string{"MyTest.java", "AnotherTest.java", "HelperTest.java"}
	for _, tf := range testFiles {
		path := filepath.Join(pkgDir, tf)
		if err := os.WriteFile(path, []byte("package com.example;"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Also create a non-test file that should be ignored
	nonTest := filepath.Join(pkgDir, "MyClass.java")
	if err := os.WriteFile(nonTest, []byte("package com.example;"), 0644); err != nil {
		t.Fatal(err)
	}

	driver := &MavenDriver{}
	ctx := context.Background()
	files, err := driver.DetectTestFiles(ctx, tmpDir)

	if err != nil {
		t.Errorf("DetectTestFiles() should not return error, got %v", err)
	}

	if len(files) != 3 {
		t.Errorf("DetectTestFiles() returned %d files, want 3", len(files))
	}

	// All files should be relative paths starting with src/
	for _, f := range files {
		if !strings.HasPrefix(f, "src/") {
			t.Errorf("DetectTestFiles() file %q should start with 'src/'", f)
		}
		if !strings.HasSuffix(f, "Test.java") {
			t.Errorf("DetectTestFiles() file %q should end with 'Test.java'", f)
		}
	}
}

func TestMavenDriver_DetectTestFiles_NoTestDir(t *testing.T) {
	tmpDir := t.TempDir()

	driver := &MavenDriver{}
	ctx := context.Background()
	files, err := driver.DetectTestFiles(ctx, tmpDir)

	if err != nil {
		t.Errorf("DetectTestFiles() should not return error for missing dir, got %v", err)
	}

	if len(files) != 0 {
		t.Errorf("DetectTestFiles() should return empty slice, got %d files", len(files))
	}
}

func TestMavenDriver_DetectTestFiles_NestedPackages(t *testing.T) {
	tmpDir := t.TempDir()

	// Create nested package structure
	deepPkgDir := filepath.Join(tmpDir, "src", "test", "java", "com", "example", "service")
	if err := os.MkdirAll(deepPkgDir, 0755); err != nil {
		t.Fatal(err)
	}

	testFile := filepath.Join(deepPkgDir, "UserServiceTest.java")
	if err := os.WriteFile(testFile, []byte("package com.example.service;"), 0644); err != nil {
		t.Fatal(err)
	}

	driver := &MavenDriver{}
	ctx := context.Background()
	files, err := driver.DetectTestFiles(ctx, tmpDir)

	if err != nil {
		t.Errorf("DetectTestFiles() should not return error, got %v", err)
	}

	if len(files) != 1 {
		t.Errorf("DetectTestFiles() returned %d files, want 1", len(files))
	}

	expected := "src/test/java/com/example/service/UserServiceTest.java"
	// Normalize path separators for cross-platform
	expected = filepath.FromSlash(expected)
	if len(files) == 1 && files[0] != expected {
		t.Errorf("DetectTestFiles() file = %q, want %q", files[0], expected)
	}
}

func TestMavenDriver_PathToClassName(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "simple class",
			path:     "src/test/java/com/example/MyTest.java",
			expected: "com.example.MyTest",
		},
		{
			name:     "nested package",
			path:     "src/test/java/com/example/service/UserServiceTest.java",
			expected: "com.example.service.UserServiceTest",
		},
		{
			name:     "root level test",
			path:     "src/test/java/SimpleTest.java",
			expected: "SimpleTest",
		},
		{
			name:     "single package",
			path:     "src/test/java/example/MyTest.java",
			expected: "example.MyTest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Normalize to OS-specific path separators
			path := filepath.FromSlash(tt.path)
			result := pathToClassName(path)
			if result != tt.expected {
				t.Errorf("pathToClassName(%q) = %q, want %q", tt.path, result, tt.expected)
			}
		})
	}
}

func TestMavenDriver_BuildTestCommand_Basic(t *testing.T) {
	driver := &MavenDriver{}
	ctx := context.Background()
	root := "/project"
	filePath := filepath.FromSlash("src/test/java/com/example/MyTest.java")

	cmd := driver.buildTestCommand(ctx, root, filePath)

	if cmd.Dir != root {
		t.Errorf("cmd.Dir = %v, want %v", cmd.Dir, root)
	}

	// Should contain mvn test -Dtest=com.example.MyTest -q
	expectedArgs := []string{"mvn", "test", "-Dtest=com.example.MyTest", "-q"}
	if len(cmd.Args) != len(expectedArgs) {
		t.Errorf("cmd.Args = %v, want %v", cmd.Args, expectedArgs)
		return
	}

	for i, arg := range expectedArgs {
		if cmd.Args[i] != arg {
			t.Errorf("cmd.Args[%d] = %q, want %q", i, cmd.Args[i], arg)
		}
	}
}

func TestMavenDriver_BuildTestCommand_WithContext(t *testing.T) {
	driver := &MavenDriver{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	root := "/project"
	filePath := filepath.FromSlash("src/test/java/MyTest.java")

	cmd := driver.buildTestCommand(ctx, root, filePath)

	if cmd.Path == "" {
		t.Error("cmd.Path should be set")
	}

	if cmd.Dir != root {
		t.Errorf("cmd.Dir = %v, want %v", cmd.Dir, root)
	}
}

func TestMavenDriver_BuildTestCommand_DeepPackage(t *testing.T) {
	driver := &MavenDriver{}
	ctx := context.Background()
	root := "/project"
	filePath := filepath.FromSlash("src/test/java/com/example/service/UserServiceTest.java")

	cmd := driver.buildTestCommand(ctx, root, filePath)

	expectedArgs := []string{"mvn", "test", "-Dtest=com.example.service.UserServiceTest", "-q"}
	if len(cmd.Args) != len(expectedArgs) {
		t.Errorf("cmd.Args = %v, want %v", cmd.Args, expectedArgs)
		return
	}

	for i, arg := range expectedArgs {
		if cmd.Args[i] != arg {
			t.Errorf("cmd.Args[%d] = %q, want %q", i, cmd.Args[i], arg)
		}
	}
}

func TestMavenDriver_ExecuteTestCommand_BuildSuccess(t *testing.T) {
	driver := &MavenDriver{}

	cmd := exec.Command("echo", "BUILD SUCCESS")

	status, output, err := driver.executeTestCommand(cmd)

	if err != nil {
		t.Errorf("executeTestCommand() should not return error for BUILD SUCCESS, got %v", err)
	}

	if status != types.StatusPassed {
		t.Errorf("executeTestCommand() status = %v, want %v", status, types.StatusPassed)
	}

	if !strings.Contains(output, "BUILD SUCCESS") {
		t.Error("executeTestCommand() output should contain 'BUILD SUCCESS'")
	}
}

func TestMavenDriver_ExecuteTestCommand_BuildFailure(t *testing.T) {
	driver := &MavenDriver{}

	cmd := exec.Command("sh", "-c", "echo 'BUILD FAILURE' && exit 1")

	status, output, err := driver.executeTestCommand(cmd)

	// Should not return error (test failure is expected behavior)
	if err != nil {
		t.Errorf("executeTestCommand() should not return error for test failure, got %v", err)
	}

	if status != types.StatusFailed {
		t.Errorf("executeTestCommand() status = %v, want %v", status, types.StatusFailed)
	}

	if !strings.Contains(output, "BUILD FAILURE") {
		t.Error("executeTestCommand() output should contain 'BUILD FAILURE'")
	}
}

func TestMavenDriver_ExecuteTestCommand_BuildSuccessWithSkipped(t *testing.T) {
	driver := &MavenDriver{}

	cmd := exec.Command("echo", "BUILD SUCCESS\nTests run: 3, Failures: 0, Errors: 0, Skipped: 1")

	status, output, err := driver.executeTestCommand(cmd)

	if err != nil {
		t.Errorf("executeTestCommand() should not return error, got %v", err)
	}

	if status != types.StatusSkipped {
		t.Errorf("executeTestCommand() status = %v, want %v", status, types.StatusSkipped)
	}

	if !strings.Contains(output, "Skipped: 1") {
		t.Error("executeTestCommand() output should contain 'Skipped: 1'")
	}
}

func TestMavenDriver_ExecuteTestCommand_CommandError(t *testing.T) {
	driver := &MavenDriver{}

	// Command that exits with error but no BUILD FAILURE in output
	cmd := exec.Command("false")

	status, _, err := driver.executeTestCommand(cmd)

	if err == nil {
		t.Error("executeTestCommand() should return error when command fails without BUILD FAILURE")
	}

	if status != types.StatusFailed {
		t.Errorf("executeTestCommand() status = %v, want %v", status, types.StatusFailed)
	}
}

func TestMavenDriver_ExecuteTestCommand_EmptyOutput(t *testing.T) {
	driver := &MavenDriver{}

	cmd := exec.Command("echo", "")

	status, _, err := driver.executeTestCommand(cmd)

	if err != nil {
		t.Errorf("executeTestCommand() should not return error, got %v", err)
	}

	// Default to passed when no status markers found
	if status != types.StatusPassed {
		t.Errorf("executeTestCommand() status = %v, want %v", status, types.StatusPassed)
	}
}

type wrongMavenItem struct{}

func (w *wrongMavenItem) FilterValue() string { return "wrong" }

func TestMavenDriver_RunTest_InvalidTestCase(t *testing.T) {
	driver := &MavenDriver{}
	ctx := context.Background()

	err := driver.RunTest(ctx, "/tmp", &wrongMavenItem{})

	if err == nil {
		t.Error("RunTest() should return error for invalid test case type")
	}

	if !strings.Contains(err.Error(), "can't convert to TestCase") {
		t.Errorf("RunTest() error = %v, want error containing 'can't convert to TestCase'", err)
	}
}

func TestMavenDriver_RunTest_SetsTestCaseFields(t *testing.T) {
	driver := &MavenDriver{}
	ctx := context.Background()
	tmpDir := t.TempDir()

	tc := &types.TestCase{
		Name:       "SomeTest",
		Filepath:   filepath.FromSlash("src/test/java/SomeTest.java"),
		TestStatus: types.StatusNotStarted,
	}

	// This will fail since there's no actual Maven project, but we're testing field assignment
	_ = driver.RunTest(ctx, tmpDir, tc)

	// The test status should have changed from NotStarted
	if tc.TestStatus == types.StatusNotStarted {
		t.Error("Expected test status to change after running")
	}

	// Output should be set (even if empty or contains error)
	_ = tc.Output // Just verify it's accessible
}

func TestMavenDriver_Detect_PomXmlStatError(t *testing.T) {
	// Use a path that exists but can't be read as pom.xml
	tmpFile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	driver := &MavenDriver{}
	detected, err := driver.Detect(tmpFile.Name())

	if detected {
		t.Error("Expected false when path is not a valid Maven project")
	}
	// Error may or may not occur depending on implementation
	_ = err
}

func TestMavenDriver_DetectTestFiles_ContextCancellation(t *testing.T) {
	driver := &MavenDriver{}
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	tmpDir := t.TempDir()

	// Create a test dir with many files to ensure walk gets interrupted
	pkgDir := filepath.Join(tmpDir, "src", "test", "java", "com", "example")
	if err := os.MkdirAll(pkgDir, 0755); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		path := filepath.Join(pkgDir, "Test"+string(rune('A'+i))+".java")
		if err := os.WriteFile(path, []byte("package com.example;"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Context cancellation is handled by exec.CommandContext internally
	// but DetectTestFiles uses filepath.Walk which doesn't support context.
	// We just verify it doesn't panic or hang.
	_, _ = driver.DetectTestFiles(ctx, tmpDir)
}

func TestMavenDriver_ExecuteTestCommand_SuccessWithTestsRun(t *testing.T) {
	driver := &MavenDriver{}

	// Simulate BUILD SUCCESS with test summary
	cmd := exec.Command("echo", "BUILD SUCCESS\nTests run: 5, Failures: 0, Errors: 0, Skipped: 0")

	status, output, err := driver.executeTestCommand(cmd)

	if err != nil {
		t.Errorf("executeTestCommand() should not return error, got %v", err)
	}

	if status != types.StatusPassed {
		t.Errorf("executeTestCommand() status = %v, want %v", status, types.StatusPassed)
	}

	if !strings.Contains(output, "Tests run:") {
		t.Error("executeTestCommand() output should contain 'Tests run:'")
	}
}

func TestMavenDriver_Detect_MultiModuleProject(t *testing.T) {
	tmpDir := t.TempDir()

	pomPath := filepath.Join(tmpDir, "pom.xml")
	if err := os.WriteFile(pomPath, []byte(`<project><modelVersion>4.0.0</modelVersion></project>`), 0644); err != nil {
		t.Fatal(err)
	}

	moduleDir := filepath.Join(tmpDir, "module-a")
	if err := os.MkdirAll(filepath.Join(moduleDir, "src", "test", "java", "com", "example"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(moduleDir, "pom.xml"), []byte(`<project><modelVersion>4.0.0</modelVersion></project>`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(moduleDir, "src", "test", "java", "com", "example", "ModuleTest.java"), []byte("package com.example;"), 0644); err != nil {
		t.Fatal(err)
	}

	driver := &MavenDriver{}
	detected, err := driver.Detect(tmpDir)
	if err != nil {
		t.Fatalf("Detect() unexpected error: %v", err)
	}
	if !detected {
		t.Fatal("Detect() should return true for a multi-module Maven project with tests in a submodule")
	}
}

func TestMavenDriver_DetectTestFiles_MultiModuleProject(t *testing.T) {
	tmpDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(tmpDir, "pom.xml"), []byte(`<project><modelVersion>4.0.0</modelVersion></project>`), 0644); err != nil {
		t.Fatal(err)
	}

	moduleTest := filepath.Join(tmpDir, "module-a", "src", "test", "java", "com", "example", "ModuleTest.java")
	if err := os.MkdirAll(filepath.Dir(moduleTest), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "module-a", "pom.xml"), []byte(`<project><modelVersion>4.0.0</modelVersion></project>`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(moduleTest, []byte("package com.example;"), 0644); err != nil {
		t.Fatal(err)
	}

	driver := &MavenDriver{}
	files, err := driver.DetectTestFiles(context.Background(), tmpDir)
	if err != nil {
		t.Fatalf("DetectTestFiles() unexpected error: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("DetectTestFiles() returned %d files, want 1", len(files))
	}

	expected := filepath.Join("module-a", "src", "test", "java", "com", "example", "ModuleTest.java")
	if files[0] != expected {
		t.Fatalf("DetectTestFiles() file = %q, want %q", files[0], expected)
	}
}

func TestMavenDriver_PathToClassName_WindowsPath(t *testing.T) {
	path := `src\\test\\java\\com\\example\\WindowsStyleTest.java`
	got := pathToClassName(path)
	want := "com.example.WindowsStyleTest"
	if got != want {
		t.Fatalf("pathToClassName(%q) = %q, want %q", path, got, want)
	}
}

func TestMavenDriver_BuildTestCommand_WindowsPath(t *testing.T) {
	driver := &MavenDriver{}
	cmd := driver.buildTestCommand(context.Background(), "/project", `src\\test\\java\\com\\example\\WindowsStyleTest.java`)
	if got, want := cmd.Args[2], "-Dtest=com.example.WindowsStyleTest"; got != want {
		t.Fatalf("cmd.Args[2] = %q, want %q", got, want)
	}
}

func TestMavenDriver_ExecuteTestCommand_MissingMavenBinary(t *testing.T) {
	driver := &MavenDriver{}
	cmd := exec.Command("/definitely/not/a/real/binary")

	status, _, err := driver.executeTestCommand(cmd)
	if err == nil {
		t.Fatal("executeTestCommand() should return an error when the mvn binary is missing")
	}
	if status != types.StatusFailed {
		t.Fatalf("executeTestCommand() status = %v, want %v", status, types.StatusFailed)
	}
}

func TestMavenDriver_ExecuteTestCommand_ContextTimeout(t *testing.T) {
	driver := &MavenDriver{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", "sleep 1")
	status, _, err := driver.executeTestCommand(cmd)
	if err == nil {
		t.Fatal("executeTestCommand() should return an error when the command times out")
	}
	if status != types.StatusFailed {
		t.Fatalf("executeTestCommand() status = %v, want %v", status, types.StatusFailed)
	}
}

func TestFindModuleRoot_SingleModule(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "pom.xml"), []byte("<project/>"), 0644)

	result := findModuleRoot(tmpDir, "src/test/java/com/example/MyTest.java")
	if result != tmpDir {
		t.Errorf("findModuleRoot() = %v, want %v", result, tmpDir)
	}
}

func TestFindModuleRoot_MultiModule(t *testing.T) {
	tmpDir := t.TempDir()
	// Root pom.xml
	os.WriteFile(filepath.Join(tmpDir, "pom.xml"), []byte("<project/>"), 0644)

	// Module with its own pom.xml
	moduleDir := filepath.Join(tmpDir, "module-a")
	os.MkdirAll(filepath.Join(moduleDir, "src/test/java"), 0755)
	os.WriteFile(filepath.Join(moduleDir, "pom.xml"), []byte("<project/>"), 0644)

	filePath := "module-a/src/test/java/com/example/ModuleATest.java"
	result := findModuleRoot(tmpDir, filePath)
	if result != moduleDir {
		t.Errorf("findModuleRoot() = %v, want %v", result, moduleDir)
	}
}

func TestFindModuleRoot_NoPomXmlFallsBackToRoot(t *testing.T) {
	tmpDir := t.TempDir()
	// No pom.xml at all

	result := findModuleRoot(tmpDir, "src/test/java/MyTest.java")
	if result != tmpDir {
		t.Errorf("findModuleRoot() should fall back to root, got %v", result)
	}
}

func TestMavenDriver_BuildTestCommand_MultiModule(t *testing.T) {
	driver := &MavenDriver{}
	ctx := context.Background()
	tmpDir := t.TempDir()

	// Create module structure
	moduleDir := filepath.Join(tmpDir, "module-a")
	os.MkdirAll(filepath.Join(moduleDir, "src/test/java/com/example"), 0755)
	os.WriteFile(filepath.Join(moduleDir, "pom.xml"), []byte("<project/>"), 0644)

	filePath := "module-a/src/test/java/com/example/ModuleATest.java"
	cmd := driver.buildTestCommand(ctx, tmpDir, filePath)

	if cmd.Dir != moduleDir {
		t.Errorf("buildTestCommand() Dir = %v, want %v (module root)", cmd.Dir, moduleDir)
	}

	expectedArgs := []string{"mvn", "test", "-Dtest=com.example.ModuleATest", "-q"}
	if len(cmd.Args) != len(expectedArgs) {
		t.Fatalf("buildTestCommand() args length = %d, want %d", len(cmd.Args), len(expectedArgs))
	}
	for i, arg := range expectedArgs {
		if cmd.Args[i] != arg {
			t.Errorf("buildTestCommand() arg[%d] = %v, want %v", i, cmd.Args[i], arg)
		}
	}
}

func TestMavenDriver_BuildTestCommand_SingleModule(t *testing.T) {
	driver := &MavenDriver{}
	ctx := context.Background()
	tmpDir := t.TempDir()

	cmd := driver.buildTestCommand(ctx, tmpDir, "src/test/java/com/example/MyTest.java")

	if cmd.Dir != tmpDir {
		t.Errorf("buildTestCommand() Dir = %v, want %v", cmd.Dir, tmpDir)
	}
}
