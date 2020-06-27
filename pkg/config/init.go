package config

import (
	"os"

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

func readConfig() (config *Config) {
	config = &Config{}
	logrus.Info("Loading config...")

	if err := loadYaml("simpleauth.default.yml", config); err != nil {
		logrus.Fatalf("Error loading default config file simpleauth.default.yml: %v", err)
	}
	if err := loadYaml("simpleauth.yml", config); err != nil {
		logrus.Warnf("Unable to load simpleauth.yml: %v", err)
	}

	err := envconfig.Process(envPrefix, config)
	if err != nil {
		logrus.Warnf("Unable to load environment variables: %v", err)
	}
	return
}

// Global configuration
var Global *Config

func init() {
	Global = readConfig()
}
