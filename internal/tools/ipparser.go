package tools

import (
	"regexp"
	"strings"
)

var (
	ipv4re *regexp.Regexp

	ipv6re0 *regexp.Regexp
	ipv6re  *regexp.Regexp
)

func init() {
	ipv4re = regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)

	ipv6re0 = regexp.MustCompile(`^fe80:(:[0-9a-fA-F]{1,4}){0,4}(%\w+)?|([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|(([0-9a-fA-F]{1,4}:){0,6}[0-9a-fA-F]{1,4})?::(([0-9a-fA-F]{1,4}:){0,6}[0-9a-fA-F]{1,4})?$`)
	ipv6re = regexp.MustCompile(`fe80:(:[0-9a-fA-F]{1,4}){0,4}(%\w+)?|([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|(([0-9a-fA-F]{1,4}:){0,6}[0-9a-fA-F]{1,4})?::(([0-9a-fA-F]{1,4}:){0,6}[0-9a-fA-F]{1,4})?`)
}

func ValidIP4(IP string) bool {
	arr := strings.Split(IP, ".")
	if len(arr) != 4 {
		return false
	}
	for _, elem := range arr {
		if elem == "" {
			return false
		}
		if len(elem) > 1 && elem[0] == '0' {
			return false
		}
		num := 0
		for _, c := range elem {
			if c >= '0' && c <= '9' {
				num = num*10 + int(c-'0')
			} else {
				return false
			}
		}
		if num > 255 {
			return false
		}
	}
	return true
}

func ValidIP6(str string) bool {
	str = strings.Trim(str, " ")
	return ipv6re0.MatchString(str)

	//return isIPV6(str)
}

func GetIP4FromString(str string) []string {
	str = strings.Trim(str, " ")
	return ipv4re.FindAllString(str, -1)
}

func GetIP6FromString(str string) []string {
	str = strings.Trim(str, " ")
	return ipv6re.FindAllString(str, -1)
}

//IPV6地址的判断：
//1. 用“：”分割字符串，若长度不等于8，则return Neither
//2. 遍历每一个数组的每一个元素，若元素的长度大于4，则return Neither
//3. 判断每一个元素的字符，若出现非0-9，A-F的字符，则return Neither
func isIPV6(IP string) bool {
	IP = strings.ToUpper(IP)
	arr := strings.Split(IP, ":")
	if len(arr) != 8 {
		return false
	}
	for _, elem := range arr {
		if elem == "" || len(elem) > 4 {
			return false
		}

		for _, c := range elem {
			if (c >= '0' && c <= '9') || (c >= 'A' && c <= 'F') {
				continue
			} else {
				return false
			}
		}
	}
	return true
}
