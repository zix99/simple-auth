package main

import (
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"
)

var cmdStipulation = &cli.Command{
	Name:  "stipulation",
	Usage: "Modify stipulations on an account",
	Subcommands: []*cli.Command{
		cmdStipulationRemoveAll,
	},
}

var cmdStipulationRemoveAll = &cli.Command{
	Name:      "removeall",
	ArgsUsage: "<email>",
	Action:    funcRemoveAllStipulations,
}

func funcRemoveAllStipulations(c *cli.Context) error {
	email := c.Args().First()

	if email == "" {
		return errors.New("Missing email")
	}

	db := getDB()

	account, err := db.FindAccountByEmail(email)
	if err != nil {
		return err
	}

	fmt.Printf("Account: %s\n", account.UUID)

	err = db.ForceSatisfyStipulations(account)
	if err != nil {
		return err
	}

	fmt.Println("Done.")
	return nil
}
