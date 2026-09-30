package geoip

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/oschwald/maxminddb-golang"
	"github.com/zu1k/nali/pkg/ipinfo"

	"github.com/oschwald/geoip2-golang"
	"github.com/spf13/viper"
)

// GeoIP2
type GeoIP struct {
	db     *geoip2.Reader
	ipinfo *ipinfo.DB
}

// new geoip from database file
func NewGeoIP(filePath string) (*GeoIP, error) {
	reader, err := maxminddb.Open(filePath)
	if err != nil {
		return nil, err
	}
	isIPinfo := strings.Contains(strings.ToLower(reader.Metadata.DatabaseType), "ipinfo")
	if isIPinfo {
		return &GeoIP{ipinfo: &ipinfo.DB{Reader: reader}}, nil
	}
	reader.Close()
	db, err := geoip2.Open(filePath)
	if err != nil {
		return nil, err
	}
	return &GeoIP{db: db}, nil
}

func (g GeoIP) Find(query string, params ...string) (result fmt.Stringer, err error) {
	if g.ipinfo != nil {
		return g.ipinfo.Find(query, params...)
	}
	ip := net.ParseIP(query)
	if ip == nil {
		return nil, errors.New("Query should be valid IP")
	}
	record, err := g.db.City(ip)
	if err != nil {
		return
	}

	lang := viper.GetString("selected.lang")
	if lang == "" {
		lang = "zh-CN"
	}

	result = Result{
		Country:     getMapLang(record.Country.Names, lang),
		CountryCode: record.Country.IsoCode,
		Area:        getMapLang(record.City.Names, lang),
	}
	return
}

func (db GeoIP) Name() string {
	if db.ipinfo != nil {
		return db.ipinfo.Name()
	}
	return "geoip"
}

type Result struct {
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	Area        string `json:"area"`
}

func (r Result) String() string {
	if r.Area == "" {
		return r.Country
	} else {
		return fmt.Sprintf("%s %s", r.Country, r.Area)
	}
}

const DefaultLang = "en"

func getMapLang(data map[string]string, lang string) string {
	res, found := data[lang]
	if found {
		return res
	}
	return data[DefaultLang]
}
