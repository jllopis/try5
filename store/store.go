package store

import "github.com/jllopis/try5/user"

type Storer interface {
	LoadUser(uuid string) (*user.User, error)
	SaveUser(user *user.User) (*user.User, error)
	DeleteUser(uuid string) (int, error)
}
