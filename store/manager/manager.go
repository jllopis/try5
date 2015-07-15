/*
Package manager expose an API for the apps to call given access to the resources
of the try5 system.

It allows to manage all resources: accounts, keys, JWT, certs.
*/
package manager

import (
	"github.com/jllopis/try5/account"
	"github.com/jllopis/try5/log"
	"github.com/jllopis/try5/tryerr"
)

// Storer es la interfaz que deben implementar los diferentes backends que realizen
// la persistencia de la aplicación.
type Storer interface {
	Dial(options map[string]interface{}) error
	Status() (int, string)
	Close() error
	// Account
	LoadAllAccounts() ([]*account.Account, error)
	LoadAccount(tuuid string) (*account.Account, error)
	SaveAccount(account *account.Account) error
	DeleteAccount(uuid string) error
	GetAccountByEmail(email string) (*account.Account, error)
	// Keys
	//	LoadAllKeys() ([]*keys.Key, error)
	//	LoadKey(kid string) (*keys.Key, error)
	//	SaveKey(key *keys.Key) (*keys.Key, error)
	//	DeleteKey(kid string) error
	//	GetKeyByAccountID(uid string) (*keys.Key, error)
	//	GetKeyByEmail(email string) (*keys.Key, error)
	//	GetKeyByPub(pubkey []byte) (*keys.Key, error)
	// Tokens
	//	LoadToken(kid string) (string, error)
	//	GetTokenByEmail(email string) (string, error)
	//	GetTokenByAccountID(uid string) (string, error)
	//	SaveToken(string) (string, error)
	//	DeleteToken(kid string) error
}

const (
	// DISCONNECTED indicates that there is no connection with the Storer
	DISCONNECTED = iota
	// CONNECTED indicate that the connection with the Storer is up and running
	CONNECTED
)

var (
	// StatusStr is a string representation of the status of the connections with the Storer
	StatusStr        = []string{"Disconnected", "Connected"}
	registeredStores = make(map[string]Storer)
)

// Manager holds the instance of the actual Storer provider to operate
type Manager struct {
	provider       Storer
	providerConfig map[string]interface{}
}

// Register registers the provided StoreProvider and allows to be used afterwards.
// The StoreProvider can be called usign the name provided.
func Register(name string, s Storer) error {
	if s == nil {
		return tryerr.ErrNilStore
	}
	if _, dup := registeredStores[name]; dup {
		return tryerr.ErrStorerAlreadyRegistered
	}
	log.LogD("registering store backend", "p", "manager", "f", "Register()", "name", name)
	registeredStores[name] = s
	return nil
}

// Unregister delete the registered StoreProvider. If the StoreProvider is not registered,
// an ErrStoreNotRegistered error is returned.
func Unregister(name string) error {
	if _, ok := registeredStores[name]; ok {
		delete(registeredStores, name)
		log.LogD("unregistering store backend", "p", "manager", "f", "Unregister()", "name", name)
		return nil
	}
	log.LogE("unregister store backend", "p", "manager", "f", "Unregister()", "name", name, "error", tryerr.ErrStoreNotRegistered)
	return tryerr.ErrStoreNotRegistered
}

// NewManager return a new instance of store.Manager
func NewManager(config map[string]interface{}, storeName string) (*Manager, error) {
	log.LogD("initiating new manager", "p", "manager", "f", "NewManager()", "store", storeName)
	if s, ok := registeredStores[storeName]; ok {
		return &Manager{provider: s, providerConfig: config}, nil
	}
	return nil, tryerr.ErrStoreNotRegistered
}

// Init initialize the structure and connect to the declared store
func (m *Manager) Init() error {
	log.LogD("initializing manager", "p", "manager", "f", "Init()")
	if m.provider == nil {
		return tryerr.ErrNilStore
	}
	err := m.provider.Dial(m.providerConfig)
	if err != nil {
		log.LogE("error dialing store", "p", "manager", "f", "Init()", "error", err)
		return err
	}
	return nil
}

// GetStore return a valid Storer. If there is none registered an error is returned
func (m *Manager) GetStore(name string) (Storer, error) {
	if s, ok := registeredStores[name]; ok {
		return s, nil
	}
	return nil, tryerr.ErrStoreNotRegistered
}
