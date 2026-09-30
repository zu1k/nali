package ip2region

import (
	"encoding/binary"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

type segment struct {
	start, end string
	region     string
}

// buildXDB builds a minimal xdb database: header, a vector index whose every
// slot points at the full segment index, the region data and the segment index.
func buildXDB(t *testing.T, structure uint16, ipVersion int, segments []segment) []byte {
	t.Helper()
	ipBytes, segSize := 4, xdb.IPv4.SegmentIndexSize
	if ipVersion == 6 {
		ipBytes, segSize = 16, xdb.IPv6.SegmentIndexSize
	}

	vectorLen := xdb.VectorIndexRows * xdb.VectorIndexCols * xdb.VectorIndexSize
	buf := make([]byte, xdb.HeaderInfoLength+vectorLen)

	regionPtr := make([]uint32, len(segments))
	for i, s := range segments {
		regionPtr[i] = uint32(len(buf))
		buf = append(buf, s.region...)
	}

	encodeIP := func(s string) []byte {
		ip := net.ParseIP(s)
		if ipVersion == 4 {
			// IPv4 addresses are stored little endian
			b := slices.Clone(ip.To4())
			slices.Reverse(b)
			return b
		}
		return ip.To16()
	}

	startPtr := uint32(len(buf))
	for i, s := range segments {
		entry := make([]byte, segSize)
		copy(entry, encodeIP(s.start))
		copy(entry[ipBytes:], encodeIP(s.end))
		binary.LittleEndian.PutUint16(entry[2*ipBytes:], uint16(len(s.region)))
		binary.LittleEndian.PutUint32(entry[2*ipBytes+2:], regionPtr[i])
		buf = append(buf, entry...)
	}
	endPtr := startPtr + uint32((len(segments)-1)*segSize)

	for i := 0; i < xdb.VectorIndexRows*xdb.VectorIndexCols; i++ {
		off := xdb.HeaderInfoLength + i*xdb.VectorIndexSize
		binary.LittleEndian.PutUint32(buf[off:], startPtr)
		binary.LittleEndian.PutUint32(buf[off+4:], endPtr)
	}

	binary.LittleEndian.PutUint16(buf[0:], structure)
	binary.LittleEndian.PutUint16(buf[2:], uint16(xdb.VectorIndexPolicy))
	binary.LittleEndian.PutUint32(buf[8:], startPtr)
	binary.LittleEndian.PutUint32(buf[12:], endPtr)
	if structure == xdb.Structure30 {
		binary.LittleEndian.PutUint16(buf[16:], uint16(ipVersion))
		binary.LittleEndian.PutUint16(buf[18:], 4)
	}
	return buf
}

func writeXDB(t *testing.T, data []byte) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.xdb")
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

var v4Segments = []segment{
	{"0.0.0.0", "1.255.255.255", "0|0|0|内网IP|内网IP"},
	{"2.0.0.0", "223.5.5.4", "中国|0|浙江省|杭州市|阿里"},
	{"223.5.5.5", "255.255.255.255", "中国|0|0|0|0"},
}

var v6Segments = []segment{
	{"::", "240d:ffff:ffff:ffff:ffff:ffff:ffff:ffff", "0|0|0|0|0"},
	{"240e::", "240e:ffff:ffff:ffff:ffff:ffff:ffff:ffff", "中国|北京市|北京市|电信|CN"},
	{"240f::", "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff", "0|0|0|0|0"},
}

func TestFind(t *testing.T) {
	tests := []struct {
		name      string
		structure uint16
		ipVersion int
		segments  []segment
		query     string
		want      string
		wantErr   bool
	}{
		{"legacy v2", xdb.Structure20, 4, v4Segments, "3.3.3.3", "中国|浙江省|杭州市|阿里", false},
		{"v3 ipv4", xdb.Structure30, 4, v4Segments, "223.5.5.5", "中国", false},
		{"v3 ipv4 first segment", xdb.Structure30, 4, v4Segments, "1.1.1.1", "内网IP|内网IP", false},
		{"v3 ipv6", xdb.Structure30, 6, v6Segments, "240e::1", "中国|北京市|北京市|电信|CN", false},
		{"ipv6 query on ipv4 db", xdb.Structure30, 4, v4Segments, "240e::1", "", true},
		{"ipv4 query on ipv6 db", xdb.Structure30, 6, v6Segments, "1.1.1.1", "", true},
		{"invalid query", xdb.Structure30, 4, v4Segments, "not-an-ip", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := NewIp2Region(writeXDB(t, buildXDB(t, tt.structure, tt.ipVersion, tt.segments)), nil)
			if err != nil {
				t.Fatal(err)
			}
			got, err := db.Find(tt.query)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got.String() != tt.want {
				t.Fatalf("got %q, want %q", got.String(), tt.want)
			}
		})
	}
}

func TestCheckFile(t *testing.T) {
	for _, data := range [][]byte{
		buildXDB(t, xdb.Structure20, 4, v4Segments),
		buildXDB(t, xdb.Structure30, 4, v4Segments),
		buildXDB(t, xdb.Structure30, 6, v6Segments),
	} {
		if !CheckFile(data) {
			t.Error("valid xdb rejected")
		}
		if CheckFile(data[:len(data)-10]) {
			t.Error("truncated xdb accepted")
		}
	}

	badVersion := buildXDB(t, xdb.Structure30, 4, v4Segments)
	binary.LittleEndian.PutUint16(badVersion[0:], 9)
	for _, data := range [][]byte{nil, []byte("<html>404</html>"), badVersion} {
		if CheckFile(data) {
			t.Errorf("invalid xdb accepted: %.20q", data)
		}
	}
}

func TestNewIp2RegionDownloadsMissingFile(t *testing.T) {
	data := buildXDB(t, xdb.Structure30, 6, v6Segments)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(data)
	}))
	defer srv.Close()

	path := filepath.Join(t.TempDir(), "ip2region_v6.xdb")
	db, err := NewIp2Region(path, []string{srv.URL})
	if err != nil {
		t.Fatal(err)
	}
	if got, err := db.Find("240e::1"); err != nil || got.String() != "中国|北京市|北京市|电信|CN" {
		t.Fatalf("got %v, %v", got, err)
	}
}

func TestNewIp2RegionInvalidFile(t *testing.T) {
	if _, err := NewIp2Region(writeXDB(t, []byte("not an xdb file")), nil); err == nil {
		t.Fatal("invalid file should fail")
	}
}

func TestFormatRegion(t *testing.T) {
	for in, want := range map[string]string{
		"中国|0|浙江省|杭州市|阿里":         "中国|浙江省|杭州市|阿里",
		"0|0|0|内网IP|内网IP":         "内网IP|内网IP",
		"中国|北京市|北京市|电信|CN":        "中国|北京市|北京市|电信|CN",
		"美国|0|0|0|0":              "美国",
		"0|0|0|0|0":               "",
		"Reserved|10.0.0.0|0|0|0": "Reserved|10.0.0.0",
	} {
		if got := formatRegion(in); got != want {
			t.Errorf("formatRegion(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestCheckFileRejectsOutOfRangeIndex(t *testing.T) {
	data := buildXDB(t, xdb.Structure30, 4, v4Segments)
	// an end pointer past the end of the file, large enough to overflow a
	// 32-bit int
	binary.LittleEndian.PutUint32(data[12:], 0xFFFFFFF0)
	if CheckFile(data) {
		t.Fatal("out of range index accepted")
	}
}
