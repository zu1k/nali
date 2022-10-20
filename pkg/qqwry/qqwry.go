package qqwry

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/zu1k/nali/pkg/download"
	"github.com/zu1k/nali/pkg/wry"
)

var DownloadUrls = []string{
	"https://99wry.cf/qqwry.dat",
}

type QQwry struct {
	wry.IPDB[uint32]
}

// NewQQwry new database from path
func NewQQwry(filePath string) (*QQwry, error) {
	var fileData []byte

	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		log.Println("文件不存在，尝试从网络获取最新纯真 IP 库")
		fileData, err = download.Download(filePath, DownloadUrls...)
		if err != nil {
			return nil, err
		}
	} else {
		fileBase, err := os.OpenFile(filePath, os.O_RDONLY, 0400)
		if err != nil {
			return nil, err
		}
		defer fileBase.Close()

		fileData, err = io.ReadAll(fileBase)
		if err != nil {
			return nil, err
		}
	}

	if len(fileData) < 8 {
		log.Fatalln("纯真 IP 库存在错误，请重新下载")
	}

	header := fileData[0:8]
	start := binary.LittleEndian.Uint32(header[:4])
	end := binary.LittleEndian.Uint32(header[4:])

	if uint32(len(fileData)) < end+7 {
		log.Fatalln("纯真 IP 库存在错误，请重新下载")
	}

	return &QQwry{
		IPDB: wry.IPDB[uint32]{
			Data: fileData,

			OffLen:   3,
			IPLen:    4,
			IPCnt:    (end-start)/7 + 1,
			IdxStart: start,
			IdxEnd:   end,
		},
	}, nil
}

func (db QQwry) Find(query string, params ...string) (result fmt.Stringer, err error) {
	ip := net.ParseIP(query)
	if ip == nil {
		return nil, errors.New("query should be IPv4")
	}
	ip4 := ip.To4()
	if ip4 == nil {
		return nil, errors.New("query should be IPv4")
	}
	ip4uint := binary.BigEndian.Uint32(ip4)

	offset := db.SearchIndexV4(ip4uint)
	if offset <= 0 {
		return nil, errors.New("query not valid")
	}

	reader := wry.NewReader(db.Data)
	reader.Parse(offset + 4)
	return reader.Result.DecodeGBK(), nil
}
