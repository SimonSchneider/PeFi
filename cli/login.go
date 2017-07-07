package main

import (
	"github.com/urfave/cli"
)

type (
	loginFlags struct {
		Name   string
		Passwd string
		Host   string
	}
)

//LoginCommand return the urfave.cli command of the transaction
func loginCommand() cli.Command {
	var ()

	subcmds := []cli.Command{
		{
			Name:  "ls",
			Usage: "print all logins",
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add a login",
		},
		{
			Name:    "use",
			Aliases: []string{"d"},
			Usage:   "set a login as the currently used one",
		},
		{
			Name:    "rm",
			Aliases: []string{"d"},
			Usage:   "delete a login",
		},
	}
	return cli.Command{
		Name:        "login",
		Usage:       "login interface",
		Subcommands: subcmds,
	}
}
