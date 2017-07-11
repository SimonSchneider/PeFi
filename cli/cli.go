//cli client for the pefi api
package main

import (
	"github.com/urfave/cli"
	"os"
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
			Destination: &Conn.Host,
		},
		cli.IntFlag{
			Name:        "port, p",
			Value:       22400,
			Usage:       "port of the server",
			EnvVar:      "PEFI_PORT",
			Destination: &Conn.Port,
		},
		cli.IntFlag{
			Name:        "user, u",
			Value:       1,
			Usage:       "user on the server",
			EnvVar:      "PEFI_USER",
			Destination: &Conn.User,
		},
	}
	app.Commands = []cli.Command{
		pingCommand(),
		accountCommand(),
		transactionCommand(),
		loginCommand(),
		labelCommand(),
		categoryCommand(),
	}
	app.Run(os.Args)
}
