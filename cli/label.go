package main

import (
	"github.com/urfave/cli"
	"pefi/api/models"
)

func labelCommand() cli.Command {
	return cli.Command{
		Name:    "label",
		Aliases: []string{"lab", "l"},
		Usage:   "label interface",
		Subcommands: GetAPISubCmd(
			"/labels",
			new(models.Label),
			new([]models.Label),
			createLabel,
			labelFlags,
			nil,
		),
	}
}

var (
	labelFlags = APIFlags{
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
				Name:  "categorie,c",
				Usage: "Name of categorie",
			},
		},
	}
)

func createLabel(c *cli.Context) (nl interface{}, err error) {
	return &models.Label{
		Name:        c.String("name"),
		Description: c.String("description"),
		CategoryID:  c.Int64("categorie"),
	}, nil
}
