package qqwry

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
)

type FileInfo struct {
	Data     []byte
	FilePath string
	FileBase *os.File
	IPNum    int64
}

type QQwry struct {
	data   *FileInfo
	offset int64
}

const (
	IndexLen      = 7
	RedirectMode1 = 0x01
	RedirectMode2 = 0x02
)

func NewQQwry(filePath string) QQwry {
	var tmpData []byte
	var fileInfo FileInfo

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

	fileInfo.IPNum = int64((end-start)/IndexLen + 1)

	return QQwry{
		data: &fileInfo,
	}
}

// setOffset 设置偏移量
func (q *QQwry) setOffset(offset int64) {
	q.offset = offset
}

// Find ip地址查询对应归属地信息
func (q QQwry) Find(ip string) (res string) {
	if strings.Count(ip, ".") != 3 {
		return
	}
	offset := q.searchIndex(binary.BigEndian.Uint32(net.ParseIP(ip).To4()))
	if offset <= 0 {
		return
	}

	var gbkCountry []byte
	var gbkArea []byte

	mode := q.readMode(offset + 4)
	if mode == RedirectMode1 {
		countryOffset := q.readUInt24()
		mode = q.readMode(countryOffset)
		if mode == RedirectMode2 {
			c := q.readUInt24()
			gbkCountry = q.readString(c)
			countryOffset += 4
		} else {
			gbkCountry = q.readString(countryOffset)
			countryOffset += uint32(len(gbkCountry) + 1)
		}
		gbkArea = q.readArea(countryOffset)
	} else if mode == RedirectMode2 {
		countryOffset := q.readUInt24()
		gbkCountry = q.readString(countryOffset)
		gbkArea = q.readArea(offset + 8)
	} else {
		gbkCountry = q.readString(offset + 4)
		gbkArea = q.readArea(offset + uint32(5+len(gbkCountry)))
	}

	enc := simplifiedchinese.GBK.NewDecoder()
	country, _ := enc.String(string(gbkCountry))
	area, _ := enc.String(string(gbkArea))

	return fmt.Sprintf("%s %s", country, area)
}

// readMode 获取偏移值类型
func (q *QQwry) readMode(offset uint32) byte {
	mode := q.readData(1, int64(offset))
	return mode[0]
}

// readArea 读取区域
func (q *QQwry) readArea(offset uint32) []byte {
	mode := q.readMode(offset)
	if mode == RedirectMode1 || mode == RedirectMode2 {
		areaOffset := q.readUInt24()
		if areaOffset == 0 {
			return []byte("")
		}
		return q.readString(areaOffset)
	}
	return q.readString(offset)
}

// readString 获取字符串
func (q *QQwry) readString(offset uint32) []byte {
	q.setOffset(int64(offset))
	data := make([]byte, 0, 30)
	buf := make([]byte, 1)
	for {
		buf = q.readData(1)
		if buf[0] == 0 {
			break
		}
		data = append(data, buf[0])
	}
	return data
}

// searchIndex 查找索引位置
func (q *QQwry) searchIndex(ip uint32) uint32 {
	header := q.readData(8, 0)

	start := binary.LittleEndian.Uint32(header[:4])
	end := binary.LittleEndian.Uint32(header[4:])

	buf := make([]byte, IndexLen)
	mid := uint32(0)
	_ip := uint32(0)

	for {
		mid = q.getMiddleOffset(start, end)
		buf = q.readData(IndexLen, int64(mid))
		_ip = binary.LittleEndian.Uint32(buf[:4])

		if end-start == IndexLen {
			offset := byteToUInt32(buf[4:])
			buf = q.readData(IndexLen)
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
			return byteToUInt32(buf[4:])
		}
	}
}

// readUInt24
func (q *QQwry) readUInt24() uint32 {
	buf := q.readData(3)
	return byteToUInt32(buf)
}

// getMiddleOffset
func (q *QQwry) getMiddleOffset(start uint32, end uint32) uint32 {
	records := ((end - start) / IndexLen) >> 1
	return start + records*IndexLen
}

// byteToUInt32 将 byte 转换为uint32
func byteToUInt32(data []byte) uint32 {
	i := uint32(data[0]) & 0xff
	i |= (uint32(data[1]) << 8) & 0xff00
	i |= (uint32(data[2]) << 16) & 0xff0000
	return i
}

// readData 从文件中读取数据
func (q *QQwry) readData(num int, offset ...int64) (rs []byte) {
	if len(offset) > 0 {
		q.setOffset(offset[0])
	}
	nums := int64(num)
	end := q.offset + nums
	dataNum := int64(len(q.data.Data))
	if q.offset > dataNum {
		return nil
	}

	if end > dataNum {
		end = dataNum
	}
	rs = q.data.Data[q.offset:end]
	q.offset = end
	return
}
