package main

import (
	"github.com/urfave/cli"
	"pefi/model"
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

type (
	categories []categorie

	categorie struct {
		model.Categorie
	}
)

var (
	categorieHeader = []string{
		"id",
		"name",
		"description",
	}

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

func (cs *categories) Header() (s []string) {
	return categorieHeader
}

func (cs *categories) Body() (s [][]string) {
	for _, c := range *cs {
		s = append(s, c.Table())
	}
	return s
}

func (cs *categories) Footer() (s []string) {
	return []string{}
}

func (c *categorie) Header() (s []string) {
	return categorieHeader
}

func (c *categorie) Body() (s [][]string) {
	return [][]string{c.Table()}
}

func (c *categorie) Footer() (s []string) {
	return []string{}
}

func (c *categorie) Table() (s []string) {
	s = []string{
		strconv.Itoa(int(c.ID)),
		c.Name,
		c.Description,
	}
	//labelIds := []string{}
	//for _, id := range c.LabelIds {
	//labelIds = append(labelIds, strconv.Itoa(int(id)))
	//}
	//s = append(s, strings.Join(labelIds, ","))
	//childIds := []string{}
	//for _, id := range c.ChildrenIds {
	//childIds = append(childIds, strconv.Itoa(int(id)))
	//}
	//s = append(s, strings.Join(childIds, ","))
	return s
}

func createCategorie(c *cli.Context) (nc tabular, err error) {
	return &categorie{
		Categorie: model.Categorie{
			Name:        c.String("name"),
			Description: c.String("description"),
			//LabelIds:    c.Int64Slice("labels"),
			//ChildrenIds: c.Int64Slice("children"),
		},
	}, nil
}
