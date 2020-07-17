package app

import (
	"fmt"

	"github.com/zu1k/nali/internal/ipdb"

	geoip2 "github.com/zu1k/nali/pkg/geoip"
	"github.com/zu1k/nali/pkg/qqwry"
)

var (
	db    ipdb.IPDB
	qqip  qqwry.QQwry
	geoip geoip2.GeoIP
)

func init() {
	qqip = qqwry.NewQQwry("db/qqwry.dat")
	geoip = geoip2.NewGeoIP("db/GeoLite2-City.mmdb")
	db = qqip
}

func SetDB(dbName ipdb.IPDBType) {
	switch dbName {
	case ipdb.GEOIP2:
		db = geoip
	case ipdb.QQIP:
		db = qqip
	}
}

func ParseIPs(ips []string) {
	for _, ip := range ips {
		ParseIP(ip)
	}
}

func ParseIP(ip string) {
	result := db.Find(ip)
	fmt.Println(formatResult(ip, result))
}

func formatResult(ip string, result string) string {
	if result == "" {
		result = "未找到"
	}
	return fmt.Sprintf("%s [%s]", ip, result)
}
