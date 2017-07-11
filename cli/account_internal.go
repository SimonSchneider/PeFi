package main

import (
	"errors"
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
			new(models.InternalAccount),
			new(models.InternalAccount),
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

func createInternalAccount(c *cli.Context) (t interface{}, err error) {
	tmp, err := createExternalAccount(c)
	if err != nil {
		return nil, err
	}
	exAcc, ok := tmp.(*models.ExternalAccount)
	if !ok {
		return nil, errors.New("not possible to get external account")
	}
	return &models.InternalAccount{
		ExternalAccount: *exAcc,
		Balance:         c.Float64("balance"),
	}, nil
}
