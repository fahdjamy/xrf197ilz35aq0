package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	xrf "xrf197ilz35aq0"
	"xrf197ilz35aq0/internal/constants"
)

type dataResponse struct {
	Code       int         `json:"code"`
	Data       interface{} `json:"data,omitempty"`
	Pagination Pagination  `json:"pagination,omitempty"`
}

type errorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

type Pagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Total  int `json:"total"`
	Start  int `json:"start"`
}

func writeResponse(data dataResponse, w http.ResponseWriter, logger xrf.Logger) {
	w.Header().Set(constants.ContentType, constants.ContentTypeJson)
	w.WriteHeader(data.Code)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error(fmt.Sprintf("event=writeResponse :: error encoding response: %v", err))
	}
}

func writeErrorResponse(error httpErr, w http.ResponseWriter, logger xrf.Logger) {
	w.Header().Set(constants.ContentType, constants.ContentTypeJson)
	w.WriteHeader(error.status)

	errResp := errorResponse{Error: error.msg, Code: error.status}

	err := json.NewEncoder(w).Encode(errResp)
	if err != nil {
		logger.Error(fmt.Sprintf("error writing error response: %s", err))
	}
}
