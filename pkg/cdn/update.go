package cdn

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/zu1k/nali/pkg/common"
)

func Download(filePath string) (data []byte, err error) {
	data, err = getData()
	if err != nil {
		log.Printf("CDN数据库下载失败，请手动下载解压后保存到本地: %s \n", filePath)
		log.Println("下载链接： https://cdn.jsdelivr.net/gh/SukkaLab/cdn/dist/cdn.json")
		return
	}

	common.ExistThenRemove(filePath)
	if err := ioutil.WriteFile(filePath, data, 0644); err == nil {
		log.Printf("已将最新的 CDN数据库 保存到本地: %s \n", filePath)
	}
	return
}

func getData() (data []byte, err error) {
	//resp, err := http.Get("https://raw.githubusercontent.com/SukkaLab/cdn/master/dist/cdn.json")
	resp, err := http.Get("https://cdn.jsdelivr.net/gh/SukkaLab/cdn/dist/cdn.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
