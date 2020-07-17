package qqwry

import (
	"bytes"
"compress/zlib"
"encoding/binary"
"io/ioutil"
"net/http"
)

// @ref https://zhangzifan.com/update-qqwry-dat.html

func getKey() (uint32, error) {
	resp, err := http.Get("http://update.cz88.net/ip/copywrite.rar")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		return 0, err
	} else {
		// @see https://stackoverflow.com/questions/34078427/how-to-read-packed-binary-data-in-go
		return binary.LittleEndian.Uint32(body[5*4:]), nil
	}
}

func GetOnline() ([]byte, error) {
	resp, err := http.Get("http://update.cz88.net/ip/qqwry.rar")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	} else {
		if key, err := getKey(); err != nil {
			return nil, err
		} else {
			for i := 0; i < 0x200; i++ {
				key = key * 0x805
				key++
				key = key & 0xff

				body[i] = byte(uint32(body[i]) ^ key)
			}

			reader, err := zlib.NewReader(bytes.NewReader(body))
			if err != nil {
				return nil, err
			}

			return ioutil.ReadAll(reader)
		}
	}
}
