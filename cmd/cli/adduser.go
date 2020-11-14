package main

import (
	"errors"
	"fmt"
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"

	"github.com/urfave/cli/v2"
)

var cmdAddUser = &cli.Command{
	Name:      "adduser",
	Category:  "user",
	Usage:     "Add a new user to simple-auth DB",
	ArgsUsage: "<email> <username>",
	Action:    funcAddUser,
}

func funcAddUser(c *cli.Context) error {
	config := config.Load(false)

	email := c.Args().Get(0)
	username := c.Args().Get(1)
	password := c.Args().Get(2)

	if email == "" || username == "" {
		return errors.New("please specify <email> <username> [password]")
	}

	fmt.Printf("Email:    %s\n", email)
	fmt.Printf("Username: %s\n", username)

	if password == "" {
		readPass, err := readPasswordTwice()
		if err != nil {
			return err
		}
		password = readPass
	}

	fmt.Println("Creating account...")
	db := db.New(config.Db.Driver, config.Db.URL)
	account, err := db.CreateAccount(email)
	if err != nil {
		return err
	}
	fmt.Printf("Account %s created\n", account.UUID)

	fmt.Println("Creating simple auth...")
	_, err = db.CreateAuthLocal(account, username, password)
	if err != nil {
		return err
	}

	fmt.Println("Success!")

	return nil
}
