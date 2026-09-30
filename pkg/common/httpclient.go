package common

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36"

type HttpClient struct {
	*http.Client
}

var httpClient *HttpClient

func init() {
	httpClient = &HttpClient{http.DefaultClient}
	httpClient.Timeout = time.Second * 300
	httpClient.Transport = &http.Transport{
		TLSHandshakeTimeout:   time.Second * 5,
		IdleConnTimeout:       time.Second * 10,
		ResponseHeaderTimeout: time.Second * 10,
		ExpectContinueTimeout: time.Second * 20,
		Proxy:                 http.ProxyFromEnvironment,
	}
}

func GetHttpClient() *HttpClient {
	c := *httpClient
	return &c
}

// Get downloads the first of urls that answers 200 OK.
func (c *HttpClient) Get(urls ...string) (body []byte, err error) {
	body, _, err = c.GetIfModified(time.Time{}, urls...)
	return body, err
}

// GetIfModified is like Get, but when since is not zero it sends an
// If-Modified-Since header and reports notModified when a server answers
// 304 Not Modified.
func (c *HttpClient) GetIfModified(since time.Time, urls ...string) (body []byte, notModified bool, err error) {
	var req *http.Request
	var resp *http.Response

	err = errors.New("no download url")
	for _, url := range urls {
		req, err = http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			log.Println(err)
			continue
		}
		req.Header.Set("User-Agent", UserAgent)
		if !since.IsZero() {
			req.Header.Set("If-Modified-Since", since.UTC().Format(http.TimeFormat))
		}
		resp, err = c.Do(req)
		if err != nil {
			continue
		}

		if resp.StatusCode == http.StatusNotModified {
			_ = resp.Body.Close()
			return nil, true, nil
		}
		if resp.StatusCode != http.StatusOK {
			_ = resp.Body.Close()
			err = fmt.Errorf("%s: unexpected HTTP status %s", url, resp.Status)
			continue
		}

		body, err = io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			continue
		}
		return body, false, nil
	}

	return nil, false, err
}
