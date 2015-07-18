package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jllopis/try5/account"
	"github.com/jllopis/try5/keys"
	logger "github.com/jllopis/try5/log"
	"github.com/jllopis/try5/store/manager"
	"github.com/jllopis/try5/tryerr"
	"github.com/labstack/echo"
	"github.com/lib/pq"
)

// GetAllAccounts devuelve una lista con todos los accounts de la base de datos
//   curl -ks https://b2d:8000/v1/accounts | jp -
func GetAllAccounts(m *manager.Manager) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		var res []*account.Account
		var err error
		if res, err = m.LoadAllAccounts(); err != nil {
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		return ctx.JSON(http.StatusOK, res)
	}
}

// GetAccountByID devuelve el account de la base de datos que coincide con el ID suministrado
//   curl -ks https://b2d:8000/v1/accounts/342947fd-6c4b-4d2b-85ab-da14b37d047a | jp -
func GetAccountByID(m *manager.Manager) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		var res *account.Account
		var err error
		var uid string
		if uid = ctx.Param("uid"); uid == "" {
			return ctx.JSON(http.StatusInternalServerError, &logMessage{Status: "error", Action: "get", Info: "uid cannot be nil"})
		}
		if res, err = m.LoadAccount(uid); err != nil {
			logger.LogE("GetAccountByID", "error", "account not found", "uid", uid)
			return ctx.JSON(http.StatusNotFound, &logMessage{Status: "error", Action: "get", Info: err.Error(), Table: "accounts", UID: uid})
		}
		return ctx.JSON(http.StatusOK, res)
	}
}

// NewAccount crea un nuevo account.
//   curl -k https://b2d:8000/v1/accounts -X POST -d '{"email":"tu2@test.com","name":"test user 2","password":"1234","active":true}'
func NewAccount(m *manager.Manager) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		var data account.Account
		err := json.NewDecoder(ctx.Request().Body).Decode(&data)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, &logMessage{Status: "error", Action: "create", Info: err.Error(), Table: "accounts"})
		}
		if err = data.ValidateFields(); err != nil {
			return ctx.JSON(http.StatusBadRequest, &logMessage{Status: "error", Action: "create", Info: err.Error(), Table: "accounts"})
		}
		if err := m.SaveAccount(&data); err != nil {
			logger.LogE("func NewAccount", "error", err)
			switch err {
			case tryerr.ErrDupEmail:
				return ctx.JSON(http.StatusConflict, &logMessage{Status: "error", Action: "create", Info: err.Error(), Table: "accounts", Code: "209"})
			case err.(*pq.Error):
				return ctx.JSON(http.StatusInternalServerError, &logMessage{Status: "error", Action: "create", Info: err.(*pq.Error).Detail, Table: err.(*pq.Error).Table, Code: string(err.(*pq.Error).Code)})
			default:
				return ctx.JSON(http.StatusInternalServerError, &logMessage{Status: "error", Action: "create", Info: err.Error(), Table: "accounts", Code: "500"})
			}
		} else {
			logger.LogD("creating keypair for account", "pkg", "api", "func", "NewAccount()", "uid", *data.UID)
			// Generate and assign a key pair
			if creat := ctx.Form("nokey"); creat == "" {
				k := keys.New(*data.UID)
				if err := m.SaveKey(k); err != nil {
					logger.LogE("func NewAccount/keys.New", "error", err)
					return ctx.JSON(http.StatusInternalServerError, &logMessage{Status: "error", Action: "create", Info: err.Error(), Table: "keys", Code: "500"})
				}
				logger.LogD("keys created", "pkg", "api", "func", "NewAccount()", "kid", *k.KID)
			}
			// and return the created account
			return ctx.JSON(http.StatusCreated, data)
		}
	}
}

