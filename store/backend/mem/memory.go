package mem

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/jllopis/try5/account"
)

type MemStore struct {
	accounts map[string]*account.account
}

func NewMemStore() *MemStore {
	return &MemStore{make(map[string]*account.account, 10)}
}

func (s *MemStore) LoadAllaccounts() ([]*account.account, error) {
	accounts := make([]*account.account, len(s.accounts))
	for _, v := range s.accounts {
		accounts = append(accounts, v)
	}
	return accounts, nil
}

func (s *MemStore) Loadaccount(uuid string) (*account.account, error) {
	return s.accounts[uuid], nil
}

func (s *MemStore) Saveaccount(account *account.account) (*account.account, error) {
	if account.UID == "" {
		account.UID = uuid.New()
	}
	s.accounts[account.UID] = account
	return account, nil
}

func (s *MemStore) Deleteaccount(uuid string) (int, error) {
	delete(s.accounts, uuid)
	return 1, nil
}
