package common

import "golang.org/x/sys/windows"

// ConsoleUsesGBK reports whether the console (or, without a console, the
// system ANSI code page) uses GBK/GB18030, as on Simplified Chinese Windows.
func ConsoleUsesGBK() bool {
	cp, err := windows.GetConsoleOutputCP()
	if err != nil || cp == 0 {
		cp = windows.GetACP()
	}
	return cp == 936 || cp == 54936
}
