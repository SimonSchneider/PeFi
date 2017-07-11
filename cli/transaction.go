package main

import (
	"fmt"
	tm "github.com/buger/goterm"
	"github.com/urfave/cli"
	"pefi/api/models"
	"strconv"
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
			createGraph,
		),
	}
}

var (
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
				Usage: "Sender ID",
			},
			cli.Int64Flag{
				Name:  "receiver,r",
				Usage: "Receiver ID",
			},
			cli.Int64Flag{
				Name:  "label,l",
				Usage: "Label ID",
			},
		},
		LsFlags: []cli.Flag{
			cli.BoolFlag{
				Name:  "graph, g",
				Usage: "view a graph",
			},
		},
	}
)

func createTransaction(c *cli.Context) (t tabular, err error) {
	timeT, err := time.Parse(time.RFC3339, c.String("time"))
	if err != nil {
		return nil, err
	}
	return model.Transaction{
		Time:       timeT,
		Amount:     c.Float64("amount"),
		SenderID:   c.Int64("sender"),
		ReceiverID: c.Int64("receiver"),
		LabelID:    c.Int64("label"),
	}, nil
}

func createGraph(c *cli.Context, t tabular) error {
	if !c.Bool("graph") {
		return nil
	}
	//tm.Clear()
	//tm.MoveCursor(0, 0)
	chart := tm.NewLineChart(70, 20)
	data := new(tm.DataTable)
	data.AddColumn("Time")
	data.AddColumn("Transactions")
	for i := 0.0; i < 10; i += 1 {
		data.AddRow(i, i*10)
	}
	tm.Println(chart.Draw(data))
	tm.Flush()
	return nil
}
