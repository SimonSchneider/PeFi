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
	}

	type account {
		id: String!
		name: String!
		description: String!
		owner_id: String!
		balance: Int!
	}
`

//type user {
//id: ID!
//name: String!
//}
//owner: user!
//type MonetaryAmount {
//Amount: Int!
//Currency: String!
//}
//`

//type userResolver struct {
//u *pefi.User
//}

//func (r userResolver) ID() graphql.ID {
//return graphql.ID(r.u.ID)
//}

//func (r userResolver) Name() string {
//return r.u.Name
//}

type accountResolver struct {
	a *pefi.Account
}

func (r accountResolver) ID() string {
	return string(r.a.ID)
}

func (r accountResolver) Name() string {
	return r.a.Name
}

//func (r accountResolver) Owner() *userResolver {
//u, _ := Services.US.Get(r.a.OwnerID)
//return &userResolver{u}
//}

func (r accountResolver) OwnerID() string {
	return string(r.a.OwnerID)
}

func (r accountResolver) Description() string {
	return r.a.Description
}

func (r accountResolver) Balance() int32 {
	return int32(r.a.Balance.Amount)
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
