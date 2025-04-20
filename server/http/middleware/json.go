package middleware

import (
	"mime"
	"net/http"
	"xrf197ilz35aq0/internal"
	"xrf197ilz35aq0/internal/constants"
	xrfErr "xrf197ilz35aq0/internal/error"
	"xrf197ilz35aq0/server/http/response"
)

func EnforceJSONMiddleware(logger internal.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get(constants.ContentType)
		mediaType, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			externalErr := &xrfErr.External{Message: err.Error(), Code: 400}
			response.WriteErrorResponse(externalErr, w, logger)
			return
		}

		if mediaType != "application/json" {
			externalErr := &xrfErr.External{Message: "expected Content-Type header must be application/json", Code: 400}
			response.WriteErrorResponse(externalErr, w, logger)
			return
		}

		next.ServeHTTP(w, r)
	})
}