// UpdateAccount actualiza los datos del account y devuelve el objeto actualizado.
//   curl -ks https://b2d:8000/v1/accounts/342947fd-6c4b-4d2b-85ab-da14b37d047a -X PUT -d '{}' | jp -
func UpdateAccount(m *manager.Manager) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		var newdata account.Account
		var err error
		var uid string
		if uid = ctx.Param("uid"); uid == "" {
			return ctx.JSON(http.StatusInternalServerError, &logMessage{Status: "error", Action: "update", Info: "uid cannot be nil", Table: "accounts"})
		}
		err = json.NewDecoder(ctx.Request().Body).Decode(&newdata)
		if err != nil {
			logger.LogE("func UpdateAccount", "p", "api", "f", "UpdateAccount()", "error", err.Error())
			return ctx.JSON(http.StatusInternalServerError, &logMessage{Status: "error", Action: "update", Info: err.Error(), Table: "accounts"})
		}
		if err = newdata.ValidateFields(); err != nil {
			return ctx.JSON(http.StatusBadRequest, &logMessage{Status: "error", Action: "update", Info: err.Error(), Table: "accounts"})
		}
		logger.LogD("updated register", "p", "api", "f", "UpdateAccount()", uid)

		// TODO if newdata.UID is nil check if register exist in db _before_ update.
		// If do not exist, return error and quit
		if newdata.UID == nil {
			newdata.UID = &uid
		} else {
			if *newdata.UID != uid {
				logger.LogE("uid's does not match", "p", "api", "f", "UpdateAccount()", "body", *newdata.UID, "path", uid)
				return ctx.JSON(http.StatusInternalServerError, &logMessage{Status: "error",
					Action: "update",
					Info:   fmt.Sprintf("provided uid's does not match: body: %v - path: %v", *newdata.UID, uid),
					Table:  "accounts"})
			}
		}
		if err := m.SaveAccount(&newdata); err != nil {
			logger.LogE("func UpdateAccount", "p", "api", "f", "UpdateAccount()", "error", err.Error())
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		logger.LogD("updated account", "p", "api", "f", "UpdateAccount()", "uid", *newdata.UID)
		return ctx.JSON(http.StatusOK, newdata)
	}
}

// DeleteAccount elimina el account solicitado.
//   curl -ks https://b2d:8000/v1/accounts/3 -X DELETE | jp -
func DeleteAccount(m *manager.Manager) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		var uid string
		if uid = ctx.Param("uid"); uid == "" {
			return ctx.JSON(http.StatusBadRequest, &logMessage{Status: "error", Action: "delete", Info: "uid cannot be nil", Table: "accounts"})
		}
		if err := m.DeleteAccount(uid); err != nil {
			switch err {
			case tryerr.ErrAccountNotFound:
				logger.LogE("account delete error", "p", "api", "f", "DeleteAccount()", "error", err.Error(), "uid", uid)
				return ctx.JSON(http.StatusOK, &logMessage{Status: "error", Action: "delete", Info: "no se ha encontrado el registro", Table: "accounts", Code: "RNF-11", UID: uid})
			default:
				logger.LogE("error deleting account", "p", "api", "f", "DeleteAccount()", "error", err)
				return ctx.JSON(http.StatusInternalServerError, err.Error())
			}
		} else {
			logger.LogD("func DeleteAccount", "registro eliminado", uid)
			return ctx.JSON(http.StatusOK, &logMessage{Status: "ok", Action: "delete", Info: uid, Table: "accounts", UID: uid})
		}
	}
}

// GetAccountKeys obtiene la pareja de claves RSA pública y privada de la cuenta
//   curl -ksi https://b2d:8000/v1/accounts/3/keys | jp -
func GetAccountKeys(m *manager.Manager) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		var uid string
		if uid = ctx.Param("uid"); uid == "" {
			return ctx.JSON(http.StatusInternalServerError, &logMessage{Status: "error", Action: "get", Info: "uid cannot be nil"})
		}
		var res *keys.Key
		var err error
		if res, err = m.GetKeyByAccountID(uid); err != nil {
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		return ctx.JSON(http.StatusOK, res)
	}
}

// GetAccountTokens obtiene la pareja de claves RSA pública y privada de la cuenta
//   curl -ksi https://b2d:8000/v1/accounts/3/tokens | jp -
func GetAccountTokens(m *manager.Manager) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		var uid string
		if uid = ctx.Param("uid"); uid == "" {
			return ctx.JSON(http.StatusInternalServerError, &logMessage{Status: "error", Action: "get", Info: "uid cannot be nil"})
		}
		var res string
		var err error
		if res, err = m.GetTokenByAccountID(uid); err != nil {
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		return ctx.JSON(http.StatusOK, res)
	}
}
