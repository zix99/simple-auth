package main

import (
	"log"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type configDatabase struct {
	Driver string
	URL    string
}

type configWeb struct {
	Host string
}

type config struct {
	Db         configDatabase
	Web        configWeb
	Production bool
}

func readConfig() (config *config) {
	v := viper.New()
	v.SetConfigName("config")
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
