package db

import (
	"path/filepath"

	"github.com/zu1k/nali/internal/constant"
	"github.com/zu1k/nali/pkg/cdn"
	"github.com/zu1k/nali/pkg/dbif"
	"github.com/zu1k/nali/pkg/qqwry"
	"github.com/zu1k/nali/pkg/zxipv6wry"
)

var (
	QQWryPath        = filepath.Join(constant.HomePath, "qqwry.dat")
	ZXIPv6WryPath    = filepath.Join(constant.HomePath, "zxipv6wry.db")
	GeoLite2CityPath = filepath.Join(constant.HomePath, "GeoLite2-City.mmdb")
	IPIPFreePath     = filepath.Join(constant.HomePath, "ipipfree.ipdb")
	CDNPath          = filepath.Join(constant.HomePath, "cdn.json")
)

func GetDB(typ dbif.QueryType) (db dbif.DB) {
	if db, found := dbCache[typ]; found {
		return db
	}

	switch typ {
	case dbif.TypeIPv4:
		db = qqwry.NewQQwry(QQWryPath)
	case dbif.TypeIPv6:
		db = zxipv6wry.NewZXwry(ZXIPv6WryPath)
		// geoip2.NewGeoIP(GeoLite2CityPath)
		// ipip.NewIPIPFree(IPIPFreePath)
	case dbif.TypeDomain:
		db = cdn.NewCDN(CDNPath)
	default:
		panic("Query type not supported!")
	}
	return
}

func Update() {
	qqwry.Download(QQWryPath)
	zxipv6wry.Download(ZXIPv6WryPath)
	cdn.Download(CDNPath)
}
