package db

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/zu1k/nali/pkg/download"
	"github.com/zu1k/nali/pkg/qqwry"
	"github.com/zu1k/nali/pkg/zxipv6wry"
)

func UpdateDB(dbNames ...string) {
	if len(dbNames) == 0 {
		dbNames = DbNameListForUpdate
	}

	done := make(map[string]struct{})
	for _, dbName := range dbNames {
		update, name := getUpdateFuncByName(dbName)
		if _, found := done[name]; !found {
			done[name] = struct{}{}
			if err := update(); err != nil {
				continue
			}
		}
	}
}

var DbNameListForUpdate = []string{
	"qqwry",
	"zxipv6wry",
	"ip2region",
	"cdn",
}

var DbCheckFunc = map[Format]func([]byte) bool{
	FormatQQWry:     qqwry.CheckFile,
	FormatZXIPv6Wry: zxipv6wry.CheckFile,
}

func getUpdateFuncByName(name string) (func() error, string) {
	name = strings.TrimSpace(name)
	if db := getDbByName(name); db != nil {
		// direct download if download-url not null
		if len(db.DownloadUrls) > 0 {
			return func() error {
				log.Printf("正在下载最新 %s 数据库...\n", db.Name)
				data, err := download.Download(db.File, db.DownloadUrls...)
				if err != nil {
					log.Printf("%s 数据库下载失败，请手动下载解压后保存到本地: %s \n", db.Name, db.File)
					log.Println("下载链接：", db.DownloadUrls)
					log.Println("error:", err)
					return err
				} else {
					if check, ok := DbCheckFunc[db.Format]; ok {
						if !check(data) {
							log.Printf("%s 数据库下载失败，请手动下载解压后保存到本地: %s \n", db.Name, db.File)
							log.Println("下载链接：", db.DownloadUrls)
							return errors.New("数据库内容出错")
						}
					}
					log.Printf("%s 数据库下载成功: %s\n", db.Name, db.File)
					return nil
				}
			}, string(db.Format)
		}

		// intenel download func
		switch db.Format {
		case FormatZXIPv6Wry:
			return func() error {
				log.Println("正在下载最新 ZX IPv6数据库...")
				_, err := zxipv6wry.Download(getDbByName("zxipv6wry").File)
				if err != nil {
					log.Println("数据库 ZXIPv6Wry 下载失败:", err)
				}
				return err
			}, FormatZXIPv6Wry
		default:
			return func() error {
				log.Println("暂不支持该类型数据库的自动更新")
				log.Println("可通过指定数据库的 download-urls 从特定链接下载数据库文件")
				return nil
			}, time.Now().String()
		}
	} else {
		return func() error {
			log.Fatalln("该名称的数据库未找到：", name)
			return nil
		}, time.Now().String()
	}
}
