package config

import (
	"fmt"
	"os"
	"simple-auth/pkg/box"
	"simple-auth/pkg/config/argparser"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const envPrefix = "sa"

func loadYaml(config *Config, filename string) error {
	f, err := box.Global.ReadEx(filename, true)
	if err != nil {
		logrus.Warnf("Unable to open config %s: %v", filename, err)
		return err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(config); err != nil {
		logrus.Errorf("Unable to parse config %s: %v", filename, err)
		return err
	}

	logrus.Infof("Loaded config file %s", filename)

	// recurse into include files
	processIncludes(config, false)

	return nil
}

func processIncludes(config *Config, abortOnError bool) error {
	if len(config.Include) > 0 {
		includes := config.Include
		config.Include = nil
		for _, fn := range includes {
			if err := loadYaml(config, fn); err != nil && abortOnError {
				return fmt.Errorf("Unable to load %s: %v", fn, err)
			}
		}
	}
	return nil
}

func loadEnvironment(config *Config) error {
	return envconfig.Process(envPrefix, config)
}

func readConfig(parseArgs bool) (config *Config) {
	config = &Config{}
	logrus.Info("Loading config...")

	if err := loadYaml(config, "simpleauth.default.yml"); err != nil {
		logrus.Fatalf("Error loading default config file simpleauth.default.yml: %v", err)
	}

	if err := loadEnvironment(config); err != nil {
		logrus.Warnf("Unable to load environment variables: %v", err)
		if err := processIncludes(config, true); err != nil {
			logrus.Fatal(err)
		}
	}

	if parseArgs {
		if err := argparser.LoadArgs(config, os.Args[1:]...); err != nil {
			logrus.Fatalf("Unable to load arguments: %v", err)
		}
		if err := processIncludes(config, true); err != nil {
			logrus.Fatal(err)
		}
	}

	return
}

// Global configuration
var config *Config

// Global reads global configuration
func Load(parseArgs bool) *Config {
	if config == nil {
		config = readConfig(parseArgs)
	}
	return config
}
