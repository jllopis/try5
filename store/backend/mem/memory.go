package mem

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/jllopis/try5/user"
)

type MemStore struct {
	users map[string]*user.User
}

func NewMemStore() *MemStore {
	return &MemStore{make(map[string]*user.User, 10)}
}

func (s *MemStore) LoadUser(uuid string) (*user.User, error) {
	return s.users[uuid], nil
}

func (s *MemStore) SaveUser(user *user.User) (*user.User, error) {
	if user.UID == "" {
		user.UID = uuid.New()
	}
	s.users[user.UID] = user
	return user, nil
}

func (s *MemStore) DeleteUser(uuid string) (int, error) {
	delete(s.users, uuid)
	return 1, nil
}
