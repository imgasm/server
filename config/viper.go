package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// var ViperConfigPath = "/go/src/github.com/imgasm/server/config/"
var ViperConfigPath = "/home/snow/code/go/src/github.com/imgasm/server/config/"
var RuntimeViper = viper.New()

func init() {
	RuntimeViper.SetConfigName("config")
	RuntimeViper.SetConfigType("json")
	RuntimeViper.AddConfigPath(ViperConfigPath)
	err := RuntimeViper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error reading config file: %s \n", err))
	}
}
