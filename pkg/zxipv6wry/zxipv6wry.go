package zxipv6wry

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"os"

	"github.com/zu1k/nali/pkg/common"
)

type ZXwry struct {
	common.IPDB
}

func NewZXwry(filePath string) ZXwry {
	var tmpData []byte
	var fileInfo common.FileData

	// 判断文件是否存在
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		log.Println("文件不存在，请自行下载 ZX IPV6库，并保存在", filePath)
		os.Exit(1)
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

	return ZXwry{
		IPDB: common.IPDB{
			Data:     &fileInfo,
			IndexLen: 7,
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
	return fmt.Sprintf("%s %s", country, area)
}

func (db *ZXwry) getAddr(offset uint32) (string, string) {
	mode := db.ReadMode(offset)
	if mode == common.RedirectMode1 {
		offset = db.ReadUInt24()
		return db.getAddr(offset)
	}
	_offset := db.Offset - 1
	c1 := db.ReadArea(_offset)
	if mode == common.RedirectMode2 {
		db.Offset = 4 + _offset
	} else {
		db.Offset = _offset + uint32(1+len(c1))
	}
	c2 := db.ReadArea(db.Offset)
	return string(c1), string(c2)
}

func (db *ZXwry) searchIndex(ip uint64) uint32 {
	db.ItemLen = 8
	db.IndexLen = 11

	header := db.Data.Data[8:24]
	start := binary.LittleEndian.Uint32(header[8:])
	counts := binary.LittleEndian.Uint32(header[:8])
	end := start + counts*db.IndexLen

	buf := make([]byte, db.IndexLen)

	for {
		mid := db.GetMiddleOffset(start, end)
		buf = db.Data.Data[mid : mid+db.IndexLen]
		_ip := binary.LittleEndian.Uint64(buf[:db.ItemLen])

		if end-start == db.IndexLen {
			if ip >= binary.LittleEndian.Uint64(db.Data.Data[end:end+db.ItemLen]) {
				buf = db.Data.Data[end : end+db.IndexLen]
			}
			return common.ByteToUInt32(buf[db.ItemLen:])
		}

		if _ip > ip {
			end = mid
		} else if _ip < ip {
			start = mid
		} else if _ip == ip {
			return common.ByteToUInt32(buf[db.ItemLen:])
		}
	}
}
