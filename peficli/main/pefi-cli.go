package main

import (
	"github.com/urfave/cli"
	"os"
	"pefi/peficli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Pefi Cli application"
	app.Usage = "interface with the pefi server via a cli app"
	app.Version = "0.1"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "ip, i",
			Value:       "127.0.0.1",
			Usage:       "ip of the server",
			EnvVar:      "PEFI_IP",
			Destination: &peficli.Conn.Host,
		},
		cli.IntFlag{
			Name:        "port, p",
			Value:       22400,
			Usage:       "port of the server",
			EnvVar:      "PEFI_PORT",
			Destination: &peficli.Conn.Port,
		},
	}
	app.Commands = []cli.Command{
		peficli.PingCommand(),
		peficli.AccountCommand(),
		peficli.TransactionCommand(),
		peficli.LoginCommand(),
		peficli.LabelCommand(),
		peficli.CategorieCommand(),
	}
	app.Run(os.Args)
}
