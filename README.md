# pefi

Main Repository for PErsonal FInance. This repo contains the database, model descriptions as well as the golang REST API. It also contains the docker compose build and development files for both production and development deployments.

* [pefi-cli](https://github.com/simonschneider/pefi-cli) Golang CLI client
* [pefi-web](https://github.com/simonschneider/pefi-web) Web client

## Requirements

* Accounts can be external (non tracked) or internal (tracked)
* Categories can be applied to accounts (savings-longterm, savings-shortterm) (either savings or non-savings)
* [Deprecated] Nested categories
* Transactions between accounts
* Labels can be applied to transactions
* With categories accounts and expenses can be categorized.
* Loans (transactions with other transactions as payback)
* Future transactions for predictions and planning

### [Deprecated] Categorie Tree
* Saving
  * Long-term
  * Short-term
* Daily
  * Food
  * Other
* Bills
  * Living
  * Extras
* Loans
  * Living
