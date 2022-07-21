package db

import (
	"github.com/zu1k/nali/pkg/cdn"
	"github.com/zu1k/nali/pkg/ip2region"
	"github.com/zu1k/nali/pkg/qqwry"
)

func GetDefaultDBList() List {
	return List{
		&DB{
			Name: "qqwry",
			NameAlias: []string{
				"chunzhen",
			},
			Format:       FormatQQWry,
			File:         "qqwry.dat",
			Languages:    LanguagesZH,
			Types:        TypesIPv4,
			DownloadUrls: qqwry.DownloadUrls,
		},
		&DB{
			Name: "zxipv6wry",
			NameAlias: []string{
				"zxipv6",
				"zx",
			},
			Format:    FormatZXIPv6Wry,
			File:      "zxipv6wry.db",
			Languages: LanguagesZH,
			Types:     TypesIPv6,
		},
		&DB{
			Name: "geoip",
			NameAlias: []string{
				"geoip2",
				"geolite",
				"geolite2",
			},
			Format:    FormatMMDB,
			File:      "GeoLite2-City.mmdb",
			Languages: LanguagesAll,
			Types:     TypesIP,
		},
		&DB{
			Name: "dbip",
			NameAlias: []string{
				"db-ip",
			},
			Format:    FormatMMDB,
			File:      "dbip.mmdb",
			Languages: LanguagesAll,
			Types:     TypesIP,
		},
		&DB{
			Name:      "ipip",
			Format:    FormatIPIP,
			File:      "ipipfree.ipdb",
			Languages: LanguagesZH,
			Types:     TypesIP,
		},
		&DB{
			Name: "ip2region",
			NameAlias: []string{
				"i2r",
			},
			Format:       FormatIP2Region,
			File:         "ip2region.xdb",
			Languages:    LanguagesZH,
			Types:        TypesIPv4,
			DownloadUrls: ip2region.DownloadUrls,
		},
		&DB{
			Name:      "ip2location",
			Format:    FormatIP2Location,
			File:      "IP2LOCATION-LITE-DB3.IPV6.BIN",
			Languages: LanguagesEN,
			Types:     TypesIP,
		},

		&DB{
			Name:         "cdn",
			Format:       FormatCDNYml,
			File:         "cdn.yml",
			Languages:    LanguagesZH,
			Types:        TypesCDN,
			DownloadUrls: cdn.DownloadUrls,
		},
	}
}
