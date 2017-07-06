package peficli

import (
	"strconv"
)

type (
	connection struct {
		Host string
		Port int
	}
)

var (
	//Conn a connection struct to the API server
	Conn connection
)

//GetAddr get the full address of the server endpoint
func GetAddr(endpoint string) string {
	return "http://" + Conn.Host + ":" + strconv.Itoa(Conn.Port) + endpoint
}
