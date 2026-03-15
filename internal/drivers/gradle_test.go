package drivers

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/pueblomo/lazytest/internal/types"
)

func writeGradleFile(t *testing.T, path, content string, mode os.FileMode) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), mode); err != nil {
		t.Fatal(err)
	}
}

func TestGradleDriver_Detect_BothBuildFilesExist(t *testing.T) {
	root := t.TempDir()
	writeGradleFile(t, filepath.Join(root, "build.gradle"), "plugins { id 'java' }", 0644)
	writeGradleFile(t, filepath.Join(root, "build.gradle.kts"), "plugins { java }", 0644)
	writeGradleFile(t, filepath.Join(root, "src/test/java/com/example/MyTest.java"), "package com.example;", 0644)

	got, err := (&GradleDriver{}).Detect(root)
	if err != nil || !got {
		t.Fatalf("Detect() = (%v, %v), want (true, nil)", got, err)
	}
}

func TestGradleDriver_Detect_NoSrcTestJavaDir(t *testing.T) {
	root := t.TempDir()
	writeGradleFile(t, filepath.Join(root, "build.gradle"), "plugins { id 'java' }", 0644)

	got, err := (&GradleDriver{}).Detect(root)
	if err != nil {
		t.Fatalf("Detect() unexpected error: %v", err)
	}
	if got {
		t.Fatal("Detect() = true, want false when src/test/java is missing")
	}
}

func TestGradleDriver_Detect_IgnoresNonStandardTestNames(t *testing.T) {
	root := t.TempDir()
	writeGradleFile(t, filepath.Join(root, "build.gradle"), "plugins { id 'java' }", 0644)
	writeGradleFile(t, filepath.Join(root, "src/test/java/com/example/MyTests.java"), "package com.example;", 0644)
	writeGradleFile(t, filepath.Join(root, "src/test/java/com/example/TestMyThing.java"), "package com.example;", 0644)

	got, err := (&GradleDriver{}).Detect(root)
	if err != nil {
		t.Fatalf("Detect() unexpected error: %v", err)
	}
	if got {
		t.Fatal("Detect() = true, want false for non-*Test.java names")
	}
}

func TestGradleDriver_DetectTestFiles_MultiModuleNestedModule(t *testing.T) {
	root := t.TempDir()
	writeGradleFile(t, filepath.Join(root, "build.gradle"), "plugins { id 'java' }", 0644)
	writeGradleFile(t, filepath.Join(root, "services/api/build.gradle"), "plugins { id 'java' }", 0644)
	writeGradleFile(t, filepath.Join(root, "services/api/src/test/java/com/example/NestedModuleTest.java"), "package com.example;", 0644)

	files, err := (&GradleDriver{}).DetectTestFiles(context.Background(), root)
	if err != nil {
		t.Fatalf("DetectTestFiles() unexpected error: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("DetectTestFiles() returned %d files, want 1", len(files))
	}
	want := filepath.Join("services", "api", "src", "test", "java", "com", "example", "NestedModuleTest.java")
	if files[0] != want {
		t.Fatalf("DetectTestFiles() file = %q, want %q", files[0], want)
	}
}

func TestGradleDriver_BuildTestCommand_UsesNestedModuleRoot(t *testing.T) {
	root := t.TempDir()
	moduleRoot := filepath.Join(root, "services", "api")
	writeGradleFile(t, filepath.Join(root, "build.gradle"), "plugins { id 'java' }", 0644)
	writeGradleFile(t, filepath.Join(moduleRoot, "build.gradle.kts"), "plugins { java }", 0644)

	cmd := (&GradleDriver{}).buildTestCommand(context.Background(), root, filepath.Join("services", "api", "src", "test", "java", "com", "example", "NestedModuleTest.java"))
	if cmd.Dir != moduleRoot {
		t.Fatalf("buildTestCommand() Dir = %q, want %q", cmd.Dir, moduleRoot)
	}
	if got, want := cmd.Args[3], "com.example.NestedModuleTest"; got != want {
		t.Fatalf("buildTestCommand() class = %q, want %q", got, want)
	}
}

func TestGradleDriver_BuildTestCommand_NonExecutableGradlewFallsBackToGradle(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("unix execute bits do not apply")
	}
	root := t.TempDir()
	writeGradleFile(t, filepath.Join(root, "build.gradle"), "plugins { id 'java' }", 0644)
	writeGradleFile(t, filepath.Join(root, "gradlew"), "#!/bin/sh\necho hi", 0644)

	cmd := (&GradleDriver{}).buildTestCommand(context.Background(), root, filepath.Join("src", "test", "java", "com", "example", "MyTest.java"))
	if got := cmd.Args[0]; got != "gradle" {
		t.Fatalf("buildTestCommand() command = %q, want gradle", got)
	}
}

