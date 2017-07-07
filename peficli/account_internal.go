package peficli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	"net/http"
	"pefi/model"
)

type (
	internalAccounts []internalAccount

	internalAccount struct {
		model.InternalAccount
	}
)

var (
	internalAccountHeader = append(
		externalAccountHeader,
		"balance",
	)
)

func (is *internalAccounts) Header() (s []string) {
	return internalAccountHeader
}

func (is *internalAccounts) Body() (s [][]string) {
	for _, i := range *is {
		s = append(s, i.Table())
	}
	return s
}

func (is *internalAccounts) Footer() (s []string) {
	sum := 0.0
	for _, i := range *is {
		sum += i.Balance
	}
	for i := 0; i < len(internalAccountHeader); i++ {
		s = append(s, "")
	}
	s[len(s)-1] = fmt.Sprintf("%.2f", sum)
	s[len(s)-2] = "Total"
	return s
}

func (i *internalAccount) Header() (s []string) {
	return internalAccountHeader
}

func (i *internalAccount) Body() (s [][]string) {
	return [][]string{i.Table()}
}

func (i *internalAccount) Footer() (s []string) {
	return []string{}
}

func (a *internalAccount) Table() (s []string) {
	s = a.ExternalAccount.Table()
	s = append(s, fmt.Sprintf("%.2f", a.Balance))
	return s
}

var (
	accAddFlags = []cli.Flag{
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
			Name:  "labels, l",
			Usage: "add a label to the account",
		},
		cli.Float64Flag{
			Name:  "balance, b",
			Usage: "set the balance of the account",
		},
	}
	accLsFlags = []cli.Flag{
		cli.BoolFlag{
			Name:  "json, j",
			Usage: "print in json format",
		},
	}
)

func InternalAccountCommand() cli.Command {
	return cli.Command{
		Name:    "internal",
		Aliases: []string{"in", "i"},
		Usage:   "internal account interface",
		Subcommands: []cli.Command{
			listInternalAccountsCmd(),
			addInternalAccountCmd(),
			getInternalAccountCmd(),
			delInternalAccountCmd(),
		},
	}
}

func listInternalAccountsCmd() cli.Command {
	return cli.Command{
		Name:  "ls",
		Usage: "print accounts",
		Flags: accLsFlags,
		Action: func(c *cli.Context) (err error) {
			return ListCmd(c, listInternalAccounts)
		},
	}
}

func listInternalAccounts() (accs model.Tabular, err error) {
	resp, err := http.Get(GetAddr("/accounts/internal"))
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	defer resp.Body.Close()
	accs = new(internalAccounts)
	if err = json.NewDecoder(resp.Body).Decode(accs); err != nil {
		fmt.Printf("error unmarshaling: %s\n", err)
		return
	}
	return
}

func addInternalAccountCmd() cli.Command {
	return cli.Command{
		Name:  "add",
		Usage: "add account",
		Flags: accAddFlags,
		Action: func(c *cli.Context) (err error) {
			return AddCmd(
				c,
				new(internalAccount),
				createInternalAccount,
				addInternalAccount,
			)
		},
	}
}

func createInternalAccount(c *cli.Context) (t model.Tabular, err error) {
	return &internalAccount{
		InternalAccount: model.InternalAccount{
			ExternalAccount: model.ExternalAccount{
				Name:        c.String("name"),
				Description: c.String("description"),
				LabelIds:    c.Int64Slice("labels"),
			},
			Balance: c.Float64("balance"),
		},
	}, nil
}

func addInternalAccount(acc model.Tabular) (nacc model.Tabular, err error) {
	buf, err := json.Marshal(acc)
	req, err := http.NewRequest("POST", GetAddr("/accounts/internal"), bytes.NewBuffer(buf))
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
	nacc = new(internalAccount)
	err = json.NewDecoder(resp.Body).Decode(nacc)
	return
}

func getInternalAccountCmd() cli.Command {
	return cli.Command{
		Name:  "get",
		Usage: "get account with id",
		Action: func(c *cli.Context) error {
			return GetCmd(c, getInternalAccount)
		},
	}
}

func getInternalAccount(id string) (nacc model.Tabular, err error) {
	resp, err := http.Get(GetAddr("/accounts/internal/" + id))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	nacc = new(internalAccount)
	err = json.NewDecoder(resp.Body).Decode(nacc)
	return

}

func delInternalAccountCmd() cli.Command {
	return cli.Command{
		Name:  "del",
		Usage: "delete account with id",
		Action: func(c *cli.Context) (err error) {
			if len(c.Args()) != 1 {
				return cli.NewExitError("incorrect number of args", 1)
			}
			return delInternalAccount(c.Args().First())
		},
	}
}

func delInternalAccount(id string) (err error) {
	req, err := http.NewRequest("DEL", GetAddr("/accounts/internal/"+id), nil)
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
