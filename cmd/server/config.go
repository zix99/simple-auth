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

type configAuthenticator struct {
	Enabled bool
}

type configAuthencatorSet struct {
	Exchange struct {
		configAuthenticator
	}
}

type configWeb struct {
	Host     string
	Metadata map[string]interface{}
}

type config struct {
	Db             configDatabase
	Web            configWeb
	Authenticators configAuthencatorSet
	Production     bool
}

func readConfig() (config *config) {
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
