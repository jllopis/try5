package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/jllopis/aloja"
	"github.com/jllopis/try5/account"
)

// GetAllAccounts devuelve una lista con todos los accounts de la base de datos
// curl -ks https://b2d:8000/v1/accounts | jp -
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
// curl -ks https://b2d:8000/v1/accounts/1 | jp -
func (ctx *ApiContext) GetAccountByID(w http.ResponseWriter, r *http.Request) {
	var res *account.Account
	var err error
	var id string
	if id = aloja.Params(r).ByName("id"); id == "" {
		ctx.Render.JSON(w, http.StatusInternalServerError, &logMessage{Status: "error", Action: "get", Info: "missing id"})
		return
	}
	if res, err = ctx.DB.LoadAccount(id); err != nil {
		if err == sql.ErrNoRows {
			errMsg := fmt.Sprintf("El item con id=%v no se ha encontrado", id)
			ctx.Render.JSON(w, http.StatusNotFound, errMsg)
			//ctx.Render.JSON(w, http.StatusNotFound, &logMessage{Status: "error", Action: "get", Info: err.Detail, Table: err.Table, Code: string(err.Code), ID: id})
			logger.Error("func GetAccountByID", "error", "account no encontrado", "id", id)
			return
		}
		ctx.Render.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Render.JSON(w, http.StatusOK, res)
}

//// NewAccount crea un nuevo account.
//// curl -k https://b2d:8000/v1/accounts -X POST -d '{}'
//func (ctx *ApiContext) NewAccount(w http.ResponseWriter, r *http.Request) {
//	var data account.Account
//	err := json.NewDecoder(r.Body).Decode(&data)
//	if err != nil {
//		ctx.Render.JSON(w, http.StatusInternalServerError, &logMessage{Status: "error", Action: "create", Info: err.Error(), Table: "accounts"})
//		return
//	}
//	if outdata, err := ctx.DB.NewAccount(data); err != nil {
//		ctx.Render.JSON(w, http.StatusInternalServerError, &logMessage{Status: "error", Action: "create", Info: err.Detail, Table: err.Table, Code: string(err.Code)})
//		logger.Error("func NewAccount", "error", err)
//		return
//	} else {
//		ctx.Render.JSON(w, http.StatusOK, outdata)
//	}
//}
//
//// UpdateAccount actualiza los datos del account y devuelve el objeto actualizado.
//// curl -ks https://b2d:8000/v1/accounts/3 -X PUT -d '{}' | jp -
//func (ctx *ApiContext) UpdateAccount(w http.ResponseWriter, r *http.Request) {
//	newdata := &account.Account{}
//	var err error
//	var id int64
//	if id, err = strconv.ParseInt(aloja.Params(r).ByName("id"), 10, 64); err != nil {
//		ctx.Render.JSON(w, http.StatusInternalServerError, &logMessage{Status: "error", Action: "update", Info: err.Error(), Table: "accounts"})
//		return
//	}
//	err = json.NewDecoder(r.Body).Decode(&newdata)
//	if err != nil {
//		ctx.Render.JSON(w, http.StatusInternalServerError, &logMessage{Status: "error", Action: "update", Info: err.Error(), Table: "accounts"})
//		logger.Error("func UpdateAccount", "error", err.Error())
//		return
//	}
//	if logger.IsDebug() {
//		logger.Info("func UpdateAccount", "updated register", id)
//	}
//
//	if newdata.ID == nil {
//		newdata.ID = &id
//	} else {
//		if *newdata.ID != int64(id) {
//			ctx.Render.JSON(w, http.StatusInternalServerError, fmt.Sprintf("los identificadores de registro no coindiden: body: %v - path: %v", newdata.ID, id))
//			logger.Error("func UpdateAccount", "error", err.Error())
//			return
//		}
//	}
//	if _, err := ctx.DB.Updateaccount(*newdata); err != nil {
//		ctx.Render.JSON(w, http.StatusInternalServerError, err.Error())
//		logger.Error("func UpdateAccount", "error", err.Error())
//		return
//	} else {
//		logger.Info("func UpdateAccount", "updated", "ok", "id", newdata.ID)
//		ctx.Render.JSON(w, http.StatusOK, newdata)
//		return
//	}
//}
//
//// DeleteAccount elimina el account solicitado.
//// curl -ks https://b2d:8000/v1/accounts/3 -X DELETE | jp -
//func (ctx *ApiContext) DeleteAccount(w http.ResponseWriter, r *http.Request) {
//	var id int64
//	var err error
//	if id, err = strconv.ParseInt(aloja.Params(r).ByName("id"), 10, 64); err != nil {
//		ctx.Render.JSON(w, http.StatusInternalServerError, err.Error())
//		return
//	}
//	if n, err := ctx.DB.DeleteAccount(id); err != nil {
//		ctx.Render.JSON(w, http.StatusInternalServerError, err.Error())
//		logger.Error("func DeleteAccount", "error", err)
//		return
//	} else {
//		switch n {
//		case 0:
//			logger.Info("func DeleteAccount", "error", "id no encontrado", "id", id)
//			ctx.Render.JSON(w, http.StatusOK, &logMessage{Status: "error", Action: "delete", Info: "no se ha encontrado el registro", Table: "accounts", Code: "RNF-11", ID: id})
//		default:
//			logger.Info("func DeleteAccount", "registro eliminado", id)
//			ctx.Render.JSON(w, http.StatusOK, &logMessage{Status: "ok", Action: "delete", Info: "eliminado registro", Table: "accounts", ID: id})
//		}
//		return
//	}
//}
