package common

import (
	"bytes"
	"math/rand"
	"regexp"
	"testing"
	"time"
	"unsafe"
)

var (
	n     = 100
	lines = make([][]byte, n)
	d     = []string{"\r", "\n", "\r\n"}
)

func init() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < n; i++ {
		lines[i] = []byte(RandStr(rand.Intn(50)) + d[rand.Intn(3)] + RandStr(rand.Intn(50)))
	}
}

func BenchmarkIndexByteTwice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, line := range lines {
			_ = bytes.IndexByte(line, '\n')
			_ = bytes.IndexByte(line, '\r')
		}
	}
}

func BenchmarkLastIndexByteTwice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, line := range lines {
			_ = bytes.LastIndexByte(line, '\n')
			_ = bytes.LastIndexByte(line, '\r')
		}
	}
}

func BenchmarkIndexAny(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, line := range lines {
			_ = bytes.IndexAny(line, "\r\n")
		}
	}
}

var newlineReg = regexp.MustCompile(`\r?\n|\r\n?`)

func BenchmarkRegexFindIndex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, line := range lines {
			_ = newlineReg.FindIndex(line)
		}
	}
}

//////////////////////////////////////////////

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var src = rand.NewSource(time.Now().UnixNano())

const (
	// 6 bits to represent a letter index
	letterIdBits = 6
	// All 1-bits as many as letterIdBits
	letterIdMask = 1<<letterIdBits - 1
	letterIdMax  = 63 / letterIdBits
)

func RandStr(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdMax
		}
		if idx := int(cache & letterIdMask); idx < len(letters) {
			b[i] = letters[idx]
			i--
		}
		cache >>= letterIdBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}
