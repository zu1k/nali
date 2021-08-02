package dbif

func init() {

}

type langMap map[string][]DB
type dataTypeMap map[QueryType][]DB

var (
	lang2DB = make(langMap)
	type2DB = make(dataTypeMap)
)

func RegistLang(lang string, db DB) {
	originDBs, found := lang2DB[lang]
	if !found {
		originDBs = make([]DB, 0, 1)
	}
	lang2DB[lang] = append(originDBs, db)
}

func RegistType(typ QueryType, db DB) {
	originDBs, found := type2DB[typ]
	if !found {
		originDBs = make([]DB, 0, 1)
	}
	type2DB[typ] = append(originDBs, db)
}
