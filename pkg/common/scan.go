package common

import (
	"regexp"
)

var newlineReg = regexp.MustCompile(`\r?\n|\r\n?`)

// ScanLines scan lines but keep the suffix \r and \n
func ScanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if loc := newlineReg.FindIndex(data); len(loc) > 0 {
		delimiterLen := 1
		i := loc[0]
		if i+1 < len(data) && data[i] == '\r' && data[i+1] == '\n' {
			delimiterLen = 2
		}
		return i + delimiterLen, data[:i+delimiterLen], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}
