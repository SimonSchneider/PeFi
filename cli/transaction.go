package main

import (
	"fmt"
	tm "github.com/buger/goterm"
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"pefi/api/models"
	"strconv"
	"strings"
	"time"
)

type (
	transactions []models.Transaction
)

func transactionCommand() cli.Command {
	return cli.Command{
		Name:    "transaction",
		Aliases: []string{"tran", "t"},
		Usage:   "transaction interface",
		Subcommands: GetAPISubCmd(
			"/transactions",
			new(models.Transaction),
			new(transactions),
			createTransaction,
			transactionFlags,
			createGraph,
		),
	}
}

func (t transactions) Footer() ([]string, error) {
	sum := float64(0.0)
	for _, s := range t {
		sum += s.Amount
	}
	sums := fmt.Sprintf("%.2f", sum)
	return []string{"", "Total", sums, "", "", " "}, nil
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

func createTransaction(c *cli.Context) (t interface{}, err error) {
	timeT, err := time.Parse(time.RFC3339, c.String("time"))
	if err != nil {
		return nil, err
	}
	return models.Transaction{
		Time:       timeT,
		Amount:     c.Float64("amount"),
		SenderID:   c.Int64("sender"),
		ReceiverID: c.Int64("receiver"),
		LabelID:    c.Int64("label"),
	}, nil
}

func createGraph(c *cli.Context, in interface{}) error {
	trans, _ := in.(transactions)
	if !c.Bool("graph") {
		return nil
	}
	w := termWidth()
	chart := tm.NewLineChart(w, 20)
	data := new(tm.DataTable)
	data.AddColumn("past Days")
	data.AddColumn("Daily total")
	sum := map[int]float64{}
	for _, t := range trans {
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
