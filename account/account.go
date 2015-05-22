package account

import (
	"errors"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID       *int64     `json:"-" db:"id"`
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

var (
	ErrInvalidName     = errors.New("invalid name")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidEmail    = errors.New("invalid email address")

	GravatarURI = "https://gravatar.com/avatar/%s?s=%v"

	RegexpEmail = regexp.MustCompile(`^[^@]+@[^@.]+\.[^@.]+`)
)

func NewAccount(email, name, password string) (*Account, error) {
	account := &Account{Email: &email, Name: &name}
	err := account.hashPassword([]byte(password))
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (account *Account) SetPassword(password string) error {
	if len(password) < 8 || len(password) > 256 {
		return ErrInvalidPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if account.Password == nil {
		s := string(hash)
		account.Password = &s
		return nil
	}
	*account.Password = string(hash)
	return nil
}

// Tiene que devolver nil, sino es que hay error
func (account *Account) hashPassword(password []byte) error {
	pass, err := bcrypt.GenerateFromPassword(password, 0)
	if err != nil {
		return err
	}
	if account.Password == nil {
		s := string(pass)
		account.Password = &s
		return nil
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
	if account.Deleted == nil {
		d := true
		account.Deleted = &d
		return nil
	}
	*account.Deleted = true
	return nil
}

func (a *Account) ValidateFields() error {
	switch {
	case a.Name == nil:
		return ErrInvalidName
	case a.Email == nil:
		return ErrInvalidEmail
	case len(*a.Name) == 0 || len(*a.Name) > 256:
		return ErrInvalidName
	case len(*a.Email) == 0 || len(*a.Email) > 256 || RegexpEmail.MatchString(*a.Email) == false:
		return ErrInvalidEmail
	default:
		return nil
	}
}
