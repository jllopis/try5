package api

import (
	"net/http"
	"strconv"
	"time"
)

type logMessage struct {
	Status string `json:"status"`
	Action string `json:"action"`
	Info   string `json:"info,omitempty"`
	Table  string `json:"table,omitempty"`
	Code   string `json:"code,omitempty"`
	UID    string `json:"id,omitempty"`
}

// add default service to give server time
func Time(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(strconv.FormatInt(time.Now().UTC().UnixNano(), 10)))
}
