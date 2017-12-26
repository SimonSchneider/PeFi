package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/simonschneider/dyntab"
	"github.com/simonschneider/pefi/services/pefi"
	"github.com/simonschneider/pefi/services/pefi/http"
	"os"
	"reflect"
	"strconv"
)

func MonAm2String(i interface{}) (string, error) {
	a, ok := i.(pefi.MonetaryAmount)
	s := strconv.Itoa(int(a.Amount))
	for len(s) < 3 {
		s = "0" + s
	}
	if ok {
		return "" + s[:len(s)-2] + "." + s[len(s)-2:] + " " + a.Currency, nil
	}
	return "", errors.New("not time.location")
}

func main() {
	router := http.GetRouter(&http.AccountHandler{})
	service := http.NewAccountService(router)
	accI, _ := service.Open(context.Background(), "test2", "2ae9052c-8033-477c-a7ae-ae65e6b58879", "description")
	fmt.Println(accI)
	accI2, err := service.Get(context.Background(), accI.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	iAccs := []pefi.Account{*accI, *accI2}
	accountTable().SetData(iAccs).PrintTo(os.Stdout)
}

func accountTable() *dyntab.Table {
	return dyntab.NewTable().
		Specialize([]dyntab.ToSpecialize{{
			Type:     reflect.TypeOf(pefi.MonetaryAmount{}),
			ToString: MonAm2String}})
}
