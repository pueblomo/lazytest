package drivers

import (
	"os"
	"path/filepath"
	"strings"
)

// ContainPath validates that filePath is within root and returns the cleaned
// root-relative path. If the file escapes root, it returns root as fallback.
// This prevents path traversal attacks when constructing test commands.
func ContainPath(root string, filePath string) string {
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return root
	}

	var absFile string
	if filepath.IsAbs(filePath) {
		absFile = filePath
	} else {
		absFile = filepath.Join(rootAbs, filePath)
	}
	absFile = filepath.Clean(absFile)

	if !strings.HasPrefix(absFile, rootAbs+string(filepath.Separator)) {
		// Try exact match for root itself
		if absFile != rootAbs {
			return root
		}
	}

	return absFile
}

// findModuleRootByPom walks up from filePath to find the nearest ancestor
// containing a marker file (e.g., pom.xml, build.gradle). Returns root if
// no marker is found or if the file path is outside root.
func findModuleRootByPom(root string, filePath string, marker string) string {
	absFile := filePath
	if !filepath.IsAbs(absFile) {
		absFile = filepath.Join(root, filePath)
	}
	absFile = filepath.Clean(absFile)

	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return root
	}

	if !strings.HasPrefix(absFile, rootAbs+string(filepath.Separator)) {
		return root
	}

	dir := filepath.Dir(absFile)
	for {
		if _, err := os.Stat(filepath.Join(dir, marker)); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir || !strings.HasPrefix(dir, rootAbs) {
			break
		}
		dir = parent
	}

	return root
}
