package download

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/jazzathedev/bam/internal/plugin"
	"github.com/jazzathedev/bam/internal/setup"
)

func DownloadURL(toolURL string) (string, error) {
	fileName := path.Base(toolURL)

	bamDir, err := setup.BamDir()
	if err != nil {
		return "", fmt.Errorf("Could not find ~/.bam: %w", err)
	}

	cacheDir := path.Join(bamDir, "cache")
	cacheFilePath := path.Join(cacheDir, fileName)

	toolFile, err := os.Open(cacheFilePath)

	if err != nil {

		if errors.Is(err, os.ErrNotExist) {
			// Error file not exist, make it
			toolFile, err = os.Create(cacheFilePath)

			if err != nil {
				return "", fmt.Errorf("Unable to create tool cache file: %w", err)
			}

			defer toolFile.Close()

			response, err := http.Get(toolURL)

			if err != nil {
				return "", fmt.Errorf("Unable to GET tool url %s: %w", toolURL, err)
			}

			defer response.Body.Close()

			_, err = io.Copy(toolFile, response.Body)
			if err != nil {
				return "", fmt.Errorf("Unable to write tool to cache file: %w", err)
			}
		} else {
			// Error file exist, can't open it
			return "", fmt.Errorf("Unable to open tool cache file: %w", err)
		}
	} else {
		// File exist, return path
		return cacheFilePath, nil
	}

	return cacheFilePath, nil
}

func VerifyToolHash(pluginStruct plugin.PluginConfig, toolBinaryPath, version string) (bool, error) {
	hashURL := strings.ReplaceAll(pluginStruct.Download.HashURL, "{version}", version)

	response, err := http.Get(hashURL)
	if err != nil {
		return false, fmt.Errorf("Unable to fetch hash URL: %w", err)
	}

	defer response.Body.Close()

	hashURLBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return false, fmt.Errorf("Unable to read response body: %w", err)
	}

	body := string(hashURLBytes)
	lines := strings.Split(body, "\n")

	hashMap := make(map[string]string)

	for _, hashLine := range lines {
		parts := strings.Split(hashLine, "  ")

		if len(parts) == 1 {
			continue
		}

		if len(parts) != 2 {
			return false, fmt.Errorf("Unsupported hash URL contents format")
		}

		hashMap[parts[1]] = parts[0]
	}

	toolBinaryName := path.Base(toolBinaryPath)

	toolExpectedHash, ok := hashMap[toolBinaryName]
	if !ok {
		return false, fmt.Errorf("Hash map does not contain entry for tool binary %s", toolBinaryName)
	}

	toolExpectedHash = strings.ToLower(toolExpectedHash)

	toolFile, err := os.Open(toolBinaryPath)
	if err != nil {
		return false, fmt.Errorf("Unable to open tool binary: %w", err)
	}

	hasher := sha256.New()

	_, err = io.Copy(hasher, toolFile)
	if err != nil {
		return false, fmt.Errorf("Unable to copy tool binary into hasher: %w", err)
	}

	toolRealHash := strings.ToLower(hex.EncodeToString(hasher.Sum(nil)))

	return toolRealHash == toolExpectedHash, nil
}
