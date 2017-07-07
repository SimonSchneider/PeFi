package main

import (
	"github.com/urfave/cli"
)

func accountCommand() cli.Command {
	return cli.Command{
		Name:    "account",
		Aliases: []string{"acc", "a"},
		Usage:   "account interface",
		Subcommands: []cli.Command{
			externalAccountCommand(),
			internalAccountCommand(),
		},
	}
}
