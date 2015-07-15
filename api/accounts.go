package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jllopis/aloja"
	"github.com/jllopis/try5/account"
	"github.com/jllopis/try5/keys"
	"github.com/jllopis/try5/store"
	"github.com/lib/pq"
)

// GetAllAccounts devuelve una lista con todos los accounts de la base de datos
//   curl -ks https://b2d:8000/v1/accounts | jp -
func (ctx *ApiContext) GetAllAccounts(w http.ResponseWriter, r *http.Request) {
	var res []*account.Account
	var err error
	if res, err = ctx.DB.LoadAllAccounts(); err != nil {
		ctx.Render.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Render.JSON(w, http.StatusOK, res)
}

// GetAccountByID devuelve el account de la base de datos que coincide con el ID suministrado
//   curl -ks https://b2d:8000/v1/accounts/342947fd-6c4b-4d2b-85ab-da14b37d047a | jp -
func (ctx *ApiContext) GetAccountByID(w http.ResponseWriter, r *http.Request) {
	var res *account.Account
	var err error
	var uid string
	if uid = aloja.Params(r).ByName("uid"); uid == "" {
		ctx.Render.JSON(w, http.StatusInternalServerError, &logMessage{Status: "error", Action: "get", Info: "uid cannot be nil"})
		return
	}
	if res, err = ctx.DB.LoadAccount(uid); err != nil {
		logger.Info("GetAccountByID", "error", "account not found", "uid", uid)
		ctx.Render.JSON(w, http.StatusNotFound, &logMessage{Status: "error", Action: "get", Info: err.Error(), Table: "accounts", UID: uid})
		return
	}
	ctx.Render.JSON(w, http.StatusOK, res)
}

// NewAccount crea un nuevo account.
//   curl -k https://b2d:8000/v1/accounts -X POST -d '{"email":"tu2@test.com","name":"test user 2","password":"1234","active":true}'
func (ctx *ApiContext) NewAccount(w http.ResponseWriter, r *http.Request) {
	var data account.Account
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		ctx.Render.JSON(w, http.StatusInternalServerError, &logMessage{Status: "error", Action: "create", Info: err.Error(), Table: "accounts"})
		return
	}
	if err = data.ValidateFields(); err != nil {
		ctx.Render.JSON(w, http.StatusBadRequest, &logMessage{Status: "error", Action: "create", Info: err.Error(), Table: "accounts"})
		return
	}
	if outdata, err := ctx.DB.SaveAccount(&data); err != nil {
		switch err {
		case store.ErrDupEmail:
			ctx.Render.JSON(w, http.StatusConflict, &logMessage{Status: "error", Action: "create", Info: err.Error(), Table: "accounts", Code: "209"})
		case err.(*pq.Error):
			ctx.Render.JSON(w, http.StatusInternalServerError, &logMessage{Status: "error", Action: "create", Info: err.(*pq.Error).Detail, Table: err.(*pq.Error).Table, Code: string(err.(*pq.Error).Code)})
		default:
			ctx.Render.JSON(w, http.StatusInternalServerError, &logMessage{Status: "error", Action: "create", Info: err.Error(), Table: "accounts", Code: "500"})
		}
		logger.Error("func NewAccount", "error", err)
		return
	} else {
		// Generate and assign a key pair
		if creat := r.FormValue("wkey"); creat == "" {
			k := keys.New(*outdata.UID)
			if _, err := ctx.DB.SaveKey(k); err != nil {
				ctx.Render.JSON(w, http.StatusInternalServerError, &logMessage{Status: "error", Action: "create", Info: err.Error(), Table: "keys", Code: "500"})
				logger.Error("func NewAccount/keys.New", "error", err)
				return
			}
		}
		// and return the created account
		ctx.Render.JSON(w, http.StatusCreated, outdata)
	}
}

