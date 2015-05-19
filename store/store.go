package store

import "github.com/jllopis/try5/account"

type Storer interface {
	Status() (int, string)
	Close() error
	LoadAllAccounts() ([]*account.Account, error)
	LoadAccount(uuid string) (*account.Account, error)
	SaveAccount(account *account.Account) (*account.Account, error)
	DeleteAccount(uuid string) (int, error)
}

const (
	DISCONNECTED = iota
	CONNECTED
)

var (
	StatusStr = []string{"Disconnected", "Connected"}
)
