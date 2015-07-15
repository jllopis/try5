package bolt

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/boltdb/bolt"
	"github.com/jllopis/try5/account"
	"github.com/jllopis/try5/log"
	"github.com/jllopis/try5/tryerr"
)

func (s *Store) LoadAllAccounts() ([]*account.Account, error) {
	var accounts []*account.Account
	err := s.C.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("accounts"))
		bs := bucket.Stats()
		fmt.Printf("\nBoltDB STATS:\n%#v\n", bs)
		bucket.ForEach(func(k, v []byte) error {
			var a *account.Account
			dec := gob.NewDecoder(bytes.NewBuffer(v))
			err := dec.Decode(&a)
			if err == nil && a != nil {
				accounts = append(accounts, a)
			}
			return nil
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (s *Store) LoadAccount(uuid string) (*account.Account, error) {
	var a *account.Account
	err := s.C.View(func(tx *bolt.Tx) error {
		data := tx.Bucket([]byte("accounts")).Get([]byte(uuid))
		if data == nil {
			return errors.New("account not found")
		}
		dec := gob.NewDecoder(bytes.NewBuffer(data))
		return dec.Decode(&a)
	})
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Store) GetAccountByEmail(email string) (*account.Account, error) {
	var found *account.Account
	err := s.C.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("accounts"))
		bucket.ForEach(func(k, v []byte) error {
			var a *account.Account
			dec := gob.NewDecoder(bytes.NewBuffer(v))
			err := dec.Decode(&a)
			if err == nil && a != nil {
				if email == *a.Email {
					found = a
					return nil
				}
			}
			return nil
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	if found != nil {
		return found, nil
	}
	return nil, tryerr.ErrEmailNotFound
}

func (s *Store) SaveAccount(acc *account.Account) error {
	now := time.Now().UTC()
	acc.Updated = &now
	// Check if we have an id. If we do, it "could" be an update (check if account exist first)
	// If don't, its a new account
	if acc.UID == nil {
		if _, err := s.GetAccountByEmail(*acc.Email); err == nil {
			// must get an ErrEmailNotFound err
			return tryerr.ErrDupEmail
		}
		u := uuid.New()
		acc.UID = &u
		acc.Created = &now
		if ok := acc.Password; ok == nil {
			return errors.New("nil password")
		}
		acc.UpdatePassword(*acc.Password)
		if acc.Active == nil {
			t := true
			acc.Active = &t
		}
	} else {
		savedAcc, err := s.LoadAccount(*acc.UID)
		if err != nil {
			return err
		}
		if savedAcc == nil {
			log.LogE("cant retrieve account from db", "p", "account", "f", "SaveAccount()", "uid", *acc.UID)
			return errors.New("error getting account from db")
		}
		// copy immutable data, that we are not allowed to modify
		acc.Created = savedAcc.Created
		if acc.Password != nil {
			if savedAcc.Password != nil {
				log.LogD("SaveAccount", "change passord required", *acc.UID, "old", *savedAcc.Password, "new", *acc.Password)
				if *acc.Password != *savedAcc.Password {
					acc.UpdatePassword(*acc.Password)
					log.LogD("SaveAccount", "passord changed", *acc.UID, "old", *savedAcc.Password, "new", *acc.Password)
				}
			}
		} else {
			if savedAcc.Password != nil {
				acc.Password = savedAcc.Password
			}
		}
		if acc.Active == nil {
			if savedAcc.Active != nil {
				acc.Active = savedAcc.Active
			} else {
				// active is true by default
				t := true
				acc.Active = &t
			}
		}
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(acc)
	if err != nil {
		return err
	}
	err = s.C.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("accounts")).Put([]byte(*acc.UID), buf.Bytes())
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) DeleteAccount(uuid string) error {
	err := s.C.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("accounts")).Delete([]byte(uuid))
	})
	if err != nil {
		return err
	}
	return nil
}
