package db

import (
	"log"
	"net"

	"github.com/spf13/viper"

	"github.com/zu1k/nali/pkg/cdn"
	"github.com/zu1k/nali/pkg/dbif"
	"github.com/zu1k/nali/pkg/geoip"
	"github.com/zu1k/nali/pkg/qqwry"
	"github.com/zu1k/nali/pkg/zxipv6wry"
)

func GetDB(typ dbif.QueryType) (db dbif.DB) {
	if db, found := dbTypeCache[typ]; found {
		return db
	}

	lang := viper.GetString("selected.lang")
	if lang == "" {
		lang = "zh-CN"
	}

	var err error
	switch typ {
	case dbif.TypeIPv4:
		selected := viper.GetString("selected.ipv4")
		if selected != "" {
			db = getDbByName(selected).get()
			break
		}

		if lang == "zh-CN" {
			db, err = qqwry.NewQQwry(getDbByName("qqwry").File)
		} else {
			db, err = geoip.NewGeoIP(getDbByName("geoip").File)
		}
	case dbif.TypeIPv6:
		selected := viper.GetString("selected.ipv6")
		if selected != "" {
			db = getDbByName(selected).get()
			break
		}

		if lang == "zh-CN" {
			db, err = zxipv6wry.NewZXwry(getDbByName("zxipv6wry").File)
		} else {
			db, err = geoip.NewGeoIP(getDbByName("geoip").File)
		}
	case dbif.TypeDomain:
		selected := viper.GetString("selected.cdn")
		if selected != "" {
			db = getDbByName(selected).get()
			break
		}

		db, err = cdn.NewCDN(getDbByName("cdn").File)
	default:
		panic("Query type not supported!")
	}

	if err != nil || db == nil {
		log.Fatalln("Database init failed:", err)
	}

	dbTypeCache[typ] = db
	return
}

func Find(typ dbif.QueryType, query string) *Result {
	if result, found := queryCache.Load(query); found {
		return result.(*Result)
	}
	// Convert NAT64 64:ff9b::/96 to IPv4
	if typ == dbif.TypeIPv6 {
		ip := net.ParseIP(query)
		if ip != nil {
			_, NAT64, _ := net.ParseCIDR("64:ff9b::/96")
			if NAT64.Contains(ip) {
				ip4 := make(net.IP, 4)
				copy(ip4, ip[12:16])
				query = ip4.String()
				typ = dbif.TypeIPv4
			}
		}
	}
	db := GetDB(typ)
	result, err := db.Find(query)
	if err != nil {
		return nil
	}
	res := &Result{db.Name(), result}
	queryCache.Store(query, res)
	return res
}
