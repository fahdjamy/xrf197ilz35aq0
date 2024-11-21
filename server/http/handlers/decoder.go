package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"xrf197ilz35aq0/internal/constants"
)

type httpErr struct {
	status int
	msg    string
	err    error
}

func (e *httpErr) Error() string {
	return fmt.Sprintf("message=%s :: \n\terr=%s", e.msg, e.err)
}

func decodeJSONBody[T any](r *http.Request, dst *T) *httpErr {
	ct := r.Header.Get(constants.ContentType)
	if ct != constants.EMPTY {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != constants.ContentTypeJson {
			msg := fmt.Sprintf("Content-Type header is not %s", constants.ContentTypeJson)
			return &httpErr{status: http.StatusUnsupportedMediaType, msg: msg}
		}
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(dst)
	if err != nil {
		return parseBodyError(err)
	}

	return nil
}

func parseBodyError(err error) *httpErr {
	var syntaxError *json.SyntaxError
	var maxBytesError *http.MaxBytesError
	var unmarshalTypeError *json.UnmarshalTypeError
	var invalidUnmarshalError *json.InvalidUnmarshalError

	switch {
	// Syntax errors in the JSON
	case errors.As(err, &syntaxError):
		msg := fmt.Sprintf("Request contains badly-formed JSON (at position %d)", syntaxError.Offset)
		return &httpErr{status: http.StatusBadRequest, msg: msg}

	// In some circumstances Decode() may return an
	// io.ErrUnexpectedEOF error for syntax errors in the JSON. https://github.com/golang/go/issues/25956.
	case errors.Is(err, io.ErrUnexpectedEOF):
		msg := fmt.Sprintf("Request contains badly-formed JSON")
		return &httpErr{status: http.StatusBadRequest, msg: msg}

	// Catching error types like trying to assign a string in the
	// JSON request body to a int field.
	// interpolate the relevant field name and position into the error message
	case errors.Is(err, unmarshalTypeError):
		msg := fmt.Sprintf("Request contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
		return &httpErr{status: http.StatusBadRequest, msg: msg}

	// Catch the error caused by extra unexpected fields in the request
	// body. https://github.com/golang/go/issues/29035 regarding turning this into a sentinel error.
	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
		return &httpErr{status: http.StatusBadRequest, msg: msg}

	// An io.EOF error is returned by Decode() if the request body is empty.
	case errors.Is(err, io.EOF):
		msg := "Request body must not be empty"
		return &httpErr{status: http.StatusBadRequest, msg: msg}

	// Catch any error caused by the request body being too large.
	case errors.Is(err, maxBytesError):
		msg := "Request body must not be larger than 1MB"
		return &httpErr{status: http.StatusBadRequest, msg: msg}

	case errors.As(err, &invalidUnmarshalError):
		msg := "Request body must contain a valid JSON pointer"
		return &httpErr{status: http.StatusBadRequest, msg: msg}

	default:
		return &httpErr{
			status: http.StatusBadRequest,
			msg:    fmt.Sprintf("Internal server :: err=%s", err.Error()),
		}
	}
}
