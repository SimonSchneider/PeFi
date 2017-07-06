package peficli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
	"net/http"
	"os"
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
			labs, err := listLabels()
			if err != nil {
				fmt.Println(err)
				return err
			}
			if c.Bool("json") {
				json.NewEncoder(os.Stdout).Encode(labs)
				return
			}
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{
				"id",
				"name",
				"desc",
			})
			for _, l := range *labs {
				table.Append(l.Table())
			}
			table.Render()
			return err
		},
	}
}

func listLabels() (labs *[]model.Label, err error) {
	resp, err := http.Get(GetAddr("/labels"))
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	defer resp.Body.Close()
	labs = new([]model.Label)
	if err = json.NewDecoder(resp.Body).Decode(labs); err != nil {
		return
	}
	return
}

func addLabelCmd() cli.Command {
	return cli.Command{
		Name:  "add",
		Usage: "add label",
		Flags: labAddFlags,
		Action: func(c *cli.Context) (err error) {
			var nlab *model.Label
			if path := c.String("file"); path != "" {
				lab := new(model.Label)
				file, err := os.Open(path)
				if err != nil {
					return err
				}
				if err = json.NewDecoder(file).Decode(lab); err != nil {
					return err
				}
				nlab, err = addLabel(*lab)
			} else if c.String("name") != "" {
				lab := model.Label{
					Name:        c.String("name"),
					Description: c.String("description"),
				}
				nlab, err = addLabel(lab)
			} else {
				return nil
			}
			err = json.NewEncoder(os.Stdout).Encode(nlab)
			return err
		},
	}
}

func addLabel(lab model.Label) (nlab *model.Label, err error) {
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
			if len(c.Args()) != 1 {
				return cli.NewExitError("incorrect number of args", 1)
			}
			lab, err := getLabel(c.Args().First())
			json.NewEncoder(os.Stdout).Encode(lab)
			return err

		},
	}
}

func getLabel(id string) (nlab *model.Label, err error) {
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
			if len(c.Args()) != 1 {
				return cli.NewExitError("incorrect number of args", 1)
			}
			return delLabel(c.Args().First())
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
