package tools

import (
	"fmt"
	"testing"
)

func TestIP4Re(t *testing.T) {
	str := "aaa1.1.11.23a36.36.32.200"
	fmt.Println(GetIP4FromString(str))
	fmt.Println(ValidIP4(str))
}

func TestValidIP6(t *testing.T) {
	ipv6Valid := []string{
		"1:2:3:4:5:6:7::",
		"1:2:3:4:5:6:7:8",

		"1:2:3:4:5:6::",
		"1:2:3:4:5:6::8",

		"1:2:3:4:5::",
		"1:2:3:4:5::8",

		"1:2:3:4::",
		"1:2:3:4::8",

		"1:2:3::",
		"1:2:3::8",

		"1:2::",
		"1:2::8",

		"1::",
		"1::8",

		"::",
		"::8",
		"::7:8",
		"::6:7:8",
		"::5:6:7:8",
		"::4:5:6:7:8",
		"::3:4:5:6:7:8",
		"::2:3:4:5:6:7:8",

		"::192.168.1.1",
		"::ffff:135.75.43.52",
		"A:0f:0F:FFFF:5:6:7:8",
	}

	ipv6Invalid := []string{
		"A:0f:0F:FFFF1:5:6:7:8",
		"G:0f:0F:FFFF:5:6:7:8",
		"2001::25de::cade",
	}

	for _, i := range ipv6Valid {
		if !ValidIP6(i) {
			t.Log("valid:", i)
		}
	}

	for _, i := range ipv6Invalid {
		if ValidIP6(i) {
			t.Log("invalid:", i)
		}
	}
}