// UpdateAccount actualiza los datos del account y devuelve el objeto actualizado.
//   curl -ks https://b2d:8000/v1/accounts/342947fd-6c4b-4d2b-85ab-da14b37d047a -X PUT -d '{}' | jp -
func (ctx *ApiContext) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	var newdata account.Account
	var err error
	var uid string
	if uid = aloja.Params(r).ByName("uid"); uid == "" {
		ctx.Render.JSON(w, http.StatusInternalServerError, &logMessage{Status: "error", Action: "update", Info: "uid cannot be nil", Table: "accounts"})
		return
	}
	err = json.NewDecoder(r.Body).Decode(&newdata)
	if err != nil {
		ctx.Render.JSON(w, http.StatusInternalServerError, &logMessage{Status: "error", Action: "update", Info: err.Error(), Table: "accounts"})
		logger.Error("func UpdateAccount", "error", err.Error())
		return
	}
	if err = newdata.ValidateFields(); err != nil {
		ctx.Render.JSON(w, http.StatusBadRequest, &logMessage{Status: "error", Action: "update", Info: err.Error(), Table: "accounts"})
		return
	}
	if logger.IsDebug() {
		logger.Info("func UpdateAccount", "updated register", uid)
	}

	// TODO if newdata.UID is nil check if register exist in db _before_ update.
	// If do not exist, return error and quit
	if newdata.UID == nil {
		newdata.UID = &uid
	} else {
		if *newdata.UID != uid {
			ctx.Render.JSON(w, http.StatusInternalServerError, &logMessage{Status: "error",
				Action: "update",
				Info:   fmt.Sprintf("provided uid's does not match: body: %v - path: %v", *newdata.UID, uid),
				Table:  "accounts"})
			logger.Error("func UpdateAccount", "error", "uid's does not match", "body", *newdata.UID, "path", uid)
			return
		}
	}
	if _, err := ctx.DB.SaveAccount(&newdata); err != nil {
		ctx.Render.JSON(w, http.StatusInternalServerError, err.Error())
		logger.Error("func UpdateAccount", "error", err.Error())
		return
	} else {
		logger.Info("func UpdateAccount", "updated", "ok", "uid", *newdata.UID)
		ctx.Render.JSON(w, http.StatusOK, newdata)
		return
	}
}

// DeleteAccount elimina el account solicitado.
//   curl -ks https://b2d:8000/v1/accounts/3 -X DELETE | jp -
func (ctx *ApiContext) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	var uid string
	if uid = aloja.Params(r).ByName("uid"); uid == "" {
		ctx.Render.JSON(w, http.StatusInternalServerError, &logMessage{Status: "error", Action: "delete", Info: "uid cannot be nil", Table: "accounts"})
		return
	}
	if n, err := ctx.DB.DeleteAccount(uid); err != nil {
		ctx.Render.JSON(w, http.StatusInternalServerError, err.Error())
		logger.Error("func DeleteAccount", "error", err)
		return
	} else {
		switch n {
		case 0:
			logger.Info("func DeleteAccount", "error", "uid no encontrado", "uid", uid)
			ctx.Render.JSON(w, http.StatusOK, &logMessage{Status: "error", Action: "delete", Info: "no se ha encontrado el registro", Table: "accounts", Code: "RNF-11", UID: uid})
		default:
			logger.Info("func DeleteAccount", "registro eliminado", uid)
			ctx.Render.JSON(w, http.StatusOK, &logMessage{Status: "ok", Action: "delete", Info: uid, Table: "accounts", UID: uid})
		}
		return
	}
}

// GetAccountKeys obtiene la pareja de claves RSA pública y privada de la cuenta
//   curl -ksi https://b2d:8000/v1/accounts/3/keys | jp -
func (ctx *ApiContext) GetAccountKeys(w http.ResponseWriter, r *http.Request) {
	var uid string
	if uid = aloja.Params(r).ByName("uid"); uid == "" {
		ctx.Render.JSON(w, http.StatusInternalServerError, &logMessage{Status: "error", Action: "get", Info: "uid cannot be nil"})
		return
	}
	var res *keys.Key
	var err error
	if res, err = ctx.DB.GetKeyByAccountID(uid); err != nil {
		ctx.Render.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Render.JSON(w, http.StatusOK, res)
}

// GetAccountTokens obtiene la pareja de claves RSA pública y privada de la cuenta
//   curl -ksi https://b2d:8000/v1/accounts/3/tokens | jp -
func (ctx *ApiContext) GetAccountTokens(w http.ResponseWriter, r *http.Request) {
	var uid string
	if uid = aloja.Params(r).ByName("uid"); uid == "" {
		ctx.Render.JSON(w, http.StatusInternalServerError, &logMessage{Status: "error", Action: "get", Info: "uid cannot be nil"})
		return
	}
	var res string
	var err error
	if res, err = ctx.DB.GetTokenByAccountID(uid); err != nil {
		ctx.Render.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Render.JSON(w, http.StatusOK, res)
}
