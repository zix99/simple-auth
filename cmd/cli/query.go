package main

import (
	"fmt"
	"simple-auth/pkg/db"

	"github.com/urfave/cli/v2"
)

var cmdQuery = &cli.Command{
	Name:  "query",
	Usage: "Query information from DB",
	Subcommands: []*cli.Command{
		&cli.Command{
			Name:   "accounts",
			Action: funcQueryAccounts,
		},
	},
}

func funcQueryAccounts(c *cli.Context) error {
	sadb := getDB()

	return sadb.GetAllAccounts(func(account *db.Account) bool {
		fmt.Printf("%s\t%s\t%s\t%s\n", account.UUID, account.Email, account.Name, account.CreatedAt)
		return false
	})
}
