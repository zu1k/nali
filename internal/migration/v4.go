package migration

import (
	"log"

	"github.com/spf13/viper"
	"github.com/zu1k/nali/internal/constant"
	"github.com/zu1k/nali/internal/db"
	"github.com/zu1k/nali/pkg/cdn"
	"github.com/zu1k/nali/pkg/ip2region"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(constant.WorkDirPath)

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
		if adb.Name == "ip2region" && adb.File != "ip2region.xdb" {
			needOverwrite = true
			adb.File = "ip2region.xdb"
			adb.DownloadUrls = ip2region.DownloadUrls
		}

		if adb.Name == "cdn" && adb.Format != "cdn-yml" {
			needOverwrite = true
			adb.Format = "cdn-yml"
			adb.DownloadUrls = cdn.DownloadUrls
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
