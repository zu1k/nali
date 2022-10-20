package ip2region

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/zu1k/nali/pkg/download"
	"github.com/zu1k/nali/pkg/wry"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

var DownloadUrls = []string{
	"https://cdn.jsdelivr.net/gh/lionsoul2014/ip2region/data/ip2region.xdb",
	"https://raw.githubusercontent.com/lionsoul2014/ip2region/master/data/ip2region.xdb",
}

type Ip2Region struct {
	seacher *xdb.Searcher
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
}

func (db Ip2Region) Find(query string, params ...string) (result fmt.Stringer, err error) {
	if db.seacher != nil {
		res, err := db.seacher.SearchByStr(query)
		if err != nil {
			return nil, err
		} else {
			return wry.Result{
				Country: strings.ReplaceAll(res, "|0", ""),
			}, nil
		}
	}

	return nil, errors.New("ip2region 未初始化")
}
