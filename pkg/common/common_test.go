package common

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
)

func TestGetFallsBackToNextURL(t *testing.T) {
	var hits []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits = append(hits, r.URL.Path)
		if ua := r.Header.Get("User-Agent"); ua != UserAgent {
			t.Errorf("unexpected User-Agent %q", ua)
		}
		switch r.URL.Path {
		case "/ok":
			_, _ = w.Write([]byte("data"))
		default:
			http.Error(w, "boom", http.StatusInternalServerError)
		}
	}))
	defer srv.Close()

	body, err := GetHttpClient().Get(srv.URL+"/fail", "://bad-url", srv.URL+"/ok")
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "data" {
		t.Fatalf("got %q", body)
	}
	if len(hits) != 2 || hits[0] != "/fail" || hits[1] != "/ok" {
		t.Fatalf("unexpected requests %v", hits)
	}
}

func TestGetReturnsErrorWhenAllURLsFail(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer srv.Close()

	body, err := GetHttpClient().Get(srv.URL+"/a", srv.URL+"/b")
	if err == nil {
		t.Fatalf("expected an error, got body %q", body)
	}

	if _, err := GetHttpClient().Get(); err == nil {
		t.Fatal("no URLs should be an error")
	}
}

func TestSaveFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "db.dat")

	for _, content := range []string{"first", "second"} {
		if err := SaveFile(path, []byte(content)); err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile(path)
		if err != nil || string(data) != content {
			t.Fatalf("got %q, %v", data, err)
		}
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("temporary files left behind: %v", entries)
	}

	if err := SaveFile(filepath.Join(dir, "missing", "db.dat"), []byte("x")); err == nil {
		t.Fatal("saving into a missing directory should fail")
	}
}

func TestReadFileMapped(t *testing.T) {
	dir := t.TempDir()

	path := filepath.Join(dir, "db.dat")
	want := make([]byte, 3*4096+17)
	for i := range want {
		want[i] = byte(i * 7)
	}
	if err := os.WriteFile(path, want, 0644); err != nil {
		t.Fatal(err)
	}
	got, err := ReadFileMapped(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Fatal("mapped content differs from the file")
	}

	// replacing the file with SaveFile must not affect an existing mapping
	if err := SaveFile(path, []byte("new content")); err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Fatal("mapping changed after the file was replaced")
	}

	empty := filepath.Join(dir, "empty.dat")
	if err := os.WriteFile(empty, nil, 0644); err != nil {
		t.Fatal(err)
	}
	if got, err := ReadFileMapped(empty); err != nil || len(got) != 0 {
		t.Fatalf("empty file: got %d bytes, %v", len(got), err)
	}

	if _, err := ReadFileMapped(filepath.Join(dir, "missing.dat")); err == nil {
		t.Fatal("missing file should fail")
	}
}

func TestDecodeInput(t *testing.T) {
	gbk, err := simplifiedchinese.GBK.NewEncoder().String("跟踪路由 8.8.8.8 中国")
	if err != nil {
		t.Fatal(err)
	}
	const utf8Line = "跟踪路由 8.8.8.8 中国"

	tests := []struct {
		name              string
		in                string
		forceGBK, autoGBK bool
		want              string
	}{
		{"utf8 untouched", utf8Line, false, false, utf8Line},
		{"gbk untouched without auto detection", gbk, false, false, gbk},
		{"gbk decoded when auto detected", gbk, false, true, utf8Line},
		{"valid utf8 kept when auto detected", utf8Line, false, true, utf8Line},
		{"ascii kept when auto detected", "1.1.1.1\n", false, true, "1.1.1.1\n"},
		{"forced gbk", gbk, true, false, utf8Line},
	}
	for _, tt := range tests {
		if got := DecodeInput(tt.in, tt.forceGBK, tt.autoGBK); got != tt.want {
			t.Errorf("%s: got %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestGetIfModified(t *testing.T) {
	modTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeContent(w, r, "db", modTime, strings.NewReader("data"))
	}))
	defer srv.Close()

	body, notModified, err := GetHttpClient().GetIfModified(modTime.Add(time.Hour), srv.URL)
	if err != nil || !notModified || body != nil {
		t.Fatalf("newer local copy: got %q, %v, %v", body, notModified, err)
	}

	body, notModified, err = GetHttpClient().GetIfModified(modTime.Add(-time.Hour), srv.URL)
	if err != nil || notModified || string(body) != "data" {
		t.Fatalf("older local copy: got %q, %v, %v", body, notModified, err)
	}

	body, notModified, err = GetHttpClient().GetIfModified(time.Time{}, srv.URL)
	if err != nil || notModified || string(body) != "data" {
		t.Fatalf("no local copy: got %q, %v, %v", body, notModified, err)
	}
}
