package geoip

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/oschwald/geoip2-golang"
)

// GeoIP2
type GeoIP struct {
	db *geoip2.Reader
}

// new geoip from db file
func NewGeoIP(filePath string) (geoip GeoIP) {
	// 判断文件是否存在
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		log.Println("文件不存在，请自行下载 Geoip2 City库，并保存在", filePath)
		os.Exit(1)
	} else {
		db, err := geoip2.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		geoip = GeoIP{db: db}
	}
	return
}

// find ip info
func (g GeoIP) Find(ip string) string {
	ipData := net.ParseIP(ip)
	record, err := g.db.City(ipData)
	if err != nil {
		log.Fatal(err)
	}
	country := record.Country.Names["zh-CN"]
	city := record.City.Names["zh-CN"]
	if city == "" {
		return country
	} else {
		return fmt.Sprintf("%s %s", country, city)
	}
}
