package entity

import (
	"net/netip"
	"sort"

	"github.com/zu1k/nali/internal/db"
	"github.com/zu1k/nali/pkg/dbif"
	"github.com/zu1k/nali/pkg/re"
)

// ParseLine parse a line into entities
func ParseLine(line string) Entities {
	ip4sLoc := re.IPv4Re.FindAllStringIndex(line, -1)
	ip6sLoc := re.IPv6Re.FindAllStringIndex(line, -1)
	domainsLoc := re.DomainRe.FindAllStringIndex(line, -1)

	tmp := make(Entities, 0, len(ip4sLoc)+len(ip6sLoc)+len(domainsLoc))
	for _, e := range ip4sLoc {
		tmp = append(tmp, &Entity{
			Loc:  *(*[2]int)(e),
			Type: TypeIPv4,
			Text: line[e[0]:e[1]],
		})
	}
	for _, e := range ip6sLoc {
		text := line[e[0]:e[1]]
		if ip, _ := netip.ParseAddr(text); !ip.Is4In6() {
			tmp = append(tmp, &Entity{
				Loc:  *(*[2]int)(e),
				Type: TypeIPv6,
				Text: text,
			})
		}
	}
	for _, e := range domainsLoc {
		tmp = append(tmp, &Entity{
			Loc:  *(*[2]int)(e),
			Type: TypeDomain,
			Text: line[e[0]:e[1]],
		})
	}

	sort.Sort(tmp)
	var es Entities

	idx := 0
	for _, e := range tmp {
		start := e.Loc[0]
		if start >= idx {
			if start > idx {
				es = append(es, &Entity{
					Loc:  [2]int{idx, start},
					Type: TypePlain,
					Text: line[idx:start],
				})
			}
			results := db.Find(dbif.QueryType(e.Type), e.Text)
			if len(results) > 0 {
				e.InfoText = results[0].String()
				e.Info = results[0].Result
				e.Source = results[0].Source
				if len(results) > 1 {
					for _, r := range results {
						e.Results = append(e.Results, Result{Source: r.Source, Text: r.String(), Info: r.Result})
					}
				}
			} else {
				e.Type = TypePlain
			}
			idx = e.Loc[1]
			es = append(es, e)
		}
	}
	if total := len(line); idx < total {
		es = append(es, &Entity{
			Loc:  [2]int{idx, total},
			Type: TypePlain,
			Text: line[idx:total],
		})
	}
	return es
}
