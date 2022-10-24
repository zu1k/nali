package download

import (
	"errors"
	"github.com/zu1k/nali/pkg/common"
)

func Download(filePath string, urls ...string) (data []byte, err error) {
	if len(urls) == 0 {
		return nil, errors.New("未指定下载 url")
	}

	data, err = common.GetHttpClient().Get(urls...)
	if err != nil {
		return
	}

	err = common.SaveFile(filePath, data)
	return
}
