package main

import (
	"fmt"
	"log"

	"github.com/jazzathedev/bam/internal/setup"
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
}
