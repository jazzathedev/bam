package download

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jazzathedev/bam/internal/plugin"
)

// destPath is absolute to the file but should usually be in ~/.bam/cache
func DownloadURL(url, destPath string, ttl time.Duration) (string, error) {
	infoPath := destPath + ".info.json"

	infoFile, err := os.Open(infoPath)

	if errors.Is(err, os.ErrNotExist) {
		infoFile, err = os.Create(infoPath)

		if err != nil {
			return "", fmt.Errorf("Unable to create cache info file: %w", err)
		}
	}

	decoder := json.NewDecoder(infoFile)

	type infoJson struct {
		TTL time.Time `json:"ttl"`
	}

	var infoTtl infoJson
	decoder.Decode(&infoTtl)
	sinceExpired := time.Since(infoTtl.TTL)
	expired := sinceExpired > ttl

	if ttl == 0 {
		expired = false
	}

	infoFile.Close()

	cacheFile, err := os.Open(destPath)

	// Cache file exist and it's not expired, return path
	if (err == nil) && !expired {
		return destPath, nil
	}

	// Error cache file not exist or it is expired
	if errors.Is(err, os.ErrNotExist) || expired {
		cacheFile.Close()
		// It doesn't exist, or it does but it's expired
		// Create truncates for us, so both paths are ok
		cacheFile, err = os.Create(destPath)

		if err != nil {
			return "", fmt.Errorf("Unable to create cache file: %w", err)
		}

		defer cacheFile.Close()

		// Cache file made, GET url
		response, err := http.Get(url)
		if err != nil {
			return "", fmt.Errorf("Unable to GET url %s: %w", url, err)
		}

		defer response.Body.Close()

		// URL got, write to cache file
		_, err = io.Copy(cacheFile, response.Body)
		if err != nil {
			return "", fmt.Errorf("Unable to write to cache file: %w", err)
		}

		infoTtl.TTL = time.Now()
		newInfoJson, err := json.Marshal(&infoTtl)

		if err != nil {
			return "", fmt.Errorf("Unable to write info.json due to json encoding: %w", err)
		}

		infoFile, err = os.Create(infoPath)
		if err != nil {
			return "", fmt.Errorf("Unable to open and truncate info.json for writing: %w", err)
		}

		_, err = io.Copy(infoFile, bytes.NewBuffer(newInfoJson))

		if err != nil {
			return "", fmt.Errorf("Unable to write info.json: %w", err)
		}
	} else {
		// Error cache file exist, can't open it
		return "", fmt.Errorf("Unable to open cache file: %w", err)
	}

	return destPath, nil
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

	toolBinaryName := filepath.Base(toolBinaryPath)

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
