package main

import (
	"fmt"
	"log"

	"github.com/jazzathedev/bam/internal/download"
	"github.com/jazzathedev/bam/internal/setup"
	"github.com/jazzathedev/bam/internal/version"
	"github.com/jazzathedev/bam/plugins"
)

func main() {
	goOs, arch := setup.DetectOS()
	fmt.Println(goOs, arch)

	if err := setup.MakeDirs(); err != nil {
		log.Fatalf("Error making critical dirs: %s", err)
	}

	pluginConfig, err := plugins.LoadBuiltinPlugins()
	if err != nil {
		log.Fatalf("Error reading plugin configs")
	}

	fmt.Println(pluginConfig[0].Name)

	resolvedVersion, err := version.ResolveVersion("22.x", pluginConfig[0])
	if err != nil {
		log.Fatalf("Resolver error: %s", err)
	}
	fmt.Println(resolvedVersion)

	toolURL, _ := download.ConstructURL(pluginConfig[0], resolvedVersion)
	toolPath, _ := download.DownloadURL(toolURL)

	fmt.Println(toolURL)
	fmt.Println(toolPath)

	fmt.Println(download.VerifyToolHash(pluginConfig[0], toolPath, resolvedVersion))
}
