package image

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const imageSourceTypeHttp imageSourceType = "http"

type httpImageSource struct {
	Config *sourceConfig
}

func init() {
	registerImageSource(imageSourceTypeHttp, newHttpImageSource)
}

func newHttpImageSource(config *sourceConfig) imageSource {
	return &httpImageSource{config}
}

func (_ *httpImageSource) matches(r *http.Request) bool {
	fmt.Println("method:", r.Method)
	fmt.Println("url query:", r.URL.Query().Get("url"))
	return (r.Method == "GET" || r.Method == "POST") && r.URL.Query().Get("url") != ""
}

func (s *httpImageSource) getImage(r *http.Request) ([]byte, error) {
	url, err := parseURL(r)
	if err != nil {
		return nil, err
	}
	return s.fetchImage(url)
}

func parseURL(r *http.Request) (*url.URL, error) {
	queryUrl := r.URL.Query().Get("url")
	return url.Parse(queryUrl)
}

func (s *httpImageSource) fetchImage(url *url.URL) ([]byte, error) {
	if s.Config.MaxAllowedSize > 0 {
		req := newHTTPRequest(s, "HEAD", url)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("Error fetching image http headers: %v", err)
		}
		res.Body.Close()
		if res.StatusCode < 200 && res.StatusCode > 206 {
			return nil, fmt.Errorf("Error fetching image http headers: (status=%d) (url=%s)", res.StatusCode, req.URL.String())
		}

		contentLength, _ := strconv.Atoi(res.Header.Get("Content-Length"))
		if contentLength > s.Config.MaxAllowedSize {
			return nil, fmt.Errorf("Content-Length %d exceeds maximum allowed %d bytes", contentLength, s.Config.MaxAllowedSize)
		}
	}

	req := newHTTPRequest(s, "GET", url)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error downloading image: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Error downloading image: (status=%d) (url=%s)", res.StatusCode, req.URL.String())
	}

	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Unable to create image from response body: %s (url=%s)", req.URL.String(), err)
	}
	return buf, nil
}

func newHTTPRequest(s *httpImageSource, method string, url *url.URL) *http.Request {
	req, _ := http.NewRequest(method, url.String(), nil)
	req.Header.Set("User-Agent", "imgasm")
	req.URL = url
	return req
}
