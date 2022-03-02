package ip2region

import (
	"log"

	"github.com/zu1k/nali/pkg/common"
)

const (
	githubUrl   = "https://raw.githubusercontent.com/lionsoul2014/ip2region/master/data/ip2region.db"
	jsdelivrUrl = "https://cdn.jsdelivr.net/gh/lionsoul2014/ip2region/data/ip2region.db"
)

func Download(filePath ...string) (data []byte, err error) {
	data, err = common.GetHttpClient().Get(jsdelivrUrl, githubUrl)
	if err != nil {
		log.Printf("CDN数据库下载失败，请手动下载解压后保存到本地: %s \n", filePath)
		log.Println("下载链接：", githubUrl)
		return
	}

	if len(filePath) == 1 {
		if err := common.SaveFile(filePath[0], data); err == nil {
			log.Println("已将最新的 ip2region 保存到本地:", filePath)
		}
	}
	return
}
