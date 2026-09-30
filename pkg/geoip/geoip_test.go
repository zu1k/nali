package geoip

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"

	"github.com/zu1k/nali/pkg/ipinfo"
)

const (
	geoliteFixture = "../../testdata/mmdb/geolite2-city-test.mmdb"
	ipinfoFixture  = "../../testdata/mmdb/ipinfo-lite-test.mmdb"
)

func setLang(t *testing.T, lang string) {
	t.Helper()
	old := viper.GetString("selected.lang")
	viper.Set("selected.lang", lang)
	t.Cleanup(func() { viper.Set("selected.lang", old) })
}

func TestGeoLite2Find(t *testing.T) {
	db, err := NewGeoIP(geoliteFixture, nil)
	if err != nil {
		t.Fatal(err)
	}
	if db.Name() != "geoip" || db.ipinfo != nil {
		t.Fatalf("GeoLite2 database detected as %q", db.Name())
	}

	tests := []struct {
		lang, query string
		want        Result
	}{
		{"zh-CN", "81.2.69.160", Result{Country: "英国", CountryCode: "GB", Area: "伦敦"}},
		{"en", "81.2.69.160", Result{Country: "United Kingdom", CountryCode: "GB", Area: "London"}},
		{"", "81.2.69.160", Result{Country: "英国", CountryCode: "GB", Area: "伦敦"}},
		// no zh-CN / ja names in the record: fall back to English
		{"zh-CN", "2001:218::1", Result{Country: "Japan", CountryCode: "JP"}},
		{"ja", "81.2.69.160", Result{Country: "United Kingdom", CountryCode: "GB", Area: "London"}},
		// not in the database
		{"en", "192.0.2.1", Result{}},
	}
	for _, tt := range tests {
		t.Run(tt.lang+"/"+tt.query, func(t *testing.T) {
			setLang(t, tt.lang)
			got, err := db.Find(tt.query)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Fatalf("got %+v, want %+v", got, tt.want)
			}
		})
	}

	if _, err := db.Find("not-an-ip"); err == nil {
		t.Error("invalid IP should fail")
	}
}

func TestResultString(t *testing.T) {
	if got := (Result{Country: "英国", Area: "伦敦"}).String(); got != "英国 伦敦" {
		t.Errorf("got %q", got)
	}
	if got := (Result{Country: "日本"}).String(); got != "日本" {
		t.Errorf("got %q", got)
	}
}

func TestIPinfoDetection(t *testing.T) {
	db, err := NewGeoIP(ipinfoFixture, nil)
	if err != nil {
		t.Fatal(err)
	}
	if db.Name() != "ipinfo" || db.ipinfo == nil {
		t.Fatalf("IPinfo database detected as %q", db.Name())
	}
	got, err := db.Find("8.8.8.8")
	if err != nil {
		t.Fatal(err)
	}
	if r, ok := got.(ipinfo.Result); !ok || r.ASN != "AS15169" {
		t.Fatalf("unexpected result %#v", got)
	}
}

func TestCheckFile(t *testing.T) {
	for _, path := range []string{geoliteFixture, ipinfoFixture} {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if !CheckFile(data) {
			t.Errorf("%s should be valid", path)
		}
		if CheckFile(data[:len(data)/2]) {
			t.Errorf("truncated %s should be invalid", path)
		}
	}
	for _, data := range [][]byte{nil, []byte("<html>404</html>")} {
		if CheckFile(data) {
			t.Errorf("%q should be invalid", data)
		}
	}
}

func TestNewGeoIPMissingFile(t *testing.T) {
	if _, err := NewGeoIP(filepath.Join(t.TempDir(), "missing.mmdb"), nil); err == nil {
		t.Fatal("missing file without download URLs should fail")
	}
}

func TestNewGeoIPDownloadsMissingFile(t *testing.T) {
	fixture, err := os.ReadFile(ipinfoFixture)
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(fixture)
	}))
	defer srv.Close()

	path := filepath.Join(t.TempDir(), "ipinfo_lite.mmdb")
	db, err := NewGeoIP(path, []string{srv.URL})
	if err != nil {
		t.Fatal(err)
	}
	if db.Name() != "ipinfo" {
		t.Fatalf("downloaded database detected as %q", db.Name())
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("downloaded database was not saved: %v", err)
	}
}

func TestNewGeoIPRejectsInvalidDownload(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("<html>rate limited</html>"))
	}))
	defer srv.Close()

	path := filepath.Join(t.TempDir(), "ipinfo_lite.mmdb")
	if _, err := NewGeoIP(path, []string{srv.URL}); err == nil {
		t.Fatal("invalid download should fail")
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("invalid download should not be kept: %v", err)
	}
}
