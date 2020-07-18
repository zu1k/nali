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
	var fileInfo common.FileInfo

	// 判断文件是否存在
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		log.Println("文件不存在，尝试从网络获取最新ZX IPv6 库")
		tmpData, err = GetOnline()
		if err != nil {
			panic(err)
		} else {
			if err := ioutil.WriteFile(filePath, tmpData, 0644); err == nil {
				log.Printf("已将最新的ZX IPv6 库保存到本地 %s ", filePath)
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

	return ZXwry{
		IPDB: common.IPDB{
			Data:     &fileInfo,
			IndexLen: 7,
		},
	}
}

func (q ZXwry) Find(ip string) (result string) {
	q.Offset = 0

	tp := big.NewInt(0)
	op := big.NewInt(0)
	tp.SetBytes(net.ParseIP(ip).To16())
	op.SetString("18446744073709551616", 10)
	op.Div(tp, op)
	tp.SetString("FFFFFFFFFFFFFFFF", 16)
	op.And(op, tp)

	v6ip := op.Uint64()
	offset := q.searchIndex(v6ip)

	country, area := q.getAddr(offset)

	if area == "ZX" {
		area = ""
	}

	return fmt.Sprintf("%s %s", country, area)
}

func (q *ZXwry) getAddr(offset uint32) (string, string) {
	mode := q.ReadMode(offset)
	if mode == common.RedirectMode1 {
		// [IP][0x01][国家和地区信息的绝对偏移地址]
		offset = q.ReadUInt24()
		return q.getAddr(offset)
	}
	// [IP][0x02][信息的绝对偏移][...] or [IP][国家][...]
	_offset := q.Offset - 1
	c1 := q.ReadArea(_offset)
	if mode == common.RedirectMode2 {
		q.Offset = 4 + _offset
	} else {
		q.Offset = _offset + uint32(1+len(c1))
	}
	c2 := q.ReadArea(q.Offset)
	return string(c1), string(c2)
}

func (q *ZXwry) searchIndex(ip uint64) uint32 {
	q.ItemLen = 8
	q.IndexLen = 11

	header := q.Data.Data[8:24]
	start := binary.LittleEndian.Uint32(header[8:])
	counts := binary.LittleEndian.Uint32(header[:8])
	end := start + counts*q.IndexLen

	buf := make([]byte, q.IndexLen)

	for {
		mid := q.GetMiddleOffset(start, end)
		buf = q.Data.Data[mid : mid+q.IndexLen]
		_ip := binary.LittleEndian.Uint64(buf[:q.ItemLen])

		if end-start == q.IndexLen {
			if ip >= binary.LittleEndian.Uint64(q.Data.Data[end:end+q.ItemLen]) {
				buf = q.Data.Data[end : end+q.IndexLen]
			}
			return common.ByteToUInt32(buf[q.ItemLen:])
		}

		if _ip > ip {
			end = mid
		} else if _ip < ip {
			start = mid
		} else if _ip == ip {
			return common.ByteToUInt32(buf[q.ItemLen:])
		}
	}
}
