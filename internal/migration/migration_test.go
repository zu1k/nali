package migration

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/spf13/viper"

	"github.com/zu1k/nali/internal/constant"
	"github.com/zu1k/nali/internal/db"
	"github.com/zu1k/nali/pkg/ip2region"
	"github.com/zu1k/nali/pkg/qqwry"
)

// withConfig points the migrations at a temporary config directory holding
// the given config.yaml (or none, if content is empty).
func withConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	if content != "" {
		if err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	old := constant.ConfigDirPath
	constant.ConfigDirPath = dir
	viper.Reset()
	t.Cleanup(func() {
		constant.ConfigDirPath = old
		viper.Reset()
	})
	return dir
}

func readDBs(t *testing.T, dir string) map[string]*db.DB {
	t.Helper()
	v := viper.New()
	v.SetConfigFile(filepath.Join(dir, "config.yaml"))
	if err := v.ReadInConfig(); err != nil {
		t.Fatal(err)
	}
	var list db.List
	if err := v.UnmarshalKey("databases", &list); err != nil {
		t.Fatal(err)
	}
	dbs := make(map[string]*db.DB)
	for _, d := range list {
		dbs[d.Name] = d
	}
	return dbs
}

const oldConfig = `databases:
- name: qqwry
  name-alias: [chunzhen]
  format: qqwry
  file: qqwry.dat
  languages: [zh-CN]
  types: [IPv4]
  download-urls:
  - https://gh-release.zu1k.com/HMBSbige/qqwry/qqwry.dat
- name: ip2region
  name-alias: [i2r]
  format: ip2region
  file: ip2region.xdb
  languages: [zh-CN]
  types: [IPv4]
  download-urls:
  - https://cdn.jsdelivr.net/gh/lionsoul2014/ip2region/data/ip2region.xdb
  - https://raw.githubusercontent.com/lionsoul2014/ip2region/master/data/ip2region.xdb
- name: cdn
  format: cdn-yml
  file: cdn.yml
  languages: [zh-CN]
  types: [CDN]
  download-urls:
  - https://example.com/my-cdn.yml
selected:
  ipv4: ip2region
  ipv6: zxipv6wry
`

func TestMigrationsUpdateStaleDownloadURLs(t *testing.T) {
	dir := withConfig(t, oldConfig)

	migration2v7()
	viper.Reset()
	migration2v8()

	dbs := readDBs(t, dir)
	if got := dbs["qqwry"].DownloadUrls; !slices.Equal(got, qqwry.DownloadUrls) {
		t.Errorf("qqwry urls not migrated: %v", got)
	}
	if got := dbs["ip2region"].DownloadUrls; !slices.Equal(got, ip2region.DownloadUrls) {
		t.Errorf("ip2region urls not migrated: %v", got)
	}
	// unrelated entries and settings are preserved
	if got := dbs["cdn"].DownloadUrls; !slices.Equal(got, []string{"https://example.com/my-cdn.yml"}) {
		t.Errorf("cdn urls changed: %v", got)
	}
	if got := dbs["ip2region"].NameAlias; !slices.Equal(got, []string{"i2r"}) {
		t.Errorf("ip2region alias changed: %v", got)
	}
	v := viper.New()
	v.SetConfigFile(filepath.Join(dir, "config.yaml"))
	if err := v.ReadInConfig(); err != nil {
		t.Fatal(err)
	}
	if got := v.GetString("selected.ipv4"); got != "ip2region" {
		t.Errorf("selected.ipv4 changed to %q", got)
	}
}

func TestMigrationsKeepCustomDownloadURLs(t *testing.T) {
	const custom = `databases:
- name: qqwry
  format: qqwry
  file: qqwry.dat
  download-urls: [https://mirror.example.com/qqwry.dat]
- name: ip2region
  format: ip2region
  file: ip2region.xdb
  download-urls: [https://mirror.example.com/ip2region_v4.xdb]
`
	dir := withConfig(t, custom)
	before, err := os.ReadFile(filepath.Join(dir, "config.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	migration2v7()
	viper.Reset()
	migration2v8()

	after, err := os.ReadFile(filepath.Join(dir, "config.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Fatalf("config rewritten although nothing needed migrating:\n%s", after)
	}
}

func TestMigrationsWithoutConfig(t *testing.T) {
	dir := withConfig(t, "")

	migration2v7()
	migration2v8()

	if _, err := os.Stat(filepath.Join(dir, "config.yaml")); !os.IsNotExist(err) {
		t.Fatalf("migrations must not create a config file: %v", err)
	}
}
