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

	// Determine package manager
	if _, err := os.Stat(filepath.Join(root, "pnpm-lock.yaml")); err == nil {
		d.packageManager = "pnpm"
	} else if _, err := os.Stat(filepath.Join(root, "package-lock.json")); err == nil {
		d.packageManager = "npm"
	} else {
		d.packageManager = "npm"
	}

	data, err := os.ReadFile(pkgPath)
	if err != nil {
		return false, err
	}

	var pkg struct {
		Scripts    map[string]string `json:"scripts"`
		Workspaces []interface{}     `json:"workspaces"`
	}

	if err := json.Unmarshal(data, &pkg); err != nil {
		return false, err
	}

	testScript, ok := pkg.Scripts["test"]
	if ok && strings.Contains(testScript, "vitest") {
		return true, nil
	}

	// If workspaces are defined, check each workspace for vitest
	if len(pkg.Workspaces) > 0 {
		workspaceDirs := []string{}
		for _, ws := range pkg.Workspaces {
			if wsStr, ok := ws.(string); ok {
				// Expand glob patterns (simple match)
				matches, _ := filepath.Glob(filepath.Join(root, wsStr))
				workspaceDirs = append(workspaceDirs, matches...)
			}
		}
		for _, wsDir := range workspaceDirs {
			wsPkgPath := filepath.Join(wsDir, "package.json")
			if data, err := os.ReadFile(wsPkgPath); err == nil {
				var wsPkg struct {
					Scripts map[string]string `json:"scripts"`
				}
				if err := json.Unmarshal(data, &wsPkg); err == nil {
					if testScript, ok := wsPkg.Scripts["test"]; ok && strings.Contains(testScript, "vitest") {
						return true, nil
					}
				}
			}
		}
	}

	return false, nil
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
	// Sanitize the file path to ensure it stays within root
	safeAbs := ContainPath(root, filePath)

	// Find the package.json root (module root) from the file's directory
	moduleRoot := d.findPkgRoot(root, safeAbs)

	// Compute relative path from module root to the test file
	relPath, err := filepath.Rel(moduleRoot, safeAbs)
	if err != nil {
		relPath = filePath // fallback to original (shouldn't happen)
	}

	// If the sanitized path is the module root itself (e.g., filePath escaped), run all tests
	if relPath == "." {
		// Run all tests in the module
		cmd := exec.CommandContext(ctx, d.packageManager, "vitest", "--run")
		cmd.Dir = moduleRoot
		return cmd
	}

	cmd := exec.CommandContext(ctx, d.packageManager, "vitest", "--run", relPath)
	cmd.Dir = moduleRoot
	return cmd
}

// findPkgRoot locates the nearest ancestor directory containing a package.json,
// starting from the given file's directory and walking upward until within root.
// If none found, returns root.
func (d *VitestDriver) findPkgRoot(rootAbs string, fileAbs string) string {
	absRoot, err := filepath.Abs(rootAbs)
	if err != nil {
		return rootAbs
	}
	dir := filepath.Dir(fileAbs)
	for {
		if _, err := os.Stat(filepath.Join(dir, "package.json")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir || !strings.HasPrefix(filepath.Clean(dir), absRoot+string(filepath.Separator)) {
			break
		}
		dir = parent
	}
	return absRoot
}

func (d *VitestDriver) executeTestCommand(cmd *exec.Cmd) (types.TestStatus, string, error) {
	output, err := cmd.CombinedOutput()
	outputString := string(output)

	if err != nil {
		// If the command started and exited with non-zero, it's a test failure
		if _, ok := err.(*exec.ExitError); ok {
			if strings.Contains(outputString, "FAIL") || strings.Contains(outputString, "failed") {
				return types.StatusFailed, outputString, nil
			}
			// Other exit errors still indicate failure
			return types.StatusFailed, outputString, nil
		}
		// Command didn't start or was canceled
		return types.StatusFailed, outputString, err
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
