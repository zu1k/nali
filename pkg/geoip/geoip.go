package geoip

import (
	"fmt"
	"log"
	"net"

	"github.com/oschwald/geoip2-golang"
)

type GeoIP struct {
	db *geoip2.Reader
}

func NewGeoIP(filePath string) GeoIP {
	db, err := geoip2.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return GeoIP{db: db}
}

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
