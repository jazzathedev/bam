//go:build windows

package bamexec

import (
	"errors"
	"os"
	"os/exec"
)

// Unfortunately windows does not have support for execve so we need to handle sigs and exit codes ourselves
// Unlike execve, windows DOES return stuff on exit so lets just make it not return on success
func Execute(execPath string, args []string) error {

	cmd := exec.Command("")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	cmd.Run()
	cmd.ProcessState.ExitCode() // Yummy exit code if its non 0 exit

	return errors.New("Shut up go")

}
