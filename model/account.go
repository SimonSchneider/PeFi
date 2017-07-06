package model

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"io"
	"pefi/model/redis"
	"strconv"
)

type (
	Account interface {
		Table() []string
	}

	Accounts struct {
		InternalAccounts []InternalAccount `json:"ia"`
		ExternalAccounts []ExternalAccount `json:"ea"`
	}

	ExternalAccount struct {
		Id          int64   `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		LabelIds    []int64 `json:"label_ids"`
	}

	InternalAccount struct {
		ExternalAccount
		Balance float64 `json:"balance"`
	}
)

//Todo change to proper format type
func (a *Accounts) Print(w io.Writer, format string) {
	if format == "json" {
		json.NewEncoder(w).Encode(a)
	} else if format == "table" {
		fmt.Println("External Accounts")
		extTable := tablewriter.NewWriter(w)
		extTable.SetHeader([]string{"id", "name", "description"})
		for _, a := range a.ExternalAccounts {
			extTable.Append(a.Table())
		}
		extTable.Render()
		fmt.Println("Internal Accounts")
		intSum := 0.0
		intTable := tablewriter.NewWriter(w)
		intTable.SetHeader([]string{"id", "name", "description", "balance"})
		for _, a := range a.InternalAccounts {
			intTable.Append(a.Table())
			intSum += a.Balance
		}
		strSum := fmt.Sprintf("%.2f", intSum)
		intTable.SetFooter([]string{"", "", "Total", strSum})
		intTable.Render()
	}
}

func (a *InternalAccount) Table() (s []string) {
	s = a.ExternalAccount.Table()
	s = append(s, fmt.Sprintf("%.2f", a.Balance))
	return s
}

func (a *ExternalAccount) Table() (s []string) {
	s = []string{
		strconv.Itoa(int(a.Id)),
		a.Name,
		a.Description,
	}
	return s
}

func (a *InternalAccount) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		InternalAccount
	}{
		InternalAccount: *a,
	})
}
func (a *ExternalAccount) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ExternalAccount
	}{
		ExternalAccount: *a,
	})
}

func GetExternalAccounts() (accs []ExternalAccount, err error) {
	vals, err := redis.HGetAll("ExternalAccount")
	if err != nil {
		return
	}
	for _, val := range vals {
		a := new(ExternalAccount)
		if err = json.Unmarshal([]byte(val), a); err != nil {
			return
		}
		accs = append(accs, *a)
	}
	return
}

// Returns the account that has the ID
func GetAccounts(ids []int64) (a Accounts, err error) {
	for _, id := range ids {
		var tmp string
		if tmp, err = redis.HGet("InternalAccount", string(id)); err == nil {
			ta := new(InternalAccount)
			if err = json.Unmarshal([]byte(tmp), ta); err != nil {
				return
			}
			a.InternalAccounts = append(a.InternalAccounts, *ta)
		} else if tmp, err = redis.HGet("ExternalAccount", string(id)); err == nil {
			ta := new(ExternalAccount)
			if err = json.Unmarshal([]byte(tmp), ta); err != nil {
				return
			}
			a.ExternalAccounts = append(a.ExternalAccounts, *ta)
		}
	}
	return
}

//NewExternalAccount create a new External Account and return it
func NewExternalAccount(data io.Reader) (a *ExternalAccount, err error) {
	a = new(ExternalAccount)
	id, err := redis.HIncrBy("unique_ids", "Account", 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = json.NewDecoder(data).Decode(a); err != nil {
		fmt.Println(err)
		return
	}
	a.Id = id
	ma, err := json.Marshal(a)
	if err != nil {
		return
	}
	redis.HSet("ExternalAccount", string(a.Id), string(ma))
	return
}

//NewInternalAccount create a new External Account and return it
func NewInternalAccount(data io.Reader) (a *InternalAccount, err error) {
	a = new(InternalAccount)
	id, err := redis.HIncrBy("unique_ids", "Account", 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = json.NewDecoder(data).Decode(a); err != nil {
		fmt.Println(err)
		return
	}
	a.Id = id
	ma, err := json.Marshal(a)
	if err != nil {
		return
	}
	redis.HSet("InternalAccount", string(a.Id), string(ma))
	return
}

func GetAllAccounts() (as Accounts, err error) {
	ias, err := getAllInternalAccounts()
	eas, err := getAllExternalAccounts()
	as.InternalAccounts = ias
	as.ExternalAccounts = eas
	return
}

func getAllInternalAccounts() (as []InternalAccount, err error) {
	vals, err := redis.HGetAll("InternalAccount")
	if err != nil {
		return
	}
	for _, val := range vals {
		a := new(InternalAccount)
		if err = json.Unmarshal([]byte(val), a); err != nil {
			return
		}
		as = append(as, *a)
	}
	return
}

func getAllExternalAccounts() (as []ExternalAccount, err error) {
	vals, err := redis.HGetAll("ExternalAccount")
	if err != nil {
		return
	}
	for _, val := range vals {
		a := new(ExternalAccount)
		if err = json.Unmarshal([]byte(val), a); err != nil {
			return
		}
		as = append(as, *a)
	}
	return
}

func DelAccount(id int64) (err error) {
	err = redis.HDel("InternalAccount", string(id))
	err = redis.HDel("ExternalAccount", string(id))
	return
}
