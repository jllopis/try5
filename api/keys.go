package api

import (
	"net/http"

	"github.com/jllopis/try5/keys"
	logger "github.com/jllopis/try5/log"
	"github.com/jllopis/try5/store/manager"
	"github.com/jllopis/try5/tryerr"
	"github.com/labstack/echo"
)

// GetKey return a HandlerFunc that returns the key obtanined by specifying its
// key id (kid)
func GetKey(m *manager.Manager) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		var res *keys.Key
		var err error
		var kid string
		if kid = ctx.Param("kid"); kid == "" {
			return ctx.JSON(http.StatusInternalServerError, &logMessage{Status: "error", Action: "get", Info: "kid cannot be nil"})
		}
		if res, err = m.LoadKey(kid); err != nil {
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		return ctx.JSON(http.StatusOK, res)
	}
}

// GetKey return a HandlerFunc that returns the keys that exists in the store.
//
// If a public key (pubkey) param is specified as query param, the key matching
// the public key (pubkey) is returned.
//    curl -k -H "Origin: http://b2d" https://localhost:9000/api/v1/keys -G -d @pk.txt
func GetAllKeys(m *manager.Manager) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		if pubkey := ctx.Query("pubkey"); pubkey != "" {
			var res *keys.Key
			var err error
			if res, err = m.GetKeyByPub([]byte(pubkey)); err != nil {
				return ctx.JSON(http.StatusInternalServerError, err.Error())
			}
			return ctx.JSON(http.StatusOK, res)
		}
		var res []*keys.Key
		var err error
		if res, err = m.LoadAllKeys(); err != nil {
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		return ctx.JSON(http.StatusOK, res)
	}
}

// DeleteKey will delete the specified key identified by its key id (kid) from
// the store.
func DeleteKey(m *manager.Manager) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		var kid string
		if kid = ctx.Param("kid"); kid == "" {
			return ctx.JSON(http.StatusBadRequest, &logMessage{Status: "error", Action: "delete", Info: "kid cannot be nil", Table: "keys"})
		}
		if err := m.DeleteKey(kid); err != nil {
			switch err {
			case tryerr.ErrKeyNotFound:
				logger.LogE("key delete error", "p", "api", "f", "DeleteKey()", "error", err.Error(), "kid", kid)
				return ctx.JSON(http.StatusNotFound, &logMessage{Status: "error", Action: "delete", Info: "no se ha encontrado el registro", Table: "keys", UID: kid})
			default:
				logger.LogE("error deleting key", "p", "api", "f", "DeleteKey()", "error", err)
				return ctx.JSON(http.StatusInternalServerError, err.Error())
			}
		} else {
			logger.LogD("func DeleteKey", "registro eliminado", kid)
			return ctx.JSON(http.StatusOK, &logMessage{Status: "ok", Action: "delete", Table: "keys", UID: kid})
		}
	}
}
