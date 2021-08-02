package entity

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	fmt.Println(ParseLine("2001:0db8:85a3:0000:0000:8a2e:0370:7334 baidu.com 1.2.3.4 baidu.com"))
	fmt.Println(ParseLine("a.cn.b.com.c.org d.com"))
}

func TestColorPrint(t *testing.T) {
	fmt.Println(ParseLine("2001:0db8:85a3:0000:0000:8a2e:0370:7334 baidu.com 1.2.3.4 baidu.com").ColorString())
}
