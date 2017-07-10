package main

import (
	"github.com/urfave/cli"
	"pefi/model"
	"strconv"
)

func externalAccountCommand() cli.Command {
	return cli.Command{
		Name:    "external",
		Aliases: []string{"ex", "e"},
		Usage:   "external account interface",
		Subcommands: GetAPISubCmd(
			"/accounts/external",
			new(externalAccount),
			new(externalAccounts),
			createExternalAccount,
			externalAccountFlags,
			nil,
		),
	}
}

type (
	externalAccounts []externalAccount

	externalAccount struct {
		model.ExternalAccount
	}
)

var (
	externalAccountHeader = []string{
		"id",
		"name",
		"description",
		"category",
	}

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

func (es *externalAccounts) Header() (s []string) {
	return externalAccountHeader
}

func (es *externalAccounts) Body() (s [][]string) {
	for _, e := range *es {
		s = append(s, e.Table())
	}
	return s
}

func (es *externalAccounts) Footer() (s []string) {
	return []string{}
}

func (e *externalAccount) Header() (s []string) {
	return externalAccountHeader
}

func (e *externalAccount) Body() (s [][]string) {
	return [][]string{e.Table()}
}

func (e *externalAccount) Footer() (s []string) {
	return []string{}
}

func (a *externalAccount) Table() (s []string) {
	s = []string{
		strconv.Itoa(int(a.ID)),
		a.Name,
		a.Description,
		strconv.Itoa(int(a.CategorieID)),
	}
	return s
}

func createExternalAccount(c *cli.Context) (t tabular, err error) {
	return &externalAccount{
		ExternalAccount: model.ExternalAccount{
			Name:        c.String("name"),
			Description: c.String("description"),
			CategorieID: c.Int64("category"),
		},
	}, nil
}
