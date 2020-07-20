package ipip

import (
	"fmt"
	"log"
	"os"

	"github.com/ipipdotnet/ipdb-go"
)

type IPIPFree struct {
	*ipdb.City
}

func NewIPIPFree(filePath string) IPIPFree {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		log.Printf("IPIP数据库不存在，请手动下载解压后保存到本地: %s \n", filePath)
		log.Println("下载链接： https://www.ipip.net/product/ip.html")
		os.Exit(1)
		return IPIPFree{}
	} else {
		db, err := ipdb.NewCity(filePath)
		if err != nil {
			log.Fatalln("IPIP 数据库 初始化失败")
			log.Fatal(err)
			os.Exit(1)
		}
		return IPIPFree{City: db}
	}
}

func (db IPIPFree) Find(ip string) string {
	info, err := db.FindInfo(ip, "CN")
	if err != nil {
		log.Fatalln("IPIP 查询失败：", err.Error())
		return ""
	} else {
		if info.CityName == "" {
			return fmt.Sprintf("%s %s", info.CountryName, info.RegionName)
		}
		return fmt.Sprintf("%s %s %s", info.CountryName, info.RegionName, info.CityName)
	}
}
