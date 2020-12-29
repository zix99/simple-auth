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
	decoder.SetStrict(true)
	if err := decoder.Decode(config); err != nil {
		logrus.Errorf("Unable to parse config %s: %v", filename, err)
		return err
	}

	logrus.Infof("Loaded config file %s", filename)

	return processIncludes(config, false)
}

func processIncludes(config *Config, abortOnNotFound bool) error {
	if len(config.Include) > 0 {
		includes := config.Include
		config.Include = nil
		for _, fn := range includes {
			if err := loadYaml(config, fn); err != nil {
				if err == os.ErrNotExist {
					if abortOnNotFound {
						return fmt.Errorf("unable to load %s: %w", fn, err)
					}
				} else {
					return fmt.Errorf("unable to load %s: %w", fn, err)
				}
			}
		}
	}
	return nil
}

func loadEnvironment(config *Config) error {
	return envconfig.Process(envPrefix, config)
}

func readConfig(args []string) (config *Config) {
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

	if len(args) > 0 {
		if err := argparser.LoadArgs(config, args...); err != nil {
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
func Load(args ...string) *Config {
	if config == nil {
		config = readConfig(args)
	}
	return config
}
