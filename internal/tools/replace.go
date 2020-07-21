package tools

import "strings"

func ReplaceAdd(origin string, old string, new string) string {
	subLen := len(new) - len(old)
	wanted := old + strings.Repeat(" ", subLen)
	if strings.Contains(origin, wanted) {
		return strings.ReplaceAll(origin, wanted, new)
	}
	return strings.ReplaceAll(origin, old, new)
}
