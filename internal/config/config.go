package config

import (
	"log"

	"github.com/spf13/viper"
	"github.com/zu1k/nali/internal/db"
)

func ReadConfig(basePath string) {
	viper.SetDefault("databases", db.GetDefaultDBList())
	viper.SetDefault("selected.ipv4", "qqwry")
	viper.SetDefault("selected.ipv6", "zxipv6wry")
	viper.SetDefault("selected.cdn", "cdn")
	viper.SetDefault("selected.lang", "zh-CN")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(basePath)
	err := viper.ReadInConfig()
	if err != nil {
		err = viper.SafeWriteConfig()
		if err != nil {
			panic(err)
		}
	}

	_ = viper.BindEnv("selected.ipv4", "NALI_DB_IP4")
	_ = viper.BindEnv("selected.ipv6", "NALI_DB_IP6")
	_ = viper.BindEnv("selected.cdn", "NALI_DB_CDN")
	_ = viper.BindEnv("selected.lang", "NALI_LANG")

	dbList := db.List{}
	err = viper.UnmarshalKey("databases", &dbList)
	if err != nil {
		log.Fatalln("Config invalid:", err)
	}

	db.NameDBMap.From(dbList)
	db.TypeDBMap.From(dbList)
}
