package api

import (
	"github.com/gorilla/securecookie"
	"github.com/jllopis/try5/store"
	"github.com/mgutz/logxi/v1"
	"github.com/unrolled/render"
)

type ApiContext struct {
	DB            store.Storer
	Render        *render.Render
	CookieHandler *securecookie.SecureCookie
}

type logMessage struct {
	Status string `json:"status"`
	Action string `json:"action"`
	Info   string `json:"info,omitempty"`
	Table  string `json:"table,omitempty"`
	Code   string `json:"code,omitempty"`
	UID    string `json:"id,omitempty"`
}

var logger log.Logger

func init() {
	logger = log.New("try5")
}
