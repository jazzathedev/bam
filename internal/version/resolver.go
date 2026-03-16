package version

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/mod/semver"

	"github.com/jazzathedev/bam/internal/plugin"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
)

func needsExpanding(versionString string) bool {
	expandHints := []string{"x", "latest"}

	for _, hint := range expandHints {
		if strings.Contains(versionString, hint) {
			return true
		}
	}

	return false
}

func filter(array []string, filterFunction func(string) bool) []string {
	var filteredArray []string

	for _, item := range array {
		if filterFunction(item) {
			filteredArray = append(filteredArray, item)
		}
	}

	return filteredArray
}

func ResolveVersion(rawVersionString string, pluginConfig plugin.PluginConfig) (string, error) {
	if !needsExpanding(rawVersionString) {
		return rawVersionString, nil
	}

	// TODO: cache the version list response to ~/.bam/cache/<tool>/versions.json
	// with a TTL (e.g. 1 hour). Check cache before hitting the network.
	// Implement after Component 4 (downloader) has cache infrastructure in place.
	response, err := http.Get(pluginConfig.Versions.ListURL)
	if err != nil {
		return "", fmt.Errorf("Error fetching tool versions list: %w", err)
	}

	defer response.Body.Close()

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading response body: %w", err)
	}

	responseJson, err := oj.Parse(responseBytes)
	if err != nil {
		return "", fmt.Errorf("Error parsing JSON response: %w", err)
	}

	jsonExpression, err := jp.ParseString(pluginConfig.Versions.ListPath)
	if err != nil {
		return "", fmt.Errorf("Error parsing JSON expression: %w", err)
	}

	responseVersions := jsonExpression.Get(responseJson)
	var versionStrings []string

	for _, responseVersion := range responseVersions {
		if versionString, ok := responseVersion.(string); ok {
			versionStrings = append(versionStrings, versionString)
		}
	}

	// Sorts ASCENDING, so must retrieve LAST item in array
	semver.Sort(versionStrings)

	var strippedVersions []string

	for _, versionString := range versionStrings {
		strippedVersion := strings.TrimPrefix(versionString, "v")
		strippedVersions = append(strippedVersions, strippedVersion)
	}

	if len(strippedVersions) == 0 {
		return "", fmt.Errorf("no versions found for %s", rawVersionString)
	}

	var majorVersions []string

	if strings.Contains(rawVersionString, "x") {
		majorVersions = filter(strippedVersions, func(str string) bool {
			return strings.HasPrefix(str, strings.Split(rawVersionString, ".")[0])
		})

		return majorVersions[len(majorVersions)-1], nil
	}

	if strings.Contains(rawVersionString, "latest") {
		return strippedVersions[len(strippedVersions)-1], nil
	}

	return "", fmt.Errorf("could not find a matching version for %s", rawVersionString)
}
