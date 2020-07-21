package common

import (
	"log"
	"os"
)

func ByteToUInt32(data []byte) uint32 {
	i := uint32(data[0]) & 0xff
	i |= (uint32(data[1]) << 8) & 0xff00
	i |= (uint32(data[2]) << 16) & 0xff0000
	return i
}

func ExistThenRemove(filePath string) {
	_, err := os.Stat(filePath)
	if err == nil {
		err = os.Remove(filePath)
		if err != nil {
			log.Fatalln("旧文件删除失败", err.Error())
			os.Exit(1)
		}
	}
}
