package manager

import (
	"github.com/jllopis/try5/jwt"
	"github.com/jllopis/try5/keys"
	"github.com/jllopis/try5/tryerr"
)

// CreateToken will create a JWT with the private RSA Key from the given account id (uid).
func (m *Manager) CreateToken(uid string) (string, error) {
	var keypair *keys.Key
	var err error
	if keypair, err = m.GetKeyByAccountID(uid); err != nil {
		logger.LogE("key not found", "error", err.Error(), "uid", uid)
		return nil, err
	}
	token, err := jwt.GenerateToken(uid, keypair.PrivKey)
	if err != nil {
		logger.LogE("can not create token", "error", err.Error(), "uid", uid)
		return nil, err
	}
	token, err = m.SaveToken(token)
	if err != nil {
		logger.LogE("can not save token", "pkg", "api", "func", "NewJWTToken()", "error", err.Error())
		return nil, err
	}
	return token, nil
}

// LoadToken will return the JWT identified by token id (kid).
func (m *Manager) LoadToken(kid string) (string, error) {
	tok, err := m.provider.LoadToken(kid)
	if err != nil {
		logger.LogD("error loading JWT", "pkg", "manager", "func", "LoadToken()", "error", err.Error())
		return nil, err
	}
	return tok, nil
}

// GetTokenByEmail will return the JWT that belongs to the account identified by email.
func (m *Manager) GetTokenByEmail(email string) (string, error) {
	return "", tryerr.ErrNotImplemented
}

// GetTokenByAccountID return the JWT associated with the account identified by
// the uid. If no tokens are found an error is returned.
func (m *Manager) GetTokenByAccountID(uid string) (string, error) {
	tok, err := m.provider.GetTokenByAccountID(uid)
	if err != nil {
		logger.LogE("error getting JWT for account", "pkg", "manager", "func", "GetTokenByAccountID()", "uid", uid, "error", err.Error())
		return nil, err
	}
	return tok, nil
}

// SaveToken will save the given JWT on store.
func (m *Manager) SaveToken(tok *string) error {
	err := m.provider.SaveToken(tok)
	if err != nil {
		logger.LogD("error saving JWT", "pkg", "manager", "func", "SaveToken()", "error", err.Error())
		return err
	}
	return nil
}

// DeleteToken will remove the JWT identified by token id (kid) from the store.
func (m *Manager) DeleteToken(kid string) error {
	err := m.provider.DeleteToken(kid)
	if err != nil {
		logger.LogD("error deleting JWT", "pkg", "manager", "func", "DeleteToken()", "kid", kid, "error", err.Error())
		return err
	}
	return nil
}
