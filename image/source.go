package image

import (
	"fmt"
	"github.com/imgasm/server/config"
	"github.com/spf13/viper"
	"net/http"
	// "net/url"
)

// todo: get rid of factory function and make it more readable

type imageSourceType string
type imageSourceFactoryFunc func(*sourceConfig) imageSource

type sourceConfig struct {
	Type           imageSourceType
	MountPath      string
	AuthForwarding bool
	Authorization  string
	// AllowedOrigings []*url.URL
	MaxAllowedSize int
}

var imageSourceMap = make(map[imageSourceType]imageSource)
var imageSourceFactoryMap = make(map[imageSourceType]imageSourceFactoryFunc)

type imageSource interface {
	matches(*http.Request) bool
	getImage(*http.Request) ([]byte, error)
}

func registerImageSource(sourceType imageSourceType, factory imageSourceFactoryFunc) {
	imageSourceFactoryMap[sourceType] = factory
}

func matchImageSource(r *http.Request) imageSource {
	for _, source := range imageSourceMap {
		if source.matches(r) {
			return source
		}
	}
	return nil
}

func LoadImageSources() {
	viper.SetConfigName("config")
	viper.AddConfigPath(config.ViperConfigPath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error reading config file: %s \n", err))
	}
	for name, factory := range imageSourceFactoryMap {
		imageSourceMap[name] = factory(&sourceConfig{
			Type:           name,
			MountPath:      viper.GetString("imgsource.mountpath"),
			AuthForwarding: viper.GetBool("imgsource.authforwarding"),
			Authorization:  viper.GetString("imgsource.authorization"),
			// AllowedOrigings: viper.GetString("imgsource.allowedorigins"),
			MaxAllowedSize: viper.GetInt("imgsource.maxallowedsize"),
		})
	}
}
