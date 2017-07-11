package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/simonschneider/gentab"
	"github.com/urfave/cli"
	"io"
	"net/http"
	"os"
	"strconv"
)

type (
	tabular interface {
		Header() []string
		Body() [][]string
		Footer() []string
	}

	connection struct {
		Host string
		Port int
		User int
	}

	APIFlags struct {
		LsFlags  []cli.Flag
		AddFlags []cli.Flag
		DelFlags []cli.Flag
		GetFlags []cli.Flag
	}
)

var (
	Conn connection

	lsFlags = []cli.Flag{
		cli.BoolFlag{
			Name:  "json, j",
			Usage: "print in json format",
		},
	}

	addFlags = []cli.Flag{
		cli.StringFlag{
			Name:  "file, f",
			Usage: "add from file",
		},
		cli.BoolFlag{
			Name:  "json, j",
			Usage: "print in json format",
		},
	}

	getFlags = []cli.Flag{
		cli.BoolFlag{
			Name:  "json, j",
			Usage: "print in json format",
		},
	}

	delFlags = []cli.Flag{}
)

//GetAddr get the full address of the server endpoint
func GetAddr(endpoint string) string {
	return "http://" + Conn.Host + ":" + strconv.Itoa(Conn.Port) + endpoint
}

//ToTable create a table from the content of t
//func ToTable(t tabular, w io.Writer) {
//table := tablewriter.NewWriter(w)
//table.SetHeader(t.Header())
//table.AppendBulk(t.Body())
//table.SetFooter(t.Footer())
//table.Render()
//}

func GetAPISubCmd(endpoint string, mod tabular, mods tabular, cF func(*cli.Context) (tabular, error), flags APIFlags, finalF func(*cli.Context, tabular) error) []cli.Command {
	return []cli.Command{
		cli.Command{
			Name:  "ls",
			Usage: "list all",
			Flags: append(flags.LsFlags, lsFlags...),
			Action: func(c *cli.Context) (err error) {
				return ListCmd(
					c,
					GetReq(
						mods,
						endpoint,
					),
					finalF,
				)
			},
		},
		cli.Command{
			Name:  "add",
			Usage: "add",
			Flags: append(flags.AddFlags, addFlags...),
			Action: func(c *cli.Context) (err error) {
				return AddCmd(
					c,
					mod,
					cF,
					AddReq(endpoint),
				)
			},
		},
		cli.Command{
			Name:  "mod",
			Usage: "mod",
			Flags: append(flags.AddFlags, addFlags...),
			Action: func(c *cli.Context) (err error) {
				return ModCmd(
					c,
					mod,
					cF,
					ModReq(endpoint),
				)
			},
		},
		cli.Command{
			Name:  "get",
			Usage: "get id",
			Flags: append(flags.GetFlags, getFlags...),
			Action: func(c *cli.Context) error {
				return GetCmd(
					c,
					GetReq(
						mod,
						endpoint,
					),
				)
			},
		},
		cli.Command{
			Name:  "del",
			Usage: "delete id",
			Flags: append(flags.DelFlags, delFlags...),
			Action: func(c *cli.Context) (err error) {
				return DelCmd(c, DelReq(endpoint))
			},
		},
	}
}

//ListCmd list the content retreived by f
func ListCmd(c *cli.Context, f func(string) (tabular, error), ff func(*cli.Context, tabular) error) (err error) {
	if len(c.Args()) != 0 {
		return cli.NewExitError("incorrect number of args", 1)
	}
	out := os.Stdout
	t, err := f("")
	if err != nil {
		s := fmt.Sprintf("%s", err)
		return cli.NewExitError("error listing:"+s, 1)
	}
	if c.Bool("json") {
		if err = json.NewEncoder(out).Encode(t); err != nil {
			s := fmt.Sprintf("%s", err)
			return cli.NewExitError("error marshaling:"+s, 1)
		}
		return nil
	}
	gentab.PrintTable(out, t)
	if ff != nil {
		ff(c, t)
	}
	return nil
}

//GetCmd Meta for getting and printing against the API
func GetCmd(c *cli.Context, f func(string) (tabular, error)) error {
	out := os.Stdout
	if len(c.Args()) != 1 {
		return cli.NewExitError("incorrect number of args", 1)
	}
	t, err := f(c.Args().First())
	if err != nil {
		s := fmt.Sprintf("%s", err)
		return cli.NewExitError("error getting:"+s, 1)
	}
	if c.Bool("json") {
		if err = json.NewEncoder(out).Encode(t); err != nil {
			s := fmt.Sprintf("%s", err)
			return cli.NewExitError("error unmarshaling"+s, 1)
		}
		return nil
	}
	gentab.PrintTable(out, t)
	return nil
}

