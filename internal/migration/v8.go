package migration

import (
	"log"
	"strings"

	"github.com/spf13/viper"
	"github.com/zu1k/nali/internal/constant"
	"github.com/zu1k/nali/internal/db"
	"github.com/zu1k/nali/pkg/ip2region"
)

func migration2v8() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(constant.ConfigDirPath)

	err := viper.ReadInConfig()
	if err != nil {
		return
	}

	dbList := db.List{}
	err = viper.UnmarshalKey("databases", &dbList)
	if err != nil {
		log.Fatalln("Config invalid:", err)
	}

	needOverwrite := false
	for _, adb := range dbList {
		if adb.Name == "ip2region" {
			if len(adb.DownloadUrls) > 0 && strings.Contains(adb.DownloadUrls[0], "ip2region.xdb") {
				needOverwrite = true
				adb.DownloadUrls = ip2region.DownloadUrls
			}
		}
	}

	if needOverwrite {
		viper.Set("databases", dbList)
		err = viper.WriteConfig()
		if err != nil {
			log.Println(err)
		}
	}
}
