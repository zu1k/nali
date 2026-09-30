package common

import (
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// DecodeInput converts a line read from stdin to UTF-8. With forceGBK the
// line is always decoded as GBK. With autoGBK (see ConsoleUsesGBK) only lines
// that are not valid UTF-8 are decoded, so output of legacy Windows tools such
// as tracert or nslookup is shown correctly while UTF-8 input is untouched.
func DecodeInput(line string, forceGBK, autoGBK bool) string {
	if !forceGBK && (!autoGBK || utf8.ValidString(line)) {
		return line
	}
	decoded, _, err := transform.String(simplifiedchinese.GBK.NewDecoder(), line)
	if err != nil {
		return line
	}
	return decoded
}
