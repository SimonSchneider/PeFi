package peficli

import (
	"fmt"
	"github.com/urfave/cli"
)

type (
	transactionFlags struct{}
	lsFlags          struct {
		Days int
	}
	addFlags struct {
		DontCommit bool
	}
	getFlags struct{}
	delFlags struct{}
)

//TransactionCommand return the urfave.cli command of the transaction
func TransactionCommand() cli.Command {
	var (
		lsF  lsFlags
		addF addFlags
		delF delFlags
		getF getFlags
	)

	subcmds := []cli.Command{
		{
			Name:  "ls",
			Usage: "print all transactions",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:        "days, d",
					Usage:       "only show transactions from the last `days`",
					Destination: &lsF.Days,
				},
			},
			Action: func(c *cli.Context) {
				listTransactions(c.Args(), lsF)
			},
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add a transaction",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:        "dont-commit",
					Usage:       "dont-commit addition",
					Destination: &addF.DontCommit,
				},
			},
			Action: func(c *cli.Context) error {
				return addTransaction(c.Args(), addF)
			},
		},
		{
			Name:    "rm",
			Aliases: []string{"d"},
			Usage:   "delete a transaction",
			Action: func(c *cli.Context) error {
				return delTransaction(c.Args(), delF)
			},
		},
		{
			Name:    "get",
			Aliases: []string{"d"},
			Usage:   "get a transaction",
			Action: func(c *cli.Context) error {
				return getTransaction(c.Args(), getF)
			},
		},
	}
	return cli.Command{
		Name:        "transaction",
		Aliases:     []string{"t"},
		Usage:       "transaction interface",
		Subcommands: subcmds,
	}
}

func listTransactions(args []string, flags lsFlags) error {
	fmt.Printf("Listing transactions with flags %+v\n", flags)
	return nil
}

func addTransaction(args []string, flags addFlags) error {
	fmt.Printf("Adding transaction \"%s\" with flags %+v\n", args[0], flags)
	//for _, a := range args {
	//fmt.Printf("%s\n", a)
	//}
	return nil
}

func delTransaction(args []string, flags delFlags) error {
	fmt.Printf("Deleting transaction \"%s\" with flags %+v\n", args[0], flags)
	return nil
}

func getTransaction(args []string, flags getFlags) error {
	fmt.Printf("Getting transaction \"%s\" with flags %+v\n", args[0], flags)
	return nil
}
