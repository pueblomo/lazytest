package drivers

import (
	"path/filepath"
	"strings"
)

func pathToClassName(filePath string) string {
	normalized := strings.ReplaceAll(filePath, "\\", "/")
	dir := filepath.Dir(normalized)
	base := strings.TrimSuffix(filepath.Base(normalized), ".java")

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
