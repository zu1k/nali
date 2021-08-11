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
