package cdn

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

const testYAML = `
cloudflare.net:
  name: Cloudflare
  link: https://www.cloudflare.com
akamaiedge.net:
  name: Akamai CDN
  link: https://www.akamai.com
kunlun[^.]+.com:
  name: 阿里云 CDN
  link: https://www.aliyun.com/product/cdn
"broken[regexp":
  name: Broken
  link: https://example.com
`

func writeYAML(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "cdn.yml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestFind(t *testing.T) {
	db, err := NewCDN(writeYAML(t, testYAML), nil)
	if err != nil {
		t.Fatal(err)
	}
	// the invalid regexp entry must be skipped, not stored as a nil *Regexp
	for _, entry := range db.ReMap {
		if entry.Regexp == nil {
			t.Fatal("nil regexp stored in ReMap")
		}
	}

	tests := []struct {
		query, want string
	}{
		{"abc.cdn.cloudflare.net", "Cloudflare"},
		{"cloudflare.net", "Cloudflare"},
		{"e1234.a.akamaiedge.net", "Akamai CDN"},
		{"img.kunlunsl.com", "阿里云 CDN"},
	}
	for _, tt := range tests {
		got, err := db.Find(tt.query)
		if err != nil {
			t.Errorf("%s: %v", tt.query, err)
			continue
		}
		if got.String() != tt.want {
			t.Errorf("%s: got %q, want %q", tt.query, got, tt.want)
		}
	}

	for _, query := range []string{"example.org", "notcloudflare.net.example", ""} {
		if got, err := db.Find(query); err == nil {
			t.Errorf("%q should not match, got %q", query, got)
		}
	}
}

func TestParseBaseCname(t *testing.T) {
	got := parseBaseCname("a.b.example.com")
	want := []string{"com", "example.com", "b.example.com", "a.b.example.com"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestCheckFile(t *testing.T) {
	if !CheckFile([]byte(testYAML)) {
		t.Error("valid yaml rejected")
	}
	for _, data := range []string{
		"", "<html>404</html>", "{}", "- a\n- b\n",
		// entries that NewCDN would skip or that carry no name
		"\"broken[regexp\":\n  name: Broken\n",
		"example.com:\n  link: https://example.com\n",
	} {
		if CheckFile([]byte(data)) {
			t.Errorf("%q accepted", data)
		}
	}
}

func TestNewCDNDownloadsMissingFile(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(testYAML))
	}))
	defer srv.Close()

	path := filepath.Join(t.TempDir(), "cdn.yml")
	db, err := NewCDN(path, []string{srv.URL})
	if err != nil {
		t.Fatal(err)
	}
	if got, err := db.Find("x.cloudflare.net"); err != nil || got.String() != "Cloudflare" {
		t.Fatalf("got %v, %v", got, err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("downloaded file was not saved: %v", err)
	}
}
