package bolt

import (
	"bytes"
	"encoding/gob"

	"github.com/boltdb/bolt"
	"github.com/jllopis/try5/tryerr"
)

// LoadToken will load the requested token from the database. The token is
// retrieved if its key id (kid) match with the param kid.
// If not found or an error occurs the error is returned
func (s *Store) LoadToken(kid string) (string, error) {
	var tok string
	err := s.C.View(func(tx *bolt.Tx) error {
		data := tx.Bucket([]byte("tokens")).Get([]byte(kid))
		if data == nil {
			return tryerr.ErrTokenNotFound
		}
		return gob.NewDecoder(bytes.NewBuffer(data)).Decode(&tok)
	})
	if err != nil {
		return "", err
	}
	return tok, nil
}

// GetTokenByEmail return the token associated with the account that match the provided email.
func (s *Store) GetTokenByEmail(email string) (string, error) {
	a, err := s.GetAccountByEmail(email)
	if err != nil {
		return "", err
	}
	return s.GetTokenByAccountID(*a.UID)
}

// GetTokenByAccountID will load the token from the database that belongs to the account
// which uid match with the param uid.
// If not found or an error occurs the error is returned.
func (s *Store) GetTokenByAccountID(uid string) (string, error) {
	return s.LoadToken(uid)
}

// SaveToken will save the provided token to the store associating it with the
// uid. If the key exist it will be updated.
func (s *Store) SaveToken(uid string, tok *string) error {
	if tok == nil {
		return tryerr.ErrNilToken
	}
	if uid == "" {
		return tryerr.ErrNilUID
	}
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(tok)
	if err != nil {
		return err
	}
	err = s.C.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("tokens")).Put([]byte(uid), buf.Bytes())
	})
	if err != nil {
		return err
	}
	return nil
}

// ExistToken return true if the token was found in the store and false otherwise.
func (s *Store) ExistToken(kid string) bool {
	found := false
	s.C.View(func(tx *bolt.Tx) error {
		data := tx.Bucket([]byte("tokens")).Get([]byte(kid))
		if data == nil {
			return nil
		}
		found = true
		return nil
	})
	return found
}

// DeleteToken will delete the token associated with the account id (uid) from the store.
func (s *Store) DeleteToken(uid string) error {
	if !s.ExistToken(uid) {
		return tryerr.ErrTokenNotFound
	}
	err := s.C.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("tokens")).Delete([]byte(uid))
	})
	if err != nil {
		return err
	}
	return nil
}
