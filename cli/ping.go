package main

import (
	"fmt"
	"github.com/urfave/cli"
	"net/http"
)

func pingCommand() cli.Command {
	return cli.Command{
		Name:    "ping",
		Aliases: []string{"p"},
		Usage:   "ping the server to test connection",
		Action: func(c *cli.Context) error {
			return pingServer()
		},
	}
}

func pingServer() error {
	fmt.Printf("pinging server at: %s:%d - ", Conn.Host, Conn.Port)
	_, err := http.Get(GetAddr("/"))
	if err != nil {
		fmt.Println("Fail")
		errors := fmt.Sprintf("Couldnt conect to server:\n%s", err)
		return cli.NewExitError(errors, 1)
	}
	fmt.Printf("Successful\n")
	return nil
}
