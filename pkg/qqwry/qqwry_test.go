package qqwry

import (
	"encoding/binary"
	"net"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/text/encoding/simplifiedchinese"
)

func gbk(t *testing.T, s string) []byte {
	t.Helper()
	b, err := simplifiedchinese.GBK.NewEncoder().Bytes([]byte(s))
	if err != nil {
		t.Fatal(err)
	}
	return append(b, 0)
}

func ip4(s string) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, binary.BigEndian.Uint32(net.ParseIP(s).To4()))
	return b
}

func off3(v int) []byte { return []byte{byte(v), byte(v >> 8), byte(v >> 16)} }

// buildQQwry builds a small qqwry.dat exercising the plain record layout and
// both redirect modes used by the real database.
func buildQQwry(t *testing.T) []byte {
	t.Helper()
	buf := make([]byte, 8)
	cat := func(parts ...[]byte) {
		for _, p := range parts {
			buf = append(buf, p...)
		}
	}

	// shared strings targeted by redirects
	usOff := len(buf)
	cat(gbk(t, "美国"), gbk(t, "加利福尼亚州"))
	cnOff := len(buf)
	cat(gbk(t, "中国"))
	areaOff := len(buf)
	cat(gbk(t, "东京都"))

	// records: [end ip][...]
	recA := len(buf) // plain: [country][area]
	cat(ip4("1.255.255.255"), gbk(t, "本机地址"), gbk(t, " CZ88.NET"))
	recB := len(buf) // mode 1: country and area both redirected
	cat(ip4("2.255.255.255"), []byte{0x01}, off3(usOff))
	recC := len(buf) // mode 2: country redirected, area inline
	cat(ip4("3.255.255.255"), []byte{0x02}, off3(cnOff), gbk(t, "浙江省杭州市"))
	recD := len(buf) // inline country, area redirected
	cat(ip4("255.255.255.255"), gbk(t, "日本"), []byte{0x02}, off3(areaOff))

	// index: [start ip][record offset]
	idxStart := len(buf)
	for _, e := range []struct {
		ip  string
		rec int
	}{{"0.0.0.0", recA}, {"2.0.0.0", recB}, {"3.0.0.0", recC}, {"4.0.0.0", recD}} {
		cat(ip4(e.ip), off3(e.rec))
	}
	idxEnd := len(buf) - 7

	binary.LittleEndian.PutUint32(buf[0:], uint32(idxStart))
	binary.LittleEndian.PutUint32(buf[4:], uint32(idxEnd))
	return buf
}

func TestFind(t *testing.T) {
	path := filepath.Join(t.TempDir(), "qqwry.dat")
	if err := os.WriteFile(path, buildQQwry(t), 0644); err != nil {
		t.Fatal(err)
	}
	db, err := NewQQwry(path, nil)
	if err != nil {
		t.Fatal(err)
	}

	tests := map[string]string{
		"0.0.0.0":         "本机地址",
		"1.2.3.4":         "本机地址",
		"2.0.0.0":         "美国 加利福尼亚州",
		"2.200.1.1":       "美国 加利福尼亚州",
		"3.3.3.3":         "中国 浙江省杭州市",
		"8.8.8.8":         "日本 东京都",
		"255.255.255.255": "日本 东京都",
		"::ffff:3.3.3.3":  "中国 浙江省杭州市",
	}
	for query, want := range tests {
		got, err := db.Find(query)
		if err != nil {
			t.Errorf("%s: %v", query, err)
			continue
		}
		if got.String() != want {
			t.Errorf("%s: got %q, want %q", query, got.String(), want)
		}
	}

	for _, query := range []string{"not-an-ip", "2001:db8::1"} {
		if _, err := db.Find(query); err == nil {
			t.Errorf("%s should fail", query)
		}
	}
}

func TestCheckFile(t *testing.T) {
	data := buildQQwry(t)
	if !CheckFile(data) {
		t.Fatal("valid database rejected")
	}
	for _, bad := range [][]byte{nil, []byte("<html>404</html>"), data[:len(data)-1]} {
		if CheckFile(bad) {
			t.Errorf("invalid database accepted: %.20q", bad)
		}
	}
}
