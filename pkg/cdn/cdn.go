package cdn

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type CDN struct {
	Data CDNDist
}

type CDNDist map[string]CDNResult

type CDNResult struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

func NewCDN(filePath string) CDN {
	cdnDist := make(CDNDist)
	cdnData := make([]byte, 0)

	// 判断文件是否存在
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		log.Println("文件不存在，尝试从网络获取最新CDN数据库")
		cdnData, err = Download()
		if err != nil {
			panic(err)
		} else {
			if err := ioutil.WriteFile(filePath, cdnData, 0644); err == nil {
				log.Printf("已将最新的 CDN数据库 保存到本地: %s ", filePath)
			}
		}
	} else {
		// 打开文件句柄
		cdnFile, err := os.OpenFile(filePath, os.O_RDONLY, 0400)
		if err != nil {
			panic(err)
		}
		defer cdnFile.Close()

		cdnData, err = ioutil.ReadAll(cdnFile)
		if err != nil {
			panic(err)
		}
	}

	err = json.Unmarshal(cdnData, &cdnDist)
	if err != nil {
		panic("cdn data parse failed!")
	}
	return CDN{Data: cdnDist}
}
