package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jazzathedev/bam/shim/bamexec"
)

type shimEntry struct {
	Tool string   `json:"tool"`
	Run  []string `json:"run"`
}

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
	shimName := strings.TrimSuffix(execName, filepath.Ext(execName))
	parentEntry, err := checkManifest(shimName)
	if err != nil {
		log.Fatalf("Unable to check shim manifest: %s", err)
	}

	parentVersionPath := filepath.Join(BamDir(), "versions", parentEntry.Tool)
	parentVersionBytes, err := os.ReadFile(parentVersionPath)
	if err != nil {
		log.Fatalf("Unable to find pinned version file: %s", err)
	}

	parentVersion := string(parentVersionBytes)
	parentVersion = strings.TrimSpace(parentVersion)

	installDir := filepath.Join(BamDir(), "installs", parentEntry.Tool, parentVersion)

	parts := []string{}
	for _, part := range parentEntry.Run {
		parts = append(parts, filepath.Join(installDir, part))
	}

	argv := append(parts, os.Args[1:]...)

	err = bamexec.Execute(argv)
	if err != nil {
		log.Fatalf("Error starting process: %s", err)
	}
}

func checkManifest(shimName string) (shimEntry, error) {
	manifestPath := filepath.Join(BamDir(), "shims", "manifest.json")

	manifestFile, err := os.Open(manifestPath)

	if err != nil {
		return shimEntry{}, fmt.Errorf("Unable to open %s: %w", manifestPath, err)
	}

	defer manifestFile.Close()

	var entries map[string]shimEntry
	decoder := json.NewDecoder(manifestFile)
	err = decoder.Decode(&entries)
	if err != nil {
		return shimEntry{}, fmt.Errorf("Invalid shim manifest %s: %w", manifestPath, err)
	}

	entry, ok := entries[shimName]
	if !ok {
		return shimEntry{}, fmt.Errorf("Can not find tool %s in shim manifest", shimName)
	}

	return entry, nil
}
