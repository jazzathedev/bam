package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jazzathedev/bam/shim/bamexec"
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

func main() {
	execPath, err := os.Executable()
	if err != nil {
		log.Fatalf("Unable to find executable path: %s", err)
	}

	execName := filepath.Base(execPath)
	toolName := strings.TrimSuffix(execName, filepath.Ext(execName))

	toolVersionPath := filepath.Join(BamDir(), "versions", toolName)
	toolVersion, err := os.ReadFile(toolVersionPath)
	if err != nil {
		log.Fatalf("Unable to find pinned version file: %s", err)
	}

	toolDir := filepath.Join(BamDir(), "installs", toolName, string(toolVersion))
	// Somehow need to go from the toolDir to the various executables/scripts in it...?

	// THIS DOES NOT WORK!!! IT IS A WIP
	err = bamexec.Execute(toolDir, os.Args)
}
