package main

import (
	"fmt"
	tm "github.com/buger/goterm"
	"github.com/urfave/cli"
	"os"
	"os/exec"
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
			createGraph,
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
		"label",
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
		strconv.Itoa(int(t.ID)),
		t.Time.Format("2006-01-02"),
		fmt.Sprintf("%.2f", t.Amount),
		strconv.Itoa(int(t.SenderID)),
		strconv.Itoa(int(t.ReceiverID)),
		strconv.Itoa(int(t.LabelID)),
	}
	return s
}

func createTransaction(c *cli.Context) (t tabular, err error) {
	timeT, err := time.Parse(time.RFC3339, c.String("time"))
	cherr(err)
	return &transaction{
		Transaction: model.Transaction{
			Time:       timeT,
			Amount:     c.Float64("amount"),
			SenderID:   c.Int64("sender"),
			ReceiverID: c.Int64("receiver"),
			LabelID:    c.Int64("label"),
		},
	}, nil
}

func createGraph(c *cli.Context, in tabular) error {
	trans, _ := in.(*transactions)
	if !c.Bool("graph") {
		return nil
	}
	w := termWidth()
	chart := tm.NewLineChart(w, 20)
	data := new(tm.DataTable)
	data.AddColumn("past Days")
	data.AddColumn("Daily total")
	sum := map[int]float64{}
	for _, t := range *trans {
		days := time.Since(t.Time)
		sum[int(days.Hours()/24)] += t.Amount
	}
	for i := 0; i <= 30; i++ {
		data.AddRow(float64(i), sum[i])
		delete(sum, i)
		if len(sum) == 0 {
			break
		}
	}
	tm.Println(chart.Draw(data))
	tm.Flush()
	return nil
}

func termWidth() int {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, _ := cmd.Output()
	widths := strings.Split(string(out), " ")
	width := strings.Trim(widths[1], "\n")
	w, _ := strconv.Atoi(width)
	return w
}
