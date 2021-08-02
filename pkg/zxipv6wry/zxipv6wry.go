package zxipv6wry

import (
	"encoding/binary"
	"errors"
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

	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		log.Println("文件不存在，尝试从网络获取最新ZX IPv6数据库")
		fileData, err = Download(filePath)
		if err != nil {
			os.Exit(1)
		}
	} else {
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
			Data: &fileInfo,
		},
	}
}

func (db ZXwry) Find(query string, params ...string) (result fmt.Stringer, err error) {
	ip := net.ParseIP(query)
	if ip == nil {
		return nil, errors.New("Query should be IPv6")
	}
	ip6 := ip.To16()
	if ip6 == nil {
		return nil, errors.New("Query should be IPv6")
	}

	tp := big.NewInt(0)
	op := big.NewInt(0)
	tp.SetBytes(ip6)
	op.SetString("18446744073709551616", 10)
	op.Div(tp, op)
	tp.SetString("FFFFFFFFFFFFFFFF", 16)
	op.And(op, tp)

	ipv6 := op.Uint64()
	offset := db.searchIndex(ipv6)
	country, area := db.getAddr(offset)

	result = common.Result{
		Country: strings.ReplaceAll(country, " CZ88.NET", ""),
		Area:    strings.ReplaceAll(area, " CZ88.NET", ""),
	}
	return result, nil
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
	end := start + counts*11

	buf := make([]byte, 11)

	for {
		mid := common.GetMiddleOffset(start, end, 11)
		buf = db.ReadData(11, mid)
		ipBytes := binary.LittleEndian.Uint64(buf[:8])

		if end-start == 11 {
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
