package peficli

import (
	"fmt"
	"github.com/urfave/cli"
	"net/http"
)

func PingFlags() []cli.Flag {
	return nil
}

func PingCommand() cli.Command {
	return cli.Command{
		Name:    "ping",
		Aliases: []string{"p"},
		Usage:   "ping the server to test connection",
		Action: func(c *cli.Context) error {
			return pingServer()
		},
		Flags: PingFlags(),
	}
}

func pingServer() error {
	fmt.Printf("pinging server at: %s:%d - ", Conn.Host, Conn.Port)
	_, err := http.Get(GetAddr("/account/get/all"))
	if err != nil {
		fmt.Println("Fail")
		errors := fmt.Sprintf("Couldnt conect to server:\n%s", err)
		return cli.NewExitError(errors, 1)
	}
	fmt.Printf("Successful\n")
	return nil
}
