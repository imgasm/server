package image

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

const imageSourceTypeBody imageSourceType = "payload"
const maxMemory int64 = 1024 * 1024 * 64

type bodyImageSource struct {
	Config *sourceConfig
}

func init() {
	registerImageSource(imageSourceTypeBody, newBodyImageSource)
}

func newBodyImageSource(config *sourceConfig) imageSource {
	return &bodyImageSource{config}
}

func (_ *bodyImageSource) matches(r *http.Request) bool {
	return r.Method == "POST" || r.Method == "PUT"
}

func (_ *bodyImageSource) getImage(r *http.Request) ([]byte, error) {
	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/") {
		return readFormBody(r)
	}
	return readRawBody(r)
}

func readFormBody(r *http.Request) ([]byte, error) {
	err := r.ParseMultipartForm(maxMemory)
	if err != nil {
		return nil, err
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if len(buf) == 0 {
		err = errors.New("empty body")
	}

	return buf, err
}

func readRawBody(r *http.Request) ([]byte, error) {
	return ioutil.ReadAll(r.Body)
}
