package bolt

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"time"

	"code.google.com/p/go-uuid/uuid"

	"github.com/boltdb/bolt"
	"github.com/jllopis/try5/account"
	"github.com/jllopis/try5/store"
)

type BoltStore struct {
	C      *bolt.DB
	status int
	BoltStoreOptions
}

type BoltStoreOptions struct {
	Dbpath  string
	Timeout time.Duration
}

func NewBoltStore(options *BoltStoreOptions) *BoltStore {
	db, err := bolt.Open(options.Dbpath, 0600, &bolt.Options{Timeout: options.Timeout * time.Second})
	if err != nil {
		log.Fatal(err)
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
	b := &BoltStore{C: db, status: store.CONNECTED}
	b.BoltStoreOptions = *options
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
		accounts = make([]*account.Account, bs.KeyN)
		bucket.ForEach(func(k, v []byte) error {
			var a *account.Account
			dec := gob.NewDecoder(bytes.NewBuffer(v))
			err := dec.Decode(&a)
			if err == nil {
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
			return nil
		}
		dec := gob.NewDecoder(bytes.NewBuffer(data))
		return dec.Decode(&a)
	})
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (s *BoltStore) SaveAccount(account *account.Account) (*account.Account, error) {
	now := time.Now().UTC()
	account.Updated = &now
	if account.Password != nil {
		account.UpdatePassword(*account.Password)
	}
	if account.UID == nil {
		u := uuid.New()
		account.UID = &u
		account.Created = &now
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(account)
	if err != nil {
		return nil, err
	}
	err = s.C.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("accounts")).Put([]byte(*account.UID), buf.Bytes())
	})
	if err != nil {
		return nil, err
	}
	return account, nil
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
