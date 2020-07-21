package zxipv6wry

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/zu1k/nali/pkg/common"

	"github.com/saracen/go7z"
)

func Download(filePath string) (data []byte, err error) {
	data, err = getData()
	if err != nil {
		log.Printf("ZX IPv6数据库下载失败，请手动下载解压后保存到本地: %s \n", filePath)
		log.Println("下载链接： https://www.zxinc.org/ip.7z")
		return
	}
	common.ExistThenRemove(filePath)
	if err = ioutil.WriteFile(filePath, data, 0644); err == nil {
		log.Printf("已将最新的 ZX IPv6数据库 保存到本地: %s ", filePath)
	}
	return

}

func getData() (data []byte, err error) {
	resp, err := http.Get("https://www.zxinc.org/ip.7z")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	file7z, err := ioutil.TempFile("", "*")
	if err != nil {
		panic(err)
	}
	defer os.Remove(file7z.Name())

	if err := ioutil.WriteFile(file7z.Name(), body, 0644); err == nil {
		return Un7z(file7z.Name())
	}
	return
}

func Un7z(filePath string) (data []byte, err error) {
	sz, err := go7z.OpenReader(filePath)
	if err != nil {
		panic(err)
	}
	defer sz.Close()

	fileNoNeed, err := ioutil.TempFile("", "*")
	if err != nil {
		panic(err)
	}
	fileNeed, err := ioutil.TempFile("", "*")
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}
	for {
		hdr, err := sz.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			panic(err)
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
		panic(err)
	}
	defer os.Remove(fileNoNeed.Name())
	defer os.Remove(fileNeed.Name())
	return ioutil.ReadFile(fileNeed.Name())
}
