package main

import (
	"fmt"
	"io/ioutil"
	"simple-auth/pkg/box"
	"simple-auth/pkg/config"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

var cmdConfig = &cli.Command{
	Name:  "config",
	Usage: "See default config",
	Subcommands: []*cli.Command{
		{
			Name:   "dump-default",
			Usage:  "Dump defualt configuration (embedded)",
			Action: funcDumpDefaultConfig,
		},
		{
			Name:   "dump",
			Usage:  "Dump current flattened configuration",
			Action: funcDumpCurrentConfig,
		},
	},
}

func funcDumpDefaultConfig(c *cli.Context) error {
	f, err := box.Read("simpleauth.default.yml")
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil
	}

	fmt.Print(string(b))
	return nil
}

func funcDumpCurrentConfig(c *cli.Context) error {
	cfg := config.Load()
	b, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	fmt.Print(string(b))
	return nil
}
