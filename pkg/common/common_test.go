package common

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
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
