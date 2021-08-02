package db

import "github.com/zu1k/nali/pkg/dbif"

var dbCache = make(map[dbif.QueryType]dbif.DB)
var queryCache = make(map[string]string)
