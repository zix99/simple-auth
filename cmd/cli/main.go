package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func cliMain(args ...string) error {
	app := &cli.App{
		Usage:                  "CLI Tool for simple-auth",
		Description:            `CLI Tool for inspecting, testing, and modifying data for simple-auth`,
		Version:                fmt.Sprintf("%s, %s", version, buildSha),
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			cmdAddUser,
			cmdPasswd,
			cmdOneTime,
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
