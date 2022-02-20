package ip2region

import (
	"fmt"
	"log"
	"os"

	"github.com/lionsoul2014/ip2region/binding/golang/ip2region"

	"github.com/zu1k/nali/pkg/common"
)

type Ip2Region struct {
	db *ip2region.Ip2Region
}

func NewIp2Region(filePath string) Ip2Region {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		log.Println("文件不存在，尝试从网络获取最新 ip2region 库")
		_, err = Download(filePath)
		if err != nil {
			os.Exit(1)
		}
	}

	region, err := ip2region.New(filePath)
	if err != nil {
		panic(err)
	}

	return Ip2Region{
		db: region,
	}
}

func (db Ip2Region) Find(query string, params ...string) (result fmt.Stringer, err error) {
	ip, err := db.db.MemorySearch(query)
	if err != nil {
		return nil, err
	}

	fmt.Println(ip)

	result = common.Result{
		Country: ip.Country,
		Area:    ip.Province,
	}
	return result, nil
}
