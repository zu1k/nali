package tools

import (
	"net"
	"strings"

	"github.com/zu1k/nali/internal/re"
)

func GetIP4FromString(str string) []string {
	str = strings.Trim(str, " ")
	return re.IPv4Re.FindAllString(str, -1)
}

func GetIP6FromString(str string) []string {
	str = strings.Trim(str, " ")
	return re.IPv6Re.FindAllString(str, -1)
}

const (
	ValidIPv4 = iota
	ValidIPv6
	InvalidIP
)

type Valid int

func ValidIP(IP string) (v Valid) {
	for i := 0; i < len(IP); i++ {
		switch IP[i] {
		case '.':
			v = ValidIPv4
		case ':':
			v = ValidIPv6
		}
	}
	if ip := net.ParseIP(IP); ip != nil {
		return v
	}
	return InvalidIP
}
