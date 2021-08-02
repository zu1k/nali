package db

import (
	"github.com/zu1k/nali/pkg/dbif"
)

func Find(typ dbif.QueryType, query string) string {
	if result, found := queryCache[query]; found {
		return result
	}
	result, err := GetDB(typ).Find(query)
	if err != nil {
		//log.Printf("Query [%s] error: %s\n", query, err)
		return ""
	}
	return result.String()
}
