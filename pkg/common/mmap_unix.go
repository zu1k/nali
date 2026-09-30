//go:build linux || darwin || freebsd

package common

import (
	"os"
	"syscall"
)

// ReadFileMapped returns the contents of a read-only database file. On this
// platform the file is memory mapped, so opening a large database does not
// read and copy the whole file: only the pages a lookup touches are loaded.
// The returned slice must not be modified, and stays valid for the life of
// the process. Database updates replace files by renaming (see SaveFile), so
// an existing mapping is never affected by an update.
func ReadFileMapped(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	size := info.Size()
	if size <= 0 || size != int64(int(size)) {
		return os.ReadFile(path)
	}

	data, err := syscall.Mmap(int(f.Fd()), 0, int(size), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return os.ReadFile(path)
	}
	return data, nil
}
