package main

import (
	"github.com/urfave/cli"
	"pefi/api/models"
)

func categoryCommand() cli.Command {
	return cli.Command{
		Name:    "category",
		Aliases: []string{"cat", "c"},
		Usage:   "category interface",
		Subcommands: GetAPISubCmd(
			"/categories",
			new(models.Category),
			new([]models.Category),
			createCategory,
			categoryFlags,
			nil,
		),
	}
}

var (
	categoryFlags = APIFlags{
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

func createCategory(c *cli.Context) (nc interface{}, err error) {
	return &models.Category{
		Name:        c.String("name"),
		Description: c.String("description"),
		//ChildrenIds: c.Int64Slice("children"),
	}, nil
}
