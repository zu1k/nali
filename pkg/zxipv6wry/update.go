package zxipv6wry

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/saracen/go7z"
	"github.com/zu1k/nali/pkg/common"
)

func Download(filePath ...string) (data []byte, err error) {
	data, err = getData()
	if err != nil {
		log.Printf("ZX IPv6数据库下载失败，请手动下载解压后保存到本地: %s \n", filePath)
		log.Println("下载链接： https://ip.zxinc.org/ip.7z")
		return
	}

	if !CheckFile(data) {
		log.Printf("ZX IPv6数据库下载出错，请手动下载解压后保存到本地: %s \n", filePath)
		log.Println("下载链接： https://ip.zxinc.org/ip.7z")
		return nil, errors.New("数据库下载内容出错")
	}

	if len(filePath) == 1 {
		if err := common.SaveFile(filePath[0], data); err == nil {
			log.Println("已将最新的 ZX IPv6数据库 保存到本地:", filePath)
		}
	}
	return
}

const (
	zx = "https://ip.zxinc.org/ip.7z"
)

func getData() (data []byte, err error) {
	data, err = common.GetHttpClient().Get(zx)

	file7z, err := os.CreateTemp("", "*")
	if err != nil {
		return nil, err
	}
	defer os.Remove(file7z.Name())
	if err := os.WriteFile(file7z.Name(), data, 0644); err == nil {
		return Un7z(file7z.Name())
	}
	return
}

func Un7z(filePath string) (data []byte, err error) {
	sz, err := go7z.OpenReader(filePath)
	if err != nil {
		return nil, err
	}
	defer sz.Close()

	fileNoNeed, err := os.CreateTemp("", "*")
	if err != nil {
		return nil, err
	}
	fileNeed, err := os.CreateTemp("", "*")
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	for {
		hdr, err := sz.Next()
		if err == io.EOF {
			break // IdxEnd of archive
		}
		if err != nil {
			return nil, err
		}

		if hdr.Name == "ipv6wry.db" {
			if _, err := io.Copy(fileNeed, sz); err != nil {
				log.Fatalln("ZX ipv6数据库解压出错：", err.Error())
			}
		} else {
			if _, err := io.Copy(fileNoNeed, sz); err != nil {
				log.Fatalln("ZX ipv6数据库解压出错：", err.Error())
			}
		}
	}
	err = fileNoNeed.Close()
	if err != nil {
		return nil, err
	}
	defer os.Remove(fileNoNeed.Name())
	defer os.Remove(fileNeed.Name())
	return ioutil.ReadFile(fileNeed.Name())
}
