package zxipv6wry

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"net"

	"github.com/zu1k/nali/pkg/common"
)

var (
	header []byte
	v6ip   uint64
	offset uint32
	start  uint32
	end    uint32
)

type ZXwry struct {
	common.IPDB
}

func (q *ZXwry) Find(ip string) (result string) {
	q.Offset = 0

	tp := big.NewInt(0)
	op := big.NewInt(0)
	tp.SetBytes(net.ParseIP(ip).To16())
	op.SetString("18446744073709551616", 10)
	op.Div(tp, op)
	tp.SetString("FFFFFFFFFFFFFFFF", 16)
	op.And(op, tp)

	v6ip = op.Uint64()
	offset = q.searchIndexV6(v6ip)

	country, area := q.getAddr(offset)

	if area == "ZX" {
		area = ""
	}

	return fmt.Sprintf("%s %s", country, area)
}

func (q *ZXwry) getAddr(offset uint32) (string, string) {
	mode := q.ReadMode(offset)
	if mode == 0x01 {
		// [IP][0x01][国家和地区信息的绝对偏移地址]
		offset = q.ReadUInt24()
		return q.getAddr(offset)
	}
	// [IP][0x02][信息的绝对偏移][...] or [IP][国家][...]
	_offset := q.Offset - 1
	c1 := q.ReadArea(_offset)
	if mode == 0x02 {
		q.Offset = 4 + _offset
	} else {
		q.Offset = _offset + uint32(1+len(c1))
	}
	c2 := q.ReadArea(q.Offset)
	return string(c1), string(c2)
}

func (q *ZXwry) searchIndexV4(ip uint32) uint32 {
	q.ItemLen = 4
	q.IndexLen = 7
	header = q.Data.Data[0:8]
	start = binary.LittleEndian.Uint32(header[:4])
	end = binary.LittleEndian.Uint32(header[4:])

	buf := make([]byte, q.IndexLen)

	for {
		mid := start + q.IndexLen*(((end-start)/q.IndexLen)>>1)
		buf = q.Data.Data[mid : mid+q.IndexLen]
		_ip := binary.LittleEndian.Uint32(buf[:q.ItemLen])

		if end-start == q.IndexLen {
			if ip >= binary.LittleEndian.Uint32(q.Data.Data[end:end+q.ItemLen]) {
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

func (q *ZXwry) searchIndexV6(ip uint64) uint32 {

	q.ItemLen = 8
	q.IndexLen = 11

	header = q.Data.Data[8:24]
	start = binary.LittleEndian.Uint32(header[8:])
	counts := binary.LittleEndian.Uint32(header[:8])
	end = start + counts*q.IndexLen

	buf := make([]byte, q.IndexLen)

	for {
		mid := start + q.IndexLen*(((end-start)/q.IndexLen)>>1)
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
