package main

import (
	"github.com/urfave/cli"
	"pefi/api/models"
	"reflect"
	"strconv"
)

func labelCommand() cli.Command {
	return cli.Command{
		Name:    "label",
		Aliases: []string{"lab", "l"},
		Usage:   "label interface",
		Subcommands: GetAPISubCmd(
			"/labels",
			new(label),
			new(labels),
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

func createLabel(c *cli.Context) (nl tabular, err error) {
	return model.Label{
		Name:        c.String("name"),
		Description: c.String("description"),
		CategorieID: c.Int64("categorie"),
	}, nil
}
