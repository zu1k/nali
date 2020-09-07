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
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334",

		"::192.168.1.1",
		"::ffff:135.75.43.52",
		"A:0f:0F:FFFF:5:6:7:8",
	}

	ipv6Invalid := []string{
		"A:0f:0F:FFFF1:5:6:7:8",
		"G:0f:0F:FFFF:5:6:7:8",
		"2001::25de::cade",
		"2001:0db8:85a3:0:0:8A2E:0370:73341",
		"a1:a2:a3:a4::b1:b2:b3:b4",
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

func BenchmarkValidIP6Re(b *testing.B) {
	b.ResetTimer()
	origin := "::ffff:135.75.43.52"
	for i:=0; i<b.N;i++ {
		ValidIP6Re(origin)
	}
}

func BenchmarkValidIP6(b *testing.B) {
	b.ResetTimer()
	origin := "::ffff:135.75.43.52"
	for i:=0; i<b.N;i++ {
		ValidIP6(origin)
	}
}

/*

   IPv6-addr      = IPv6-full / IPv6-comp / IPv6v4-full / IPv6v4-comp

   IPv6-hex       = 1*4HEXDIG

   IPv6-full      = IPv6-hex 7(":" IPv6-hex)

   IPv6-comp      = [IPv6-hex *5(":" IPv6-hex)] "::"
                  [IPv6-hex *5(":" IPv6-hex)]
                  ; The "::" represents at least 2 16-bit groups of
                  ; zeros.  No more than 6 groups in addition to the
                  ; "::" may be present.

   IPv6v4-full    = IPv6-hex 5(":" IPv6-hex) ":" IPv4-address-literal

   IPv6v4-comp    = [IPv6-hex *3(":" IPv6-hex)] "::"
                  [IPv6-hex *3(":" IPv6-hex) ":"]
                  IPv4-address-literal
                  ; The "::" represents at least 2 16-bit groups of
                  ; zeros.  No more than 4 groups in addition to the
                  ; "::" and IPv4-address-literal may be present.
*/