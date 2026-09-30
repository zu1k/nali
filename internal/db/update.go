package db

import (
	"bytes"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/zu1k/nali/pkg/cdn"
	"github.com/zu1k/nali/pkg/common"
	"github.com/zu1k/nali/pkg/geoip"
	"github.com/zu1k/nali/pkg/ip2region"
	"github.com/zu1k/nali/pkg/qqwry"
	"github.com/zu1k/nali/pkg/zxipv6wry"
)

func UpdateDB(dbNames ...string) {
	if len(dbNames) == 0 {
		dbNames = defaultUpdateList()
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

// DbNameListForUpdate is always refreshed by a plain `nali update`.
var DbNameListForUpdate = []string{
	"qqwry",
	"zxipv6wry",
	"ip2region",
	"cdn",
}

// DbNameListForUpdateIfUsed holds large optional databases. A plain
// `nali update` only refreshes them when they are selected or already
// downloaded; `nali update --db <name>` always updates them.
var DbNameListForUpdateIfUsed = []string{
	"ip2region-ipv6",
	"ipinfo",
}

func defaultUpdateList() []string {
	names := append([]string{}, DbNameListForUpdate...)

	selected := make(map[string]bool)
	for _, key := range []string{"selected.ipv4", "selected.ipv6", "selected.cdn"} {
		for _, name := range strings.Split(viper.GetString(key), ",") {
			if db, found := lookupDb(strings.TrimSpace(name)); found {
				selected[db.Name] = true
			}
		}
	}

	for _, name := range DbNameListForUpdateIfUsed {
		db, found := lookupDb(name)
		if !found {
			continue
		}
		if _, err := os.Stat(db.File); selected[db.Name] || err == nil {
			names = append(names, name)
		}
	}
	return names
}

var DbCheckFunc = map[Format]func([]byte) bool{
	FormatQQWry:     qqwry.CheckFile,
	FormatZXIPv6Wry: zxipv6wry.CheckFile,
	FormatMMDB:      geoip.CheckFile,
	FormatIP2Region: ip2region.CheckFile,
	FormatCDNYml:    cdn.CheckFile,
}

func getUpdateFuncByName(name string) (func() error, string) {
	name = strings.TrimSpace(name)
	if db := getDbByName(name); db != nil {
		// direct download if download-url not null
		if len(db.DownloadUrls) > 0 {
			return func() error {
				log.Printf("正在下载最新 %s 数据库...\n", db.Name)
				data, err := common.GetHttpClient().Get(db.DownloadUrls...)
				if err == nil && sameContent(db.File, data) {
					log.Printf("%s 数据库已是最新版本: %s\n", db.Name, db.File)
					return nil
				}
				if err == nil {
					// validate before saving so a bad download never replaces a working database
					if check, ok := DbCheckFunc[db.Format]; ok && !check(data) {
						err = errors.New("数据库内容出错")
					}
				}
				if err == nil {
					err = common.SaveFile(db.File, data)
				}
				if err != nil {
					log.Printf("%s 数据库下载失败，请手动下载解压后保存到本地: %s \n", db.Name, db.File)
					log.Println("下载链接：", db.DownloadUrls)
					log.Println("error:", err)
					return err
				}
				log.Printf("%s 数据库下载成功: %s\n", db.Name, db.File)
				return nil
			}, db.Name
		}

		// intenel download func
		switch db.Format {
		case FormatZXIPv6Wry:
			return func() error {
				log.Println("正在下载最新 ZX IPv6数据库...")
				// download and validate without saving, so an unchanged
				// database is not rewritten
				data, err := zxipv6wry.Download()
				if err == nil && sameContent(db.File, data) {
					log.Printf("%s 数据库已是最新版本: %s\n", db.Name, db.File)
					return nil
				}
				if err == nil {
					err = common.SaveFile(db.File, data)
				}
				if err != nil {
					log.Println("数据库 ZXIPv6Wry 下载失败:", err)
					return err
				}
				log.Printf("%s 数据库下载成功: %s\n", db.Name, db.File)
				return nil
			}, db.Name
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

// sameContent reports whether the file at path holds exactly data.
func sameContent(path string, data []byte) bool {
	info, err := os.Stat(path)
	if err != nil || info.Size() != int64(len(data)) {
		return false
	}
	old, err := os.ReadFile(path)
	return err == nil && bytes.Equal(old, data)
}
