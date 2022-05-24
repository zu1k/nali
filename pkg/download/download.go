package download

import (
	"log"

	"github.com/zu1k/nali/pkg/common"
)

func Download(filePath string, urls ...string) (data []byte, err error) {
	_ = urls[0]

	data, err = common.GetHttpClient().Get(urls...)
	if err != nil {
		log.Printf("文件下载失败，请手动下载解压后保存到本地: %s \n", filePath)
		log.Println("下载链接：", urls)
		return
	}

	if len(filePath) == 1 {
		if err := common.SaveFile(filePath, data); err == nil {
			log.Println("文件下载成功:", filePath)
		}
	}
	return
}
