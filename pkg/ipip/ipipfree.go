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

type Result struct {
	Country string
	Region  string
	City    string
}

func (r Result) String() string {
	if r.City == "" {
		return fmt.Sprintf("%s %s", r.Country, r.Region)
	}
	return fmt.Sprintf("%s %s %s", r.Country, r.Region, r.City)
}

func (db IPIPFree) Find(query string, params ...string) (result fmt.Stringer, err error) {
	info, err := db.FindInfo(query, "CN")
	if err != nil || info == nil {
		return nil, err
	} else {
		// info contains more info
		result = Result{
			Country: info.CountryName,
			Region:  info.RegionName,
			City:    info.CityName,
		}
		return
	}
}
