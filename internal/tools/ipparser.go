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
	// ipv4-compatible address
	ipv6head :=
	`(` +
		`(`+
			`(0{1,4}:){6}` + `|` +
			`(0{1,4}:){1,5}:(0{1,4}:){0}` + `|`+
			`(0{1,4}:){1,4}:(0{1,4}:){1}` + `|`+
			`(0{1,4}:){1,3}:(0{1,4}:){2}` + `|`+
			`(0{1,4}:){1,2}:(0{1,4}:){3}` + `|`+
			`(0{1,4}:){1}:(0{1,4}:){4}` + `|`+
			`::(0{1,4}:){0,5}` +
		`)`+`|`+
		`(`+
			`(0{1,4}:){5}` + `|` +
			`(0{1,4}:){1,4}:(0{1,4}:){0}` + `|`+
			`(0{1,4}:){1,3}:(0{1,4}:){1}` + `|`+
			`(0{1,4}:){1,2}:(0{1,4}:){2}` + `|`+
			`(0{1,4}:){1}:(0{1,4}:){3}` + `|`+
			`::(0{1,4}:){0,4}` +
		`)`+ `[fF]{4}:`+
	`)`
	ipv6head = `(`+ipv6head+`)?`

	ipv4re = regexp.MustCompile(ipv6head+`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)

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

func ValidIP6Re(str string) bool {
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
func ValidIP6(IP string) bool {
	IP = strings.ToUpper(IP)
	n := len(IP)
	if n > 39 || n == 0 {
		return false
	}

	// 以 ":" 结尾 但是只有一个
	if strings.HasSuffix(IP, ":") && !strings.HasSuffix(IP, "::") {
		return false
	}
	// 如果"::" 有两个以上
	if strings.Count(IP, "::") > 1 {
		return false
	}
	// 如果 ":" 有8个以上
	if strings.Count(IP, ":") > 8 {
		return false
	}
	tmp := strings.Split(IP, ":")
	// 如果有ipv4， 则返回真， 前面的部分未校验。
	if ValidIP4(tmp[len(tmp)-1]) {
		return true
	}
	if strings.Contains(IP,"::") {
		var count int
		for _, v := range tmp {
			if v != "" {
				count++
				continue
			}
		}
		if count == 8 {
			return false
		}
	}

	// 对每个元素进行遍历
	for k := 0; k < n-1; {
		if IP[k] == ':' {
			k++
			continue
		} else if valid(IP[k]) {
			var bits int
			for valid(IP[k]) {
				k++
				bits++
				if bits > 4 {
					return false
				}
				if k == n {
					break
				}
			}

		} else {
			return false
		}
	}

	// 到了这一步， 可以确定是ipv6
	return true
}

func valid(i uint8) bool {
	return (i >= 'A' && i <= 'F') || (i >= '0' && i <= '9')
}