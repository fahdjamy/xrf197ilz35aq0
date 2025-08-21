package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	xrf "xrf197ilz35aq0/internal"
	"xrf197ilz35aq0/internal/constants"
	xrfErr "xrf197ilz35aq0/internal/error"
	"xrf197ilz35aq0/server/http/decoder"
)

type DataResponse struct {
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

func WriteResponse(data DataResponse, w http.ResponseWriter, logger xrf.Logger) {
	WritePaginatedResponse(data, nil, w, logger)
}

func WritePaginatedResponse(data DataResponse, pag *pagination, w http.ResponseWriter, logger xrf.Logger) {
	w.Header().Set(constants.ContentType, constants.ApplicationJson)
	w.WriteHeader(data.Code)

	if pag == nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			logger.Error(fmt.Sprintf("event=writeResponseFailure :: error encoding response: %v", err))
		}
	} else {
		err := json.NewEncoder(w).Encode(struct {
			*pagination
			DataResponse
		}{})
		if err != nil {
			logger.Error(fmt.Sprintf("event=writePaginatedResponseFailure :: error encoding response: %v", err))
		}
	}
}

func WriteErrorResponse(error error, w http.ResponseWriter, logger xrf.Logger) {
	msg := "Something went wrong"
	statusCode := http.StatusInternalServerError

	var decoderError *decoder.Err
	var internalError *xrfErr.Internal
	var externalError *xrfErr.External

	switch {
	case errors.As(error, &decoderError):
		var decErr *decoder.Err
		errors.As(error, &decErr)
		statusCode = decErr.Status
		msg = decErr.Msg
	case errors.As(error, &internalError):
		var internalErr *xrfErr.Internal
		errors.As(error, &internalErr)
	case errors.As(error, &externalError):
		if externalError.Code != 0 && externalError.Code >= 400 {
			statusCode = externalError.Code
		} else {
			statusCode = externalErrorCode(externalError.Message)
		}
		msg = externalError.Message
	case errors.Is(error, xrfErr.InvalidXrfToXrfTokenError):
		statusCode = http.StatusNotFound
		msg = "Not found"
	default:
		statusCode = http.StatusInternalServerError
		msg = "Something went wrong"
	}

	w.Header().Set(constants.ContentType, constants.ApplicationJson)
	w.WriteHeader(statusCode)

	errResp := errorResponse{Error: msg, Code: statusCode}

	err := json.NewEncoder(w).Encode(errResp)
	if err != nil {
		logger.Error(fmt.Sprintf("error writing error response: %s", err))
	}
}

func externalErrorCode(errorMessage string) int {
	switch errorMessage {
	case constants.NotFoundOrgErrMsg:
		return http.StatusNotFound
	default:
		return http.StatusBadRequest
	}
}
