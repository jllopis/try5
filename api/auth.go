package api

import (
	"net/http"

	"github.com/jllopis/try5/account"
)

func (ctx *ApiContext) Authenticate(w http.ResponseWriter, r *http.Request) {
	var res *account.Account
	var err error
	var email, password string
	if email = r.FormValue("email"); email == "" {
		ctx.Render.JSON(w, http.StatusBadRequest, &logMessage{Status: "error", Action: "authenticate", Info: "email cannot be nil"})
		return
	}
	if password = r.FormValue("password"); password == "" {
		ctx.Render.JSON(w, http.StatusBadRequest, &logMessage{Status: "error", Action: "authenticate", Info: "password cannot be nil"})
		return
	}
	if res, err = ctx.DB.GetAccountByEmail(email); err != nil {
		logger.Error("func Authenticate", "error", "account no encontrado", "email", email)
		ctx.Render.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = res.MatchPassword(password)
	if err != nil {
		ctx.Render.JSON(w, http.StatusForbidden, map[string]interface{}{"status": "fail", "reason": err.Error()})
		return
	}
	res.Password = nil
	ctx.Render.JSON(w, http.StatusOK, map[string]interface{}{"status": "ok", "account": res})
}
