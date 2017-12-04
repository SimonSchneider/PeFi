package pefi

type (
	AccountService interface {
		OpenExternal(name, owner, description string) (*ExternalAccount, error)
		OpenInternal(name, owner, description string) (*InternalAccount, error)
		Update(name string, new interface{}) error
		Delete(name string) error
		Get(name string) (interface{}, error)
		Transfer(sender, receiver string) (string, error)
		Deposit(name string, amount uint64) (string, error)
		Withdraw(name string, amount uint64) (string, error)
	}

	MonetaryAmount struct {
		Amount   int64  `json:amount`
		Currency string `json:currency`
	}

	ExternalAccount struct {
		Name        string `json:name`
		OwnerId     string `json:owner-id`
		Description string `json:description`
	}

	InternalAccount struct {
		ExternalAccount
		Amount MonetaryAmount `json:amount`
	}
)
