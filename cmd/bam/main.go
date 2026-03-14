package main

import (
	"fmt"
	"log"

	"github.com/jazzathedev/bam/internal/setup"
)

func main() {
	goOs, arch := setup.DetectOS()
	fmt.Println(goOs, arch)

	if err := setup.MakeDirs(); err != nil {
		log.Fatalf("Error making critical dirs: %s", err)
	}
}
