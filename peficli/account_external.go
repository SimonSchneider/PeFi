package peficli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	"net/http"
	"pefi/model"
	"strings"
)

type (
	externalAccounts []externalAccount

	externalAccount struct {
		model.ExternalAccount
	}
)

var (
	externalAccountHeader = []string{
		"id",
		"name",
		"description",
		"labels",
	}
)

func (es *externalAccounts) Header() (s []string) {
	return externalAccountHeader
}

func (es *externalAccounts) Body() (s [][]string) {
	for _, e := range *es {
		s = append(s, e.Table())
	}
	return s
}

func (es *externalAccounts) Footer() (s []string) {
	return []string{}
}

func (e *externalAccount) Header() (s []string) {
	return externalAccountHeader
}

func (e *externalAccount) Body() (s [][]string) {
	return [][]string{e.Table()}
}

func (e *externalAccount) Footer() (s []string) {
	return []string{}
}

func (a *externalAccount) Table() (s []string) {
	s = []string{
		strconv.Itoa(int(a.Id)),
		a.Name,
		a.Description,
	}
	labelIds := []string{}
	for _, id := range a.LabelIds {
		labelIds = append(labelIds, strconv.Itoa(int(id)))
	}
	s = append(s, strings.Join(labelIds, ","))
	return s
}

func ExternalAccountCommand() cli.Command {
	return cli.Command{
		Name:    "external",
		Aliases: []string{"ex", "e"},
		Usage:   "external account interface",
		Subcommands: []cli.Command{
			listExternalAccountsCmd(),
			addExternalAccountCmd(),
			getExternalAccountCmd(),
			delExternalAccountCmd(),
		},
	}
}

func listExternalAccountsCmd() cli.Command {
	return cli.Command{
		Name:  "ls",
		Usage: "print accounts",
		Flags: accLsFlags,
		Action: func(c *cli.Context) (err error) {
			return ListCmd(c, listExternalAccounts)
		},
	}
}

func listExternalAccounts() (accs model.Tabular, err error) {
	resp, err := http.Get(GetAddr("/accounts/external"))
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	defer resp.Body.Close()
	accs = new(externalAccounts)
	if err = json.NewDecoder(resp.Body).Decode(accs); err != nil {
		fmt.Printf("error unmarshaling: %s\n", err)
		return
	}
	return
}

func addExternalAccountCmd() cli.Command {
	return cli.Command{
		Name:  "add",
		Usage: "add account",
		Flags: accAddFlags,
		Action: func(c *cli.Context) (err error) {
			return AddCmd(
				c,
				new(externalAccount),
				createExternalAccount,
				addExternalAccount,
			)
		},
	}
}

func createExternalAccount(c *cli.Context) (t model.Tabular, err error) {
	return &externalAccount{
		ExternalAccount: model.ExternalAccount{
			Name:        c.String("name"),
			Description: c.String("description"),
			LabelIds:    c.Int64Slice("labels"),
		},
	}, nil
}

func addExternalAccount(acc model.Tabular) (nacc model.Tabular, err error) {
	buf, err := json.Marshal(acc)
	req, err := http.NewRequest("POST", GetAddr("/accounts/external"), bytes.NewBuffer(buf))
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
	nacc = new(externalAccount)
	err = json.NewDecoder(resp.Body).Decode(nacc)
	return
}

func getExternalAccountCmd() cli.Command {
	return cli.Command{
		Name:  "get",
		Usage: "get account with id",
		Action: func(c *cli.Context) error {
			return GetCmd(c, getExternalAccount)
		},
	}
}

func getExternalAccount(id string) (nacc model.Tabular, err error) {
	resp, err := http.Get(GetAddr("/accounts/external/" + id))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	nacc = new(externalAccount)
	err = json.NewDecoder(resp.Body).Decode(nacc)
	return

}

func delExternalAccountCmd() cli.Command {
	return cli.Command{
		Name:  "del",
		Usage: "delete account with id",
		Action: func(c *cli.Context) (err error) {
			if len(c.Args()) != 1 {
				return cli.NewExitError("incorrect number of args", 1)
			}
			return delExternalAccount(c.Args().First())
		},
	}
}

func delExternalAccount(id string) (err error) {
	req, err := http.NewRequest("DEL", GetAddr("/accounts/external/"+id), nil)
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
