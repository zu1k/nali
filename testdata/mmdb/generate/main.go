// Command generate writes the small MMDB fixtures used by nali's tests.
//
// It lives in its own module so mmdbwriter never becomes a dependency of nali.
// Regenerate the fixtures with:
//
//	cd testdata/mmdb/generate && go run .
package main

import (
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
)

type entry struct {
	cidr   string
	record mmdbtype.Map
}

func write(name, dbType string, entries []entry) {
	w, err := mmdbwriter.New(mmdbwriter.Options{
		DatabaseType:            dbType,
		Description:             map[string]string{"en": "nali test fixture"},
		Languages:               []string{"en", "zh-CN"},
		IPVersion:               6,
		RecordSize:              24,
		IncludeReservedNetworks: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, e := range entries {
		_, network, err := net.ParseCIDR(e.cidr)
		if err != nil {
			log.Fatal(err)
		}
		if err := w.Insert(network, e.record); err != nil {
			log.Fatal(err)
		}
	}
	f, err := os.Create(filepath.Join("..", name))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if _, err := w.WriteTo(f); err != nil {
		log.Fatal(err)
	}
}

func names(en, zh string) mmdbtype.Map {
	m := mmdbtype.Map{"en": mmdbtype.String(en)}
	if zh != "" {
		m["zh-CN"] = mmdbtype.String(zh)
	}
	return m
}

func ipinfo(network, country, cc, continent, cont, asn, asName, asDomain string) mmdbtype.Map {
	m := mmdbtype.Map{
		"country":        mmdbtype.String(country),
		"country_code":   mmdbtype.String(cc),
		"continent":      mmdbtype.String(continent),
		"continent_code": mmdbtype.String(cont),
		"asn":            mmdbtype.String(asn),
		"as_name":        mmdbtype.String(asName),
		"as_domain":      mmdbtype.String(asDomain),
	}
	if network != "" {
		m["network"] = mmdbtype.String(network)
	}
	return m
}

func main() {
	write("ipinfo-lite-test.mmdb", "ipinfo bundle_location_lite.mmdb", []entry{
		{"8.8.8.0/24", ipinfo("8.8.8.0/24", "United States", "US", "North America", "NA", "AS15169", "Google LLC", "google.com")},
		// no "network" field: the reader must fall back to the matched CIDR
		{"2001:4860::/32", ipinfo("", "United States", "US", "North America", "NA", "AS15169", "Google LLC", "google.com")},
		// country only, no ASN
		{"203.0.113.0/24", ipinfo("", "Australia", "AU", "Oceania", "OC", "", "", "")},
	})

	write("geolite2-city-test.mmdb", "GeoLite2-City", []entry{
		{"81.2.69.0/24", mmdbtype.Map{
			"country": mmdbtype.Map{"iso_code": mmdbtype.String("GB"), "names": names("United Kingdom", "英国")},
			"city":    mmdbtype.Map{"names": names("London", "伦敦")},
		}},
		// English-only names: other languages must fall back to English
		{"2001:218::/32", mmdbtype.Map{
			"country": mmdbtype.Map{"iso_code": mmdbtype.String("JP"), "names": names("Japan", "")},
		}},
	})
}
