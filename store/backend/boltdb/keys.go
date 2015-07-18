package bolt

import (
	"bytes"
	"encoding/gob"
	"time"

	"code.google.com/p/go-uuid/uuid"

	"github.com/boltdb/bolt"
	"github.com/jllopis/try5/keys"
	"github.com/jllopis/try5/log"
	"github.com/jllopis/try5/tryerr"
)

// LoadAllKeys query the database and returns all the keys found.
// If some error could happen, it is returned
func (s *Store) LoadAllKeys() ([]*keys.Key, error) {
	var kk []*keys.Key
	err := s.C.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("keys"))
		bucket.ForEach(func(k, v []byte) error {
			var a *keys.Key
			err := gob.NewDecoder(bytes.NewBuffer(v)).Decode(&a)
			if err == nil && a != nil {
				kk = append(kk, a)
			}
			return nil
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return kk, nil
}

// LoadKey will load the requested key from the database. The account is
// retrieved if its key id (kid) match with the param kid.
// If not found or an error occurs the error is returned
func (s *Store) LoadKey(kid string) (*keys.Key, error) {
	var a *keys.Key
	err := s.C.View(func(tx *bolt.Tx) error {
		data := tx.Bucket([]byte("keys")).Get([]byte(kid))
		if data == nil {
			return tryerr.ErrAccountNotFound
		}
		return gob.NewDecoder(bytes.NewBuffer(data)).Decode(&a)
	})
	if err != nil {
		return nil, err
	}
	return a, nil
}

// SaveKey will save the provided key to the store. If the key exist it will be
// updated.
func (s *Store) SaveKey(key *keys.Key) error {
	if key.PrivKey == nil || key.PubKey == nil {
		return tryerr.ErrNilKey
	}
	now := time.Now().UTC()
	key.Updated = &now
	// Check if we have an id. If we do, it "could" be an update (check if key exist first)
	// If don't, its a new key
	if key.KID == nil {
		if _, err := s.GetKeyByAccountID(*key.AccountID); err == nil {
			return tryerr.ErrKeyExists
		}
		u := uuid.New()
		key.KID = &u
		key.Created = &now
		if key.Active == nil {
			t := true
			key.Active = &t
		}
	} else {
		SavedKey, err := s.LoadKey(*key.KID)
		if err != nil {
			return err
		}
		if SavedKey == nil {
			log.LogE("cant retrieve key from db", "p", "bolt", "f", "SaveKey()", "kid", *key.KID)
			return tryerr.ErrKeyNotFound
		}
		// copy immutable data, that we are not allowed to modify
		key.Created = SavedKey.Created
		if key.Active == nil {
			if SavedKey.Active != nil {
				key.Active = SavedKey.Active
			} else {
				// active is true by default
				t := true
				key.Active = &t
			}
		}
		if key.AccountID == nil {
			key.AccountID = SavedKey.AccountID
		}
	}
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(key)
	if err != nil {
		return err
	}
	err = s.C.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("keys")).Put([]byte(*key.KID), buf.Bytes())
	})
	if err != nil {
		return err
	}
	return nil
}

// DeleteKey will delete the kye identified with key id (kid) from the store.
func (s *Store) DeleteKey(kid string) error {
	if !s.ExistKey(kid) {
		return tryerr.ErrKeyNotFound
	}
	err := s.C.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("keys")).Delete([]byte(kid))
	})
	if err != nil {
		return err
	}
	return nil
}

// ExistKey return true if the key was found in the store and false otherwise.
func (s *Store) ExistKey(kid string) bool {
	found := false
	s.C.View(func(tx *bolt.Tx) error {
		data := tx.Bucket([]byte("keys")).Get([]byte(kid))
		if data == nil {
			return nil
		}
		found = true
		return nil
	})
	return found
}

// GetKeyByAccountID will load the key from the database that belongs to the account
// which uid match with the param uid.
// If not found or an error occurs the error is returned
func (s *Store) GetKeyByAccountID(uid string) (*keys.Key, error) {
	var found *keys.Key
	err := s.C.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("keys"))
		bucket.ForEach(func(k, v []byte) error {
			var kp *keys.Key
			dec := gob.NewDecoder(bytes.NewBuffer(v))
			err := dec.Decode(&kp)
			if err == nil && kp != nil {
				if uid == *kp.AccountID {
					found = kp
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
	return nil, tryerr.ErrKeyNotFound
}

// GetKeyByEmail return the key associated with the account that match the provided email.
func (s *Store) GetKeyByEmail(email string) (*keys.Key, error) {
	a, err := s.GetAccountByEmail(email)
	if err != nil {
		return nil, err
	}
	return s.GetKeyByAccountID(*a.UID)
}

// GetKeyByPub returns the key that matches the given pubkey.
func (s *Store) GetKeyByPub(pubkey []byte) (*keys.Key, error) {
	var found *keys.Key
	err := s.C.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("keys"))
		bucket.ForEach(func(k, v []byte) error {
			var kp *keys.Key
			err := gob.NewDecoder(bytes.NewBuffer(v)).Decode(&kp)
			if err == nil && kp != nil {
				if bytes.Equal(pubkey, kp.PubKey) {
					found = kp
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
	return nil, tryerr.ErrKeyNotFound
}
