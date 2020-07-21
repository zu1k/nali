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

	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		log.Println("文件不存在，尝试从网络获取最新CDN数据库")
		cdnData, err = Download(filePath)
		if err != nil {
			os.Exit(1)
		}
	} else {
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
