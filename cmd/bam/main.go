package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jazzathedev/bam/internal/download"
	"github.com/jazzathedev/bam/internal/extract"
	"github.com/jazzathedev/bam/internal/setup"
	"github.com/jazzathedev/bam/internal/version"
	"github.com/jazzathedev/bam/plugins"
)

func main() {
	fmt.Print(extract.ExtractZip("C:\\Users\\jazza\\Downloads\\node-v22.22.1-win-x64.zip", "C:\\Users\\jazza\\Downloads\\temp\\", true))

	return

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

	toolURL, err := download.ConstructURL(pluginConfig[0], resolvedVersion)
	if err != nil {
		log.Fatalf("Error constructing plugin download URL %s", err)
	}
	fmt.Printf("toolURL: %s\n", toolURL)

	toolDest := "C:\\Users\\jazza\\.bam\\cache\\node-v22.22.1-win-x64.zip"

	toolPath, err := download.DownloadURL(toolURL, toolDest, time.Hour)
	if err != nil {
		log.Fatalf("Error downloading toolURL %s: %s", toolURL, err)
	}
	fmt.Printf("toolPath: %s\n", toolPath)

	matchedHash, err := download.VerifyToolFile(pluginConfig[0], toolPath, resolvedVersion)
	if err != nil {
		log.Fatalf("Error verifying tool hash: %s", err)
	}
	fmt.Printf("tool hash verified? %t\n", matchedHash)
}
