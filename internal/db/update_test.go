package db

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/spf13/viper"
)

// formatRaw has no DbCheckFunc, so any downloaded content is accepted.
const formatRaw Format = "raw"

type fakeServer struct {
	*httptest.Server
	mu   sync.Mutex
	hits map[string]int
}

func newFakeServer(t *testing.T, handler func(w http.ResponseWriter, r *http.Request)) *fakeServer {
	t.Helper()
	s := &fakeServer{hits: make(map[string]int)}
	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		s.hits[r.URL.Path]++
		s.mu.Unlock()
		handler(w, r)
	}))
	t.Cleanup(s.Close)
	return s
}

func (s *fakeServer) hitCount(path string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.hits[path]
}

// withNameDBMap replaces the configured databases for the duration of a test.
func withNameDBMap(t *testing.T, dbs ...*DB) {
	t.Helper()
	original := NameDBMap
	NameDBMap = make(NameMap)
	NameDBMap.From(dbs)
	t.Cleanup(func() { NameDBMap = original })
}

func withViper(t *testing.T, kv map[string]string) {
	t.Helper()
	for k, v := range kv {
		old := viper.GetString(k)
		viper.Set(k, v)
		t.Cleanup(func() { viper.Set(k, old) })
	}
}

// mockAllUpdatableDBs points every database nali may update by default at the
// fake server so no test ever reaches the real network.
func mockAllUpdatableDBs(t *testing.T, serverURL, dir string) {
	t.Helper()
	var dbs []*DB
	for _, name := range append(append([]string{}, DbNameListForUpdate...), DbNameListForUpdateIfUsed...) {
		dbs = append(dbs, &DB{
			Name:         name,
			Format:       formatRaw,
			File:         filepath.Join(dir, name),
			DownloadUrls: []string{serverURL + "/" + name},
		})
	}
	withNameDBMap(t, dbs...)
}

func downloaded(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	sort.Strings(names)
	return names
}

func TestUpdateDBDefaultListSkipsUnusedOptionalDBs(t *testing.T) {
	srv := newFakeServer(t, func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte(r.URL.Path)) })
	dir := t.TempDir()
	mockAllUpdatableDBs(t, srv.URL, dir)
	withViper(t, map[string]string{"selected.ipv4": "qqwry", "selected.ipv6": "zxipv6wry", "selected.cdn": "cdn"})

	UpdateDB()

	want := append([]string{}, DbNameListForUpdate...)
	sort.Strings(want)
	if got := downloaded(t, dir); !slices.Equal(got, want) {
		t.Fatalf("downloaded %v, want %v", got, want)
	}
	for _, name := range want {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil || string(data) != "/"+name {
			t.Errorf("unexpected content for %s: %q, %v", name, data, err)
		}
	}
}

func TestUpdateDBDefaultListIncludesSelectedOptionalDB(t *testing.T) {
	srv := newFakeServer(t, func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte("ok")) })
	dir := t.TempDir()
	mockAllUpdatableDBs(t, srv.URL, dir)
	withViper(t, map[string]string{"selected.ipv4": "ipinfo", "selected.ipv6": "zxipv6wry", "selected.cdn": "cdn"})

	UpdateDB()

	if srv.hitCount("/ipinfo") != 1 {
		t.Errorf("selected ipinfo should be updated")
	}
	if srv.hitCount("/ip2region-ipv6") != 0 {
		t.Errorf("unused ip2region-ipv6 should not be updated")
	}
}

func TestUpdateDBDefaultListIncludesDownloadedOptionalDB(t *testing.T) {
	srv := newFakeServer(t, func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte("new")) })
	dir := t.TempDir()
	mockAllUpdatableDBs(t, srv.URL, dir)
	withViper(t, map[string]string{"selected.ipv4": "qqwry", "selected.ipv6": "zxipv6wry", "selected.cdn": "cdn"})
	if err := os.WriteFile(filepath.Join(dir, "ip2region-ipv6"), []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}

	UpdateDB()

	data, _ := os.ReadFile(filepath.Join(dir, "ip2region-ipv6"))
	if string(data) != "new" {
		t.Errorf("already downloaded ip2region-ipv6 should be refreshed, got %q", data)
	}
	if srv.hitCount("/ipinfo") != 0 {
		t.Errorf("unused ipinfo should not be updated")
	}
}

func TestUpdateDBExplicitNameAlwaysUpdates(t *testing.T) {
	srv := newFakeServer(t, func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte("ok")) })
	dir := t.TempDir()
	mockAllUpdatableDBs(t, srv.URL, dir)

	UpdateDB("ipinfo")

	if got := downloaded(t, dir); !slices.Equal(got, []string{"ipinfo"}) {
		t.Fatalf("downloaded %v, want [ipinfo]", got)
	}
}

func TestUpdateDBDeduplicatesByName(t *testing.T) {
	srv := newFakeServer(t, func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte("ok")) })
	dir := t.TempDir()
	// Two databases with the same format must both be updated, while an
	// alias of an already updated database must not trigger a second download.
	withNameDBMap(t,
		&DB{Name: "a", NameAlias: []string{"a-alias"}, Format: formatRaw, File: filepath.Join(dir, "a"), DownloadUrls: []string{srv.URL + "/a"}},
		&DB{Name: "b", Format: formatRaw, File: filepath.Join(dir, "b"), DownloadUrls: []string{srv.URL + "/b"}},
	)

	UpdateDB("a", "b", "a-alias", " a ")

	if srv.hitCount("/a") != 1 || srv.hitCount("/b") != 1 {
		t.Fatalf("unexpected hits: a=%d b=%d", srv.hitCount("/a"), srv.hitCount("/b"))
	}
}

