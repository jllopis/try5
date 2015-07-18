package api

import (
	"net/http"

	"github.com/jllopis/try5/account"
	logger "github.com/jllopis/try5/log"
	"github.com/jllopis/try5/store/manager"
	"github.com/labstack/echo"
)

// Authenticate get the email and passowrd as query or form params and return
// a JSON object indicating if the authentication succed.
func Authenticate(m *manager.Manager) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		var res *account.Account
		var err error
		var email, password string
		if email = ctx.Form("email"); email == "" {
			return ctx.JSON(http.StatusBadRequest, &logMessage{Status: "error", Action: "authenticate", Info: "email cannot be nil"})
		}
		if password = ctx.Form("password"); password == "" {
			return ctx.JSON(http.StatusBadRequest, &logMessage{Status: "error", Action: "authenticate", Info: "password cannot be nil"})
		}
		if res, err = m.GetAccountByEmail(email); err != nil {
			logger.LogE("missing account", "error", err.Error(), "email", email)
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}

		err = res.MatchPassword(password)
		if err != nil {
			return ctx.JSON(http.StatusForbidden, map[string]interface{}{"status": "fail", "reason": err.Error()})
		}
		res.Password = nil
		return ctx.JSON(http.StatusOK, map[string]interface{}{"status": "ok", "account": res})
	}
}
