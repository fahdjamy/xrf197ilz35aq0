package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
)

const (
	GET    = "GET"
	PUT    = "PUT"
	POST   = "POST"
	DELETE = "DELETE"
)

func getAndValidateId(req *http.Request, reqIdKey string) (string, bool) {
	vars := mux.Vars(req)
	idVal, ok := vars[reqIdKey]
	if !ok || idVal == "" {
		return "", false
	}

	return idVal, true
}