func TestUpdateDBKeepsExistingFileOnInvalidDownload(t *testing.T) {
	srv := newFakeServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("<html>rate limited</html>"))
	})
	dir := t.TempDir()
	for _, format := range []Format{FormatMMDB, FormatQQWry, FormatZXIPv6Wry, FormatIP2Region, FormatCDNYml} {
		t.Run(string(format), func(t *testing.T) {
			file := filepath.Join(dir, string(format))
			if err := os.WriteFile(file, []byte("good"), 0644); err != nil {
				t.Fatal(err)
			}
			withNameDBMap(t, &DB{Name: "x", Format: format, File: file, DownloadUrls: []string{srv.URL + "/x"}})

			update, _ := getUpdateFuncByName("x")
			if err := update(); err == nil {
				t.Fatal("invalid download should fail")
			}
			if data, _ := os.ReadFile(file); string(data) != "good" {
				t.Fatalf("existing database was overwritten with %q", data)
			}
		})
	}
}

func TestUpdateDBKeepsExistingFileOnHTTPError(t *testing.T) {
	srv := newFakeServer(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	})
	file := filepath.Join(t.TempDir(), "db")
	if err := os.WriteFile(file, []byte("good"), 0644); err != nil {
		t.Fatal(err)
	}
	withNameDBMap(t, &DB{Name: "x", Format: formatRaw, File: file, DownloadUrls: []string{srv.URL + "/x"}})

	update, _ := getUpdateFuncByName("x")
	if err := update(); err == nil {
		t.Fatal("HTTP 404 should fail")
	}
	if data, _ := os.ReadFile(file); string(data) != "good" {
		t.Fatalf("existing database was overwritten with %q", data)
	}
}

func TestUpdateListsResolveToDefaultDBs(t *testing.T) {
	withNameDBMap(t)
	for _, name := range append(append([]string{}, DbNameListForUpdate...), DbNameListForUpdateIfUsed...) {
		if _, found := lookupDb(name); !found {
			t.Errorf("update list entry %q has no default database", name)
		}
	}
}

func TestDefaultDBList(t *testing.T) {
	defaults := NameMap{}
	defaults.From(GetDefaultDBList())

	ipinfo := defaults["ipinfo"]
	const want = "https://github.com/NetworkCats/IPinfoLite-Download/releases/latest/download/ipinfo_lite.mmdb"
	if ipinfo == nil || ipinfo.Format != FormatMMDB || len(ipinfo.DownloadUrls) != 1 || ipinfo.DownloadUrls[0] != want {
		t.Fatalf("incorrect IPinfo defaults: %+v", ipinfo)
	}

	v6 := defaults["i2r-ipv6"]
	if v6 == nil || v6.Name != "ip2region-ipv6" || v6.File != "ip2region_v6.xdb" || len(v6.DownloadUrls) == 0 {
		t.Fatalf("incorrect ip2region-ipv6 defaults: %+v", v6)
	}

	for alias, name := range map[string]string{"chunzhen": "qqwry", "zx": "zxipv6wry", "geolite2": "geoip", "db-ip": "dbip", "i2r": "ip2region"} {
		if db := defaults[alias]; db == nil || db.Name != name {
			t.Errorf("alias %s should resolve to %s", alias, name)
		}
	}

	files := make(map[string]string)
	for _, db := range GetDefaultDBList() {
		if other, dup := files[db.File]; dup {
			t.Errorf("%s and %s share the file %s", db.Name, other, db.File)
		}
		files[db.File] = db.Name
	}
}

func TestUpdateDBSkipsUnchangedDatabase(t *testing.T) {
	const content = "database v1"

	tests := []struct {
		name       string
		body       string
		wantUpdate bool
	}{
		{name: "identical content", body: content},
		{name: "changed content", body: "database v2", wantUpdate: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := newFakeServer(t, func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(tt.body))
			})
			file := filepath.Join(t.TempDir(), "db")
			if err := os.WriteFile(file, []byte(content), 0644); err != nil {
				t.Fatal(err)
			}
			old := time.Now().Add(-time.Minute).Truncate(time.Second)
			if err := os.Chtimes(file, old, old); err != nil {
				t.Fatal(err)
			}
			withNameDBMap(t, &DB{Name: "x", Format: formatRaw, File: file, DownloadUrls: []string{srv.URL + "/x"}})

			update, _ := getUpdateFuncByName("x")
			if err := update(); err != nil {
				t.Fatal(err)
			}

			data, _ := os.ReadFile(file)
			info, _ := os.Stat(file)
			if tt.wantUpdate {
				if string(data) != tt.body {
					t.Fatalf("database not updated: %q", data)
				}
			} else if string(data) != content || !info.ModTime().Equal(old) {
				t.Fatalf("unchanged database was rewritten: %q, mtime %v", data, info.ModTime())
			}
		})
	}
}
