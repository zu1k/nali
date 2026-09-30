package geoip

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/zu1k/nali/pkg/download"
	"github.com/zu1k/nali/pkg/ipinfo"

	"github.com/oschwald/geoip2-golang"
	"github.com/oschwald/maxminddb-golang"
	"github.com/spf13/viper"
)

// GeoIP2
type GeoIP struct {
	db     *geoip2.Reader
	ipinfo *ipinfo.DB
}

// NewGeoIP opens a MaxMind DB file. IPinfo MMDB files are detected by their
// metadata and served by the ipinfo reader; everything else is read as a
// GeoIP2/GeoLite2 City compatible database. If the file does not exist and
// downloadUrls is not empty, the database is downloaded first.
func NewGeoIP(filePath string, downloadUrls []string) (*GeoIP, error) {
	if _, err := os.Stat(filePath); err != nil && os.IsNotExist(err) {
		if len(downloadUrls) == 0 {
			log.Println("文件不存在，请自行下载 MMDB 数据库（如 GeoLite2-City），并保存在", filePath)
			return nil, err
		}
		log.Println("文件不存在，尝试从网络获取最新 MMDB 数据库")
		data, err := download.Download(filePath, downloadUrls...)
		if err != nil {
			return nil, err
		}
		if !CheckFile(data) {
			_ = os.Remove(filePath)
			return nil, errors.New("下载的 MMDB 数据库无效")
		}
	}

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

// CheckFile reports whether data is a readable MaxMind DB file.
func CheckFile(data []byte) bool {
	reader, err := maxminddb.FromBytes(data)
	if err != nil {
		return false
	}
	_ = reader.Close()
	return true
}

func (g GeoIP) Find(query string, params ...string) (result fmt.Stringer, err error) {
	if g.ipinfo != nil {
		return g.ipinfo.Find(query, params...)
	}
	ip := net.ParseIP(query)
	if ip == nil {
		return nil, errors.New("query should be valid IP")
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

// Close releases the underlying database file.
func (g GeoIP) Close() error {
	if g.ipinfo != nil {
		return g.ipinfo.Reader.Close()
	}
	return g.db.Close()
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
