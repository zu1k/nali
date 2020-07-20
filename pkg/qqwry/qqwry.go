package qqwry

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"

	"github.com/zu1k/nali/pkg/common"
	"golang.org/x/text/encoding/simplifiedchinese"
)

type QQwry struct {
	common.IPDB
}

// NewQQwry new db from path
func NewQQwry(filePath string) QQwry {
	var fileData []byte
	var fileInfo common.FileData

	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		log.Println("文件不存在，尝试从网络获取最新纯真 IP 库")
		fileData, err = Download()
		if err != nil {
			log.Printf("纯真IP库下载失败，请手动下载解压后保存到本地: %s \n", filePath)
			log.Println("下载链接： https://qqwry.mirror.noc.one/qqwry.rar")
			os.Exit(1)
		} else {
			if err := ioutil.WriteFile(filePath, fileData, 0644); err == nil {
				log.Printf("已将最新的 纯真IP库 保存到本地: %s ", filePath)
			}
		}
	} else {
		// 打开文件句柄
		fileInfo.FileBase, err = os.OpenFile(filePath, os.O_RDONLY, 0400)
		if err != nil {
			panic(err)
		}
		defer fileInfo.FileBase.Close()

		fileData, err = ioutil.ReadAll(fileInfo.FileBase)
		if err != nil {
			panic(err)
		}
	}
	fileInfo.Data = fileData

	buf := fileInfo.Data[0:8]
	start := binary.LittleEndian.Uint32(buf[:4])
	end := binary.LittleEndian.Uint32(buf[4:])

	return QQwry{
		IPDB: common.IPDB{
			Data:     &fileInfo,
			IndexLen: 7,
			IPNum:    (end-start)/7 + 1,
		},
	}
}

// Find ip地址查询对应归属地信息
func (db QQwry) Find(ip string) (res string) {
	if strings.Count(ip, ".") != 3 {
		return
	}

	ip4 := binary.BigEndian.Uint32(net.ParseIP(ip).To4())

	offset := db.searchIndex(ip4)
	if offset <= 0 {
		return
	}

	var gbkCountry []byte
	var gbkArea []byte

	mode := db.ReadMode(offset + 4)
	// [IP][0x01][国家和地区信息的绝对偏移地址]
	if mode == common.RedirectMode1 {
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
	} else if mode == common.RedirectMode2 {
		countryOffset := db.ReadUInt24()
		gbkCountry = db.ReadString(countryOffset)
		gbkArea = db.ReadArea(offset + 8)
	} else {
		gbkCountry = db.ReadString(offset + 4)
		gbkArea = db.ReadArea(offset + uint32(5+len(gbkCountry)))
	}

	enc := simplifiedchinese.GBK.NewDecoder()
	country, _ := enc.String(string(gbkCountry))
	area, _ := enc.String(string(gbkArea))

	country = strings.ReplaceAll(country, " CZ88.NET", "")
	area = strings.ReplaceAll(area, " CZ88.NET", "")

	return fmt.Sprintf("%s %s", country, area)
}

// searchIndex 查找索引位置
func (db *QQwry) searchIndex(ip uint32) uint32 {
	header := db.ReadData(8, 0)

	start := binary.LittleEndian.Uint32(header[:4])
	end := binary.LittleEndian.Uint32(header[4:])

	buf := make([]byte, db.IndexLen)
	mid := uint32(0)
	_ip := uint32(0)

	for {
		mid = db.GetMiddleOffset(start, end)
		buf = db.ReadData(db.IndexLen, mid)
		_ip = binary.LittleEndian.Uint32(buf[:4])

		if end-start == db.IndexLen {
			offset := common.ByteToUInt32(buf[4:])
			buf = db.ReadData(db.IndexLen)
			if ip < binary.LittleEndian.Uint32(buf[:4]) {
				return offset
			}
			return 0
		}

		// 找到的比较大，向前移
		if _ip > ip {
			end = mid
		} else if _ip < ip { // 找到的比较小，向后移
			start = mid
		} else if _ip == ip {
			return common.ByteToUInt32(buf[4:])
		}
	}
}
