package entity

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/fatih/color"
	"github.com/zu1k/nali/pkg/dbif"
)

type EntityType uint

const (
	TypeIPv4   = dbif.TypeIPv4
	TypeIPv6   = dbif.TypeIPv6
	TypeDomain = dbif.TypeDomain

	TypePlain = 100
)

type Entity struct {
	Loc  [2]int     `json:"-"` // s[Loc[0]:Loc[1]]
	Type EntityType `json:"type"`

	Text     string      `json:"ip"`
	InfoText string      `json:"text"`
	Source   string      `json:"source"`
	Info     interface{} `json:"info"`

	// Results holds every database's result when several databases are
	// selected (e.g. NALI_DB_IP4=qqwry,ipinfo); the fields above then carry
	// the first one.
	Results []Result `json:"results,omitempty"`
}

// Result is the answer of one database for an entity.
type Result struct {
	Source string      `json:"source"`
	Text   string      `json:"text"`
	Info   interface{} `json:"info"`
}

// infoTexts returns the non-empty result texts to show after the entity.
func (e *Entity) infoTexts() []string {
	if e.Type == TypePlain {
		return nil
	}
	if len(e.Results) == 0 {
		if e.InfoText == "" {
			return nil
		}
		return []string{e.InfoText}
	}
	var texts []string
	for _, r := range e.Results {
		if r.Text != "" {
			texts = append(texts, r.Text)
		}
	}
	return texts
}

func (e Entity) ParseInfo() error {
	return nil
}

func (e Entity) Json() string {
	jsonResult, err := json.Marshal(e)
	if err != nil {
		log.Fatal(err.Error())
	}
	return string(jsonResult)
}

type Entities []*Entity

func (es Entities) Len() int {
	return len(es)
}

func (es Entities) Less(i, j int) bool {
	return es[i].Loc[0] < es[j].Loc[0]
}

func (es Entities) Swap(i, j int) {
	es[i], es[j] = es[j], es[i]
}

func (es Entities) String() string {
	var result strings.Builder
	for _, entity := range es {
		result.WriteString(entity.Text)
		for _, text := range entity.infoTexts() {
			result.WriteString("[" + text + "] ")
		}
	}
	return result.String()
}

func (es Entities) ColorString() string {
	var line strings.Builder
	for _, e := range es {
		s := e.Text
		switch e.Type {
		case TypeIPv4:
			s = color.GreenString(e.Text)
		case TypeIPv6:
			s = color.BlueString(e.Text)
		case TypeDomain:
			s = color.YellowString(e.Text)
		}
		if texts := e.infoTexts(); len(texts) > 0 {
			for _, text := range texts {
				s += " [" + color.RedString(text) + "]"
			}
			s += " "
		}
		line.WriteString(s)
	}
	return line.String()
}

func (es Entities) Json() string {
	var s strings.Builder
	for _, e := range es {
		if e.Type == TypePlain {
			continue
		}
		s.WriteString(e.Json() + "\n")
	}
	return s.String()
}
