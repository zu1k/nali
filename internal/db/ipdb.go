package db

import (
	"os"
	"strings"
)

// ip database type
type IPDBType int

const (
	GEOIP2 = iota // geoip2
	QQIP          // chunzhen
	IPIP          // ipip.net
)

func GetIPDBType() IPDBType {
	dbname := os.Getenv("NALI_DB")
	dbname = strings.ToLower(dbname)
	switch dbname {
	case "geo", "geoip", "geoip2":
		return GEOIP2
	case "chunzhen", "qqip", "qqwry":
		return QQIP
	case "ipip", "ipipfree", "ipip.net":
		return IPIP
	default:
		return QQIP
	}
}
