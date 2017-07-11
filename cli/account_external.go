package main

import (
	"github.com/urfave/cli"
	"pefi/api/models"
	"strconv"
)

func externalAccountCommand() cli.Command {
	return cli.Command{
		Name:    "external",
		Aliases: []string{"ex", "e"},
		Usage:   "external account interface",
		Subcommands: GetAPISubCmd(
			"/accounts/external",
			new(models.externalAccount),
			new(models.externalAccounts),
			createExternalAccount,
			externalAccountFlags,
			nil,
		),
	}
}

var (
	externalAccountFlags = APIFlags{
		AddFlags: []cli.Flag{
			cli.StringFlag{
				Name:  "name,n",
				Usage: "Name of account",
			},
			cli.StringFlag{
				Name:  "description,d",
				Usage: "Name of account",
			},
			cli.Int64Flag{
				Name:  "category, c",
				Usage: "add a label to the account",
			},
		},
	}
)

func createExternalAccount(c *cli.Context) (t tabular, err error) {
	return &models.externalAccount{
		Name:        c.String("name"),
		Description: c.String("description"),
		CategorieID: c.Int64("category"),
	}, nil
}
