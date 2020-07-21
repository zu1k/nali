package app

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/zu1k/nali/internal/tools"

	"github.com/zu1k/nali/constant"

	"github.com/zu1k/nali/pkg/cdn"
)

var (
	cdnDB    cdn.CDN
	domainRe *regexp.Regexp
)

func init() {
	domainRe = regexp.MustCompile(`[0-9A-Za-z]{2,}\.[0-9A-Za-z]{2,3}\.[0-9A-Za-z]{2,3}|[0-9A-Za-z]{2,}\.[0-9A-Za-z]{2,3}`)
}

func InitCDNDB() {
	cdnDB = cdn.NewCDN(filepath.Join(constant.HomePath, "cdn.json"))
}

func ParseCDNs(str []string) {
	for _, cname := range str {
		name := find(cname)
		fmt.Printf("%s [%s]\n", cname, name)
	}
}

func find(cname string) string {
	baseCname := parseBaseCname(cname)
	if baseCname == "" {
		return "无法解析"
	}
	cdnResult, found := cdnDB.Data[baseCname]
	if found {
		return cdnResult.Name
	}
	return "未找到"
}

func ReplaceCDNInString(str string) (result string) {
	done := make(map[string]bool)

	cnames := domainRe.FindAllString(str, -1)
	result = str
	for _, cname := range cnames {
		name := find(cname)
		if name != "未找到" && name != "无法解析" {
			if _, found := done[cname]; found {
				continue
			}
			result = tools.ReplaceAdd(result, cname, fmt.Sprintf("%s [%s]", cname, name))
			done[cname] = true
		}
	}
	return
}

func parseBaseCname(domain string) string {
	hostParts := strings.Split(domain, ".")
	if len(hostParts) < 2 {
		return domain
	}
	baseCname := hostParts[len(hostParts)-2] + "." + hostParts[len(hostParts)-1]
	return baseCname
}

func UpdateDB() {
	filePath := filepath.Join(constant.HomePath, "cdn.json")

	log.Println("正在下载最新 CDN数据库...")
	_, err := cdn.Download(filePath)
	if err != nil {
		log.Fatalln("下载失败", err.Error())
		return
	}
}
