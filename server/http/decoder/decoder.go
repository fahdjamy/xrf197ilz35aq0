package decoder

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"xrf197ilz35aq0/internal/constants"
	xrfErr "xrf197ilz35aq0/internal/error"
)

type Err struct {
	Status int
	Msg    string
	Err    error
}

func (e *Err) Error() string {
	return fmt.Sprintf("message=%s :: \n\tErr=%s", e.Msg, e.Err)
}

func DecodeJSONBody[T any](r *http.Request, dst *T) error {
	ct := r.Header.Get(constants.ContentType)
	if ct != constants.EMPTY {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != constants.ApplicationJson {
			msg := fmt.Sprintf("Content-Type header is not %s", constants.ApplicationJson)
			return &Err{Status: http.StatusUnsupportedMediaType, Msg: msg}
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

func parseBodyError(err error) *Err {
	var syntaxError *json.SyntaxError
	var externalError *xrfErr.External
	var internalError *xrfErr.Internal
	var maxBytesError *http.MaxBytesError
	var unmarshalTypeError *json.UnmarshalTypeError
	var invalidUnmarshalError *json.InvalidUnmarshalError

	switch {
	// Syntax errors in the JSON
	case errors.As(err, &syntaxError):
		msg := fmt.Sprintf("Request contains badly-formed JSON (at position %d)", syntaxError.Offset)
		return &Err{Status: http.StatusBadRequest, Msg: msg}

	// In some circumstances Decode() may return an
	// io.ErrUnexpectedEOF error for syntax errors in the JSON. https://github.com/golang/go/issues/25956.
	case errors.Is(err, io.ErrUnexpectedEOF):
		msg := fmt.Sprintf("Request contains badly-formed JSON")
		return &Err{Status: http.StatusBadRequest, Msg: msg}

	// Catching error types like trying to assign a string in the
	// JSON request body to an int field.
	// interpolate the relevant field name and position into the error message
	case errors.Is(err, unmarshalTypeError):
		msg := fmt.Sprintf("Request contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
		return &Err{Status: http.StatusBadRequest, Msg: msg}

	// Catch the error caused by extra unexpected fields in the request
	// body. https://github.com/golang/go/issues/29035 regarding turning this into a sentinel error.
	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
		return &Err{Status: http.StatusBadRequest, Msg: msg}

	// An io.EOF error is returned by Decode() if the request body is empty.
	case errors.Is(err, io.EOF):
		msg := "Request body must not be empty"
		return &Err{Status: http.StatusBadRequest, Msg: msg}

	// Catch any error caused by the request body being too large.
	case errors.Is(err, maxBytesError):
		msg := "Request body must not be larger than 1MB"
		return &Err{Status: http.StatusBadRequest, Msg: msg}

	case errors.As(err, &invalidUnmarshalError):
		msg := "Request body must contain a valid JSON pointer"
		return &Err{Status: http.StatusBadRequest, Msg: msg}

	case errors.As(err, &externalError):
		return &Err{Status: http.StatusBadRequest, Msg: err.Error()}

	case errors.As(err, &internalError):
		var internalErr *xrfErr.Internal
		errors.As(err, &internalErr)
		return &Err{Status: http.StatusInternalServerError, Msg: internalErr.String()}

	default:
		return &Err{
			Status: http.StatusBadRequest,
			Msg:    fmt.Sprintf("Internal :: Err=%s", err.Error()),
		}
	}
}
