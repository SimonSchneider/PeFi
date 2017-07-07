package peficli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	"net/http"
	"pefi/model"
	"strconv"
	"strings"
)

type (
	categories []categorie

	categorie struct {
		model.Categorie
	}
)

var (
	categorieHeader = []string{
		"id",
		"name",
		"description",
		"Labels",
		"Children",
	}
)

func (cs *categories) Header() (s []string) {
	return categorieHeader
}

func (cs *categories) Body() (s [][]string) {
	for _, c := range *cs {
		s = append(s, c.Table())
	}
	return s
}

func (cs *categories) Footer() (s []string) {
	return []string{}
}

func (c *categorie) Header() (s []string) {
	return categorieHeader
}

func (c *categorie) Body() (s [][]string) {
	return [][]string{c.Table()}
}

func (c *categorie) Footer() (s []string) {
	return []string{}
}

func (c *categorie) Table() (s []string) {
	s = []string{
		strconv.Itoa(int(c.Id)),
		c.Name,
		c.Description,
	}
	labelIds := []string{}
	for _, id := range c.LabelIds {
		labelIds = append(labelIds, strconv.Itoa(int(id)))
	}
	s = append(s, strings.Join(labelIds, ","))
	childIds := []string{}
	for _, id := range c.ChildrenIds {
		childIds = append(childIds, strconv.Itoa(int(id)))
	}
	s = append(s, strings.Join(childIds, ","))
	return s
}

var (
	catAddFlags = []cli.Flag{
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
		cli.Int64SliceFlag{
			Name:  "labels,l",
			Usage: "Label set",
		},
		cli.Int64SliceFlag{
			Name:  "children,c",
			Usage: "Children set",
		},
	}
	catLsFlags = []cli.Flag{
		cli.BoolFlag{
			Name:  "json, j",
			Usage: "print in json format",
		},
	}
)

func CategorieCommand() cli.Command {
	subcmds := []cli.Command{
		cli.Command{
			Name:  "ls",
			Usage: "print labels",
			Flags: catLsFlags,
			Action: func(c *cli.Context) (err error) {
				return ListCmd(c, listCategories)
			},
		},
		cli.Command{
			Name:  "add",
			Usage: "add categorie",
			Flags: catAddFlags,
			Action: func(c *cli.Context) (err error) {
				return AddCmd(
					c,
					new(categorie),
					createCategorie,
					addCategorie,
				)
			},
		},
		cli.Command{
			Name:  "get",
			Usage: "get categorie with id",
			Action: func(c *cli.Context) error {
				return GetCmd(c, getCategorie)
			},
		},
		cli.Command{
			Name:  "del",
			Usage: "delete label with id",
			Action: func(c *cli.Context) (err error) {
				return DelCmd(c, delCategorie)
			},
		},
	}

	return cli.Command{
		Name:        "categorie",
		Aliases:     []string{"cat", "c"},
		Usage:       "categorie interface",
		Subcommands: subcmds,
	}
}

func listCategories() (cs model.Tabular, err error) {
	resp, err := http.Get(GetAddr("/categories"))
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	defer resp.Body.Close()
	cs = new(categories)
	if err = json.NewDecoder(resp.Body).Decode(cs); err != nil {
		return
	}
	return cs, nil
}

func createCategorie(c *cli.Context) (nc model.Tabular, err error) {
	return &categorie{
		Categorie: model.Categorie{
			Name:        c.String("name"),
			Description: c.String("description"),
			LabelIds:    c.Int64Slice("labels"),
			ChildrenIds: c.Int64Slice("children"),
		},
	}, nil
}

func addCategorie(c model.Tabular) (nc model.Tabular, err error) {
	buf, err := json.Marshal(c)
	req, err := http.NewRequest("POST", GetAddr("/categories"), bytes.NewBuffer(buf))
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
	nc = new(categorie)
	err = json.NewDecoder(resp.Body).Decode(nc)
	return nc, err
}

func getCategorie(id string) (nc model.Tabular, err error) {
	resp, err := http.Get(GetAddr("/categories/" + id))
	if err != nil {
		fmt.Println("hej")
		return
	}
	defer resp.Body.Close()
	nc = new(categorie)
	err = json.NewDecoder(resp.Body).Decode(nc)
	return

}

func delCategorie(id string) (err error) {
	req, err := http.NewRequest("DEL", GetAddr("/categories/"+id), nil)
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
