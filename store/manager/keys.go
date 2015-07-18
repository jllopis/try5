package manager

import (
	"github.com/jllopis/try5/keys"
	logger "github.com/jllopis/try5/log"
)

// CreateKey will create an RSA Key pair for the given account id (uid).
func (m *Manager) CreateKey(uid string) (*keys.Key, error) {
	k := keys.New(uid)
	if err := m.SaveKey(k); err != nil {
		logger.LogE("error creating rsa keys", "p", "manager", "f", "CreateKey()", "error", err)
		return nil, err
	}
	return k, nil
}

// LoadAllKeys will return all the keys found in the store.
func (m *Manager) LoadAllKeys() ([]*keys.Key, error) {
	k, err := m.provider.LoadAllKeys()
	if err != nil {
		logger.LogD("error loading keys", "pkg", "manager", "func", "LoadAllKeys()", "error", err.Error())
		return nil, err
	}
	return k, nil
}

// LoadKey will return the key identified by key id (kid).
func (m *Manager) LoadKey(kid string) (*keys.Key, error) {
	k, err := m.provider.LoadKey(kid)
	if err != nil {
		logger.LogD("error loading key", "pkg", "manager", "func", "LoadKey()", "error", err.Error())
		return nil, err
	}
	return k, nil
}

// SaveKey will save the given key to the store.
func (m *Manager) SaveKey(key *keys.Key) error {
	err := m.provider.SaveKey(key)
	if err != nil {
		logger.LogD("error saving key", "pkg", "manager", "func", "SaveKey()", "error", err.Error())
		return err
	}
	return nil
}

// DeleteKey key will delete the key identified by kid from the store.
func (m *Manager) DeleteKey(kid string) error {
	err := m.provider.DeleteKey(kid)
	if err != nil {
		logger.LogD("error deleting key", "pkg", "manager", "func", "DeleteKey()", "kid", kid, "error", err.Error())
		return err
	}
	return nil
}

// GetKeyByAccountID return the RSA key for the account uid.
// If no key is found an error is returned.
func (m *Manager) GetKeyByAccountID(uid string) (*keys.Key, error) {
	a, err := m.provider.GetKeyByAccountID(uid)
	if err != nil {
		logger.LogE("error getting keys for account", "pkg", "manager", "func", "GetKeyByAccountID()", "uid", uid, "error", err.Error())
		return nil, err
	}
	return a, nil
}

// GetKeyByEmail return the key wich matches its associated email with the one
// provided.
func (m *Manager) GetKeyByEmail(email string) (*keys.Key, error) {
	a, err := m.provider.GetKeyByEmail(email)
	if err != nil {
		logger.LogE("error getting keys for account", "pkg", "manager", "func", "GetKeyByEmail()", "email", email, "error", err.Error())
		return nil, err
	}
	return a, nil
}

// GetKeyByPub return the keys that match the public key provided.
func (m *Manager) GetKeyByPub(pubkey []byte) (*keys.Key, error) {
	a, err := m.provider.GetKeyByPub(pubkey)
	if err != nil {
		logger.LogE("error getting keys for account", "pkg", "manager", "func", "GetKeyByPub()", "pubkey", string(pubkey), "error", err.Error())
		return nil, err
	}
	return a, nil
}
