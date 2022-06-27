package ip2region

import (
	"log"

	"github.com/zu1k/nali/pkg/common"
)

var DownloadUrls = []string{
	"https://cdn.jsdelivr.net/gh/lionsoul2014/ip2region/data/ip2region.xdb",
	"https://raw.githubusercontent.com/lionsoul2014/ip2region/master/data/ip2region.xdb",
}

// Deprecated: This will be removed from 0.5.0, use package download instead
func Download(filePath ...string) (data []byte, err error) {
	data, err = common.GetHttpClient().Get(DownloadUrls...)
	if err != nil {
		log.Printf("CDN数据库下载失败，请手动下载解压后保存到本地: %s \n", filePath)
		log.Println("下载链接：", DownloadUrls)
		return
	}

	if len(filePath) == 1 {
		if err := common.SaveFile(filePath[0], data); err == nil {
			log.Println("已将最新的 ip2region 保存到本地:", filePath)
		}
	}
	return
}
