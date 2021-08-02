package app

import (
	"path/filepath"

	db2 "github.com/zu1k/nali/internal/db"

	"github.com/zu1k/nali/internal/constant"

	geoip2 "github.com/zu1k/nali/pkg/geoip"
	"github.com/zu1k/nali/pkg/ipip"
	"github.com/zu1k/nali/pkg/qqwry"
	"github.com/zu1k/nali/pkg/zxipv6wry"
)

var (
	db    []db2.IPDB
	qqip  qqwry.QQwry
	geoip geoip2.GeoIP
)

// InitIPDB init ip database content
func InitIPDB(ipdbtype db2.IPDBType) {
	db = make([]db2.IPDB, 1)
	switch ipdbtype {
	case db2.GEOIP2:
		db[0] = geoip2.NewGeoIP(filepath.Join(constant.HomePath, "GeoLite2-City.mmdb"))
	case db2.QQIP:
		db[0] = qqwry.NewQQwry(filepath.Join(constant.HomePath, "qqwry.dat"))
		db = append(db, zxipv6wry.NewZXwry(filepath.Join(constant.HomePath, "ipv6wry.database")))
	case db2.IPIP:
		db[0] = ipip.NewIPIPFree(filepath.Join(constant.HomePath, "ipipfree.ipdb"))
		db = append(db, zxipv6wry.NewZXwry(filepath.Join(constant.HomePath, "ipv6wry.database")))
	}
}
