package db

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultUpdateIncludesIPinfo(t *testing.T) {
	original := NameDBMap
	NameDBMap = make(NameMap)
	t.Cleanup(func() { NameDBMap = original })
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, r.URL.Path)
	}))
	defer server.Close()
	dir := t.TempDir()
	// Use multiple MMDB entries to ensure updates are deduplicated by database name.
	for _, name := range []string{"qqwry", "zxipv6wry", "ip2region", "cdn", "ipinfo"} {
		NameDBMap[name] = &DB{Name: name, Format: FormatMMDB, File: filepath.Join(dir, name), DownloadUrls: []string{server.URL + "/" + name}}
	}
	UpdateDB()
	for _, name := range []string{"qqwry", "zxipv6wry", "ip2region", "cdn", "ipinfo"} {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Errorf("%s was not downloaded: %v", name, err)
			continue
		}
		if string(data) != "/"+name {
			t.Errorf("unexpected data for %s: %q", name, data)
		}
	}
}

func TestIPinfoDefaultDownloadURL(t *testing.T) {
	defaults := NameMap{}
	defaults.From(GetDefaultDBList())
	info := defaults["ipinfo"]
	const want = "https://github.com/NetworkCats/IPinfoLite-Download/releases/latest/download/ipinfo_lite.mmdb"
	if info == nil || len(info.DownloadUrls) != 1 || info.DownloadUrls[0] != want {
		t.Fatalf("incorrect IPinfo defaults: %+v", info)
	}
}
