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
	"github.com/jllopis/try5/store"
	"github.com/mgutz/logxi/v1"
)

type BoltStore struct {
	C      *bolt.DB
	status int
	BoltStoreOptions
	logger log.Logger
}

type BoltStoreOptions struct {
	Dbpath  string
	Timeout time.Duration
}

func NewBoltStore(options *BoltStoreOptions) *BoltStore {
	b := &BoltStore{logger: log.New("bolt")}
	b.BoltStoreOptions = *options
	db, err := bolt.Open(options.Dbpath, 0600, &bolt.Options{Timeout: options.Timeout * time.Second})
	if err != nil {
		b.logger.Fatal("NewBoltStore", "error", err.Error())
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("accounts"))
		return err
	})
	if err != nil {
		return nil
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("sessions"))
		return err
	})
	if err != nil {
		return nil
	}
	b.C = db
	b.status = store.CONNECTED
	return b
}

func (s *BoltStore) Status() (int, string) {
	return s.status, store.StatusStr[s.status]
}

func (s *BoltStore) LoadAllAccounts() ([]*account.Account, error) {
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

func (s *BoltStore) LoadAccount(uuid string) (*account.Account, error) {
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

func (s *BoltStore) GetAccountByEmail(email string) (*account.Account, error) {
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
	return nil, errors.New("email not found")
}

func (s *BoltStore) SaveAccount(acc *account.Account) (*account.Account, error) {
	now := time.Now().UTC()
	acc.Updated = &now
	// Check if we have an id. If we do, it "could" be an update (check if account exist first)
	// If don't, its a new account
	if acc.UID == nil {
		u := uuid.New()
		acc.UID = &u
		acc.Created = &now
		if ok := acc.Password; ok == nil {
			return nil, errors.New("nil password")
		}
		acc.UpdatePassword(*acc.Password)
		if acc.Active == nil {
			t := true
			acc.Active = &t
		}
	} else {
		savedAcc, err := s.LoadAccount(*acc.UID)
		if err != nil {
			return nil, err
		}
		if savedAcc == nil {
			s.logger.Info("SaveAccount", "cant retrieve account from db", "uid", *acc.UID)
			return nil, errors.New("uid mismatch")
		}
		// copy immutable data, that we are not allowed to modify
		acc.Created = savedAcc.Created
		if acc.Password != nil {
			if savedAcc.Password != nil {
				s.logger.Info("SaveAccount", "change passord required", *acc.UID, "old", *savedAcc.Password, "new", *acc.Password)
				if *acc.Password != *savedAcc.Password {
					acc.UpdatePassword(*acc.Password)
					s.logger.Info("SaveAccount", "passord changed", *acc.UID, "old", *savedAcc.Password, "new", *acc.Password)
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
		return nil, err
	}
	err = s.C.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("accounts")).Put([]byte(*acc.UID), buf.Bytes())
	})
	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (s *BoltStore) DeleteAccount(uuid string) (int, error) {
	err := s.C.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("accounts")).Delete([]byte(uuid))
	})
	if err != nil {
		return 0, err
	}
	return 1, nil
}

func (s *BoltStore) Close() error {
	s.status = store.DISCONNECTED
	return s.C.Close()
}
