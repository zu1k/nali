package cdn

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type CDN struct {
	Data CDNDist
}

type CDNDist map[string]CDNResult

type CDNResult struct {
	Name string `yaml:"name"`
	Link string `yaml:"link"`
}

func (r CDNResult) String() string {
	return r.Name
}

func NewCDN(filePath string) (*CDN, error) {
	cdnDist := make(CDNDist)
	cdnData := make([]byte, 0)

	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		log.Println("文件不存在，尝试从网络获取最新CDN数据库")
		cdnData, err = Download(filePath)
		if err != nil {
			return nil, err
		}
	} else {
		cdnFile, err := os.OpenFile(filePath, os.O_RDONLY, 0400)
		if err != nil {
			return nil, err
		}
		defer cdnFile.Close()

		cdnData, err = ioutil.ReadAll(cdnFile)
		if err != nil {
			return nil, err
		}
	}

	err = yaml.Unmarshal(cdnData, &cdnDist)
	if err != nil {
		return nil, err
	}
	return &CDN{Data: cdnDist}, nil
}

func (db CDN) Find(query string, params ...string) (result fmt.Stringer, err error) {
	baseCname := parseBaseCname(query)
	for _, domain := range baseCname {
		if domain != "" {
			cdnResult, found := db.Data[domain]
			if found {
				return cdnResult, nil
			}
		}

		if strings.Contains(domain, "kunlun") {
			return CDNResult{
				Name: "阿里云 CDN",
			}, nil
		}
	}

	return nil, errors.New("not found")
}

func parseBaseCname(domain string) (result []string) {
	parts := strings.Split(domain, ".")
	size := len(parts)
	if size == 0 {
		return []string{}
	}
	domain = parts[size-1]
	result = append(result, domain)
	for i := len(parts) - 2; i >= 0; i-- {
		domain = parts[i] + "." + domain
		result = append(result, domain)
	}
	return result
}
