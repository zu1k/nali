package ip2region

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/zu1k/nali/pkg/download"
	"github.com/zu1k/nali/pkg/wry"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

var DownloadUrls = []string{
	"https://cdn.jsdelivr.net/gh/lionsoul2014/ip2region/data/ip2region_v4.xdb",
	"https://raw.githubusercontent.com/lionsoul2014/ip2region/master/data/ip2region_v4.xdb",
}

var DownloadUrlsV6 = []string{
	"https://raw.githubusercontent.com/lionsoul2014/ip2region/master/data/ip2region_v6.xdb",
}

type Ip2Region struct {
	seacher *xdb.Searcher
}

func NewIp2Region(filePath string, downloadUrls []string) (*Ip2Region, error) {
	if len(downloadUrls) == 0 {
		downloadUrls = DownloadUrls
	}
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		log.Println("文件不存在，尝试从网络获取最新 ip2region 库")
		_, err = download.Download(filePath, downloadUrls...)
		if err != nil {
			return nil, err
		}
	}

	f, err := os.OpenFile(filePath, os.O_RDONLY, 0400)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	version, err := detectVersion(data)
	if err != nil {
		fmt.Printf("无法解析 ip2region xdb 数据库版本: %s\n", err)
		return nil, err
	}

	searcher, err := xdb.NewWithBuffer(version, data)
	if err != nil {
		fmt.Printf("无法解析 ip2region xdb 数据库: %s\n", err)
		return nil, err
	}
	return &Ip2Region{
		seacher: searcher,
	}, nil
}

// CheckFile reports whether data looks like a valid ip2region xdb database.
func CheckFile(data []byte) bool {
	if len(data) < xdb.HeaderInfoLength {
		return false
	}
	header, err := xdb.NewHeader(data)
	if err != nil {
		return false
	}
	version, err := xdb.VersionFromHeader(header)
	if err != nil {
		return false
	}
	return header.StartIndexPtr >= xdb.HeaderInfoLength &&
		header.StartIndexPtr <= header.EndIndexPtr &&
		uint64(header.EndIndexPtr)+uint64(version.SegmentIndexSize) <= uint64(len(data))
}

func detectVersion(data []byte) (*xdb.Version, error) {
	header, err := xdb.NewHeader(data)
	if err != nil {
		return nil, err
	}
	return xdb.VersionFromHeader(header)
}

func (db Ip2Region) Find(query string, params ...string) (result fmt.Stringer, err error) {
	if db.seacher != nil {
		res, err := db.seacher.Search(query)
		if err != nil {
			return nil, err
		} else {
			return wry.Result{
				Country: formatRegion(res),
			}, nil
		}
	}

	return nil, errors.New("ip2region 未初始化")
}

func (db Ip2Region) Name() string {
	return "ip2region"
}

// formatRegion drops the "0" placeholders ip2region uses for unknown fields,
// e.g. "中国|0|浙江省|杭州市|阿里" -> "中国|浙江省|杭州市|阿里".
func formatRegion(region string) string {
	fields := strings.Split(region, "|")
	kept := fields[:0]
	for _, f := range fields {
		if f != "0" && f != "" {
			kept = append(kept, f)
		}
	}
	return strings.Join(kept, "|")
}
