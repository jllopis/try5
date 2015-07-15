package manager

import (
	"github.com/jllopis/try5/account"
	"github.com/jllopis/try5/log"
	"github.com/jllopis/try5/tryerr"
)

// CreateAccount get the basic data for an account (email, name and passord) and
// call the registered provider to persist it.
// It returns the UID for the account created of an error if could not be saved.
func (m *Manager) CreateAccount(email, name, pass string) (string, error) {
	a := &account.Account{
		Email:    &email,
		Name:     &name,
		Password: &pass,
	}
	err := m.provider.SaveAccount(a)
	if err != nil {
		log.LogD("error saving account", "pkg", "manager", "func", "CreateAccount()", "error", err.Error())
		return "", err
	}
	log.LogD("account created", "pkg", "manager", "func", "CreateAccount()", "uid", *a.UID)
	return *a.UID, nil
}

// SaveAccount save the account provided using the provider registered. If an error
// occurs, it is returned.
func (m *Manager) SaveAccount(a *account.Account) error {
	err := m.provider.SaveAccount(a)
	if err != nil {
		log.LogD("error saving account", "pkg", "manager", "func", "SaveAccount()", "error", err.Error())
		return err
	}
	return nil
}

// Close method will clean the resources taken by this Manager and free them.
func (m *Manager) Close() error {
	log.LogD("closing manager", "pkg", "manager", "func", "Close()")
	return m.provider.Close()
}

// LoadAllAccounts will load all accounts present in the store.
// If the accounts can't be loaded, an error is returned.
func (m *Manager) LoadAllAccounts() ([]*account.Account, error) {
	l, err := m.provider.LoadAllAccounts()
	if err != nil {
		log.LogD("error loading accounts", "pkg", "manager", "func", "LoadAllAccounts()", "error", err.Error())
		return nil, err
	}
	return l, nil
}

// LoadAccount will load the account identified by uid if it exists in the database.
// If the account can't be loaded, an error is returned.
func (m *Manager) LoadAccount(uid string) (*account.Account, error) {
	a, err := m.provider.LoadAccount(uid)
	if err != nil {
		log.LogD("error loading account", "pkg", "manager", "func", "LoadAccount()", "error", err.Error())
		return nil, err
	}
	return a, nil
}

// DeleteAccount deletes de account referenced by uid from the store.
// If an error ocurres, it is returned.
func (m *Manager) DeleteAccount(uid string) error {
	err := m.provider.DeleteAccount(uid)
	if err != nil {
		log.LogD("error deleting account", "pkg", "manager", "func", "DeleteAccount()", "uid", uid, "error", err.Error())
		return err
	}
	return nil
}

// ExistAccount check the store for account presence. Returns a bool true if account
// is found and false otherwise.
// The account can be checked by uid, email or name. Search by name does not grant
// uniqueness as can be name conflicts. OTOH, UID and email must be unique.
// If we found an error querying the provider, an error is returned.
func (m *Manager) ExistAccount(q string) (bool, error) {
	return false, tryerr.ErrNotImplemented
}

func (m *Manager) GetAccountByEmail(email string) (*account.Account, error) {
	return nil, tryerr.ErrNotImplemented
}
