package main

import (
	"errors"
	"fmt"
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"time"

	"github.com/urfave/cli/v2"
)

var cmdOneTime = &cli.Command{
	Name:      "onetime",
	Usage:     "Create one-time use token for an account",
	ArgsUsage: "<email>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "duration",
			Aliases: []string{"d"},
			Value:   "24h",
		},
	},
	Action: funcOneTime,
}

func funcOneTime(c *cli.Context) error {
	email := c.Args().First()
	durationArg := c.String("duration")
	if email == "" {
		return errors.New("Missing email")
	}
	duration, err := time.ParseDuration(durationArg)
	if err != nil {
		return err
	}

	config := config.Load(false)
	db := db.New(config.Db.Driver, config.Db.URL)
	account, err := db.FindAccountByEmail(email)
	if err != nil {
		return fmt.Errorf("Unable to find account for %s: %w", email, err)
	}

	token, err := db.CreateAccountOneTimeToken(account, duration)
	if err != nil {
		return err
	}

	fmt.Printf("Token: %s\n", token)
	fmt.Printf("URL:   %s\n", config.Web.GetBaseURL()+"/api/ui/onetime?token="+token)
	return nil
}
