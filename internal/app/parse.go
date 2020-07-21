package app

import (
	"fmt"
	"path/filepath"

	"github.com/zu1k/nali/constant"
	"github.com/zu1k/nali/internal/ipdb"
	"github.com/zu1k/nali/internal/tools"
	geoip2 "github.com/zu1k/nali/pkg/geoip"
	"github.com/zu1k/nali/pkg/ipip"
	"github.com/zu1k/nali/pkg/qqwry"
	"github.com/zu1k/nali/pkg/zxipv6wry"
)

var (
	db    []ipdb.IPDB
	qqip  qqwry.QQwry
	geoip geoip2.GeoIP
)

// init ip db content
func InitIPDB(ipdbtype ipdb.IPDBType) {
	db = make([]ipdb.IPDB, 1)
	switch ipdbtype {
	case ipdb.GEOIP2:
		db[0] = geoip2.NewGeoIP(filepath.Join(constant.HomePath, "GeoLite2-City.mmdb"))
	case ipdb.QQIP:
		db[0] = qqwry.NewQQwry(filepath.Join(constant.HomePath, "qqwry.dat"))
		db = append(db, zxipv6wry.NewZXwry(filepath.Join(constant.HomePath, "ipv6wry.db")))
	case ipdb.IPIP:
		db[0] = ipip.NewIPIPFree(filepath.Join(constant.HomePath, "ipipfree.ipdb"))
		db = append(db, zxipv6wry.NewZXwry(filepath.Join(constant.HomePath, "ipv6wry.db")))
	}
}

// parse several ips
func ParseIPs(ips []string) {
	db0 := db[0]
	var db1 ipdb.IPDB
	if len(db) > 1 {
		db1 = db[1]
	} else {
		db1 = nil
	}
	for _, ip := range ips {
		if tools.ValidIP4(ip) {
			result := db0.Find(ip)
			fmt.Println(formatResult(ip, result))
		} else if tools.ValidIP6(ip) && db1 != nil {
			result := db1.Find(ip)
			fmt.Println(formatResult(ip, result))
		} else {
			fmt.Println(ReplaceIPInString(ip))
		}
	}
}

func ReplaceIPInString(str string) (result string) {
	db0 := db[0]
	var db1 ipdb.IPDB
	if len(db) > 1 {
		db1 = db[1]
	} else {
		db1 = nil
	}

	result = str
	ip4s := tools.GetIP4FromString(str)
	for _, ip := range ip4s {
		info := db0.Find(ip)
		result = tools.ReplaceAdd(result, ip, formatResult(ip, info))
	}

	ip6s := tools.GetIP6FromString(str)
	for _, ip := range ip6s {
		info := db1.Find(ip)
		result = tools.ReplaceAdd(result, ip, formatResult(ip, info))
	}
	return
}

func formatResult(ip string, result string) string {
	if result == "" {
		result = "未找到"
	}
	return fmt.Sprintf("%s [%s]", ip, result)
}
