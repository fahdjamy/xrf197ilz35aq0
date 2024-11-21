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
	Error      string      `json:"error,omitempty"`
	Pagination Pagination  `json:"pagination,omitempty"`
}

type Pagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Total  int `json:"total"`
	Start  int `json:"start"`
}

func writeResponse(response Response, w http.ResponseWriter, logger xrf.Logger) {
	w.WriteHeader(response.Code)
	w.Header().Set(constants.ContentType, constants.ContentTypeJson)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		logger.Error(fmt.Sprintf("error encoding response: %v", err))
	}
}

func writeErrorResponse(error httpErr, w http.ResponseWriter, logger xrf.Logger) {
	response := Response{Error: error.msg, Code: error.status}

	data, _ := json.Marshal(response)

	w.WriteHeader(error.status)
	w.Header().Set(constants.ContentType, constants.ContentTypeJson)

	_, err := w.Write(data)
	if err != nil {
		logger.Error(fmt.Sprintf("error writing response: %s", err))
	}
}
