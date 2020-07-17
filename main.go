package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/zu1k/nali/internal/ipdb"

	"github.com/zu1k/nali/internal/app"

	"github.com/zu1k/nali/cmd"
	"github.com/zu1k/nali/constant"
)

func main() {
	setHomePath()
	app.InitIPDB(getIPDBType())
	cmd.Execute()
}

func setHomePath() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	homePath := filepath.Join(homeDir, ".nali")
	constant.HomePath = homePath
	if _, err := os.Stat(homePath); os.IsNotExist(err) {
		if err := os.MkdirAll(homePath, 0777); err != nil {
			log.Fatal("can not create", homePath, ", use bin dir instead")
		}
	}
}

func getIPDBType() ipdb.IPDBType {
	dbname := os.Getenv("NALI_DB")
	dbname = strings.ToLower(dbname)
	switch dbname {
	case "geo", "geoip", "geoip2":
		return ipdb.GEOIP2
	case "chunzhen", "qqip", "qqwry":
		return ipdb.QQIP
	default:
		return ipdb.QQIP
	}
}
