package install

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/jazzathedev/bam/internal/download"
	"github.com/jazzathedev/bam/internal/extract"
	"github.com/jazzathedev/bam/internal/plugin"
	"github.com/jazzathedev/bam/internal/setup"
	"github.com/jazzathedev/bam/internal/version"
	"github.com/jazzathedev/bam/plugins"
)

func Install(tool, rawVersion string) (string, error) {
	cfg, err := findPlugin(tool)
	if err != nil {
		return "", fmt.Errorf("Error finding matching plugin: %w", err)
	}

	bam := setup.BamDir()

	versionString, err := version.ResolveVersion(rawVersion, cfg)
	if err != nil {
		return "", fmt.Errorf("Error resolving version string %s: %w", rawVersion, err)
	}

	url, err := download.ConstructURL(cfg, versionString)
	if err != nil {
		return "", fmt.Errorf("Error constructing download URL: %w", err)
	}

	fileName := path.Base(url)
	dest := filepath.Join(bam, "cache", cfg.Name, fileName)

	archivePath, err := download.DownloadURL(url, dest, time.Hour)
	if err != nil {
		return "", fmt.Errorf("Error downloading URL: %w", err)
	}

	valid, err := download.VerifyToolFile(cfg, archivePath, versionString)
	if err != nil {
		return "", fmt.Errorf("Error checking tool file validity: %w", err)
	}

	if !valid {
		err := os.Remove(archivePath)
		if err != nil {
			return "", fmt.Errorf("Error removing invalid tool file, please remove %s manually: %w", archivePath, err)
		}

		archiveJson := archivePath + ".info.json"

		err = os.Remove(archiveJson)
		if err != nil {
			return "", fmt.Errorf("Error removing invalid tool file info.json, please remove %s manually: %w", archiveJson, err)
		}

		return "", fmt.Errorf("Invalid tool file removed. Please try your command again.")
	}

	extractDest := filepath.Join(bam, "installs", cfg.Name, versionString)
	err = extract.Extractor(cfg, archivePath, extractDest)
	if err != nil {
		return "", fmt.Errorf("Error extracting tool archive: %w", err)
	}

	return versionString, nil
}

func normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func findPlugin(tool string) (plugin.PluginConfig, error) {
	configs, err := plugins.LoadBuiltinPlugins()
	if err != nil {
		return plugin.PluginConfig{}, fmt.Errorf("Error reading plugin configs: %w", err)
	}

	// Lowered and trimmed whitespace
	targetTool := normalize(tool)
	for _, cfg := range configs {
		if normalize(cfg.Name) == targetTool ||
			slices.ContainsFunc(cfg.Aliases, func(a string) bool { return normalize(a) == targetTool }) {
			return cfg, nil
		}
	}

	return plugin.PluginConfig{}, fmt.Errorf("Unable to find plugin matching tool '%s'", tool)
}

func SetGlobal(toolName, version string) error {
	bam := setup.BamDir()

	versionPinPath := filepath.Join(bam, "versions", toolName)
	err := os.MkdirAll(filepath.Dir(versionPinPath), 0755)
	if err != nil {
		return fmt.Errorf("Error making version pinning folder: %w", err)
	}

	err = os.WriteFile(versionPinPath, []byte(version), 0644)
	if err != nil {
		return fmt.Errorf("Unable to write version pinning file: %w", err)
	}

	return nil
}
