package zxipv6wry

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/zu1k/nali/constant"

	"github.com/saracen/go7z"
)

func Download(filePath string) (data []byte, err error) {
	resp, err := http.Get("https://www.zxinc.org/ip.7z")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	file7z := filepath.Join(constant.HomePath, "ip.7z")
	_, err = os.Stat(file7z)
	if err != nil && os.IsNotExist(err) {
		if err := ioutil.WriteFile(file7z, body, 0644); err == nil {
			Un7z(file7z)
			err = os.Remove(file7z)
		}
	} else {
		Un7z(file7z)
		err = os.Remove(file7z)
	}

	return ioutil.ReadFile(filePath)
}

func Un7z(filePath string) {
	sz, err := go7z.OpenReader(filePath)
	if err != nil {
		panic(err)
	}
	defer sz.Close()

	f, err := os.Create("tmp")
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
			f, err := os.Create(filepath.Join(constant.HomePath, "ipv6wry.db"))
			if err != nil {
				panic(err)
			}
			defer f.Close()

			if _, err := io.Copy(f, sz); err != nil {
				log.Fatalln("ZX ipv6数据库解压出错：", err.Error())
			}
		} else {
			if _, err := io.Copy(f, sz); err != nil {
				log.Fatalln("ZX ipv6数据库解压出错：", err.Error())
			}
		}
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}
	_ = os.Remove("tmp")
}
