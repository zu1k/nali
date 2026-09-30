package ipinfo

import (
	"encoding/json"
	"testing"

	"github.com/oschwald/maxminddb-golang"
)

const fixture = "../../testdata/mmdb/ipinfo-lite-test.mmdb"

func openFixture(t *testing.T) *DB {
	t.Helper()
	reader, err := maxminddb.Open(fixture)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = reader.Close() })
	return &DB{Reader: reader}
}

func TestFind(t *testing.T) {
	db := openFixture(t)

	tests := []struct {
		query string
		want  Result
		text  string
	}{
		{
			query: "8.8.8.8",
			want:  Result{Network: "8.8.8.0/24", Country: "United States", CountryCode: "US", Continent: "North America", ContinentCode: "NA", ASN: "AS15169", ASName: "Google LLC", ASDomain: "google.com"},
			text:  "United States AS15169 Google LLC",
		},
		{
			// record without a "network" field uses the matched CIDR
			query: "2001:4860:4860::8888",
			want:  Result{Network: "2001:4860::/32", Country: "United States", CountryCode: "US", Continent: "North America", ContinentCode: "NA", ASN: "AS15169", ASName: "Google LLC", ASDomain: "google.com"},
			text:  "United States AS15169 Google LLC",
		},
		{
			query: "203.0.113.7",
			want:  Result{Network: "203.0.113.0/24", Country: "Australia", CountryCode: "AU", Continent: "Oceania", ContinentCode: "OC"},
			text:  "Australia",
		},
		{
			// not found: empty result, empty text
			query: "192.0.2.1",
			want:  Result{},
			text:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			got, err := db.Find(tt.query)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Fatalf("got %+v, want %+v", got, tt.want)
			}
			if got.String() != tt.text {
				t.Fatalf("String() = %q, want %q", got.String(), tt.text)
			}
		})
	}

	if _, err := db.Find("not-an-ip"); err == nil {
		t.Error("invalid IP should fail")
	}
	if db.Name() != "ipinfo" {
		t.Errorf("Name() = %q", db.Name())
	}
}

func TestResultJSON(t *testing.T) {
	data, err := json.Marshal(Result{Network: "8.8.8.0/24", ASN: "AS15169"})
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]any
	if err := json.Unmarshal(data, &fields); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"network", "country", "country_code", "continent", "continent_code", "asn", "as_name", "as_domain"} {
		if _, ok := fields[key]; !ok {
			t.Errorf("missing JSON field %q in %s", key, data)
		}
	}
}
