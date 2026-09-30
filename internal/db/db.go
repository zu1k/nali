package db

import (
	"log"
	"net"
	"strings"

	"github.com/spf13/viper"

	"github.com/zu1k/nali/pkg/cdn"
	"github.com/zu1k/nali/pkg/dbif"
	"github.com/zu1k/nali/pkg/geoip"
	"github.com/zu1k/nali/pkg/qqwry"
	"github.com/zu1k/nali/pkg/zxipv6wry"
)

// GetDBs returns the databases selected for a query type. A selection may
// name several databases separated by commas, e.g. "qqwry,ipinfo"; all of
// them are queried and their results shown side by side.
func GetDBs(typ dbif.QueryType) (dbs []dbif.DB) {
	if dbs, found := dbTypeCache[typ]; found {
		return dbs
	}

	lang := viper.GetString("selected.lang")
	if lang == "" {
		lang = "zh-CN"
	}

	var selectedKey string
	switch typ {
	case dbif.TypeIPv4:
		selectedKey = "selected.ipv4"
	case dbif.TypeIPv6:
		selectedKey = "selected.ipv6"
	case dbif.TypeDomain:
		selectedKey = "selected.cdn"
	default:
		panic("Query type not supported!")
	}

	for _, name := range strings.Split(viper.GetString(selectedKey), ",") {
		if name = strings.TrimSpace(name); name != "" {
			dbs = append(dbs, getDbByName(name).get())
		}
	}

	if len(dbs) == 0 {
		var db dbif.DB
		var err error
		switch {
		case typ == dbif.TypeDomain:
			db, err = cdn.NewCDN(getDbByName("cdn").File, getDbByName("cdn").DownloadUrls)
		case lang != "zh-CN":
			db, err = geoip.NewGeoIP(getDbByName("geoip").File, getDbByName("geoip").DownloadUrls)
		case typ == dbif.TypeIPv4:
			db, err = qqwry.NewQQwry(getDbByName("qqwry").File, getDbByName("qqwry").DownloadUrls)
		default:
			db, err = zxipv6wry.NewZXwry(getDbByName("zxipv6wry").File)
		}
		if err != nil || db == nil {
			log.Fatalln("Database init failed:", err)
		}
		dbs = []dbif.DB{db}
	}

	dbTypeCache[typ] = dbs
	return
}

// Find queries every database selected for typ and returns the successful
// results in selection order, or nil when no database has a result.
func Find(typ dbif.QueryType, query string) []*Result {
	if results, found := queryCache.Load(query); found {
		return results.([]*Result)
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

	var results []*Result
	for _, db := range GetDBs(typ) {
		result, err := db.Find(query)
		if err != nil {
			continue
		}
		results = append(results, &Result{db.Name(), result})
	}
	if len(results) == 0 {
		return nil
	}
	queryCache.Store(query, results)
	return results
}
