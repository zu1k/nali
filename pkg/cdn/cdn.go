package cdn

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/zu1k/nali/pkg/download"
	"github.com/zu1k/nali/pkg/re"
	"gopkg.in/yaml.v2"
)

var DownloadUrls = []string{
	"https://cdn.jsdelivr.net/gh/4ft35t/cdn/src/cdn.yml",
	"https://raw.githubusercontent.com/4ft35t/cdn/master/src/cdn.yml",
	"https://raw.githubusercontent.com/SukkaLab/cdn/master/src/cdn.yml",
}

type CDN struct {
	Map   map[string]CDNResult
	ReMap []CDNReTuple
}

type CDNReTuple struct {
	*regexp.Regexp
	CDNResult
}

type CDNResult struct {
	Name string `yaml:"name"`
	Link string `yaml:"link"`
}

func (r CDNResult) String() string {
	return r.Name
}

func NewCDN(filePath string) (*CDN, error) {
	fileData := make([]byte, 0)
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		log.Println("文件不存在，尝试从网络获取最新CDN数据库")
		fileData, err = download.Download(filePath, DownloadUrls...)
		if err != nil {
			return nil, err
		}
	} else {
		cdnFile, err := os.OpenFile(filePath, os.O_RDONLY, 0400)
		if err != nil {
			return nil, err
		}
		defer cdnFile.Close()

		fileData, err = ioutil.ReadAll(cdnFile)
		if err != nil {
			return nil, err
		}
	}

	cdnMap := make(map[string]CDNResult)
	err = yaml.Unmarshal(fileData, &cdnMap)
	if err != nil {
		return nil, err
	}
	cdnReMap := make([]CDNReTuple, 0)
	for k, v := range cdnMap {
		if re.MaybeRegexp(k) {
			rex, err := regexp.Compile(k)
			if err != nil {
				log.Printf("[CDN Database] entry %s not a valid regexp", k)
			}
			cdnReMap = append(cdnReMap, CDNReTuple{
				Regexp:    rex,
				CDNResult: v,
			})
		}
	}

	return &CDN{Map: cdnMap, ReMap: cdnReMap}, nil
}

func (db CDN) Find(query string, params ...string) (result fmt.Stringer, err error) {
	baseCname := parseBaseCname(query)
	for _, domain := range baseCname {
		if domain != "" {
			cdnResult, found := db.Map[domain]
			if found {
				return cdnResult, nil
			}
		}

		for _, entry := range db.ReMap {
			if entry.Regexp.MatchString(domain) {
				return entry.CDNResult, nil
			}
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
