package geoip

import (
	"encoding/json"
	"net"
	"os"
	"testing"

	"github.com/zu1k/nali/pkg/ipinfo"
)

func TestIPinfoDatabase(t *testing.T) {
	path := os.Getenv("IPINFO_TEST_DB")
	if path == "" {
		t.Skip("set IPINFO_TEST_DB to run against an IPinfo Lite MMDB")
	}
	db, err := NewGeoIP(path, nil)
	if err != nil {
		t.Fatal(err)
	}
	if db.ipinfo == nil || db.Name() != "ipinfo" {
		t.Fatal("IPinfo database was not detected")
	}
	defer db.ipinfo.Reader.Close()
	for _, query := range []string{"8.8.8.8", "2001:4860:4860::8888"} {
		t.Run(query, func(t *testing.T) {
			value, err := db.Find(query)
			if err != nil {
				t.Fatal(err)
			}
			result, ok := value.(ipinfo.Result)
			if !ok {
				t.Fatalf("unexpected result type %T", value)
			}
			if result.CountryCode != "US" || result.ASN != "AS15169" || result.ASName == "" || result.ASDomain == "" || result.ContinentCode != "NA" || result.Country == "" || result.Continent == "" {
				t.Fatalf("incomplete record: %+v", result)
			}
			_, network, err := net.ParseCIDR(result.Network)
			if err != nil || !network.Contains(net.ParseIP(query)) {
				t.Fatalf("invalid network %q for %s", result.Network, query)
			}
			data, err := json.Marshal(result)
			if err != nil {
				t.Fatal(err)
			}
			var fields map[string]interface{}
			if err := json.Unmarshal(data, &fields); err != nil {
				t.Fatal(err)
			}
			for _, key := range []string{"network", "country", "country_code", "continent", "continent_code", "asn", "as_name", "as_domain"} {
				if _, ok := fields[key]; !ok {
					t.Errorf("missing JSON field %s", key)
				}
			}
			t.Log(string(data))
		})
	}
	if _, err := db.Find("invalid"); err == nil {
		t.Error("invalid IP should fail")
	}
	value, err := db.Find("127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	if value.String() != "" {
		t.Fatalf("unexpected loopback result: %v", value)
	}
}
