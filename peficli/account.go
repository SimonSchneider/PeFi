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
	accAddFlags = []cli.Flag{
		cli.StringFlag{
			Name:  "file",
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
		Action: func(c *cli.Context) (err error) {
			accs, err := listExternalAccounts()
			if err != nil {
				return err
			}
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"id", "name", "description", "labels"})
			for _, a := range *accs {
				table.Append(a.Table())
			}
			table.Render()
			return err
		},
	}
}

func listInternalAccountsCmd() cli.Command {
	return cli.Command{
		Name:  "ls",
		Usage: "print accounts",
		Action: func(c *cli.Context) (err error) {
			accs, err := listInternalAccounts()
			if err != nil {
				return err
			}
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"id", "name", "description", "labels", "balance"})
			for _, a := range *accs {
				table.Append(a.Table())
			}
			table.Render()
			return err
		},
	}
}

func listExternalAccounts() (accs *[]model.ExternalAccount, err error) {
	resp, err := http.Get(GetAddr("/accounts/external"))
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	defer resp.Body.Close()
	accs = new([]model.ExternalAccount)
	if err = json.NewDecoder(resp.Body).Decode(accs); err != nil {
		fmt.Printf("error unmarshaling: %s\n", err)
		return
	}
	return
}

func listInternalAccounts() (accs *[]model.InternalAccount, err error) {
	resp, err := http.Get(GetAddr("/accounts/internal"))
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	defer resp.Body.Close()
	accs = new([]model.InternalAccount)
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
			var nacc *model.ExternalAccount
			if path := c.String("file"); path != "" {
				acc := new(model.ExternalAccount)
				file, err := os.Open(path)
				if err != nil {
					return err
				}
				if err = json.NewDecoder(file).Decode(acc); err != nil {
					return err
				}
				nacc, err = addExternalAccount(*acc)
			} else if c.String("name") != "" {
				acc := model.ExternalAccount{
					Name:        c.String("name"),
					Description: c.String("description"),
					LabelIds:    c.Int64Slice("labels"),
				}
				nacc, err = addExternalAccount(acc)
			} else {
				return nil
			}
			err = json.NewEncoder(os.Stdout).Encode(nacc)
			return err
		},
	}
}

func addInternalAccountCmd() cli.Command {
	return cli.Command{
		Name:  "add",
		Usage: "add account",
		Flags: accAddFlags,
		Action: func(c *cli.Context) (err error) {
			var nacc *model.InternalAccount
			if path := c.String("file"); path != "" {
				acc := new(model.InternalAccount)
				file, err := os.Open(path)
				if err != nil {
					return err
				}
				if err = json.NewDecoder(file).Decode(acc); err != nil {
					return err
				}
				nacc, err = addInternalAccount(*acc)
			} else if c.String("name") != "" {
				acc := model.InternalAccount{
					ExternalAccount: model.ExternalAccount{
						Name:        c.String("name"),
						Description: c.String("description"),
						LabelIds:    c.Int64Slice("labels"),
					},
					Balance: c.Float64("balance"),
				}
				nacc, err = addInternalAccount(acc)
			} else {
				return nil
			}
			err = json.NewEncoder(os.Stdout).Encode(nacc)
			return err
		},
	}
}

func addExternalAccount(acc model.ExternalAccount) (nacc *model.ExternalAccount, err error) {
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
	json.NewDecoder(resp.Body).Decode(nacc)
	return
}

func addInternalAccount(acc model.InternalAccount) (nacc *model.InternalAccount, err error) {
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
	json.NewDecoder(resp.Body).Decode(nacc)
	return
}

func getExternalAccountCmd() cli.Command {
	return cli.Command{
		Name:  "get",
		Usage: "get account with id",
		Action: func(c *cli.Context) error {
			if len(c.Args()) != 1 {
				return cli.NewExitError("incorrect number of args", 1)
			}
			acc, err := getExternalAccount(c.Args().First())
			json.NewEncoder(os.Stdout).Encode(acc)
			return err

		},
	}
}

func getInternalAccountCmd() cli.Command {
	return cli.Command{
		Name:  "get",
		Usage: "get account with id",
		Action: func(c *cli.Context) error {
			if len(c.Args()) != 1 {
				return cli.NewExitError("incorrect number of args", 1)
			}
			acc, err := getInternalAccount(c.Args().First())
			json.NewEncoder(os.Stdout).Encode(acc)
			return err

		},
	}
}

func getExternalAccount(id string) (nacc *model.ExternalAccount, err error) {
	resp, err := http.Get(GetAddr("/accounts/external/" + id))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	nacc = new(model.ExternalAccount)
	err = json.NewDecoder(resp.Body).Decode(nacc)
	return

}

func getInternalAccount(id string) (nacc *model.InternalAccount, err error) {
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
