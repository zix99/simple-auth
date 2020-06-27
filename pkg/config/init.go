package config

import (
	"log"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func readConfig() (config *Config) {
	v := viper.New()
	v.SetConfigName("simpleauth")
	v.AddConfigPath(".")
	v.AddConfigPath("$HOME/")

	v.SetEnvPrefix("sa")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		switch err.(type) {
		default:
			log.Fatal(err)
		case viper.ConfigFileNotFoundError:
			logrus.Warnf("Unable to find config file")
		}
	}

	v.Unmarshal(&config)
	return
}

// Global configuration
var Global *Config

func init() {
	Global = readConfig()
}
