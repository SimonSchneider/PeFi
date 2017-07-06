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

var routes = Routes{
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
		GetInternalAccount,
	},
	Route{
		"Get internal account with Id",
		"GET",
		"/accounts/internal/{accountId}",
		GetInternalAccount,
	},
}
