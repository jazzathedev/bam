//go:build windows

package bamexec

import (
	"os"
	"os/exec"
)

// Unfortunately windows does not have support for execve so we need to handle sigs and exit codes ourselves
// Unlike execve, windows DOES return stuff on exit so lets just make it not return on success
func Execute(argv []string) error {

	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

	err := cmd.Run()
	if _, ok := err.(*exec.ExitError); err != nil && !ok {
		return err // We couldn't even launch it
	}

	os.Exit(cmd.ProcessState.ExitCode())
	return nil // Can't reach me!
}
