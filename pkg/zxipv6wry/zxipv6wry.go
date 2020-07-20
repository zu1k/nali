package zxipv6wry

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"os"
	"strings"

	"github.com/zu1k/nali/pkg/common"
)

type ZXwry struct {
	common.IPDB
}

func NewZXwry(filePath string) ZXwry {
	var fileData []byte
	var fileInfo common.FileData

	// 判断文件是否存在
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		log.Println("文件不存在，尝试从网络获取最新ZX IPv6数据库")
		fileData, err = Download(filePath)
		if err != nil {
			log.Printf("ZX IPv6数据库下载失败，请手动下载解压后保存到本地: %s \n", filePath)
			log.Println("下载链接： https://www.zxinc.org/ip.7z")
			os.Exit(1)
		} else {
			if err := ioutil.WriteFile(filePath, fileData, 0644); err == nil {
				log.Printf("已将最新的 ZX IPv6数据库 保存到本地: %s ", filePath)
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

	return ZXwry{
		IPDB: common.IPDB{
			Data:     &fileInfo,
			IndexLen: 11,
		},
	}
}

func (db ZXwry) Find(ip string) (result string) {
	db.Offset = 0

	tp := big.NewInt(0)
	op := big.NewInt(0)
	tp.SetBytes(net.ParseIP(ip).To16())
	op.SetString("18446744073709551616", 10)
	op.Div(tp, op)
	tp.SetString("FFFFFFFFFFFFFFFF", 16)
	op.And(op, tp)

	ipv6 := op.Uint64()
	offset := db.searchIndex(ipv6)
	country, area := db.getAddr(offset)

	country = strings.ReplaceAll(country, " CZ88.NET", "")
	area = strings.ReplaceAll(area, " CZ88.NET", "")

	return fmt.Sprintf("%s %s", country, area)
}

func (db *ZXwry) getAddr(offset uint32) (string, string) {
	mode := db.ReadMode(offset)
	if mode == common.RedirectMode1 {
		offset = db.ReadUInt24()
		return db.getAddr(offset)
	}
	realOffset := db.Offset - 1
	c1 := db.ReadArea(realOffset)
	if mode == common.RedirectMode2 {
		db.Offset = 4 + realOffset
	} else {
		db.Offset = realOffset + uint32(1+len(c1))
	}
	c2 := db.ReadArea(db.Offset)
	return string(c1), string(c2)
}

func (db *ZXwry) searchIndex(ip uint64) uint32 {
	header := db.ReadData(16, 8)
	start := binary.LittleEndian.Uint32(header[8:])
	counts := binary.LittleEndian.Uint32(header[:8])
	end := start + counts*db.IndexLen

	buf := make([]byte, db.IndexLen)

	for {
		mid := db.GetMiddleOffset(start, end)
		buf = db.ReadData(11, mid)
		ipBytes := binary.LittleEndian.Uint64(buf[:8])

		if end-start == db.IndexLen {
			if ip >= binary.LittleEndian.Uint64(db.ReadData(8, end)) {
				buf = db.ReadData(11, end)
			}
			return common.ByteToUInt32(buf[8:])
		}

		if ipBytes > ip {
			end = mid
		} else if ipBytes < ip {
			start = mid
		} else if ipBytes == ip {
			return common.ByteToUInt32(buf[8:])
		}
	}
}
