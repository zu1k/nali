package entity

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/spf13/viper"

	"github.com/zu1k/nali/internal/db"
)

const cdnYAML = `
cloudflare.net:
  name: Cloudflare
  link: https://www.cloudflare.com
`

// setupDBs wires the IPv4/IPv6/CDN lookups to small local fixtures. The db
// package caches opened databases, so every test in this package shares them.
func setupDBs(t *testing.T) {
	t.Helper()
	mmdb, err := filepath.Abs("../../testdata/mmdb/ipinfo-lite-test.mmdb")
	if err != nil {
		t.Fatal(err)
	}
	cdnFile := filepath.Join(t.TempDir(), "cdn.yml")
	if err := os.WriteFile(cdnFile, []byte(cdnYAML), 0644); err != nil {
		t.Fatal(err)
	}
	db.NameDBMap.From(db.List{
		{Name: "test-ipinfo", Format: db.FormatMMDB, File: mmdb, Types: db.TypesIP},
		{Name: "test-cdn", Format: db.FormatCDNYml, File: cdnFile, Types: db.TypesCDN},
	})
	viper.Set("selected.ipv4", "test-ipinfo")
	viper.Set("selected.ipv6", "test-ipinfo")
	viper.Set("selected.cdn", "test-cdn")
}

func TestParseLine(t *testing.T) {
	setupDBs(t)

	line := "dns 8.8.8.8, v6 2001:4860:4860::8888 nat64 64:ff9b::808:808 cdn a.b.cloudflare.net unknown 192.0.2.1 end"
	es := ParseLine(line)

	// entities must cover the whole line, in order, without gaps
	var rebuilt strings.Builder
	for _, e := range es {
		rebuilt.WriteString(e.Text)
	}
	if rebuilt.String() != line {
		t.Fatalf("entities do not cover the line:\n got %q\nwant %q", rebuilt.String(), line)
	}

	want := map[string]struct {
		typ    EntityType
		info   string
		source string
	}{
		"8.8.8.8":              {TypeIPv4, "United States AS15169 Google LLC", "ipinfo"},
		"2001:4860:4860::8888": {TypeIPv6, "United States AS15169 Google LLC", "ipinfo"},
		"64:ff9b::808:808":     {TypeIPv6, "United States AS15169 Google LLC", "ipinfo"},
		"a.b.cloudflare.net":   {TypeDomain, "Cloudflare", "cdn"},
		// found nothing in the database: still an IPv4 entity, with empty info
		"192.0.2.1": {TypeIPv4, "", "ipinfo"},
	}
	found := 0
	for _, e := range es {
		w, ok := want[e.Text]
		if !ok {
			if e.Type != TypePlain {
				t.Errorf("unexpected entity %q of type %d", e.Text, e.Type)
			}
			continue
		}
		found++
		if e.Type != w.typ || e.InfoText != w.info || e.Source != w.source {
			t.Errorf("%s: got type=%d info=%q source=%q, want type=%d info=%q source=%q",
				e.Text, e.Type, e.InfoText, e.Source, w.typ, w.info, w.source)
		}
	}
	if found != len(want) {
		t.Fatalf("found %d of %d expected entities in %q", found, len(want), es.String())
	}

	wantText := "dns 8.8.8.8[United States AS15169 Google LLC] , v6 2001:4860:4860::8888[United States AS15169 Google LLC]  " +
		"nat64 64:ff9b::808:808[United States AS15169 Google LLC]  cdn a.b.cloudflare.net[Cloudflare]  unknown 192.0.2.1 end"
	if got := es.String(); got != wantText {
		t.Errorf("String():\n got %q\nwant %q", got, wantText)
	}
}

func TestJSON(t *testing.T) {
	setupDBs(t)

	es := ParseLine("8.8.8.8")
	if len(es) != 1 {
		t.Fatalf("got %d entities", len(es))
	}
	var got map[string]any
	if err := json.Unmarshal([]byte(es[0].Json()), &got); err != nil {
		t.Fatal(err)
	}
	if got["ip"] != "8.8.8.8" || got["source"] != "ipinfo" || got["text"] != "United States AS15169 Google LLC" {
		t.Fatalf("unexpected JSON %v", got)
	}
	info, ok := got["info"].(map[string]any)
	if !ok || info["asn"] != "AS15169" || info["network"] != "8.8.8.0/24" {
		t.Fatalf("unexpected info %v", got["info"])
	}
}

func TestParseLinePlainText(t *testing.T) {
	es := ParseLine("no addresses here")
	if len(es) != 1 || es[0].Type != TypePlain || es[0].Text != "no addresses here" {
		t.Fatalf("unexpected entities %+v", es)
	}
	if es.String() != "no addresses here" {
		t.Fatalf("String() = %q", es.String())
	}
}

func TestRenderMultipleResults(t *testing.T) {
	old := color.NoColor
	color.NoColor = true
	t.Cleanup(func() { color.NoColor = old })

	es := Entities{
		{Type: TypePlain, Text: "dns "},
		{
			Type: TypeIPv4, Text: "8.8.8.8",
			InfoText: "美国 谷歌公司DNS服务器", Source: "qqwry", Info: "a",
			Results: []Result{
				{Source: "qqwry", Text: "美国 谷歌公司DNS服务器", Info: "a"},
				{Source: "ipinfo", Text: "", Info: "b"}, // empty texts are not shown
				{Source: "geoip", Text: "United States", Info: "c"},
			},
		},
		{Type: TypePlain, Text: "end"},
	}

	if got, want := es.String(), "dns 8.8.8.8[美国 谷歌公司DNS服务器] [United States] end"; got != want {
		t.Errorf("String():\n got %q\nwant %q", got, want)
	}
	if got, want := es.ColorString(), "dns 8.8.8.8 [美国 谷歌公司DNS服务器] [United States] end"; got != want {
		t.Errorf("ColorString():\n got %q\nwant %q", got, want)
	}

	var got map[string]any
	if err := json.Unmarshal([]byte(es[1].Json()), &got); err != nil {
		t.Fatal(err)
	}
	results, ok := got["results"].([]any)
	if got["source"] != "qqwry" || got["text"] != "美国 谷歌公司DNS服务器" || !ok || len(results) != 3 {
		t.Fatalf("unexpected JSON %v", got)
	}

	// a single result keeps the original JSON shape
	single := Entity{Type: TypeIPv4, Text: "1.1.1.1", InfoText: "x", Source: "qqwry"}
	if strings.Contains(single.Json(), "results") {
		t.Fatalf("single result JSON must not contain results: %s", single.Json())
	}
}
