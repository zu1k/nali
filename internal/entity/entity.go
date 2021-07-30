package entity

import (
	"sort"
	"strings"
)

type EntityType uint

const (
	TypePlain = iota
	TypeIPv4
	TypeIPv6
	TypeDomain
)

type Entity struct {
	Index uint
	Length uint
	Type EntityType

	Text string
	Info string
}

func (e Entity) ParseInfo() error {
	return nil
}

type Entities []*Entity

func (es Entities) Len() int {
	return len(es)
}

func (es Entities) Less(i, j int) bool {
	return es[i].Index < es[j].Index
}

func (es Entities) Swap(i, j int) {
	es[i],es[j] = es[j],es[i]
}

func (es Entities) String() string {
	sort.Sort(es)

	var result strings.Builder
	for _, entity := range es {
		result.WriteString(entity.Text)
		if entity.Type!=TypePlain && len(entity.Info)>0 {
			result.WriteString("[" + entity.Info + "] ")
		}
	}

	return result.String()
}