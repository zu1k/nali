package iptools

import (
	"regexp"
	"strings"
)

var (
	ipv4re0 *regexp.Regexp
	ipv4re  *regexp.Regexp
)

func init() {
	ipv4re0 = regexp.MustCompile(`^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}$`)
	ipv4re = regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
}

func ValidIP4(str string) bool {
	str = strings.Trim(str, " ")
	return ipv4re0.MatchString(str)
}

func GetIP4FromString(str string) []string {
	str = strings.Trim(str, " ")
	return ipv4re.FindAllString(str, -1)
}
