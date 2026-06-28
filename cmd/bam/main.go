package main

import (
	"fmt"
	"log"

	"github.com/jazzathedev/bam/internal/install"
)

func main() {
	installedTool, err := install.Install("node", "latest")
	if err != nil {
		log.Fatalf("install failed: %s", err)
	}

	err = install.SetGlobal("node", installedTool.Name)
	if err != nil {
		log.Fatalf("pinning failed: %s", err)
	}

	fmt.Println("done")
}
