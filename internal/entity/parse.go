package entity

import (
	"sort"

	"github.com/zu1k/nali/internal/re"
)

// ParseLine parse a line into entities
func ParseLine(line string) Entities {
	ip4sLoc := re.IPv4Re.FindAllStringIndex(line, -1)
	ip6sLoc := re.IPv6Re.FindAllStringIndex(line, -1)
	domainsLoc := re.DomainRe.FindAllStringIndex(line, -1)

	tmp := make(Entities, 0, len(ip4sLoc)+len(ip6sLoc)+len(domainsLoc))
	for _, e := range ip4sLoc {
		tmp = append(tmp, &Entity{
			Loc:  e,
			Type: TypeIPv4,
			Text: line[e[0]:e[1]],
		})
	}
	for _, e := range ip6sLoc {
		tmp = append(tmp, &Entity{
			Loc:  e,
			Type: TypeIPv6,
			Text: line[e[0]:e[1]],
		})
	}
	for _, e := range domainsLoc {
		tmp = append(tmp, &Entity{
			Loc:  e,
			Type: TypeDomain,
			Text: line[e[0]:e[1]],
		})
	}

	sort.Sort(tmp)
	es := make(Entities, 0, len(tmp))

	idx := 0
	for _, e := range tmp {
		start := e.Loc[0]
		if start >= idx {
			if start > idx {
				es = append(es, &Entity{
					Loc:  []int{idx, start},
					Type: TypePlain,
					Text: line[idx:start],
				})
			}
			es = append(es, e)
			idx = e.Loc[1]
		}
	}

	return es
}
