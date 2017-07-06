package peficli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
	"io/ioutil"
	"net/http"
	"os"
	"pefi/model"
	"strconv"
)

type (
	accAddFlags struct {
		Force bool
		File  string
		Type  string
	}
	accGetFlags struct {
		Format string
	}
	accDelFlags struct {
	}
	accLsFlags struct {
		Format string
	}
)

func AccountCommand() cli.Command {
	var (
		addF accAddFlags
		getF accGetFlags
		delF accDelFlags
		lsF  accLsFlags
	)

	externalAccountCmd := cli.Command{
		Name:  "external",
		Usage: "external account interface",
		Subcommands: []cli.Command{
			{
				Name:  "ls",
				Usage: "print external accounts",
				Action: func(c *cli.Context) {
					getExternalAccounts(nil, lsF)
				},
			},
			{
				Name:  "add",
				Usage: "add external account",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:        "file",
						Usage:       "add from file",
						Destination: &addF.File,
					},
				},
				Action: func(c *cli.Context) {
					addExternalAccount(c.Args(), addF)
				},
			},
		},
	}

	subcmds := []cli.Command{
		{
			Name:  "ls",
			Usage: "print all accounts",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "format, f",
					Usage:       "Printing format",
					Value:       "table",
					Destination: &lsF.Format,
				},
			},
			Action: func(c *cli.Context) {
				lsAccounts(c.Args(), lsF)
			},
		},
		externalAccountCmd,
		{
			Name:  "add",
			Usage: "add an account",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:        "force, f",
					Usage:       "force add",
					Destination: &addF.Force,
				},
				cli.StringFlag{
					Name:        "file",
					Usage:       "add from file",
					Destination: &addF.File,
				},
				cli.StringFlag{
					Name:        "type,t",
					Destination: &addF.Type,
				},
			},
			Action: func(c *cli.Context) error {
				return addAccounts(c.Args(), addF)
			},
		},
		{
			Name:    "delete",
			Aliases: []string{"del"},
			Usage:   "remove an account",
			Action: func(c *cli.Context) error {
				return delAccount(c.Args(), delF)
			},
		},
		{
			Name:  "get",
			Usage: "return an account",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "format, f",
					Usage:       "Printing format",
					Value:       "table",
					Destination: &getF.Format,
				},
			},
			Action: func(c *cli.Context) error {
				return getAccounts(c.Args(), getF)
			},
		},
	}

	return cli.Command{
		Name:        "account",
		Aliases:     []string{"a"},
		Usage:       "account interface",
		Subcommands: subcmds,
	}
}

func addAccounts(args []string, flags accAddFlags) (err error) {
	var buf []byte
	if flags.File != "" {
		buf, err = ioutil.ReadFile(flags.File)
		if err != nil {
			fmt.Printf("error reading file:\n%s", err)
			return err
		}
	}
	as := new(model.Accounts)
	if err = json.Unmarshal(buf, as); err != nil {
		return
	}
	for _, acc := range as.InternalAccounts {
		if err = addAccount(&acc, "internal", flags); err != nil {
			return
		}
	}
	for _, acc := range as.ExternalAccounts {
		if err = addAccount(&acc, "external", flags); err != nil {
			return
		}
	}
	return
}

func addAccount(acc model.Account, accType string, flags accAddFlags) (err error) {
	buf, err := json.Marshal(acc)
	req, err := http.NewRequest("POST", GetAddr("/account/add/"+accType), bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("body:", string(body))
	return
}

func getAccounts(args []string, flags accGetFlags) (err error) {
	//resp, err := http.Get(GetAddr("/account/get/" + args[0]))
	buf, err := json.Marshal(args)
	req, err := http.NewRequest("POST", GetAddr("/account/get"), bytes.NewBuffer(buf))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	defer resp.Body.Close()
	as := new(model.Accounts)
	if err = json.NewDecoder(resp.Body).Decode(as); err != nil {
		fmt.Printf("error unmarshaling: %s\n", err)
	}
	as.Print(os.Stdout, flags.Format)
	return
}

func delAccount(args []string, flags accDelFlags) (err error) {
	for _, arg := range args {
		resp, err := http.Get(GetAddr("/account/del/" + arg))
		if err != nil {
			fmt.Printf("error: %s\n", err)
		}
		resp.Body.Close()
	}
	return
}

func getExternalAccounts(args []string, flags accLsFlags) (err error) {
	resp, err := http.Get(GetAddr("/accounts/external"))
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	defer resp.Body.Close()
	accs := new([]model.ExternalAccount)
	if err = json.NewDecoder(resp.Body).Decode(accs); err != nil {
		fmt.Printf("error unmarshaling: %s\n", err)
		return
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"id", "name", "description", "labels"})
	for _, a := range *accs {
		table.Append([]string{strconv.Itoa(int(a.Id)), a.Name, a.Description, ""}) //string(a.LabelIds)})
	}
	table.Render()
	return
}

func lsAccounts(args []string, flags accLsFlags) (err error) {
	resp, err := http.Get(GetAddr("/account/ls"))
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	defer resp.Body.Close()
	as := new(model.Accounts)
	if err = json.NewDecoder(resp.Body).Decode(as); err != nil {
		fmt.Printf("error unmarshaling: %s\n", err)
		return
	}
	as.Print(os.Stdout, flags.Format)
	return
}

func addExternalAccount(args []string, flags accAddFlags) (err error) {
	var buf []byte
	if flags.File != "" {
		buf, err = ioutil.ReadFile(flags.File)
		if err != nil {
			fmt.Printf("error reading file:\n%s", err)
			return err
		}
	}
	acc := new(model.ExternalAccount)
	if err = json.Unmarshal(buf, acc); err != nil {
		return
	}
	err = addExAccount(*acc, flags)
	return
}

func addExAccount(acc model.ExternalAccount, flags accAddFlags) (err error) {
	buf, err := json.Marshal(acc)
	req, err := http.NewRequest("POST", GetAddr("/accounts/external"), bytes.NewBuffer(buf))
	if err != nil {
		fmt.Println(err)
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("body:", string(body))
	return
}
