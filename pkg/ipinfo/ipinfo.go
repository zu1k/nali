package ipinfo

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/oschwald/maxminddb-golang"
)

// DB reads IPinfo's flat MMDB records for both IPv4 and IPv6.
type DB struct{ Reader *maxminddb.Reader }

type Result struct {
	Network       string `maxminddb:"network" json:"network"`
	Country       string `maxminddb:"country" json:"country"`
	CountryCode   string `maxminddb:"country_code" json:"country_code"`
	Continent     string `maxminddb:"continent" json:"continent"`
	ContinentCode string `maxminddb:"continent_code" json:"continent_code"`
	ASN           string `maxminddb:"asn" json:"asn"`
	ASName        string `maxminddb:"as_name" json:"as_name"`
	ASDomain      string `maxminddb:"as_domain" json:"as_domain"`
}

func (d *DB) Find(query string, params ...string) (fmt.Stringer, error) {
	ip := net.ParseIP(query)
	if ip == nil {
		return nil, errors.New("query should be valid IP")
	}
	var result Result
	network, found, err := d.Reader.LookupNetwork(ip, &result)
	if err != nil {
		return nil, err
	}
	if found && result.Network == "" {
		result.Network = network.String()
	}
	return result, nil
}

func (d *DB) Name() string { return "ipinfo" }

func (r Result) String() string {
	parts := []string{}
	for _, value := range []string{r.Country, r.ASN, r.ASName} {
		if value != "" {
			parts = append(parts, value)
		}
	}
	return strings.Join(parts, " ")
}
