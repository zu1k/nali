package db

import (
	"log"
	"strings"

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

func Find(typ dbif.QueryType, query string) string {
	if result, found := queryCache.Load(query); found {
		return result.(string)
	}
	result, err := GetDB(typ).Find(query)
	if err != nil {
		return ""
	}
	r := strings.Trim(result.String(), " ")
	queryCache.Store(query, r)
	return r
}
