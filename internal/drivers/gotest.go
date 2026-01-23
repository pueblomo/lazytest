package drivers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/pueblomo/lazytest/internal/types"
)

type GoTestDriver struct{}

func (d *GoTestDriver) Name() string {
	return "go test"
}

func (d *GoTestDriver) Detect(root string) (bool, error) {
	goModPath := filepath.Join(root, "go.mod")
	if _, err := os.Stat(goModPath); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	hasTests := false
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), "_test.go") {
			hasTests = true
			return filepath.SkipAll
		}
		return nil
	})

	if err != nil {
		return false, err
	}

	return hasTests, nil
}

func (d *GoTestDriver) RunTest(ctx context.Context, root string, testCase list.Item) error {
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

func (d *GoTestDriver) buildTestCommand(ctx context.Context, root string, filePath string) *exec.Cmd {
	dir := filepath.Dir(filePath)
	packagePath := "./" + dir

	cmd := exec.CommandContext(ctx, "go", "test", "-v", packagePath)
	cmd.Dir = root

	return cmd
}

func (d *GoTestDriver) executeTestCommand(cmd *exec.Cmd) (types.TestStatus, string, error) {
	output, err := cmd.CombinedOutput()
	outputString := string(output)

	if err != nil {
		if strings.Contains(outputString, "FAIL") {
			return types.StatusFailed, outputString, nil
		}
		return types.StatusFailed, outputString, fmt.Errorf("test command %v failed: %w", cmd.Args, err)
	}

	if strings.Contains(outputString, "PASS") {
		return types.StatusPassed, outputString, nil
	}

	if strings.Contains(outputString, "SKIP") {
		return types.StatusSkipped, outputString, nil
	}

	return types.StatusPassed, outputString, nil
}

func (d *GoTestDriver) DetectTestFiles(ctx context.Context, root string) ([]string, error) {
	cmd := exec.CommandContext(ctx, "go", "list", "-json", "./...")
	cmd.Dir = root

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list packages: %w", err)
	}

	var testFiles []string
	decoder := json.NewDecoder(strings.NewReader(string(output)))

	for decoder.More() {
		var pkg struct {
			Dir          string   `json:"Dir"`
			TestGoFiles  []string `json:"TestGoFiles"`
			XTestGoFiles []string `json:"XTestGoFiles"`
		}

		if err := decoder.Decode(&pkg); err != nil {
			continue
		}

		relDir, err := filepath.Rel(root, pkg.Dir)
		if err != nil {
			relDir = pkg.Dir
		}

		for _, testFile := range pkg.TestGoFiles {
			testFiles = append(testFiles, filepath.Join(relDir, testFile))
		}

		for _, testFile := range pkg.XTestGoFiles {
			testFiles = append(testFiles, filepath.Join(relDir, testFile))
		}
	}

	return testFiles, nil
}
