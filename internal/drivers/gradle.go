package drivers

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/pueblomo/lazytest/internal/types"
)

type GradleDriver struct{}

func (d *GradleDriver) Name() string { return "gradle" }

func (d *GradleDriver) Detect(root string) (bool, error) {
	moduleRoots, err := findGradleModuleRoots(root)
	if err != nil {
		return false, err
	}

	for _, moduleRoot := range moduleRoots {
		srcTestDir := filepath.Join(moduleRoot, "src", "test", "java")
		if _, err := os.Stat(srcTestDir); os.IsNotExist(err) {
			continue
		} else if err != nil {
			return false, err
		}

		hasTests := false
		err := filepath.Walk(srcTestDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(info.Name(), "Test.java") {
				hasTests = true
				return filepath.SkipAll
			}
			return nil
		})
		if err != nil {
			return false, err
		}
		if hasTests {
			return true, nil
		}
	}

	return false, nil
}

func (d *GradleDriver) RunTest(ctx context.Context, root string, testCase list.Item) error {
	tc, ok := testCase.(*types.TestCase)
	if !ok {
		return fmt.Errorf("can't convert to TestCase")
	}
	cmd := d.buildTestCommand(ctx, root, tc.Filepath)
	testStatus, output, err := d.executeTestCommand(cmd)
	tc.TestStatus = testStatus
	tc.Output = output
	return err
}

func (d *GradleDriver) buildTestCommand(ctx context.Context, root string, filePath string) *exec.Cmd {
	className := pathToClassName(filePath)
	workingDir := findGradleModuleRoot(root, filePath)

	var cmd *exec.Cmd
	if info, err := os.Stat(filepath.Join(workingDir, "gradlew")); err == nil && info.Mode().IsRegular() && info.Mode().Perm()&0111 != 0 {
		cmd = exec.CommandContext(ctx, filepath.Join(workingDir, "gradlew"), "test", "--tests", className, "-q")
	} else if _, err := os.Stat(filepath.Join(workingDir, "gradlew.bat")); err == nil {
		cmd = exec.CommandContext(ctx, filepath.Join(workingDir, "gradlew.bat"), "test", "--tests", className, "-q")
	} else {
		cmd = exec.CommandContext(ctx, "gradle", "test", "--tests", className, "-q")
	}
	cmd.Dir = workingDir
	return cmd
}

func (d *GradleDriver) executeTestCommand(cmd *exec.Cmd) (types.TestStatus, string, error) {
	output, err := cmd.CombinedOutput()
	outputString := string(output)

	if err != nil {
		if strings.Contains(outputString, "BUILD FAILED") || strings.Contains(outputString, " Tests run:") {
			return types.StatusFailed, outputString, nil
		}
		return types.StatusFailed, outputString, fmt.Errorf("test command %v failed: %w", cmd.Args, err)
	}

	if strings.Contains(outputString, "BUILD SUCCESSFUL") {
		if strings.Contains(outputString, "Skipped:") && !strings.Contains(outputString, "Skipped: 0") {
			return types.StatusSkipped, outputString, nil
		}
		return types.StatusPassed, outputString, nil
	}
	if strings.Contains(outputString, "BUILD FAILED") {
		return types.StatusFailed, outputString, nil
	}
	return types.StatusPassed, outputString, nil
}

func (d *GradleDriver) DetectTestFiles(ctx context.Context, root string) ([]string, error) {
	var testFiles []string
	moduleRoots, err := findGradleModuleRoots(root)
	if err != nil {
		return nil, err
	}
	if len(moduleRoots) == 0 {
		moduleRoots = []string{root}
	}

	for _, moduleRoot := range moduleRoots {
		srcTestDir := filepath.Join(moduleRoot, "src", "test", "java")
		if _, err := os.Stat(srcTestDir); os.IsNotExist(err) {
			continue
		} else if err != nil {
			return nil, err
		}

		err := filepath.Walk(srcTestDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(info.Name(), "Test.java") {
				relPath, err := filepath.Rel(root, path)
				if err != nil {
					relPath = path
				}
				testFiles = append(testFiles, relPath)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return testFiles, nil
}

func findGradleModuleRoots(root string) ([]string, error) {
	buildGradlePath := filepath.Join(root, "build.gradle")
	buildGradleKtsPath := filepath.Join(root, "build.gradle.kts")
	if _, err := os.Stat(buildGradlePath); err != nil && os.IsNotExist(err) {
		if _, err := os.Stat(buildGradleKtsPath); err != nil && os.IsNotExist(err) {
			return nil, nil
		}
	}

	moduleSet := map[string]struct{}{root: {}}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == "build" || strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if info.Name() == "build.gradle" || info.Name() == "build.gradle.kts" {
			moduleSet[filepath.Dir(path)] = struct{}{}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	moduleRoots := make([]string, 0, len(moduleSet))
	for moduleRoot := range moduleSet {
		moduleRoots = append(moduleRoots, moduleRoot)
	}
	return moduleRoots, nil
}

func findGradleModuleRoot(root string, filePath string) string {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return root
	}
	absRoot = filepath.Clean(absRoot)

	absFile := filePath
	if !filepath.IsAbs(absFile) {
		absFile = filepath.Join(absRoot, filePath)
	}
	absFile = filepath.Clean(absFile)

	rootWithSep := absRoot + string(filepath.Separator)
	if !strings.HasPrefix(absFile+string(filepath.Separator), rootWithSep) {
		return root
	}

	dir := filepath.Dir(absFile)
	for {
		if _, err := os.Stat(filepath.Join(dir, "build.gradle")); err == nil {
			return dir
		}
		if _, err := os.Stat(filepath.Join(dir, "build.gradle.kts")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return root
}
