package tools

import (
	"regexp"
	"strings"
)

var (
	ipv4re0 *regexp.Regexp
	ipv4re  *regexp.Regexp

	ipv6re0 *regexp.Regexp
	ipv6re  *regexp.Regexp
)

func init() {
	ipv4re0 = regexp.MustCompile(`^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}$`)
	ipv4re = regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)

	ipv6re0 = regexp.MustCompile(`^fe80:(:[0-9a-fA-F]{1,4}){0,4}(%\w+)?|([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|(([0-9a-fA-F]{1,4}:){0,6}[0-9a-fA-F]{1,4})?::(([0-9a-fA-F]{1,4}:){0,6}[0-9a-fA-F]{1,4})?$`)
	ipv6re = regexp.MustCompile(`fe80:(:[0-9a-fA-F]{1,4}){0,4}(%\w+)?|([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|(([0-9a-fA-F]{1,4}:){0,6}[0-9a-fA-F]{1,4})?::(([0-9a-fA-F]{1,4}:){0,6}[0-9a-fA-F]{1,4})?`)
}

func ValidIP4(str string) bool {
	str = strings.Trim(str, " ")
	return ipv4re0.MatchString(str)
}

func ValidIP6(str string) bool {
	str = strings.Trim(str, " ")
	return ipv6re0.MatchString(str)
}

func GetIP4FromString(str string) []string {
	str = strings.Trim(str, " ")
	return ipv4re.FindAllString(str, -1)
}

func GetIP6FromString(str string) []string {
	str = strings.Trim(str, " ")
	return ipv6re.FindAllString(str, -1)
}
