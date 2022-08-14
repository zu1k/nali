package common

import (
	"bytes"
)

// ScanLines scan lines but keep the suffix \r and \n
func ScanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i, j := bytes.IndexByte(data, '\r'), bytes.IndexByte(data, '\n'); i >= 0 || j >= 0 {
		// case 1: TOKEN\r\nTOKEN
		if i >= 0 && j >= 0 {
			if i+1 == j {
				return i + 2, data[:i+2], nil
			}
			if i < j {
				// case 2: TOKEN\rTOKEN\nTOKEN
				return i + 1, data[:i+1], nil
			} else {
				// case 3: TOKEN\nTOKEN\rTOKEN
				return j + 1, data[:j+1], nil
			}
		} else if i >= 0 {
			// case 4: TOKEN\rTOKEN
			return i + 1, data[:i+1], nil
		} else {
			// case 5: TOKEN\nTOKEN
			return j + 1, data[:j+1], nil
		}
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}
