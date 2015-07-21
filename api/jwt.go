package api

import (
	"net/http"
	"strings"

	"github.com/jllopis/try5/jwt"
	"github.com/jllopis/try5/keys"
	logger "github.com/jllopis/try5/log"
	"github.com/jllopis/try5/store/manager"
	"github.com/labstack/echo"
)

/*
NewJWTToken genera un nuevo JSON Web Token para el `uid` suministrado.
El parámetro `uid` es obligatorio y debe pasarse como variable de path.

El token se guardará en la caché mientras dure su TTL ('exp' claim).

Ejemplo:
  curl -k https://b2d:8000/v1/jwt/token/9a77ca50-8859-4f27-9ae6-9f3b6ab92936 -X POST
*/
func NewJWTToken(m *manager.Manager) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		var uid string
		if uid = ctx.Param("uid"); uid == "" {
			return ctx.JSON(http.StatusBadRequest, &logMessage{Status: "error", Action: "NewJWTToken", Info: "user id (uid) cannot be nil"})
		}
		var res *keys.Key
		var err error
		if res, err = m.GetKeyByAccountID(uid); err != nil {
			logger.LogE("key not found", "error", err.Error(), "uid", uid)
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		token, err := jwt.GenerateToken(uid, res.PrivKey)
		if err != nil {
			logger.LogE("can not create token", "error", err.Error(), "uid", uid)
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		err = m.SaveToken(uid, &token)
		if err != nil {
			logger.LogE("can not save token", "pkg", "api", "func", "NewJWTToken()", "error", err.Error())
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		return ctx.JSON(http.StatusOK, map[string]string{"token": token})
	}
}

/*
ValidateToken recibe un JSON Web Token en la variable `jwt` y lo valida.

La variable `jwt` puede pasarse por path, cabecera authorization, form value o cookie.
Si no se suministra, se devuelve error.

EL token se valida con la clave pública RSA de la cuenta puesto que se genera y firma con
la clave privada correspondiente a la cuenta.
*/
func ValidateToken(m *manager.Manager) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		// Get token from query or form params
		tokenStr := ctx.Form("jwt")

		// Get token from URL params
		if tokenStr == "" {
			tokenStr = ctx.Query("jwt")
		}

		// Get token from authorization header
		if tokenStr == "" {
			bearer := ctx.Request().Header.Get("Authorization")
			if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
				tokenStr = bearer[7:]
			}
		}

		// Get token from cookie
		if tokenStr == "" {
			cookie, err := ctx.Request().Cookie("jwt")
			if err == nil {
				tokenStr = cookie.Value
			}
		}

		// No jwt token found
		if tokenStr == "" {
			return ctx.JSON(http.StatusInternalServerError, &logMessage{Status: "error", Action: "unauthorized", Info: "jwt token is required"})
		}

		logger.LogD("ValidateToken", "jwt", tokenStr)
		if err := jwt.Validate(m, tokenStr); err != nil {
			return ctx.JSON(http.StatusForbidden, map[string]interface{}{"status": "fail", "reason": err.Error()})
		}
		return ctx.JSON(http.StatusOK, map[string]interface{}{"status": "ok", "jwt": tokenStr})
	}
}

// GetAccountJWTToken returns the JWT associated with the account. A valid uid must be provided.
func GetAccountJWTToken(m *manager.Manager) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		var uid string
		if uid = ctx.Param("uid"); uid == "" {
			return ctx.JSON(http.StatusBadRequest, &logMessage{Status: "error", Action: "NewJWTToken", Info: "user id (uid) cannot be nil"})
		}
		var token string
		var err error
		if token, err = m.GetTokenByAccountID(uid); err != nil {
			return ctx.JSON(http.StatusInternalServerError, &logMessage{Status: "error", Info: err.Error()})
		}
		return ctx.JSON(http.StatusOK, map[string]string{"token": token})
	}
}
