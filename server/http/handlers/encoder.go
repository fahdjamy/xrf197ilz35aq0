package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	xrf "xrf197ilz35aq0"
	"xrf197ilz35aq0/internal/constants"
	xrfErr "xrf197ilz35aq0/internal/error"
)

type dataResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
}

type errorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

type pagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Total  int `json:"total"`
	Start  int `json:"start"`
}

func writeResponse(data dataResponse, w http.ResponseWriter, logger xrf.Logger) {
	writePaginatedResponse(data, nil, w, logger)
}

func writePaginatedResponse(data dataResponse, pag *pagination, w http.ResponseWriter, logger xrf.Logger) {
	w.Header().Set(constants.ContentType, constants.ContentTypeJson)
	w.WriteHeader(data.Code)

	if pag == nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			logger.Error(fmt.Sprintf("event=writeResponseFailure :: error encoding response: %v", err))
		}
	} else {
		err := json.NewEncoder(w).Encode(struct {
			*pagination
			dataResponse
		}{})
		if err != nil {
			logger.Error(fmt.Sprintf("event=writePaginatedResponseFailure :: error encoding response: %v", err))
		}
	}
}

func writeErrorResponse(error error, w http.ResponseWriter, logger xrf.Logger) {
	msg := "Something went wrong"
	statusCode := http.StatusInternalServerError

	var decoderError *decoderErr
	var internalError *xrfErr.Internal
	var externalError *xrfErr.External

	switch {
	case errors.As(error, &decoderError):
		var decErr *decoderErr
		errors.As(error, &decErr)
		statusCode = decErr.status
		msg = decErr.msg
	case errors.As(error, &internalError):
		var internalErr *xrfErr.Internal
		errors.As(error, &internalErr)
	case errors.As(error, &externalError):
		var externalErr *xrfErr.External
		errors.As(error, &externalErr)
		statusCode = http.StatusInternalServerError
		msg = externalErr.Message
	default:
		statusCode = http.StatusInternalServerError
		msg = "Something went wrong"
	}

	w.Header().Set(constants.ContentType, constants.ContentTypeJson)
	w.WriteHeader(statusCode)

	errResp := errorResponse{Error: msg, Code: statusCode}

	err := json.NewEncoder(w).Encode(errResp)
	if err != nil {
		logger.Error(fmt.Sprintf("error writing error response: %s", err))
	}
}
