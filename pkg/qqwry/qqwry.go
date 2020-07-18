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
	var tmpData []byte
	var fileInfo common.FileInfo

	// 判断文件是否存在
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		log.Println("文件不存在，尝试从网络获取最新纯真 IP 库")
		tmpData, err = GetOnline()
		if err != nil {
			panic(err)
		} else {
			if err := ioutil.WriteFile(filePath, tmpData, 0644); err == nil {
				log.Printf("已将最新的纯真 IP 库保存到本地 %s ", filePath)
			}
		}
	} else {
		// 打开文件句柄
		fileInfo.FileBase, err = os.OpenFile(filePath, os.O_RDONLY, 0400)
		if err != nil {
			panic(err)
		}
		defer fileInfo.FileBase.Close()

		tmpData, err = ioutil.ReadAll(fileInfo.FileBase)
		if err != nil {
			panic(err)
		}
	}

	fileInfo.Data = tmpData

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
func (q QQwry) Find(ip string) (res string) {
	if strings.Count(ip, ".") != 3 {
		return
	}

	ip4 := binary.BigEndian.Uint32(net.ParseIP(ip).To4())

	offset := q.searchIndex(ip4)
	if offset <= 0 {
		return
	}

	var gbkCountry []byte
	var gbkArea []byte

	mode := q.ReadMode(offset + 4)
	// [IP][0x01][国家和地区信息的绝对偏移地址]
	if mode == common.RedirectMode1 {
		countryOffset := q.ReadUInt24()
		mode = q.ReadMode(countryOffset)
		if mode == common.RedirectMode2 {
			c := q.ReadUInt24()
			gbkCountry = q.ReadString(c)
			countryOffset += 4
		} else {
			gbkCountry = q.ReadString(countryOffset)
			countryOffset += uint32(len(gbkCountry) + 1)
		}
		gbkArea = q.ReadArea(countryOffset)
	} else if mode == common.RedirectMode2 {
		countryOffset := q.ReadUInt24()
		gbkCountry = q.ReadString(countryOffset)
		gbkArea = q.ReadArea(offset + 8)
	} else {
		gbkCountry = q.ReadString(offset + 4)
		gbkArea = q.ReadArea(offset + uint32(5+len(gbkCountry)))
	}

	enc := simplifiedchinese.GBK.NewDecoder()
	country, _ := enc.String(string(gbkCountry))
	area, _ := enc.String(string(gbkArea))

	return fmt.Sprintf("%s %s", country, area)
}

// searchIndex 查找索引位置
func (q *QQwry) searchIndex(ip uint32) uint32 {
	header := q.ReadData(8, 0)

	start := binary.LittleEndian.Uint32(header[:4])
	end := binary.LittleEndian.Uint32(header[4:])

	buf := make([]byte, q.IndexLen)
	mid := uint32(0)
	_ip := uint32(0)

	for {
		mid = q.GetMiddleOffset(start, end)
		buf = q.ReadData(q.IndexLen, mid)
		_ip = binary.LittleEndian.Uint32(buf[:4])

		if end-start == q.IndexLen {
			offset := common.ByteToUInt32(buf[4:])
			buf = q.ReadData(q.IndexLen)
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
