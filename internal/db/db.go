package db

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/zu1k/nali/internal/constant"
	"github.com/zu1k/nali/pkg/cdn"
	"github.com/zu1k/nali/pkg/dbif"
	"github.com/zu1k/nali/pkg/geoip"
	"github.com/zu1k/nali/pkg/ipip"
	"github.com/zu1k/nali/pkg/qqwry"
	"github.com/zu1k/nali/pkg/zxipv6wry"
)

var (
	QQWryPath        = filepath.Join(constant.HomePath, "qqwry.dat")
	ZXIPv6WryPath    = filepath.Join(constant.HomePath, "zxipv6wry.db")
	GeoLite2CityPath = filepath.Join(constant.HomePath, "GeoLite2-City.mmdb")
	IPIPFreePath     = filepath.Join(constant.HomePath, "ipipfree.ipdb")
	CDNPath          = filepath.Join(constant.HomePath, "cdn.json")

	Language       = "zh-CN"
	IPv4DBSelected = ""
	IPv6DBSelected = ""
)

func init() {
	lang := os.Getenv("NALI_LANG")
	if lang != "" {
		Language = lang
	}

	ipv4DB := os.Getenv("NALI_DB_IP4")
	if ipv4DB != "" {
		IPv4DBSelected = ipv4DB
	}

	ipv6DB := os.Getenv("NALI_DB_IP6")
	if ipv6DB != "" {
		IPv6DBSelected = ipv6DB
	}
}

func GetDB(typ dbif.QueryType) (db dbif.DB) {
	if db, found := dbCache[typ]; found {
		return db
	}

	switch typ {
	case dbif.TypeIPv4:
		if IPv4DBSelected != "" {
			db = GetIPDBbyName(IPv4DBSelected)
		} else {
			if Language == "zh-CN" {
				db = qqwry.NewQQwry(QQWryPath)
			} else {
				db = geoip.NewGeoIP(GeoLite2CityPath)
			}
		}
	case dbif.TypeIPv6:
		if IPv6DBSelected != "" {
			db = GetIPDBbyName(IPv6DBSelected)
		} else {
			if Language == "zh-CN" {
				db = zxipv6wry.NewZXwry(ZXIPv6WryPath)
			} else {
				db = geoip.NewGeoIP(GeoLite2CityPath)
			}
		}
	case dbif.TypeDomain:
		db = cdn.NewCDN(CDNPath)
	default:
		panic("Query type not supported!")
	}

	dbCache[typ] = db
	return
}

func GetIPDBbyName(name string) (db dbif.DB) {
	switch name {
	case "geo", "geoip", "geoip2":
		return geoip.NewGeoIP(GeoLite2CityPath)
	case "chunzhen", "qqip", "qqwry":
		return qqwry.NewQQwry(QQWryPath)
	case "ipip", "ipipfree", "ipip.net":
		return ipip.NewIPIPFree(IPIPFreePath)
	default:
		return qqwry.NewQQwry(QQWryPath)
	}
}

func Update() {
	qqwry.Download(QQWryPath)
	zxipv6wry.Download(ZXIPv6WryPath)
	cdn.Download(CDNPath)
}

func Find(typ dbif.QueryType, query string) string {
	if result, found := queryCache[query]; found {
		return result
	}
	result, err := GetDB(typ).Find(query, Language)
	if err != nil {
		return ""
	}
	r := strings.Trim(result.String(), " ")
	queryCache[query] = r
	return r
}
