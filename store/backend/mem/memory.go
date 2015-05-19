package mem

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/jllopis/try5/account"
	"github.com/jllopis/try5/store"
)

type MemStore struct {
	accounts map[string]*account.Account
	status   int
}

func NewMemStore() *MemStore {
	return &MemStore{accounts: make(map[string]*account.Account, 10), status: store.CONNECTED}
}

func (s *MemStore) Status() (int, string) {
	return s.status, store.StatusStr[s.status]
}

func (s *MemStore) LoadAllaccounts() ([]*account.Account, error) {
	accounts := make([]*account.Account, len(s.accounts))
	for _, v := range s.accounts {
		accounts = append(accounts, v)
	}
	return accounts, nil
}

func (s *MemStore) Loadaccount(uuid string) (*account.Account, error) {
	return s.accounts[uuid], nil
}

func (s *MemStore) Saveaccount(account *account.Account) (*account.Account, error) {
	if account.UID == nil {
		*account.UID = uuid.New()
	}
	s.accounts[*account.UID] = account
	return account, nil
}

func (s *MemStore) Deleteaccount(uuid string) (int, error) {
	delete(s.accounts, uuid)
	return 1, nil
}
