package main

import (
	"flag"
	"fmt"
	//"pefi/model"
)

func main() {
	host := flag.String("host", "localhost", "host of the app")
	flag.Parse()
	fmt.Println(*host)
}
