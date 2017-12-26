package graphql

import (
	"context"
	"github.com/simonschneider/pefi"
)

var Services struct {
	AS pefi.AccountService
	US pefi.UserService
}

var Schema = `
	schema {
		query: Query
	}

	type Query {
		account(id: ID): account
		accounts(userID: ID): [account]!
		user(id: ID): user
	}

	type account {
		id: String!
		name: String!
		description: String!
		owner: user!
		balance: monetary_amount!
	}

	type user {
		id: String!
		name: String!
		accounts: [account]
	}

	type monetary_amount {
		amount: Int!
		currency: String!
	}
`

type monetaryAmountResolver struct {
	m *pefi.MonetaryAmount
}

type userResolver struct {
	u *pefi.User
}

func (r userResolver) ID() string {
	return string(r.u.ID)
}

func (r userResolver) Name() string {
	return r.u.Name
}

func (r userResolver) Accounts(ctx context.Context) (*[]*accountResolver, error) {
	accounts, err := Services.AS.GetAll(ctx, r.u.ID)
	if err != nil {
		return nil, err
	}
	var l []*accountResolver
	for _, account := range accounts {
		l = append(l, &accountResolver{account})
	}
	return &l, nil
}

func (r monetaryAmountResolver) Amount() int32 {
	return int32(r.m.Amount)
}

func (r monetaryAmountResolver) Currency() string {
	return r.m.Currency
}

type accountResolver struct {
	a *pefi.Account
}

func (r accountResolver) ID() string {
	return string(r.a.ID)
}

func (r accountResolver) Name() string {
	return r.a.Name
}

func (r accountResolver) Owner(ctx context.Context) (*userResolver, error) {
	u, err := Services.US.Get(ctx, r.a.OwnerID)
	if err != nil {
		return nil, err
	}
	return &userResolver{u}, nil
}

func (r accountResolver) OwnerID() string {
	return string(r.a.OwnerID)
}

func (r accountResolver) Description() string {
	return r.a.Description
}

func (r accountResolver) Balance() *monetaryAmountResolver {
	return &monetaryAmountResolver{&r.a.Balance}
}

type Resolver struct{}

func (r *Resolver) Account(ctx context.Context, args struct{ ID *string }) (*accountResolver, error) {
	a, err := Services.AS.Get(context.Background(), pefi.ID(*args.ID))
	if err != nil {
		return nil, err
	}
	return &accountResolver{a}, nil
}

func (r *Resolver) Accounts(ctx context.Context, args struct{ UserID *string }) ([]*accountResolver, error) {
	accounts, err := Services.AS.GetAll(context.Background(), pefi.ID(*args.UserID))
	if err != nil {
		return nil, err
	}
	var l []*accountResolver
	for _, account := range accounts {
		l = append(l, &accountResolver{account})
	}
	return l, nil
}

func (r *Resolver) User(ctx context.Context, args struct{ ID *string }) (*userResolver, error) {
	u, err := Services.US.Get(context.Background(), pefi.ID(*args.ID))
	if err != nil {
		return nil, err
	}
	return &userResolver{u}, nil
}
