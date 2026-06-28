package setup

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var (
	bamDir  string
	bamOnce sync.Once
)

func BamDir() string {
	bamOnce.Do(func() {
		home, err := os.UserHomeDir()
		if err != nil {
			panic("cannot resolve home dir: " + err.Error())
		}
		bamDir = filepath.Join(home, ".bam")
	})
	return bamDir
}

func MakeDirs() error {
	bam := BamDir()

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
