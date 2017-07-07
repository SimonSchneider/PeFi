package main

import (
	"fmt"
	"github.com/urfave/cli"
	"pefi/model"
)

func internalAccountCommand() cli.Command {
	return cli.Command{
		Name:    "internal",
		Aliases: []string{"in", "i"},
		Usage:   "internal account interface",
		Subcommands: GetAPISubCmd(
			"/accounts/internal",
			new(internalAccount),
			new(internalAccounts),
			createInternalAccount,
			internalAccountFlags,
		),
	}
}

type (
	internalAccounts []internalAccount

	internalAccount struct {
		model.InternalAccount
	}
)

var (
	internalAccountHeader = append(
		externalAccountHeader,
		"balance",
	)

	internalAccountFlags = APIFlags{
		AddFlags: append([]cli.Flag{
			cli.Float64Flag{
				Name:  "balance, b",
				Usage: "set the balance of the account",
			},
		}, externalAccountFlags.AddFlags...),
	}
)

func (is *internalAccounts) Header() (s []string) {
	return internalAccountHeader
}

func (is *internalAccounts) Body() (s [][]string) {
	for _, i := range *is {
		s = append(s, i.Table())
	}
	return s
}

func (is *internalAccounts) Footer() (s []string) {
	sum := 0.0
	for _, i := range *is {
		sum += i.Balance
	}
	for i := 0; i < len(internalAccountHeader); i++ {
		s = append(s, "")
	}
	s[len(s)-1] = fmt.Sprintf("%.2f", sum)
	s[len(s)-2] = "Total"
	return s
}

func (i *internalAccount) Header() (s []string) {
	return internalAccountHeader
}

func (i *internalAccount) Body() (s [][]string) {
	return [][]string{i.Table()}
}

func (i *internalAccount) Footer() (s []string) {
	return []string{}
}

func (a *internalAccount) Table() (s []string) {
	var tmp externalAccount
	tmp.ExternalAccount = a.ExternalAccount
	s = tmp.Table()
	s = append(s, fmt.Sprintf("%.2f", a.Balance))
	return s
}

func createInternalAccount(c *cli.Context) (t tabular, err error) {
	return &internalAccount{
		InternalAccount: model.InternalAccount{
			ExternalAccount: model.ExternalAccount{
				Name:        c.String("name"),
				Description: c.String("description"),
				LabelIds:    c.Int64Slice("labels"),
			},
			Balance: c.Float64("balance"),
		},
	}, nil
}
