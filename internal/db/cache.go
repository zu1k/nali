package db

import "github.com/zu1k/nali/pkg/dbif"

var (
	dbNameCache = make(map[string]dbif.DB)
	dbTypeCache = make(map[dbif.QueryType]dbif.DB)
	queryCache  = make(map[string]string)
)

var (
	NameDBMap = make(NameMap)
	TypeDBMap = make(TypeMap)
)
