package account

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID       *int64     `json:"id" db:"id"`
	UID      *string    `json:"uid" db:"uid"`
	Email    *string    `json:"email" db:"email"`
	Name     *string    `json:"name, omitempty" db:"name"`
	Password *string    `json:"password, omitempty" db:"password"`
	Active   *bool      `json:"active" db:"active"`
	Gravatar *string    `json:"gravatar" db:"gravatar"`
	Created  *time.Time `json:"created" db:"created"`
	Updated  *time.Time `json:"updated" db:"updated"`
	Deleted  *bool      `json:"deleted, omitempty" db:"deleted"`
}

func NewAccount(email, name, password string) (*Account, error) {
	account := &Account{Email: &email, Name: &name}
	err := account.hashPassword([]byte(password))
	if err != nil {
		return nil, err
	}
	return account, nil
}

// Tiene que devolver nil, sino es que hay error
func (account *Account) hashPassword(password []byte) error {
	pass, err := bcrypt.GenerateFromPassword(password, 0)
	if err != nil {
		return err
	}
	*account.Password = string(pass)
	return nil
}

func (account *Account) MatchPassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(*account.Password), []byte(password)); err != nil {
		return err
	}
	return nil
}

func (account *Account) UpdatePassword(newPassword string) error {
	err := account.hashPassword([]byte(newPassword))
	if err != nil {
		return err
	}

	return nil
}

func (account *Account) DeletePassword() error {
	account.Password = nil
	return nil
}

func (account *Account) Delete() error {
	*account.Deleted = true

	return nil
}
