package api

import "net/http"

type Route struct {
	Name   string
	Method string
	//Input       string
	//Output      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var accRoutes = Routes{
	Route{
		"Add external accounts",
		"POST",
		"/accounts/external",
		AddExternalAccount,
	},
	Route{
		"Add internal accounts",
		"POST",
		"/accounts/internal",
		AddInternalAccount,
	},
	Route{
		"Get list of external accounts",
		"GET",
		"/accounts/external",
		GetExternalAccounts,
	},
	Route{
		"Get list of internal accounts",
		"GET",
		"/accounts/internal",
		GetInternalAccounts,
	},
	Route{
		"Get external account with Id",
		"GET",
		"/accounts/external/{accountId}",
		GetExternalAccount,
	},
	Route{
		"Get internal account with Id",
		"GET",
		"/accounts/internal/{accountId}",
		GetInternalAccount,
	},
	Route{
		"Delete external account with Id",
		"DEL",
		"/accounts/external/{accountId}",
		DelExternalAccount,
	},
	Route{
		"Delete internal account with Id",
		"DEL",
		"/accounts/internal/{accountId}",
		DelInternalAccount,
	},
}

var labRoutes = Routes{
	Route{
		"Add labels",
		"POST",
		"/labels",
		AddLabel,
	},
	Route{
		"Get list of labels",
		"GET",
		"/labels",
		GetLabels,
	},
	Route{
		"Get label with Id",
		"GET",
		"/labels/{labelId}",
		GetLabel,
	},
	Route{
		"Delete label with Id",
		"DEL",
		"/labels/{labelId}",
		DelLabel,
	},
}
