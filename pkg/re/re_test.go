package re

import (
	"fmt"
	"testing"
)

var domainList = []string{
	"a.a.qiniudns.com",
	"a.com.qiniudns.com",
	"a.com.cn.qiniudns.com",
	"看这里：a.com.cn.qiniudns.com行不行",
}

func TestDomainRe(t *testing.T) {
	for _, domain := range domainList {
		if !DomainRe.MatchString(domain) {
			t.Error(domain)
			t.Fail()
		}
		fmt.Println(DomainRe.FindAllString(domain, -1))
	}
}

var validIPv6List = []string{
	"::ffff:104.26.11.119",
}

func TestIPv6Re(t *testing.T) {
	for _, ip := range validIPv6List {
		if !IPv6Re.MatchString(ip) {
			t.Error(ip)
			t.Fail()
		}
		fmt.Println(IPv6Re.FindAllString(ip, -1))
	}
}
