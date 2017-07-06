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

func AccountCommand() cli.Command {
	externalAccountCmd := cli.Command{
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

	internalAccountCmd := cli.Command{
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

	subcmds := []cli.Command{
		externalAccountCmd,
		internalAccountCmd,
	}

	return cli.Command{
		Name:        "account",
		Aliases:     []string{"acc", "a"},
		Usage:       "account interface",
		Subcommands: subcmds,
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

func listExternalAccounts() (accs model.Tabular, err error) {
	resp, err := http.Get(GetAddr("/accounts/external"))
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	defer resp.Body.Close()
	accs = new(model.ExternalAccounts)
	if err = json.NewDecoder(resp.Body).Decode(accs); err != nil {
		fmt.Printf("error unmarshaling: %s\n", err)
		return
	}
	return
}

func listInternalAccounts() (accs model.Tabular, err error) {
	resp, err := http.Get(GetAddr("/accounts/internal"))
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	defer resp.Body.Close()
	accs = new(model.InternalAccounts)
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
				new(model.ExternalAccount),
				createExternalAccount,
				addExternalAccount,
			)
		},
	}
}

func createExternalAccount(c *cli.Context) (t model.Tabular, err error) {
	return &model.ExternalAccount{
		Name:        c.String("name"),
		Description: c.String("description"),
		LabelIds:    c.Int64Slice("labels"),
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
	nacc = new(model.ExternalAccount)
	err = json.NewDecoder(resp.Body).Decode(nacc)
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
				new(model.InternalAccount),
				createInternalAccount,
				addInternalAccount,
			)
		},
	}
}

func createInternalAccount(c *cli.Context) (t model.Tabular, err error) {
	return &model.InternalAccount{
		ExternalAccount: model.ExternalAccount{
			Name:        c.String("name"),
			Description: c.String("description"),
			LabelIds:    c.Int64Slice("labels"),
		},
		Balance: c.Float64("balance"),
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
	nacc = new(model.InternalAccount)
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

func getInternalAccountCmd() cli.Command {
	return cli.Command{
		Name:  "get",
		Usage: "get account with id",
		Action: func(c *cli.Context) error {
			return GetCmd(c, getInternalAccount)
		},
	}
}

func getExternalAccount(id string) (nacc model.Tabular, err error) {
	resp, err := http.Get(GetAddr("/accounts/external/" + id))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	nacc = new(model.ExternalAccount)
	err = json.NewDecoder(resp.Body).Decode(nacc)
	return

}

func getInternalAccount(id string) (nacc model.Tabular, err error) {
	resp, err := http.Get(GetAddr("/accounts/internal/" + id))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	nacc = new(model.InternalAccount)
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
