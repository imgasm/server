package image

import (
	"errors"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
)

const imageSourceTypeFileSystem imageSourceType = "fs"

type fileSystemImageSource struct {
	Config *sourceConfig
}

func init() {
	registerImageSource(imageSourceTypeFileSystem, newFileSystemImageSource)
}

func newFileSystemImageSource(config *sourceConfig) imageSource {
	return &fileSystemImageSource{config}
}

func (s *fileSystemImageSource) matches(r *http.Request) bool {
	return (r.Method == "GET" || r.Method == "POST") && r.URL.Query().Get("file") != ""
}

func (s *fileSystemImageSource) getImage(r *http.Request) ([]byte, error) {
	file := r.URL.Query().Get("file")
	if file == "" {
		return nil, errors.New("missing param file")
	}

	file, err := s.buildPath(file)
	if err != nil {
		return nil, err
	}

	return s.read(file)
}

func (s *fileSystemImageSource) buildPath(file string) (string, error) {
	file = path.Clean(path.Join(s.Config.MountPath, file))
	if strings.HasPrefix(file, s.Config.MountPath) == false {
		return "", errors.New("invalid filepath")
	}
	return file, nil
}

func (_ *fileSystemImageSource) read(file string) ([]byte, error) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.New("invalid filepath")
	}
	return buf, nil
}
