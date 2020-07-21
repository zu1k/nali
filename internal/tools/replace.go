package tools

import (
	"strings"
)

func ReplaceAdd(origin string, old string, new string) (result string) {
	subLen := len(new) - len(old)
	wanted := old + strings.Repeat(" ", subLen)
	if strings.Contains(origin, wanted) {
		result = strings.ReplaceAll(origin, wanted, new)
	}
	result = strings.ReplaceAll(origin, old, new)
	return strings.TrimRight(result, " \t")
}
