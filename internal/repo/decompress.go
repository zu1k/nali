package repo

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
)

var (
	fileTypes = []struct {
		ext        string
		decompress func(src io.Reader) (io.Reader, error)
	}{
		{".zip", unzip},
		{".gz", gunzip},
	}
)

func decompress(src io.Reader, fileName string) (io.Reader, error) {
	for _, fileType := range fileTypes {
		if strings.HasSuffix(fileName, fileType.ext) {
			return fileType.decompress(src)
		}
	}
	return nil, fmt.Errorf("decompression algorithm not implemented")
}

func unzip(src io.Reader) (io.Reader, error) {
	// Zip format requires its file size for Decompressing.
	// So we need to read the HTTP response into a buffer at first.
	buf, err := io.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("cannot decompress zip file: %v", err)
	}

	r := bytes.NewReader(buf)
	z, err := zip.NewReader(r, r.Size())
	if err != nil {
		return nil, fmt.Errorf("cannot decompress zip file: %v", err)
	}

	for _, file := range z.File {
		if !file.FileInfo().IsDir() {
			return file.Open()
		}
	}

	return nil, fmt.Errorf("executable not found in zip file")
}

func gunzip(src io.Reader) (io.Reader, error) {
	r, err := gzip.NewReader(src)
	if err != nil {
		return nil, fmt.Errorf("cannot decompress gzip file: %v", err)
	}

	return r, nil
}
