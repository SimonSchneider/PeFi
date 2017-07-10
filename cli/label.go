package main

import (
	"github.com/urfave/cli"
	"pefi/model"
	"reflect"
	"strconv"
)

func labelCommand() cli.Command {
	return cli.Command{
		Name:    "label",
		Aliases: []string{"lab", "l"},
		Usage:   "label interface",
		Subcommands: GetAPISubCmd(
			//"/testing?orderBy=name&orderBy=description&limit=2",
			//"/testing",
			"/labels",
			new(label),
			new(labels),
			createLabel,
			labelFlags,
			nil,
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
			cli.Int64Flag{
				Name:  "categorie,c",
				Usage: "Name of categorie",
			},
		},
	}
)

func (ls *labels) Header() (s []string) {
	return (&label{}).Header()
}

//}

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
	val := reflect.ValueOf(l.Label)
	for i := 0; i < val.Type().NumField(); i++ {
		s = append(s, (val.Type().Field(i).Name))
	}
	return s
}

func (l *label) Body() (s [][]string) {
	return [][]string{l.Table()}
}

func (l *label) Footer() (s []string) {
	return []string{}
}

func (l *label) Table() (s []string) {
	return []string{
		strconv.Itoa(int(l.ID)),
		l.Name,
		l.Description,
		strconv.Itoa(int(l.CategorieID)),
	}
}

func createLabel(c *cli.Context) (nl tabular, err error) {
	return &label{
		Label: model.Label{
			Name:        c.String("name"),
			Description: c.String("description"),
			CategorieID: c.Int64("categorie"),
		},
	}, nil
}
