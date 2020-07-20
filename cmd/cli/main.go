package main

import (
	"fmt"
	"log"
	"os"
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"syscall"

	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh/terminal"
)

func funcAddUser(c *cli.Context) error {
	config := config.Load(false)

	email := c.Args().Get(0)
	username := c.Args().Get(1)

	fmt.Printf("Email:    %s\n", email)
	fmt.Printf("Username: %s\n", username)

	fmt.Print("Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return err
	}
	fmt.Print("Re-Enter: ")
	bytePassword2, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return err
	}

	pass1 := string(bytePassword)
	pass2 := string(bytePassword2)
	if pass1 != pass2 {
		return fmt.Errorf("Password do not match")
	}

	fmt.Println("Creating account...")
	db := db.New(config.Db.Driver, config.Db.URL)
	account, err := db.CreateAccount(email)
	if err != nil {
		return err
	}
	fmt.Printf("Account %s created\n", account.UUID)

	fmt.Println("Creating simple auth...")
	err = db.CreateAccountAuthSimple(account, username, pass1)
	if err != nil {
		return err
	}

	fmt.Println("Success!")

	return nil
}

func cmdAddUser() *cli.Command {
	return &cli.Command{
		Name:      "adduser",
		Category:  "user",
		Usage:     "Add a new user to simple-auth DB",
		ArgsUsage: "<email> <username>",
		Action:    funcAddUser,
	}
}

func cliMain(args ...string) error {
	app := &cli.App{
		Usage:                  "CLI Tool for simple-auth",
		Description:            `CLI Tool for inspecting, testing, and modifying data for simple-auth`,
		Version:                fmt.Sprintf("%s, %s", version, buildSha),
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			cmdAddUser(),
		},
		Copyright: `simple-auth  Copyright (C) 2020 Chris LaPointe
		This program comes with ABSOLUTELY NO WARRANTY.
		This is free software, and you are welcome to redistribute it
		under certain conditions`,
	}

	return app.Run(args)
}

func main() {
	err := cliMain(os.Args...)
	if err != nil {
		log.Fatal(err)
	}
}
