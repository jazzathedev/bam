package download_test

import (
	"testing"

	"github.com/jazzathedev/bam/internal/download"
	"github.com/jazzathedev/bam/internal/plugin"
)

var testPluginStruct = plugin.PluginConfig{Download: plugin.Download{URL: "https://nodejs.org/dist/v{version}/node-v{version}-{os}-{arch}.{ext}"}, Platform: plugin.Platform{OSMap: map[string]string{"windows": "win"}, ArchMap: map[string]string{"amd64": "x64"}, ExtMap: map[string]string{"windows": "zip"}}}

var expectedURL = "https://nodejs.org/dist/v22.22.1/node-v22.22.1-win-x64.zip"

func TestConstructUrl(t *testing.T) {
	constructedURL, err := download.ConstructURL(testPluginStruct, "22.22.1")
	if err != nil {
		t.Errorf("Failed to construct tool URL: %s", err)
	}

	if constructedURL != expectedURL {
		t.Errorf("Expected %s, got %s", expectedURL, constructedURL)
	}
}

var testPluginStructBroken = plugin.PluginConfig{Download: plugin.Download{URL: "https://nodejs.org/dist/v{version}/node-v{version}-{os}-{arch}.{ext}"}, Platform: plugin.Platform{OSMap: map[string]string{"jazzaos": "jazz"}, ArchMap: map[string]string{"amd64": "x64"}, ExtMap: map[string]string{"windows": "zip"}}}

func TestConstructUrlUnsupportedOS(t *testing.T) {
	_, err := download.ConstructURL(testPluginStructBroken, "22.22.1")
	if err == nil {
		t.Error("Expected error \"User OS {OS} not supported\"")
	}
}
