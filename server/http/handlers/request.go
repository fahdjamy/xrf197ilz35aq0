package handlers

import (
	"net/http"
)

func getAndValidateId(req *http.Request, reqIdKey string) (string, bool) {
	value := req.PathValue(reqIdKey)
	if value == "" {
		return "", false
	}
	return value, true
}
