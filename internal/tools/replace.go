package tools

import (
	"strings"
	"regexp"
)



func ReplaceLookAround(origin string, back_regex string, body_regex string, ahead_regex string, replace string) (result string) {
	sm := regexp.MustCompile(body_regex).FindAllStringIndex(origin, -1)

	result = ""
	last_end_id := 0

	for _, index := range sm {
		start_id := index[0]
		end_id := index[1]
		this_ip := origin[start_id:end_id]

		result += origin[last_end_id:start_id]

		if ( start_id == 0 || regexp.MustCompile(back_regex).MatchString(origin[start_id-1:start_id]) ) &&
		   ( end_id == len(origin) || regexp.MustCompile(ahead_regex).MatchString(origin[end_id:end_id+1]) ) {
			result += replace
		} else {
			result += this_ip
		}

		last_end_id = end_id
	}
	result += origin[last_end_id:]
	return
}

func AddInfoIp4(origin string, ip string, info string) (result string) {
	result = ReplaceLookAround(origin, `[^0-9\.:]`, strings.ReplaceAll(ip, ".", "\\."), `[^0-9\.]`, ip+" ["+info+"]")
	return strings.TrimRight(result, " \t")
}

func AddInfoIp6(origin string, ip string, info string) (result string) {
	result = ReplaceLookAround(origin, `[^0-9a-fA-F:]`, strings.ReplaceAll(ip, ".", "\\."), `[^0-9a-fA-F:\\.]`, ip+" ["+info+"]")
	return strings.TrimRight(result, " \t")
}

func AddInfoDomain(origin string, domain string, info string) (result string) {
	result = ReplaceLookAround(origin, `[^0-9a-zA-Z-]`, strings.ReplaceAll(domain, ".", "\\."), `[^0-9a-zA-Z-\\.]`, domain+" ["+info+"]")
	return strings.TrimRight(result, " \t")
}