package qqwry

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"

	"github.com/zu1k/nali/pkg/common"
	"github.com/zu1k/nali/pkg/download"
	"golang.org/x/text/encoding/simplifiedchinese"
)

var DownloadUrls = []string{
	"https://99wry.cf/qqwry.dat",
}

type QQwry struct {
	common.IPDB
}

// NewQQwry new database from path
func NewQQwry(filePath string) (*QQwry, error) {
	var fileData []byte
	var fileInfo common.FileData

	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		log.Println("文件不存在，尝试从网络获取最新纯真 IP 库")
		fileData, err = download.Download(filePath, DownloadUrls...)
		if err != nil {
			return nil, err
		}
	} else {
		fileInfo.FileBase, err = os.OpenFile(filePath, os.O_RDONLY, 0400)
		if err != nil {
			return nil, err
		}
		defer fileInfo.FileBase.Close()

		fileData, err = ioutil.ReadAll(fileInfo.FileBase)
		if err != nil {
			return nil, err
		}
	}
	fileInfo.Data = fileData

	buf := fileInfo.Data[0:8]
	start := binary.LittleEndian.Uint32(buf[:4])
	end := binary.LittleEndian.Uint32(buf[4:])

	return &QQwry{
		IPDB: common.IPDB{
			Data:  &fileInfo,
			IPNum: (end-start)/7 + 1,
		},
	}, nil
}

func (db QQwry) Find(query string, params ...string) (result fmt.Stringer, err error) {
	ip := net.ParseIP(query)
	if ip == nil {
		return nil, errors.New("Query should be IPv4")
	}
	ip4 := ip.To4()
	if ip4 == nil {
		return nil, errors.New("Query should be IPv4")
	}
	ip4uint := binary.BigEndian.Uint32(ip4)

	offset := db.searchIndex(ip4uint)
	if offset <= 0 {
		return nil, errors.New("Query not valid")
	}

	var gbkCountry []byte
	var gbkArea []byte

	mode := db.ReadMode(offset + 4)
	switch mode {
	case common.RedirectMode1: // [IP][0x01][国家和地区信息的绝对偏移地址]
		countryOffset := db.ReadUInt24()
		mode = db.ReadMode(countryOffset)
		if mode == common.RedirectMode2 {
			c := db.ReadUInt24()
			gbkCountry = db.ReadString(c)
			countryOffset += 4
		} else {
			gbkCountry = db.ReadString(countryOffset)
			countryOffset += uint32(len(gbkCountry) + 1)
		}
		gbkArea = db.ReadArea(countryOffset)
	case common.RedirectMode2:
		countryOffset := db.ReadUInt24()
		gbkCountry = db.ReadString(countryOffset)
		gbkArea = db.ReadArea(offset + 8)
	default:
		gbkCountry = db.ReadString(offset + 4)
		gbkArea = db.ReadArea(offset + uint32(5+len(gbkCountry)))
	}

	enc := simplifiedchinese.GBK.NewDecoder()
	country, _ := enc.String(string(gbkCountry))
	area, _ := enc.String(string(gbkArea))

	result = common.Result{
		Country: strings.ReplaceAll(country, " CZ88.NET", ""),
		Area:    strings.ReplaceAll(area, " CZ88.NET", ""),
	}
	return result, nil
}

// searchIndex 查找索引位置
func (db *QQwry) searchIndex(ip uint32) uint32 {
	header := db.ReadData(8, 0)

	start := binary.LittleEndian.Uint32(header[:4])
	end := binary.LittleEndian.Uint32(header[4:])

	buf := make([]byte, 7)
	mid := uint32(0)
	ipUint := uint32(0)

	for {
		mid = common.GetMiddleOffset(start, end, 7)
		buf = db.ReadData(7, mid)
		ipUint = binary.LittleEndian.Uint32(buf[:4])

		if end-start == 7 {
			offset := common.ByteToUInt32(buf[4:])
			buf = db.ReadData(7)
			if ip < binary.LittleEndian.Uint32(buf[:4]) {
				return offset
			}
			return 0
		}

		if ipUint > ip {
			end = mid
		} else if ipUint < ip {
			start = mid
		} else if ipUint == ip {
			return common.ByteToUInt32(buf[4:])
		}
	}
}
