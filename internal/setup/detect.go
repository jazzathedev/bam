package setup

import "runtime"

func DetectOS() (string, string) {
	goOS := runtime.GOOS
	arch := runtime.GOARCH

	return goOS, arch
}
