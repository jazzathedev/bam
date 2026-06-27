package main

import (
	"fmt"
	"log"

	"github.com/jazzathedev/bam/internal/install"
)

func main() {
	if err := install.Install("node", "latest"); err != nil {
		log.Fatalf("install failed: %s", err)
	}
	fmt.Println("done")
}
