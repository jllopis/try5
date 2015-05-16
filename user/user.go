package user

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int64     `json:"id" db:"id"`
	UID      string    `json:"uid" db:"uid"`
	Email    string    `json:"email" db:"email"`
	Name     string    `json:"name, omitempty" db:"name"`
	Password string    `json:"password, omitempty" db:"password"`
	Deleted  bool      `json:"deleted, omitempty" db:"deleted"`
	Created  time.Time `json:"created" db:"created"`
	Updated  time.Time `json:"updated" db:"updated"`
}

func NewUser(email, name, password string) (*User, error) {
	user := &User{Email: email, Name: name}
	err := user.hashPassword([]byte(password))
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Tiene que devolver nil, sino es que hay error
func (user *User) hashPassword(password []byte) error {
	pass, err := bcrypt.GenerateFromPassword(password, 0)
	if err != nil {
		return err
	}
	user.Password = string(pass)
	return nil
}

func (user *User) MatchPassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return err
	}
	return nil
}

func (user *User) UpdatePassword(newPassword string) error {
	err := user.hashPassword([]byte(newPassword))
	if err != nil {
		return err
	}

	return nil
}

func (user *User) DeletePassword() error {
	user.Password = ""
	return nil
}

func (user *User) Delete() error {
	user.Deleted = true

	return nil
}
