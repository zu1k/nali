package common

import (
	"io/ioutil"
	"net/http"
	"time"
)

const UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36"

type HttpClient struct {
	*http.Client
}

var httpClient *HttpClient

func init() {
	httpClient = &HttpClient{http.DefaultClient}
	httpClient.Timeout = time.Second * 30
	httpClient.Transport = &http.Transport{
		TLSHandshakeTimeout:   time.Second * 5,
		IdleConnTimeout:       time.Second * 20,
		ResponseHeaderTimeout: time.Second * 20,
		ExpectContinueTimeout: time.Second * 20,
	}
}

func GetHttpClient() *HttpClient {
	c := *httpClient
	return &c
}

func (c *HttpClient) Get(urls ...string) (body []byte, err error) {
	var req *http.Request
	var resp *http.Response

	for _, url := range urls {
		req, err = http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			continue
		}
		req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
		req.Header.Set("User-Agent", UserAgent)
		resp, err = c.Do(req)

		if err == nil && resp != nil && resp.StatusCode == 200 {
			defer resp.Body.Close()
			body, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				continue
			}
			return
		}
	}

	return nil, err
}
