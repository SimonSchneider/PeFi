package main

import (
	"pefi/model"
)

//The Routes of the API
var Routes = createRoutes([]route{
	route{
		name: "categories",
		handlers: apiHandlers{
			gets: mwGets(model.GetCategories),
			get:  mwGet(model.GetCategorie),
			add:  mwAdd(new(model.Categorie), model.NewCategorie),
			del:  mwDel(model.DelCategorie),
		},
	},
	route{
		name: "labels",
		handlers: apiHandlers{
			gets: mwGets(model.GetLabels),
			get:  mwGet(model.GetLabel),
			add:  mwAdd(new(model.Label), model.NewLabel),
			del:  mwDel(model.DelLabel),
		},
	},
	route{
		name: "accounts/external",
		handlers: apiHandlers{
			gets: mwGets(model.GetExternalAccounts),
			get:  mwGet(model.GetExternalAccount),
			add:  mwAdd(new(model.ExternalAccount), model.NewExternalAccount),
			del:  mwDel(model.DelExternalAccount),
		},
	},
	route{
		name: "accounts/internal",
		handlers: apiHandlers{
			gets: mwGets(model.GetInternalAccounts),
			get:  mwGet(model.GetInternalAccount),
			add:  mwAdd(new(model.InternalAccount), model.NewInternalAccount),
			del:  mwDel(model.DelInternalAccount),
		},
	},
	route{
		name: "transactions",
		handlers: apiHandlers{
			gets: mwGets(model.GetTransactions),
			get:  mwGet(model.GetTransaction),
			add:  mwAdd(new(model.Transaction), model.NewTransaction),
			del:  mwDel(model.DelTransaction),
		},
	},
})
