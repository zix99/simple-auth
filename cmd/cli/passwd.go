package main

import (
	"errors"
	"fmt"

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
		return errors.New("missing username")
	}

	fmt.Printf("Username: %s\n", username)
	password, err := readPasswordTwice()
	if err != nil {
		return fmt.Errorf("password error: %w", err)
	}

	// Make the modifications
	fmt.Println("Updating password...")
	db := getDB()
	err = db.UpdatePasswordForUsername(username, password)
	if err != nil {
		return fmt.Errorf("unable to update password: %w", err)
	}

	fmt.Println("Done.")

	return nil
}
