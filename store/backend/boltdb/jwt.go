package bolt

import "github.com/jllopis/try5/tryerr"

func (s *Store) LoadToken(kid string) (string, error) {
	return "", tryerr.ErrNotImplemented
}

func (s *Store) GetTokenByEmail(email string) (string, error) {
	return "", tryerr.ErrNotImplemented
}

func (s *Store) GetTokenByAccountID(uid string) (string, error) {
	return "", tryerr.ErrNotImplemented
}

func (s *Store) SaveToken(tok *string) error {
	return tryerr.ErrNotImplemented
}

func (s *Store) DeleteToken(kid string) error {
	return tryerr.ErrNotImplemented
}
