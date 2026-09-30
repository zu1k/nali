//go:build !windows

package common

// ConsoleUsesGBK reports whether the console uses a GBK code page. Outside
// Windows, terminals are expected to use UTF-8.
func ConsoleUsesGBK() bool {
	return false
}
