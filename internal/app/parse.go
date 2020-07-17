package app

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/zu1k/nali/internal/iptools"

	"github.com/zu1k/nali/constant"

	"github.com/zu1k/nali/internal/ipdb"

	geoip2 "github.com/zu1k/nali/pkg/geoip"
	"github.com/zu1k/nali/pkg/qqwry"
)

var (
	db    ipdb.IPDB
	qqip  qqwry.QQwry
	geoip geoip2.GeoIP
)

// init ip db content
func InitIPDB(ipdbtype ipdb.IPDBType) {
	switch ipdbtype {
	case ipdb.GEOIP2:
		db = geoip2.NewGeoIP(filepath.Join(constant.HomePath, "GeoLite2-City.mmdb"))
	case ipdb.QQIP:
		db = qqwry.NewQQwry(filepath.Join(constant.HomePath, "qqwry.dat"))
	}
}

// parse several ips
func ParseIPs(ips []string) {
	for _, ip := range ips {
		if iptools.ValidIP4(ip) {
			result := db.Find(ip)
			fmt.Println(formatResult(ip, result))
		} else {
			fmt.Println(ReplaceInString(ip))
		}
	}
}

func ReplaceInString(str string) (result string) {
	result = str
	ips := iptools.GetIP4FromString(str)
	for _, ip := range ips {
		info := db.Find(ip)
		result = strings.ReplaceAll(result, ip, formatResult(ip, info))
	}
	return
}

func formatResult(ip string, result string) string {
	if result == "" {
		result = "未找到"
	}
	return fmt.Sprintf("%s [%s]", ip, result)
}
