package db

import (
	"path/filepath"
	"sync"
	"testing"

	"github.com/zu1k/nali/pkg/dbif"
)

func resetCaches(t *testing.T) {
	t.Helper()
	reset := func() {
		dbNameCache = make(map[string]dbif.DB)
		dbTypeCache = make(map[dbif.QueryType][]dbif.DB)
		queryCache = sync.Map{}
	}
	reset()
	t.Cleanup(reset)
}

func fixture(t *testing.T, name string) string {
	t.Helper()
	path, err := filepath.Abs(filepath.Join("..", "..", "testdata", "mmdb", name))
	if err != nil {
		t.Fatal(err)
	}
	return path
}

func TestFindMultipleDatabases(t *testing.T) {
	resetCaches(t)
	withNameDBMap(t,
		&DB{Name: "test-ipinfo", Format: FormatMMDB, File: fixture(t, "ipinfo-lite-test.mmdb"), Types: TypesIP},
		&DB{Name: "test-geolite", Format: FormatMMDB, File: fixture(t, "geolite2-city-test.mmdb"), Types: TypesIP},
	)
	withViper(t, map[string]string{
		"selected.ipv4": "test-ipinfo, test-geolite",
		"selected.ipv6": "test-geolite",
		"selected.lang": "en",
	})

	results := Find(dbif.TypeIPv4, "81.2.69.160")
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}
	if results[0].Source != "ipinfo" || results[0].String() != "" {
		t.Errorf("first result: %s %q", results[0].Source, results[0].String())
	}
	if results[1].Source != "geoip" || results[1].String() != "United Kingdom London" {
		t.Errorf("second result: %s %q", results[1].Source, results[1].String())
	}

	// a single selected database behaves as before
	results = Find(dbif.TypeIPv6, "2001:218::1")
	if len(results) != 1 || results[0].String() != "Japan" {
		t.Fatalf("unexpected IPv6 results %+v", results)
	}

	// failing lookups (here: an invalid address) yield no results
	if results := Find(dbif.TypeIPv4, "not-an-ip"); results != nil {
		t.Fatalf("expected no results, got %+v", results)
	}
}

func TestDefaultUpdateListHonorsMultipleSelection(t *testing.T) {
	withNameDBMap(t)
	withViper(t, map[string]string{"selected.ipv4": "qqwry, ipinfo", "selected.ipv6": "zx", "selected.cdn": "cdn"})

	found := false
	for _, name := range defaultUpdateList() {
		if name == "ipinfo" {
			found = true
		}
	}
	if !found {
		t.Fatal("ipinfo selected as second IPv4 database should be updated")
	}
}
