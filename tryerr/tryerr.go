package tryerr

import "errors"

var (
	// ErrInvalidContext is returned when no context is provided or it is invalid
	ErrInvalidContext = errors.New("invalid context")
	// ErrStorerAlreadyRegistered is returned when the Storer has been previously registered
	ErrStorerAlreadyRegistered = errors.New("provider already registered")
	// ErrNilStore is returned when the provided store is nil, not a valid Storer instance
	ErrNilStore = errors.New("store cannot be nil")
	// ErrStoreNotRegistered is returned when trying to access a store that has not been registered
	ErrStoreNotRegistered = errors.New("store not registered")
	// ErrAccountNotFound is returned when the required account was not found in the store
	ErrAccountNotFound = errors.New("account not found")
	// ErrEmailNotFound is returned when no account is found in the store with the given email
	ErrEmailNotFound = errors.New("email not found")
	// ErrDupEmail is returned when the provided email is already found in the store
	ErrDupEmail = errors.New("email exists in db")
	// ErrKeyExists is returned when the provided key already exists in the store
	ErrKeyExists = errors.New("key exists in db")
	// ErrKeyNotFound is returned when the requested key is not found in the store
	ErrKeyNotFound = errors.New("key not found")
	// ErrTokenNotFound is returned when the requested token is not found in the store
	ErrTokenNotFound = errors.New("token not found")
	// ErrInvalidToken  is returned when the token is not a JWT valid token
	ErrInvalidToken = errors.New("token not valid")
	// ErrNotImplemented is returned when the functionality required is not implemented
	ErrNotImplemented = errors.New("function not implemented")
)
