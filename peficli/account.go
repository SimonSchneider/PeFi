package peficli

import (
	"github.com/urfave/cli"
)

var (
	accAddFlags = []cli.Flag{
		cli.StringFlag{
			Name:  "file, f",
			Usage: "add from file",
		},
		cli.StringFlag{
			Name:  "name,n",
			Usage: "Name of account",
		},
		cli.StringFlag{
			Name:  "description,d",
			Usage: "Name of account",
		},
		cli.Int64SliceFlag{
			Name:  "labels, l",
			Usage: "add a label to the account",
		},
		cli.Float64Flag{
			Name:  "balance, b",
			Usage: "set the balance of the account",
		},
	}
	accLsFlags = []cli.Flag{
		cli.BoolFlag{
			Name:  "json, j",
			Usage: "print in json format",
		},
	}
)

func AccountCommand() cli.Command {
	subcmds := []cli.Command{
		ExternalAccountCommand,
		InternalAccountCommand,
	}

	return cli.Command{
		Name:        "account",
		Aliases:     []string{"acc", "a"},
		Usage:       "account interface",
		Subcommands: subcmds,
	}
}
