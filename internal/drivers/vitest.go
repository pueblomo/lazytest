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

type VitestDriver struct {
	packageManager string
}

func (d *VitestDriver) Name() string {
	return "vitest"
}

func (d *VitestDriver) Detect(root string) (bool, error) {
	pkgPath := filepath.Join(root, "package.json")
	if _, err := os.Stat(pkgPath); err != nil {
		return false, err
	}

	data, err := os.ReadFile(pkgPath)
	if err != nil {
		return false, err
	}

	var pkg struct {
		Scripts map[string]string `json:"scripts"`
	}

	if err := json.Unmarshal(data, &pkg); err != nil {
		return false, err
	}

	testScript, ok := pkg.Scripts["test"]
	if !ok {
		return false, nil
	}

	if _, err := os.Stat(filepath.Join(root, "pnpm-lock.yaml")); err == nil {
		d.packageManager = "pnpm"
	} else if _, err := os.Stat(filepath.Join(root, "package-lock.json")); err == nil {
		d.packageManager = "npm"
	} else {
		d.packageManager = "npm"
	}

	return strings.Contains(testScript, "vitest"), nil
}

func (d *VitestDriver) RunTest(ctx context.Context, root string, testCase list.Item) error {
	tc, ok := testCase.(*types.TestCase)
	if !ok {
		return fmt.Errorf("Can't convert to TestCase")
	}
	cmd := d.buildTestCommand(ctx, root, tc.Filepath)
	testStatus, output, err := d.executeTestCommand(cmd)
	tc.TestStatus = testStatus
	tc.Output = output
	return err
}

func (d *VitestDriver) buildTestCommand(ctx context.Context, root string, filePath string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, d.packageManager, "vitest", "--run", filePath)
	cmd.Dir = root

	return cmd
}

func (d *VitestDriver) executeTestCommand(cmd *exec.Cmd) (types.TestStatus, string, error) {
	output, err := cmd.CombinedOutput()
	outputString := string(output)
	if err != nil {
		return types.StatusFailed, outputString, fmt.Errorf("test command %s failed: %w", cmd.Args, err)
	}

	if strings.Contains(outputString, "failed") {
		return types.StatusFailed, outputString, nil
	}

	if strings.Contains(outputString, "skipped") {
		return types.StatusSkipped, outputString, nil
	}

	return types.StatusPassed, outputString, nil
}

func (d *VitestDriver) DetectTestFiles(ctx context.Context, root string) ([]string, error) {
	cmd := exec.CommandContext(ctx, d.packageManager, "vitest", "list", "--filesOnly", "--json")
	cmd.Dir = root

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var tests []struct {
		File string `json:"file"`
	}

	if err := json.Unmarshal(output, &tests); err != nil {
		return nil, err
	}

	slice := make([]string, 0, len(tests))

	for _, t := range tests {
		slice = append(slice, t.File)
	}

	return slice, nil
}
