package qqwry

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/zu1k/nali/pkg/common"
)

func Download(filePath string) (data []byte, err error) {
	data, err = getData()
	if err != nil {
		log.Printf("纯真IP库下载失败，请手动下载解压后保存到本地: %s \n", filePath)
		log.Println("下载链接： https://qqwry.mirror.noc.one/qqwry.rar")
		return
	}
	common.ExistThenRemove(filePath)
	if err = ioutil.WriteFile(filePath, data, 0644); err == nil {
		log.Printf("已将最新的 纯真IP库 保存到本地: %s ", filePath)
	}
	return
}

func getData() (data []byte, err error) {
	resp, err := http.Get("https://qqwry.mirror.noc.one/qqwry.rar")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	key, err := getCopyWriteKey()
	if err != nil {
		return
	}

	return unRar(body, key)
}

func unRar(data []byte, key uint32) ([]byte, error) {
	for i := 0; i < 0x200; i++ {
		key = key * 0x805
		key++
		key = key & 0xff

		data[i] = byte(uint32(data[i]) ^ key)
	}

	reader, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(reader)
}

func getCopyWriteKey() (uint32, error) {
	resp, err := http.Get("https://qqwry.mirror.noc.one/copywrite.rar")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		return 0, err
	} else {
		return binary.LittleEndian.Uint32(body[5*4:]), nil
	}
}
