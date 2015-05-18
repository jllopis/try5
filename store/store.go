package store

import "github.com/jllopis/try5/account"

type Storer interface {
	LoadAllAccounts() ([]*account.Account, error)
	LoadAccount(uuid string) (*account.Account, error)
	SaveAccount(account *account.Account) (*account.Account, error)
	DeleteAccount(uuid string) (int, error)
}
