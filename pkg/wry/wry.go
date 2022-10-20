package wry

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
)

// IPDB common ip database
type IPDB[T ~uint32 | ~uint64] struct {
	Data []byte

	OffLen   uint8
	IPLen    uint8
	IPCnt    T
	IdxStart T
	IdxEnd   T
}

type Reader struct {
	s []byte
	i uint32 // current reading index
	l uint32 // last reading index

	Result Result
}

func NewReader(data []byte) Reader {
	return Reader{s: data, i: 0, l: 0, Result: Result{
		Country: "",
		Area:    "",
	}}
}

func (r *Reader) seekAbs(offset uint32) {
	r.l = r.i
	r.i = offset
}

func (r *Reader) seek(offset int64) {
	r.l = r.i
	r.i = uint32(int64(r.i) + offset)
}

// seekBack: seek to last index, can only call once
func (r *Reader) seekBack() {
	r.i = r.l
}

func (r *Reader) read(length uint32) []byte {
	rs := make([]byte, length)
	copy(rs, r.s[r.i:])
	r.l = r.i
	r.i += length
	return rs
}

func (r *Reader) readMode() (mode byte) {
	mode = r.s[r.i]
	r.l = r.i
	r.i += 1
	return
}

// readOffset: read 3 bytes as uint32 offset
func (r *Reader) readOffset(follow bool) uint32 {
	buf := r.read(3)
	offset := Bytes3ToUint32(buf)
	if follow {
		r.l = r.i
		r.i = offset
	}
	return offset
}

func (r *Reader) readString(seek bool) string {
	length := bytes.IndexByte(r.s[r.i:], 0)
	str := string(r.s[r.i : r.i+uint32(length)])
	if seek {
		r.l = r.i
		r.i += uint32(length) + 1
	}
	return str
}

type Result struct {
	Country string
	Area    string
}

func (r *Result) DecodeGBK() *Result {
	enc := simplifiedchinese.GBK.NewDecoder()
	r.Country, _ = enc.String(r.Country)
	r.Area, _ = enc.String(r.Area)
	return r
}

func (r *Result) Trim() *Result {
	r.Country = strings.TrimSpace(strings.ReplaceAll(r.Country, "CZ88.NET", ""))
	r.Area = strings.TrimSpace(strings.ReplaceAll(r.Area, "CZ88.NET", ""))
	return r
}

func (r Result) String() string {
	r.Trim()
	return strings.TrimSpace(fmt.Sprintf("%s %s", r.Country, r.Area))
}

func Bytes3ToUint32(data []byte) uint32 {
	i := uint32(data[0]) & 0xff
	i |= (uint32(data[1]) << 8) & 0xff00
	i |= (uint32(data[2]) << 16) & 0xff0000
	return i
}
