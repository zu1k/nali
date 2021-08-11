package cdn

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type CDN struct {
	Data CDNDist
}

type CDNDist map[string]CDNResult

type CDNResult struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

func (r CDNResult) String() string {
	return r.Name
}

func NewCDN(filePath string) *CDN {
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
	return &CDN{Data: cdnDist}
}

func (db CDN) Find(query string, params ...string) (result fmt.Stringer, err error) {
	baseCname := parseBaseCname(query)
	if baseCname == "" {
		return nil, errors.New("base domain parse failed")
	}
	cdnResult, found := db.Data[baseCname]
	if found {
		return cdnResult, nil
	}

	if strings.Contains(baseCname, "kunlun") {
		return CDNResult{
			Name: "阿里云 CDN",
		}, nil
	}
	return nil, errors.New("not found")
}

func parseBaseCname(domain string) string {
	hostParts := strings.Split(domain, ".")
	if len(hostParts) < 2 {
		return domain
	}
	baseCname := hostParts[len(hostParts)-2] + "." + hostParts[len(hostParts)-1]
	return baseCname
}
