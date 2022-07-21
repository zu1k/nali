package ip2region

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/lionsoul2014/ip2region/binding/golang/ip2region"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/zu1k/nali/pkg/common"
	"github.com/zu1k/nali/pkg/download"
)

var DownloadUrls = []string{
	"https://cdn.jsdelivr.net/gh/lionsoul2014/ip2region/data/ip2region.xdb",
	"https://raw.githubusercontent.com/lionsoul2014/ip2region/master/data/ip2region.xdb",
}

type Ip2Region struct {
	seacher *xdb.Searcher
	db_old  *ip2region.Ip2Region
}

func NewIp2Region(filePath string) (*Ip2Region, error) {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		log.Println("文件不存在，尝试从网络获取最新 ip2region 库")
		_, err = download.Download(filePath, DownloadUrls...)
		if err != nil {
			return nil, err
		}
	}

	switch {
	case strings.HasSuffix(filePath, ".xdb"):
		f, err := os.OpenFile(filePath, os.O_RDONLY, 0400)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		data, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}
		searcher, err := xdb.NewWithBuffer(data)
		if err != nil {
			fmt.Printf("无法解析 ip2region xdb 数据库: %s\n", err)
			return nil, err
		}
		return &Ip2Region{
			seacher: searcher,
		}, nil
	default:
		region, err := ip2region.New(filePath)
		if err != nil {
			return nil, err
		}
		return &Ip2Region{
			db_old: region,
		}, nil
	}
}

func (db Ip2Region) Find(query string, params ...string) (result fmt.Stringer, err error) {
	if db.seacher != nil {
		res, err := db.seacher.SearchByStr(query)
		if err != nil {
			return nil, err
		} else {
			return common.Result{
				Country: strings.ReplaceAll(res, "|0", ""),
			}, nil
		}
	} else if db.db_old != nil {
		ip, err := db.db_old.MemorySearch(query)
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

		return common.Result{
			Country: ip.Country,
			Area:    area,
		}, nil
	}

	return nil, errors.New("ip2region 未初始化")
}
