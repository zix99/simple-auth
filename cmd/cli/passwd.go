package main

import (
	"errors"
	"fmt"
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"

	"github.com/urfave/cli/v2"
)

var cmdPasswd = &cli.Command{
	Name:      "passwd",
	Category:  "user",
	Usage:     "Change or set password for simple-auth user",
	ArgsUsage: "<username>",
	Action:    funcPasswd,
}

func funcPasswd(c *cli.Context) error {
	username := c.Args().First()

	if username == "" {
		return errors.New("Missing username")
	}

	fmt.Printf("Username: %s\n", username)
	password, err := readPasswordTwice()
	if err != nil {
		return fmt.Errorf("Password error: %w", err)
	}

	// Make the modifications
	fmt.Println("Updating password...")
	config := config.Load(false)
	db := db.New(config.Db.Driver, config.Db.URL)
	err = db.UpdatePasswordForUsername(username, password)
	if err != nil {
		return fmt.Errorf("Unable to update password: %w", err)
	}

	fmt.Println("Done.")

	return nil
}
