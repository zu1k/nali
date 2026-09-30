//go:build !(linux || darwin || freebsd)

package common

import "os"

// ReadFileMapped returns the contents of a read-only database file. The
// returned slice must not be modified.
func ReadFileMapped(path string) ([]byte, error) {
	return os.ReadFile(path)
}
