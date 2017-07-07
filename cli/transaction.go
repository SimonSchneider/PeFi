package main

import (
	"fmt"
	"github.com/urfave/cli"
	"pefi/model"
	"strconv"
	"strings"
	"time"
)

func transactionCommand() cli.Command {
	return cli.Command{
		Name:    "transaction",
		Aliases: []string{"tran", "t"},
		Usage:   "transaction interface",
		Subcommands: GetAPISubCmd(
			"/transactions",
			new(transaction),
			new(transactions),
			createTransaction,
			transactionFlags,
		),
	}
}

type (
	transactions []transaction

	transaction struct {
		model.Transaction
	}
)

var (
	transactionHeader = []string{
		"id",
		"time",
		"amount",
		"sender",
		"receiver",
		"labels",
	}

	transactionFlags = APIFlags{
		AddFlags: []cli.Flag{
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
		},
	}
)

func (ts *transactions) Header() (s []string) {
	return transactionHeader
}

func (ts *transactions) Body() (s [][]string) {
	for _, t := range *ts {
		s = append(s, t.Table())
	}
	return s
}

func (ts *transactions) Footer() (s []string) {
	sum := 0.0
	for _, t := range *ts {
		sum += t.Amount
	}
	for i := 0; i < len(transactionHeader); i++ {
		s = append(s, "")
	}
	s[1] = "Total"
	s[2] = fmt.Sprintf("%.2f", sum)
	return s
}

func (ts *transaction) Header() (s []string) {
	return transactionHeader
}

func (t *transaction) Body() (s [][]string) {
	return [][]string{t.Table()}
}

func (t *transaction) Footer() (s []string) {
	return []string{}
}

func (t *transaction) Table() (s []string) {
	s = []string{
		strconv.Itoa(int(t.Id)),
		t.Time.Format("2006-01-02"),
		fmt.Sprintf("%.2f", t.Amount),
		strconv.Itoa(int(t.SenderId)),
		strconv.Itoa(int(t.ReceiverId)),
	}
	labelIds := []string{}
	for _, id := range t.LabelIds {
		labelIds = append(labelIds, strconv.Itoa(int(id)))
	}
	s = append(s, strings.Join(labelIds, ","))
	return s
}

func createTransaction(c *cli.Context) (t tabular, err error) {
	timeT, err := time.Parse(time.RFC3339, c.String("time"))
	if err != nil {
		return nil, err
	}
	return &transaction{
		Transaction: model.Transaction{
			Time:       timeT,
			Amount:     c.Float64("amount"),
			SenderId:   c.Int64("sender"),
			ReceiverId: c.Int64("receiver"),
			LabelIds:   c.Int64Slice("labels"),
		},
	}, nil
}
