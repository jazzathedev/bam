package download

import (
	"fmt"
	"strings"

	"github.com/jazzathedev/bam/internal/plugin"
	"github.com/jazzathedev/bam/internal/setup"
)

func ConstructURL(pluginStruct plugin.PluginConfig, version string) (string, error) {

	urlTemplate := pluginStruct.Download.URL
	// URL template contains "version", "os", "arch" and "ext". All but "version" need mapping.

	osMapping := pluginStruct.Platform.OSMap
	archMapping := pluginStruct.Platform.ArchMap
	extMapping := pluginStruct.Platform.ExtMap

	userOS, userArch := setup.DetectOS()

	mappedOS, ok := osMapping[userOS]
	if !ok {
		return "", fmt.Errorf("User OS %s not supported.", userOS)
	}

	mappedArch, ok := archMapping[userArch]
	if !ok {
		return "", fmt.Errorf("User ARCH %s not supported.", userArch)
	}

	mappedExt, ok := extMapping[userOS]
	if !ok {
		return "", fmt.Errorf("User EXT %s not supported.", userOS)
	}

	urlTemplate = strings.ReplaceAll(urlTemplate, "{version}", version)
	urlTemplate = strings.ReplaceAll(urlTemplate, "{os}", mappedOS)
	urlTemplate = strings.ReplaceAll(urlTemplate, "{arch}", mappedArch)
	urlTemplate = strings.ReplaceAll(urlTemplate, "{ext}", mappedExt)

	return urlTemplate, nil
}
