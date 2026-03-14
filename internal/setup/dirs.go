package setup

import (
	"fmt"
	"os"
	"path/filepath"
)

func MakeDirs() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	bam := filepath.Join(home, ".bam")

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
