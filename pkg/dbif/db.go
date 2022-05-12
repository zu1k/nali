package dbif

import (
	"fmt"

	"github.com/zu1k/nali/pkg/cdn"
	"github.com/zu1k/nali/pkg/geoip"
	"github.com/zu1k/nali/pkg/ip2region"
	"github.com/zu1k/nali/pkg/ipip"
	"github.com/zu1k/nali/pkg/qqwry"
	"github.com/zu1k/nali/pkg/zxipv6wry"
	"github.com/zu1k/nali/pkg/ip2location"
)

type QueryType uint

const (
	TypeIPv4 = iota
	TypeIPv6
	TypeDomain
)

type DB interface {
	Find(query string, params ...string) (result fmt.Stringer, err error)
}

var (
	_ DB = &qqwry.QQwry{}
	_ DB = &zxipv6wry.ZXwry{}
	_ DB = &ipip.IPIPFree{}
	_ DB = &geoip.GeoIP{}
	_ DB = &ip2region.Ip2Region{}
	_ DB = &ip2locationdb.IP2LocationDB{}
	_ DB = &cdn.CDN{}
)
