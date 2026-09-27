//go:build !linux

package cmd

import "os"

// isInteractiveTerminal falls back to a character-device heuristic outside
// Linux, ignoring /dev/null which is a character device but never a terminal.
func isInteractiveTerminal() bool {
	info, err := os.Stdin.Stat()
	if err != nil || info.Mode()&os.ModeCharDevice == 0 {
		return false
	}
	return info.Name() != "null"
}
