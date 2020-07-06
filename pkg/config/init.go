package config

import (
	"os"
	"simple-auth/pkg/config/argparser"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const envPrefix = "sa"

func loadYaml(filename string, config *Config) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(config); err != nil {
		return err
	}

	return nil
}

func loadEnvironment(config *Config) error {
	return envconfig.Process(envPrefix, config)
}

func readConfig() (config *Config) {
	config = &Config{}
	logrus.Info("Loading config...")

	if err := loadYaml("simpleauth.default.yml", config); err != nil {
		logrus.Fatalf("Error loading default config file simpleauth.default.yml: %v", err)
	}
	if err := loadYaml("simpleauth.yml", config); err != nil {
		logrus.Warnf("Unable to load simpleauth.yml: %v", err)
	}

	if err := loadEnvironment(config); err != nil {
		logrus.Warnf("Unable to load environment variables: %v", err)
	}

	if err := argparser.LoadArgs(config, os.Args[1:]...); err != nil {
		logrus.Fatalf("Unable to load arguments: %v", err)
	}

	return
}

// Global configuration
var Global *Config

func init() {
	Global = readConfig()
}
