package peficli

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
	"io"
	"os"
	"pefi/model"
	"strconv"
)

type (
	connection struct {
		Host string
		Port int
	}

	createF func(*cli.Context) model.Tabular
)

var (
	Conn connection
)

//GetAddr get the full address of the server endpoint
func GetAddr(endpoint string) string {
	return "http://" + Conn.Host + ":" + strconv.Itoa(Conn.Port) + endpoint
}

//ToTable create a table from the content of t
func ToTable(t model.Tabular, w io.Writer) {
	table := tablewriter.NewWriter(w)
	table.SetHeader(t.Header())
	table.AppendBulk(t.Body())
	table.Render()
}

//ListCmd list the content retreived by f
func ListCmd(c *cli.Context, f func() (model.Tabular, error)) (err error) {
	out := os.Stdout
	t, err := f()
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
	ToTable(t, out)
	return nil
}

//GetCmd Meta for getting and printing against the API
func GetCmd(c *cli.Context, f func(string) (model.Tabular, error)) error {
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
	ToTable(t, out)
	return nil
}

//DelCmd Meta for deleting against the cmd
func DelCmd(c *cli.Context, f func(string) error) error {
	//out := os.Stdout
	if len(c.Args()) != 1 {
		return cli.NewExitError("incorrect number of args", 1)
	}
	if err := f(c.Args().First()); err != nil {
		s := fmt.Sprintf("%s", err)
		return cli.NewExitError("error deleting:"+s, 1)
	}
	return nil
}

func AddCmd(c *cli.Context, t model.Tabular, cF func(*cli.Context) (model.Tabular, error), f func(model.Tabular) (model.Tabular, error)) (err error) {
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
	ToTable(nt, out)
	return nil
}
