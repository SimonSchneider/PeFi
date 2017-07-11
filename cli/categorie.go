package main

import (
	"github.com/urfave/cli"
	"pefi/api/models"
	"strconv"
)

func categorieCommand() cli.Command {
	return cli.Command{
		Name:    "categorie",
		Aliases: []string{"cat", "c"},
		Usage:   "categorie interface",
		Subcommands: GetAPISubCmd(
			"/categories",
			new(categorie),
			new(categories),
			createCategorie,
			categorieFlags,
			nil,
		),
	}
}

var (
	categorieFlags = APIFlags{
		AddFlags: []cli.Flag{
			cli.StringFlag{
				Name:  "name,n",
				Usage: "Name of account",
			},
			cli.StringFlag{
				Name:  "description,d",
				Usage: "Name of account",
			},
		},
	}
)

func createCategorie(c *cli.Context) (nc tabular, err error) {
	return model.Categorie{
		Name:        c.String("name"),
		Description: c.String("description"),
		//ChildrenIds: c.Int64Slice("children"),
	}, nil
}
