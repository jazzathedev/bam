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
	err := os.MkdirAll(filepath.Dir(destPath), 0755)
	if err != nil {
		return "", fmt.Errorf("Unable to create destPath's dir")
	}

	infoPath := destPath + ".info.json"

	infoFile, err := os.Open(infoPath)

	if errors.Is(err, os.ErrNotExist) {
		infoFile, err = os.Create(infoPath)

		if err != nil {
			return "", fmt.Errorf("Unable to create info file: %w", err)
		}
	}

	decoder := json.NewDecoder(infoFile)

	type destInfo struct {
		TTL time.Time `json:"ttl"`
	}

	var infoJson destInfo
	decoder.Decode(&infoJson)
	sinceExpired := time.Since(infoJson.TTL)
	isExpired := sinceExpired > ttl

	if ttl == 0 {
		isExpired = false
	}

	infoFile.Close()

	destFile, err := os.Open(destPath)

	// Cache file exist and it's not expired, return path
	if (err == nil) && !isExpired {
		return destPath, nil
	}

	// Error cache file not exist or it is expired
	if errors.Is(err, os.ErrNotExist) || isExpired {
		destFile.Close()
		// It doesn't exist, or it does but it's expired
		// Create truncates for us, so both paths are ok
		destFile, err = os.Create(destPath)

		if err != nil {
			return "", fmt.Errorf("Unable to create file: %w", err)
		}

		defer destFile.Close()

		// Cache file made, GET url
		response, err := http.Get(url)
		if err != nil {
			return "", fmt.Errorf("Unable to GET url %s: %w", url, err)
		}

		defer response.Body.Close()

		// URL got, write to cache file
		_, err = io.Copy(destFile, response.Body)
		if err != nil {
			return "", fmt.Errorf("Unable to write to file: %w", err)
		}

		infoJson.TTL = time.Now()
		newInfoJson, err := json.Marshal(&infoJson)

		if err != nil {
			return "", fmt.Errorf("Unable to write info.json due to json encoding: %w", err)
		}

		infoFile, err = os.Create(infoPath)
		if err != nil {
			return "", fmt.Errorf("Unable to open and truncate info.json for writing: %w", err)
		}

		defer infoFile.Close()

		_, err = io.Copy(infoFile, bytes.NewBuffer(newInfoJson))

		if err != nil {
			return "", fmt.Errorf("Unable to write info.json: %w", err)
		}
	} else {
		// Error cache file exist, can't open it
		return "", fmt.Errorf("Unable to open file: %w", err)
	}

	return destPath, nil
}

func VerifyToolFile(pluginConfig plugin.PluginConfig, toolPath, version string) (bool, error) {
	hashURL := strings.ReplaceAll(pluginConfig.Download.HashURL, "{version}", version)

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

	toolFileName := filepath.Base(toolPath)

	toolExpectedHash, ok := hashMap[toolFileName]
	if !ok {
		return false, fmt.Errorf("Hash map does not contain entry for tool file %s", toolFileName)
	}

	toolExpectedHash = strings.ToLower(toolExpectedHash)

	toolFile, err := os.Open(toolPath)
	if err != nil {
		return false, fmt.Errorf("Unable to open tool file: %w", err)
	}

	hasher := sha256.New()

	_, err = io.Copy(hasher, toolFile)
	if err != nil {
		return false, fmt.Errorf("Unable to copy tool file into hasher: %w", err)
	}

	toolRealHash := strings.ToLower(hex.EncodeToString(hasher.Sum(nil)))

	return toolRealHash == toolExpectedHash, nil
}
