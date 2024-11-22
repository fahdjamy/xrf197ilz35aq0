package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	xrf "xrf197ilz35aq0"
	"xrf197ilz35aq0/internal/constants"
)

type Response struct {
	Code       int         `json:"code"`
	Data       interface{} `json:"data,omitempty"`
	Pagination Pagination  `json:"pagination,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

type Pagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Total  int `json:"total"`
	Start  int `json:"start"`
}

func writeResponse(response Response, w http.ResponseWriter, logger xrf.Logger) {
	w.WriteHeader(response.Code)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set(constants.ContentType, constants.ContentTypeJson)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		logger.Error(fmt.Sprintf("event=writeResponse :: error encoding response: %v", err))
	}
}

func writeErrorResponse(error httpErr, w http.ResponseWriter, logger xrf.Logger) {
	errResp := ErrorResponse{Error: error.msg, Code: error.status}

	w.WriteHeader(error.status)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set(constants.ContentType, constants.ContentTypeJson)

	err := json.NewEncoder(w).Encode(errResp)
	if err != nil {
		logger.Error(fmt.Sprintf("error writing error response: %s", err))
	}
}
