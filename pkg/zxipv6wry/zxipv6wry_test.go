package zxipv6wry

import (
	"encoding/binary"
	"testing"
)

func header(start, counts uint64, size int) []byte {
	data := make([]byte, max(size, 24))
	copy(data, "IPDB")
	data[6], data[7] = 3, 8
	binary.LittleEndian.PutUint64(data[8:], counts)
	binary.LittleEndian.PutUint64(data[16:], start)
	return data[:size]
}

func TestCheckFile(t *testing.T) {
	if !CheckFile(header(24, 2, 24+2*11)) {
		t.Error("valid header rejected")
	}
	for name, data := range map[string][]byte{
		"empty":       nil,
		"html":        []byte("<html>404</html>"),
		"bad magic":   append([]byte("IPDX"), header(24, 2, 46)[4:]...),
		"short":       header(24, 2, 20),
		"truncated":   header(24, 2, 24+2*11-1),
		"zero counts": header(24, 0, 24),
	} {
		if CheckFile(data) {
			t.Errorf("%s accepted", name)
		}
	}
}
