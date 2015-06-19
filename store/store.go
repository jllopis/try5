package store

import (
	"errors"

	"github.com/jllopis/try5/account"
)

type Storer interface {
	Status() (int, string)
	Close() error
	LoadAllAccounts() ([]*account.Account, error)
	LoadAccount(uuid string) (*account.Account, error)
	SaveAccount(account *account.Account) (*account.Account, error)
	DeleteAccount(uuid string) (int, error)
	GetAccountByEmail(email string) (*account.Account, error)
}

const (
	DISCONNECTED = iota
	CONNECTED
)

var (
	StatusStr = []string{"Disconnected", "Connected"}
)

var (
	ErrEmailNotFound = errors.New("email not found")
	ErrDupEmail      = errors.New("email exists in db")
)
