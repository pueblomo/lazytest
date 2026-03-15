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

type MavenDriver struct{}

func (d *MavenDriver) Name() string {
	return "maven"
}

func (d *MavenDriver) Detect(root string) (bool, error) {
	moduleRoots, err := findMavenModuleRoots(root)
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

func (d *MavenDriver) RunTest(ctx context.Context, root string, testCase list.Item) error {
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

func (d *MavenDriver) buildTestCommand(ctx context.Context, root string, filePath string) *exec.Cmd {
	className := pathToClassName(filePath)
	workingDir := findModuleRoot(root, filePath)

	cmd := exec.CommandContext(ctx, "mvn", "test", "-Dtest="+className, "-q")
	cmd.Dir = workingDir

	return cmd
}

func (d *MavenDriver) executeTestCommand(cmd *exec.Cmd) (types.TestStatus, string, error) {
	output, err := cmd.CombinedOutput()
	outputString := string(output)

	if err != nil {
		if strings.Contains(outputString, "BUILD FAILURE") || strings.Contains(outputString, "Tests run:") {
			return types.StatusFailed, outputString, nil
		}
		return types.StatusFailed, outputString, fmt.Errorf("test command %v failed: %w", cmd.Args, err)
	}

	if strings.Contains(outputString, "BUILD SUCCESS") {
		// Check for skipped tests
		if strings.Contains(outputString, "Skipped:") && !strings.Contains(outputString, "Skipped: 0") {
			return types.StatusSkipped, outputString, nil
		}
		return types.StatusPassed, outputString, nil
	}

	if strings.Contains(outputString, "BUILD FAILURE") {
		return types.StatusFailed, outputString, nil
	}

	return types.StatusPassed, outputString, nil
}

func (d *MavenDriver) DetectTestFiles(ctx context.Context, root string) ([]string, error) {
	var testFiles []string

	moduleRoots, err := findMavenModuleRoots(root)
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

func findMavenModuleRoots(root string) ([]string, error) {
	pomPath := filepath.Join(root, "pom.xml")
	if _, err := os.Stat(pomPath); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	moduleSet := map[string]struct{}{root: {}}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == "target" || strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if info.Name() == "pom.xml" {
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

// findModuleRoot finds the nearest ancestor directory of filePath that contains a pom.xml.
// For multi-module Maven projects, this ensures mvn runs from the correct module.
// Falls back to root if no pom.xml is found along the path.
func findModuleRoot(root string, filePath string) string {
	absFile := filePath
	if !filepath.IsAbs(absFile) {
		absFile = filepath.Join(root, filePath)
	}

	dir := filepath.Dir(absFile)
	for {
		if _, err := os.Stat(filepath.Join(dir, "pom.xml")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir || !strings.HasPrefix(dir, root) {
			break
		}
		dir = parent
	}

	return root
}

// pathToClassName converts a Java test file path to its fully qualified class name.
// e.g., "src/test/java/com/example/MyTest.java" -> "com.example.MyTest"
func pathToClassName(filePath string) string {
	normalized := strings.ReplaceAll(filePath, "\\", "/")
	dir := filepath.Dir(normalized)
	base := strings.TrimSuffix(filepath.Base(normalized), ".java")

	// Find the java source root and derive the package from the path
	parts := strings.Split(strings.ReplaceAll(dir, "\\", "/"), "/")
	var pkgParts []string
	inJava := false
	for _, part := range parts {
		if part == "java" {
			inJava = true
			continue
		}
		if inJava {
			pkgParts = append(pkgParts, part)
		}
	}

	if len(pkgParts) > 0 {
		return strings.Join(pkgParts, ".") + "." + base
	}

	return base
}
