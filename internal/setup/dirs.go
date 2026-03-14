package setup

import (
	"fmt"
	"os"
	"path/filepath"
)

func BamDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".bam"), nil
}

func MakeDirs() error {
	bam, err := BamDir()
	if err != nil {
		return fmt.Errorf("Error finding .bam folder: %w", err)
	}

	if err := os.MkdirAll(bam, 0755); err != nil {
		return fmt.Errorf("Making .bam: %w", err)
	}

	pathsToMake := []string{"shims", "installs", "cache", "plugins/builtin", "plugins/user", "versions"}

	for _, pathToMake := range pathsToMake {
		if err := os.MkdirAll(filepath.Join(bam, pathToMake), 0755); err != nil {
			return fmt.Errorf("Making %s: %w", pathToMake, err)
		}
	}

	return nil
}
