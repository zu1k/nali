package qqwry

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io/ioutil"
	"log"

	"github.com/zu1k/nali/pkg/common"
)

func Download(filePath ...string) (data []byte, err error) {
	data, err = downloadAndDecrypt()
	if err != nil {
		log.Printf("纯真IP库下载失败，请手动下载解压后保存到本地: %s \n", filePath)
		log.Println("下载链接： https://qqwry.mirror.noc.one/qqwry.rar")
		return
	}

	if len(filePath) == 1 {
		if err := common.SaveFile(filePath[0], data); err == nil {
			log.Println("已将最新的 纯真IP库 保存到本地:", filePath)
		}
	}
	return
}

const (
	mirror = "https://qqwry.mirror.noc.one/qqwry.rar"
	key    = "https://qqwry.mirror.noc.one/copywrite.rar"
)

func downloadAndDecrypt() (data []byte, err error) {
	data, err = common.GetHttpClient().Get(mirror)
	if err != nil {
		return nil, err
	}

	key, err := getCopyWriteKey()
	if err != nil {
		return nil, err
	}

	return unRar(data, key)
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
	body, err := common.GetHttpClient().Get(key)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint32(body[5*4:]), nil
}
