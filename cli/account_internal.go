package main

import (
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"pefi/api/models"
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
			createGraph,
		),
	}
}

var (
	internalAccountFlags = APIFlags{
		AddFlags: append([]cli.Flag{
			cli.Float64Flag{
				Name:  "balance, b",
				Usage: "set the balance of the account",
			},
		}, externalAccountFlags.AddFlags...),
	}
)

func createInternalAccount(c *cli.Context) (t tabular, err error) {
	tmp, err := createExternalAccount(c)
	if err != nil {
		return nil, err
	}
	exAcc, ok := tmp.(*externalAccount)
	if !ok {
		return nil, errors.New("not possible to get external account")
	}
	return models.InternalAccount{
		ExternalAccount: (*exAcc).ExternalAccount,
		Balance:         c.Float64("balance"),
	}, nil
}
