package cdn

import (
	"io/ioutil"
	"net/http"
)

func Download() (data []byte, err error) {
	//resp, err := http.Get("https://raw.githubusercontent.com/SukkaLab/cdn/master/dist/cdn.json")
	resp, err := http.Get("https://cdn.jsdelivr.net/gh/SukkaLab/cdn/dist/cdn.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