func TestGradleDriver_BuildTestCommand_PathWithSpaces(t *testing.T) {
	root := filepath.Join(t.TempDir(), "project with spaces")
	moduleRoot := filepath.Join(root, "module with spaces")
	if err := os.MkdirAll(moduleRoot, 0755); err != nil {
		t.Fatal(err)
	}
	writeGradleFile(t, filepath.Join(root, "build.gradle"), "plugins { id 'java' }", 0644)
	writeGradleFile(t, filepath.Join(moduleRoot, "build.gradle"), "plugins { id 'java' }", 0644)

	cmd := (&GradleDriver{}).buildTestCommand(context.Background(), root, filepath.Join("module with spaces", "src", "test", "java", "com", "example", "SpacePathTest.java"))
	if cmd.Dir != moduleRoot {
		t.Fatalf("buildTestCommand() Dir = %q, want %q", cmd.Dir, moduleRoot)
	}
	if got, want := cmd.Args[3], "com.example.SpacePathTest"; got != want {
		t.Fatalf("buildTestCommand() class = %q, want %q", got, want)
	}
}

func TestFindGradleModuleRoots_BothBuildFilesInSameProject(t *testing.T) {
	root := t.TempDir()
	writeGradleFile(t, filepath.Join(root, "build.gradle"), "plugins { id 'java' }", 0644)
	writeGradleFile(t, filepath.Join(root, "build.gradle.kts"), "plugins { java }", 0644)
	writeGradleFile(t, filepath.Join(root, "module-a/build.gradle"), "plugins { id 'java' }", 0644)

	roots, err := findGradleModuleRoots(root)
	if err != nil {
		t.Fatalf("findGradleModuleRoots() unexpected error: %v", err)
	}
	set := map[string]bool{}
	for _, r := range roots {
		set[r] = true
	}
	if len(set) != 2 || !set[root] || !set[filepath.Join(root, "module-a")] {
		t.Fatalf("findGradleModuleRoots() unexpected roots: %v", roots)
	}
}

func TestGradleDriver_ExecuteTestCommand_BuildSuccessful(t *testing.T) {
	status, output, err := (&GradleDriver{}).executeTestCommand(exec.Command("echo", "BUILD SUCCESSFUL"))
	if err != nil {
		t.Fatalf("executeTestCommand() unexpected error: %v", err)
	}
	if status != types.StatusPassed || !strings.Contains(output, "BUILD SUCCESSFUL") {
		t.Fatalf("executeTestCommand() = (%v, %q), want passed with BUILD SUCCESSFUL", status, output)
	}
}

func TestGradleDriver_ExecuteTestCommand_BuildFailed(t *testing.T) {
	status, output, err := (&GradleDriver{}).executeTestCommand(exec.Command("sh", "-c", "echo 'BUILD FAILED' && exit 1"))
	if err != nil {
		t.Fatalf("executeTestCommand() unexpected error: %v", err)
	}
	if status != types.StatusFailed || !strings.Contains(output, "BUILD FAILED") {
		t.Fatalf("executeTestCommand() = (%v, %q), want failed with BUILD FAILED", status, output)
	}
}

func TestGradleDriver_ExecuteTestCommand_MissingBinary(t *testing.T) {
	status, _, err := (&GradleDriver{}).executeTestCommand(exec.Command("/definitely/not/a/real/binary"))
	if err == nil {
		t.Fatal("executeTestCommand() error = nil, want error")
	}
	if status != types.StatusFailed {
		t.Fatalf("executeTestCommand() status = %v, want %v", status, types.StatusFailed)
	}
}
