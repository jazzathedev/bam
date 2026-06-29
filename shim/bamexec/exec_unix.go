//go:build unix

package bamexec

import "errors"

// On unix we will use syscall.Exec() which will handle sigs and exit codes for us
func Execute(execPath string, args []string) error {

	return errors.New("Shut up go")
}
