package main

type (
	AccountServiceConfiguration struct{}
	AccountService              struct{}
)

func NewAccountService(c *AccountServiceConfiguration) *AccountService {
	return &AccountService{}
}

func (a AccountService) OpenExternal(name, owner, description string) (Account, error) {
	return &ExternalAccount{
		Name:        name,
		OwnerId:     owner,
		Description: description,
		storage:     nil,
	}, nil
}

func (a AccountService) OpenInternal(name, owner, description string) (Account, error) {
	return &InternalAccount{
		ExternalAccount: ExternalAccount{
			Name:        name,
			OwnerId:     owner,
			Description: description,
			storage:     nil,
		},
		Amount: MonetaryAmount{0, "SEK"},
	}, nil
}

func Put(name, owner, description string) (string, error) {
	return "modified acc", nil
}

func Delete(name string) error {
	return nil
}

func Get(name string) (string, error) {
	return "the acc", nil
}

func Transfer(sender, receiver string) (string, error) {
	return "transfer id", nil
}

func Deposit(name string, amount uint64) (string, error) {
	return "deposit id", nil
}

func Withdraw(name string, amount uint64) (string, error) {
	return "withdraw id", nil
}
