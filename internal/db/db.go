package db

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/zu1k/nali/internal/constant"
	"github.com/zu1k/nali/pkg/cdn"
	"github.com/zu1k/nali/pkg/dbif"
	"github.com/zu1k/nali/pkg/geoip"
	"github.com/zu1k/nali/pkg/ip2region"
	"github.com/zu1k/nali/pkg/ipip"
	"github.com/zu1k/nali/pkg/qqwry"
	"github.com/zu1k/nali/pkg/zxipv6wry"
	"github.com/zu1k/nali/pkg/ip2location"
)

var (
	QQWryPath        = filepath.Join(constant.HomePath, "qqwry.dat")
	ZXIPv6WryPath    = filepath.Join(constant.HomePath, "zxipv6wry.db")
	GeoLite2CityPath = filepath.Join(constant.HomePath, "GeoLite2-City.mmdb")
	IPIPFreePath     = filepath.Join(constant.HomePath, "ipipfree.ipdb")
	Ip2RegionPath    = filepath.Join(constant.HomePath, "ip2region.db")
	CDNPath          = filepath.Join(constant.HomePath, "cdn.yml")
	IP2LocationPath    = filepath.Join(constant.HomePath, "IP2LOCATION-LITE-DB3.IPV6.BIN")

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

	var err error

	switch typ {
	case dbif.TypeIPv4:
		if IPv4DBSelected != "" {
			db, err = GetIPDBbyName(IPv4DBSelected)
		} else {
			if Language == "zh-CN" {
				db, err = qqwry.NewQQwry(QQWryPath)
			} else {
				db, err = geoip.NewGeoIP(GeoLite2CityPath)
			}
		}
	case dbif.TypeIPv6:
		if IPv6DBSelected != "" {
			db, err = GetIPDBbyName(IPv6DBSelected)
		} else {
			if Language == "zh-CN" {
				db, err = zxipv6wry.NewZXwry(ZXIPv6WryPath)
			} else {
				db, err = geoip.NewGeoIP(GeoLite2CityPath)
			}
		}
	case dbif.TypeDomain:
		db, err = cdn.NewCDN(CDNPath)
	default:
		panic("Query type not supported!")
	}

	if err != nil || db == nil {
		log.Fatalln("Database init failed:", err)
	}

	dbCache[typ] = db
	return
}

func GetIPDBbyName(name string) (dbif.DB, error) {
	name = strings.ToLower(name)
	switch name {
	case "geo", "geoip", "geoip2":
		return geoip.NewGeoIP(GeoLite2CityPath)
	case "chunzhen", "qqip", "qqwry":
		return qqwry.NewQQwry(QQWryPath)
	case "ipip", "ipipfree", "ipip.net":
		return ipip.NewIPIPFree(IPIPFreePath)
	case "ip2region", "region", "i2r":
		return ip2region.NewIp2Region(Ip2RegionPath)
	case "ip2location", "bin", "ipl":
		return ip2locationdb.NewIP2LocationDB(IP2LocationPath)
	default:
		return qqwry.NewQQwry(QQWryPath)
	}
}

func getDBInfoMap() map[string]func() error {
	return map[string]func() error{
		"chunzhen": func() error {
			log.Println("正在下载最新 纯真 IPv4数据库...")
			_, err := qqwry.Download(QQWryPath)
			if err != nil {
				log.Fatalln("数据库 QQWry 下载失败:", err)
			}
			return err
		},
		"zxipv6": func() error {
			log.Println("正在下载最新 ZX IPv6数据库...")
			_, err := zxipv6wry.Download(ZXIPv6WryPath)
			if err != nil {
				log.Fatalln("数据库 ZXIPv6Wry 下载失败:", err)
			}
			return err
		},
		"ip2region": func() error {
			log.Println("正在下载最新 Ip2Region 数据库...")
			_, err := ip2region.Download(Ip2RegionPath)
			if err != nil {
				log.Fatalln("数据库 Ip2Region 下载失败:", err)
			}
			return err
		},
		"cdn": func() error {
			log.Println("正在下载最新 CDN服务提供商数据库...")
			_, err := cdn.Download(CDNPath)
			if err != nil {
				log.Fatalln("数据库 CDN 下载失败:", err)
			}
			return err
		},
	}
}

func UpdateDB(dbName ...string) {
	dbInfo := getDBInfoMap()
	isAll := false
	if len(dbName) == 0 {
		isAll = true
	}
	keySet := make(map[string]struct{})
	for _, v := range dbName {
		keySet[v] = struct{}{}
	}
	for key, action := range dbInfo {
		_, ok := keySet[key]
		if !isAll && !ok {
			continue
		}
		if err := action(); err != nil {
			// keep loop
			continue
		}
	}
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
