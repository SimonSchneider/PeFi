package main

import (
	"github.com/urfave/cli"
	"pefi/model"
	"strconv"
)

func labelCommand() cli.Command {
	return cli.Command{
		Name:    "label",
		Aliases: []string{"lab", "l"},
		Usage:   "label interface",
		Subcommands: GetAPISubCmd(
			"/labels",
			new(label),
			new(labels),
			createLabel,
			labelFlags,
		),
	}
}

type (
	labels []label

	label struct {
		model.Label
	}
)

var (
	labelHeader = []string{
		"id",
		"name",
		"desc",
	}

	labelFlags = APIFlags{
		AddFlags: []cli.Flag{
			cli.StringFlag{
				Name:  "name,n",
				Usage: "Name of account",
			},
			cli.StringFlag{
				Name:  "description,d",
				Usage: "Name of account",
			},
		},
	}
)

func (ls *labels) Header() (s []string) {
	return labelHeader
}

func (ls *labels) Body() (s [][]string) {
	for _, l := range *ls {
		s = append(s, l.Table())
	}
	return s
}

func (ls *labels) Footer() (s []string) {
	return []string{}
}

func (l *label) Header() (s []string) {
	return labelHeader
}

func (l *label) Body() (s [][]string) {
	return [][]string{l.Table()}
}

func (l *label) Footer() (s []string) {
	return []string{}
}

func (l *label) Table() (s []string) {
	return []string{
		strconv.Itoa(int(l.Id)),
		l.Name,
		l.Description,
	}
}

func createLabel(c *cli.Context) (nl tabular, err error) {
	return &label{
		Label: model.Label{
			Name:        c.String("name"),
			Description: c.String("description"),
		},
	}, nil
}
