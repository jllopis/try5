package account

import (
	"errors"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Account hold the information of an Account type. Variables are of type pointer
// to easily identify null variables when persist/read to/from database storage.
type Account struct {
	ID       *int64     `json:"-" db:"id"`
	UID      *string    `json:"uid" db:"uid"`
	Email    *string    `json:"email" db:"email"`
	Name     *string    `json:"name,omitempty" db:"name"`
	Password *string    `json:"password,omitempty" db:"password"`
	Active   *bool      `json:"active" db:"active"`
	Gravatar *string    `json:"gravatar,omitempty" db:"gravatar"`
	Created  *time.Time `json:"created" db:"created"`
	Updated  *time.Time `json:"updated" db:"updated"`
	Deleted  *bool      `json:"deleted,omitempty" db:"deleted"`
}

var (
	// ErrInvalidName notifies than the name is not valid
	ErrInvalidName = errors.New("invalid name")
	// ErrInvalidPassword alert about an invalid password
	ErrInvalidPassword = errors.New("invalid password")
	// ErrInvalidEmail notifies than the provided email is no valid
	ErrInvalidEmail = errors.New("invalid email address")

	// GravatarURI is the URI of the gravatar service to show the user gravatar
	GravatarURI = "https://gravatar.com/avatar/%s?s=%v"
	// RegexpEmail check that the email is compatible with an email address
	// TODO(jllopis): in go v1.5 ckeck if is of use net/mail.AddressParser.Parse()
	RegexpEmail = regexp.MustCompile(`^[^@]+@[^@.]+\.[^@.]+`)
)

/*
NewAccount creates a new instance of Account with the given email, name and
password.

email must be unique and can not exist two accounts with same email. However, the
existance check will not be done until Save account into the store.

If an error is found, nil is returned with the error.
*/
func NewAccount(email, name, password string) (*Account, error) {
	account := &Account{Email: &email, Name: &name}
	err := account.hashPassword([]byte(password))
	if err != nil {
		return nil, err
	}
	return account, nil
}

// SetPassword add the given password to the account. First it validates that
// the password length is correct and then hash it with bcrypt.
func (account *Account) SetPassword(password string) error {
	if len(password) < 8 || len(password) > 256 {
		return ErrInvalidPassword
	}
	err := account.hashPassword([]byte(password))
	if err != nil {
		return err
	}
	return nil
}

// hashPassword do the bcrypt hashing and store it to the account
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

// MatchPassword check if the given password match with the hash stored in the account.
func (account *Account) MatchPassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(*account.Password), []byte(password)); err != nil {
		return err
	}
	return nil
}

// UpdatePassword updates the hashed password in the account with the new string.
func (account *Account) UpdatePassword(newPassword string) error {
	err := account.hashPassword([]byte(newPassword))
	if err != nil {
		return err
	}

	return nil
}

// DeletePassword erase the hash from the account. Afterwards the account will
// have no password so it can not authenticate.
func (account *Account) DeletePassword() error {
	account.Password = nil
	return nil
}

// Delete marks the account as deleted so it can not be used.
func (account *Account) Delete() error {
	if account.Deleted == nil {
		d := true
		account.Deleted = &d
		return nil
	}
	*account.Deleted = true
	return nil
}

/*
ValidateFields make sure that the account fields match the requirements.

The checks performed are:

- Name: must exist and be of length between 1 and 256

- Email: must exist and be a valid email address. Valid email should match regexp:
  `^[^@]+@[^@.]+\.[^@.]+`
*/
func (account *Account) ValidateFields() error {
	switch {
	case account.Name == nil:
		return ErrInvalidName
	case account.Email == nil:
		return ErrInvalidEmail
	case len(*account.Name) == 0 || len(*account.Name) > 256:
		return ErrInvalidName
	case len(*account.Email) == 0 || len(*account.Email) > 256 || RegexpEmail.MatchString(*account.Email) == false:
		return ErrInvalidEmail
	default:
		return nil
	}
}
