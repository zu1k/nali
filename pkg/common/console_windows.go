package common

import "golang.org/x/sys/windows"

// ConsoleUsesGBK reports whether the console (or, without a console, the
// system ANSI code page) uses GBK (936) or GB18030 (54936), as on Simplified
// Chinese Windows. Input is then decoded as GB18030, which covers both.
func ConsoleUsesGBK() bool {
	cp, err := windows.GetConsoleOutputCP()
	if err != nil || cp == 0 {
		cp = windows.GetACP()
	}
	return cp == 936 || cp == 54936
}
