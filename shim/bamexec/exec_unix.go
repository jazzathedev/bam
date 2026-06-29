//go:build unix

package bamexec

import (
	"os"
	"syscall"
)

// On unix we will use syscall.Exec() which will handle sigs and exit codes for us
func Execute(argv []string) error {
	return syscall.Exec(argv[0], argv, os.Environ())
}