//DelCmd Meta for deleting against the API
func DelCmd(c *cli.Context, f func(string) error) error {
	if len(c.Args()) != 1 {
		return cli.NewExitError("incorrect number of args", 1)
	}
	if err := f(c.Args().First()); err != nil {
		s := fmt.Sprintf("%s", err)
		return cli.NewExitError("error deleting:"+s, 1)
	}
	return nil
}

//AddCmd Meta for adding against the API
func AddCmd(c *cli.Context, t tabular, cF func(*cli.Context) (tabular, error), f func(tabular) (tabular, error)) (err error) {
	if len(c.Args()) != 0 {
		return cli.NewExitError("incorrect number of args", 1)
	}
	out := os.Stdout
	if path := c.String("file"); path != "" {
		file, err := os.Open(path)
		if err != nil {
			s := fmt.Sprintf("%s", err)
			return cli.NewExitError("error reading file:"+s, 1)
		}
		if err = json.NewDecoder(file).Decode(t); err != nil {
			s := fmt.Sprintf("%s", err)
			return cli.NewExitError("error reading json file:"+s, 1)
		}
	} else {
		t, err = cF(c)
		if err != nil {
			s := fmt.Sprintf("%s", err)
			return cli.NewExitError("error creating from flags:"+s, 1)
		}
	}
	nt, err := f(t)
	if err != nil {
		s := fmt.Sprintf("%s", err)
		return cli.NewExitError("error adding:"+s, 1)
	}
	if c.Bool("json") {
		if err = json.NewEncoder(out).Encode(t); err != nil {
			s := fmt.Sprintf("%s", err)
			return cli.NewExitError("error unmarshaling"+s, 1)
		}
		return nil
	}
	//ToTable(nt, out)
	return nil
}

//AddCmd Meta for adding against the API
func ModCmd(c *cli.Context, t tabular, cF func(*cli.Context) (tabular, error), f func(string, tabular) error) (err error) {
	if len(c.Args()) != 1 {
		return cli.NewExitError("incorrect number of args", 1)
	}
	if path := c.String("file"); path != "" {
		file, err := os.Open(path)
		if err != nil {
			s := fmt.Sprintf("%s", err)
			return cli.NewExitError("error reading file:"+s, 1)
		}
		if err = json.NewDecoder(file).Decode(t); err != nil {
			s := fmt.Sprintf("%s", err)
			return cli.NewExitError("error reading json file:"+s, 1)
		}
	} else {
		t, err = cF(c)
		if err != nil {
			s := fmt.Sprintf("%s", err)
			return cli.NewExitError("error creating from flags:"+s, 1)
		}
	}
	err = f(c.Args().First(), t)
	if err != nil {
		s := fmt.Sprintf("%s", err)
		return cli.NewExitError("error adding:"+s, 1)
	}
	return nil
}

func AddReq(endpoint string) func(tabular) (tabular, error) {
	return func(mod tabular) (newMod tabular, err error) {
		buf, err := json.Marshal(mod)
		if err != nil {
			return nil, err
		}
		req, err := http.NewRequest("POST",
			GetAddr(endpoint),
			bytes.NewBuffer(buf))
		req.Header.Set("user", strconv.Itoa(Conn.User))
		if err != nil {
			return nil, err
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, errors.New(resp.Status)
		}
		if err = json.NewDecoder(resp.Body).Decode(mod); err != nil {
			return nil, err
		}
		return mod, err
	}
}

func ModReq(endpoint string) func(string, tabular) error {
	return func(id string, mod tabular) (err error) {
		buf, err := json.Marshal(mod)
		if err != nil {
			return err
		}
		req, err := http.NewRequest("PUT",
			GetAddr(endpoint+"/"+id),
			bytes.NewBuffer(buf))
		req.Header.Set("user", strconv.Itoa(Conn.User))
		if err != nil {
			return err
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return errors.New(resp.Status)
		}
		return nil
	}
}

func GetReq(mod tabular, endpoint string) func(string) (tabular, error) {
	return func(id string) (newMod tabular, err error) {
		if id != "" {
			endpoint += "/" + id
		}
		req, err := http.NewRequest("GET",
			GetAddr(endpoint), nil)
		req.Header.Set("user", strconv.Itoa(Conn.User))
		if err != nil {
			return nil, err
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, errors.New(resp.Status)
		}
		if err = json.NewDecoder(resp.Body).Decode(mod); err != nil {
			return nil, err
		}
		return mod, nil
	}
}

func DelReq(endpoint string) func(string) error {
	return func(id string) (err error) {
		req, err := http.NewRequest("DEL", GetAddr(endpoint+"/"+id), nil)
		req.Header.Set("user", strconv.Itoa(Conn.User))
		if err != nil {
			return err
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return errors.New(resp.Status)
		}
		return nil
	}
}
