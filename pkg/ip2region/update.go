package ip2region

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
		log.Println("下载链接：", githubUrl)
		return
	}

	common.ExistThenRemove(filePath)
	if err := ioutil.WriteFile(filePath, data, 0644); err == nil {
		log.Printf("已将最新的 ip2region 保存到本地: %s \n", filePath)
	}
	return
}

const (
	githubUrl   = "https://raw.githubusercontent.com/lionsoul2014/ip2region/master/data/ip2region.db"
	jsdelivrUrl = "https://cdn.jsdelivr.net/gh/lionsoul2014/ip2region/data/ip2region.db"
)

func getData() (data []byte, err error) {
	resp, err := http.Get(jsdelivrUrl)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		resp, err = http.Get(githubUrl)
		if err != nil {
			return nil, err
		}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
