package ip2region

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/lionsoul2014/ip2region/binding/golang/ip2region"
	"github.com/zu1k/nali/pkg/common"
)

type Ip2Region struct {
	db *ip2region.Ip2Region
}

func NewIp2Region(filePath string) (*Ip2Region, error) {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		log.Println("文件不存在，尝试从网络获取最新 ip2region 库")
		_, err = Download(filePath)
		if err != nil {
			return nil, err
		}
	}

	region, err := ip2region.New(filePath)
	if err != nil {
		return nil, err
	}

	return &Ip2Region{
		db: region,
	}, nil
}

func (db Ip2Region) Find(query string, params ...string) (result fmt.Stringer, err error) {
	ip, err := db.db.MemorySearch(query)
	if err != nil {
		return nil, err
	}

	area := ""
	if ip.Province != "0" {
		area = ip.Province
	}
	if ip.City != "0" && strings.EqualFold(ip.City, ip.Province) {
		area = area + " " + ip.Province
	}
	if ip.ISP != "0" {
		area = area + " " + ip.ISP
	}

	result = common.Result{
		Country: ip.Country,
		Area:    area,
	}
	return result, nil
}
