package engine

import (
	"os"
	"path/filepath"
	"strings"
)

// DetectProject finds an .xcworkspace or .xcodeproj in the given directory.
func DetectProject(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	// Prefer .xcworkspace over .xcodeproj
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".xcworkspace") {
			return filepath.Join(dir, e.Name()), nil
		}
	}
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".xcodeproj") {
			return filepath.Join(dir, e.Name()), nil
		}
	}

	return "", nil
}

// IsWorkspace returns true if path ends in .xcworkspace.
func IsWorkspace(path string) bool {
	return strings.HasSuffix(path, ".xcworkspace")
}
