package common

import (
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// DecodeInput converts a line read from stdin to UTF-8. With forceGBK the
// line is always decoded as GBK. With autoChinese (see ConsoleUsesGBK) only
// lines that are not valid UTF-8 are decoded, as GB18030 (a superset of GBK
// that covers both code pages 936 and 54936), so output of legacy Windows
// tools such as tracert or nslookup is shown correctly while UTF-8 input is
// untouched.
func DecodeInput(line string, forceGBK, autoChinese bool) string {
	var enc encoding.Encoding
	switch {
	case forceGBK:
		enc = simplifiedchinese.GBK
	case autoChinese && !utf8.ValidString(line):
		enc = simplifiedchinese.GB18030
	default:
		return line
	}
	decoded, _, err := transform.String(enc.NewDecoder(), line)
	if err != nil {
		return line
	}
	return decoded
}
