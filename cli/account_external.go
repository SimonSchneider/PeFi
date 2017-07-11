package main

import (
	"github.com/urfave/cli"
	"pefi/api/models"
)

func externalAccountCommand() cli.Command {
	return cli.Command{
		Name:    "external",
		Aliases: []string{"ex", "e"},
		Usage:   "external account interface",
		Subcommands: GetAPISubCmd(
			"/accounts/external",
			new(models.ExternalAccount),
			new([]models.ExternalAccount),
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

func createExternalAccount(c *cli.Context) (t interface{}, err error) {
	return &models.ExternalAccount{
		Name:        c.String("name"),
		Description: c.String("description"),
		CategoryID:  c.Int64("category"),
	}, nil
}
