package peficli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	"net/http"
	"pefi/model"
)

var (
	labAddFlags = []cli.Flag{
		cli.StringFlag{
			Name:  "file, f",
			Usage: "add from file",
		},
		cli.StringFlag{
			Name:  "name,n",
			Usage: "Name of account",
		},
		cli.StringFlag{
			Name:  "description,d",
			Usage: "Name of account",
		},
	}
	labLsFlags = []cli.Flag{
		cli.BoolFlag{
			Name:  "json, j",
			Usage: "print in json format",
		},
	}
)

func LabelCommand() cli.Command {
	subcmds := []cli.Command{
		listLabelsCmd(),
		addLabelCmd(),
		getLabelCmd(),
		delLabelCmd(),
	}

	return cli.Command{
		Name:        "label",
		Aliases:     []string{"lab", "l"},
		Usage:       "label interface",
		Subcommands: subcmds,
	}
}

func listLabelsCmd() cli.Command {
	return cli.Command{
		Name:  "ls",
		Usage: "print labels",
		Flags: labLsFlags,
		Action: func(c *cli.Context) (err error) {
			return ListCmd(c, listLabels)
		},
	}
}

func listLabels() (ls model.Tabular, err error) {
	resp, err := http.Get(GetAddr("/labels"))
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	defer resp.Body.Close()
	ls = new(model.Labels)
	if err = json.NewDecoder(resp.Body).Decode(ls); err != nil {
		return
	}
	return ls, nil
}

func addLabelCmd() cli.Command {
	return cli.Command{
		Name:  "add",
		Usage: "add label",
		Flags: labAddFlags,
		Action: func(c *cli.Context) (err error) {
			return AddCmd(
				c,
				new(model.Label),
				createLabel,
				addLabel,
			)
		},
	}
}

func createLabel(c *cli.Context) (nl model.Tabular, err error) {
	return &model.Label{
		Name:        c.String("name"),
		Description: c.String("description"),
	}, nil
}

func addLabel(lab model.Tabular) (nlab model.Tabular, err error) {
	buf, err := json.Marshal(lab)
	req, err := http.NewRequest("POST", GetAddr("/labels"), bytes.NewBuffer(buf))
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
	nlab = new(model.Label)
	err = json.NewDecoder(resp.Body).Decode(nlab)
	return nlab, err
}

func getLabelCmd() cli.Command {
	return cli.Command{
		Name:  "get",
		Usage: "get label with id",
		Action: func(c *cli.Context) error {
			return GetCmd(c, getLabel)
		},
	}
}

func getLabel(id string) (nlab model.Tabular, err error) {
	resp, err := http.Get(GetAddr("/labels/" + id))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	nlab = new(model.Label)
	err = json.NewDecoder(resp.Body).Decode(nlab)
	return

}

func delLabelCmd() cli.Command {
	return cli.Command{
		Name:  "del",
		Usage: "delete label with id",
		Action: func(c *cli.Context) (err error) {
			return DelCmd(c, delLabel)
		},
	}
}

func delLabel(id string) (err error) {
	req, err := http.NewRequest("DEL", GetAddr("/labels/"+id), nil)
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
