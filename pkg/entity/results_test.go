package entity

import (
	"encoding/json"
	"testing"

	"github.com/zu1k/nali/internal/db"
	"github.com/zu1k/nali/pkg/wry"
)

func found(source, text string) *db.Result {
	return &db.Result{Source: source, Result: wry.Result{Country: text}}
}

func missed(source string) *db.Result {
	return &db.Result{Source: source}
}

func TestSetResults(t *testing.T) {
	tests := []struct {
		name        string
		results     []*db.Result
		wantOK      bool
		wantSource  string
		wantText    string
		wantResults []Result
	}{
		{
			name:       "single database",
			results:    []*db.Result{found("qqwry", "美国")},
			wantOK:     true,
			wantSource: "qqwry", wantText: "美国",
		},
		{
			name:    "nothing found",
			results: nil,
		},
		{
			// the first database failed: its entry stays in results, and the
			// top-level fields come from the database that answered
			name:       "first database failed",
			results:    []*db.Result{missed("qqwry"), found("geoip", "United States")},
			wantOK:     true,
			wantSource: "geoip", wantText: "United States",
			wantResults: []Result{
				{Source: "qqwry", Text: ""},
				{Source: "geoip", Text: "United States"},
			},
		},
		{
			// an answer with empty text is skipped for the top-level fields
			name:       "first answer empty",
			results:    []*db.Result{found("ipinfo", ""), found("geoip", "United States")},
			wantOK:     true,
			wantSource: "geoip", wantText: "United States",
			wantResults: []Result{
				{Source: "ipinfo", Text: ""},
				{Source: "geoip", Text: "United States"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var e Entity
			if ok := e.setResults(tt.results); ok != tt.wantOK {
				t.Fatalf("setResults = %v, want %v", ok, tt.wantOK)
			}
			if e.Source != tt.wantSource || e.InfoText != tt.wantText {
				t.Errorf("top level: source=%q text=%q, want %q %q", e.Source, e.InfoText, tt.wantSource, tt.wantText)
			}
			if len(e.Results) != len(tt.wantResults) {
				t.Fatalf("got %d results, want %d", len(e.Results), len(tt.wantResults))
			}
			for i, want := range tt.wantResults {
				got := e.Results[i]
				if got.Source != want.Source || got.Text != want.Text {
					t.Errorf("results[%d] = %+v, want %+v", i, got, want)
				}
				if (got.Info == nil) != (tt.results[i].Result == nil) {
					t.Errorf("results[%d].Info = %v", i, got.Info)
				}
			}
		})
	}
}

func TestSetResultsJSONForFailedLookup(t *testing.T) {
	e := Entity{Type: TypeIPv4, Text: "8.8.8.8"}
	e.setResults([]*db.Result{missed("qqwry"), found("geoip", "United States")})

	var got struct {
		Source  string `json:"source"`
		Results []struct {
			Source string `json:"source"`
			Text   string `json:"text"`
			Info   any    `json:"info"`
		} `json:"results"`
	}
	if err := json.Unmarshal([]byte(e.Json()), &got); err != nil {
		t.Fatal(err)
	}
	if got.Source != "geoip" || len(got.Results) != 2 || got.Results[0].Source != "qqwry" || got.Results[0].Info != nil {
		t.Fatalf("unexpected JSON %s", e.Json())
	}
}
