//go:build linux

package cmd

import (
	"os"

	"golang.org/x/sys/unix"
)

// isInteractiveTerminal reports whether stdin is attached to a real terminal
// (as opposed to a pipe, a file or /dev/null).
func isInteractiveTerminal() bool {
	_, err := unix.IoctlGetTermios(int(os.Stdin.Fd()), unix.TCGETS)
	return err == nil
}
