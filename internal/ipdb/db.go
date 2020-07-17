package ipdb

import (
	"os"
	"strings"
)

// ip db interface
type IPDB interface {
	Find(ip string) string
}

func GetIPDBType() IPDBType {
	dbname := os.Getenv("NALI_DB")
	dbname = strings.ToLower(dbname)
	switch dbname {
	case "geo", "geoip", "geoip2":
		return GEOIP2
	case "chunzhen", "qqip", "qqwry":
		return QQIP
	default:
		return QQIP
	}
}
