package peficli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	"net/http"
	"pefi/model"
	"time"
)

var (
	traAddFlags = []cli.Flag{
		cli.StringFlag{
			Name:  "file, f",
			Usage: "add from file",
		},
		cli.StringFlag{
			Name:  "time,t",
			Value: time.Now().Format(time.RFC3339),
			Usage: "Time of transaction",
		},
		cli.Float64Flag{
			Name:  "amount,a",
			Usage: "Transaction amount",
		},
		cli.Int64Flag{
			Name:  "sender,s",
			Usage: "Sender Id",
		},
		cli.Int64Flag{
			Name:  "receiver,r",
			Usage: "Receiver Id",
		},
		cli.Int64SliceFlag{
			Name:  "labels,l",
			Usage: "Label Ids",
		},
	}
	traLsFlags = []cli.Flag{
		cli.BoolFlag{
			Name:  "json, j",
			Usage: "print in json format",
		},
	}
)

//TransactionCommand return the urfave.cli command of the transaction
func TransactionCommand() cli.Command {
	subcmds := []cli.Command{
		listTransactionsCmd(),
		addTransactionCmd(),
		getTransactionCmd(),
		delTransactionCmd(),
	}
	return cli.Command{
		Name:        "transaction",
		Aliases:     []string{"tran", "t"},
		Usage:       "transaction interface",
		Subcommands: subcmds,
	}
}

func listTransactionsCmd() cli.Command {
	return cli.Command{
		Name:  "ls",
		Usage: "print transactions",
		Flags: traLsFlags,
		Action: func(c *cli.Context) (err error) {
			return ListCmd(c, listTransactions)
		},
	}
}

func listTransactions() (ts model.Tabular, err error) {
	resp, err := http.Get(GetAddr("/transactions"))
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	defer resp.Body.Close()
	ts = new(model.Transactions)
	if err = json.NewDecoder(resp.Body).Decode(ts); err != nil {
		return
	}
	return
}

func addTransactionCmd() cli.Command {
	return cli.Command{
		Name:  "add",
		Usage: "add transaction",
		Flags: traAddFlags,
		Action: func(c *cli.Context) (err error) {
			return AddCmd(
				c,
				new(model.Transaction),
				createTransaction,
				addTransaction,
			)
		},
	}
}

func createTransaction(c *cli.Context) (t model.Tabular, err error) {
	timeT, err := time.Parse(time.RFC3339, c.String("time"))
	if err != nil {
		return nil, err
	}
	return &model.Transaction{
		Time:       timeT,
		Amount:     c.Float64("amount"),
		SenderId:   c.Int64("sender"),
		ReceiverId: c.Int64("receiver"),
		LabelIds:   c.Int64Slice("labels"),
	}, nil
}

func addTransaction(t model.Tabular) (nt model.Tabular, err error) {
	buf, err := json.Marshal(t)
	req, err := http.NewRequest("POST", GetAddr("/transactions"), bytes.NewBuffer(buf))
	if err != nil {
		fmt.Println(err)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	nt = new(model.Transaction)
	err = json.NewDecoder(resp.Body).Decode(nt)
	return nt, err
}

func getTransactionCmd() cli.Command {
	return cli.Command{
		Name:  "get",
		Usage: "get transaction with id",
		Action: func(c *cli.Context) error {
			return GetCmd(c, getTransaction)
		},
	}
}

func getTransaction(id string) (nt model.Tabular, err error) {
	resp, err := http.Get(GetAddr("/transactions/" + id))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	nt = new(model.Transaction)
	err = json.NewDecoder(resp.Body).Decode(nt)
	return

}

func delTransactionCmd() cli.Command {
	return cli.Command{
		Name:  "del",
		Usage: "delete transaction with id",
		Action: func(c *cli.Context) (err error) {
			if len(c.Args()) != 1 {
				return cli.NewExitError("incorrect number of args", 1)
			}
			return delTransaction(c.Args().First())
		},
	}
}

func delTransaction(id string) (err error) {
	req, err := http.NewRequest("DEL", GetAddr("/transactions/"+id), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	resp.Body.Close()
	return
}
