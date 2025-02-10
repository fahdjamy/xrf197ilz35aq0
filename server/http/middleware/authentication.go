package middleware

import (
	"context"
	"fmt"
	"net/http"
	"xrf197ilz35aq0/core/service"
	"xrf197ilz35aq0/internal"
	xrfErr "xrf197ilz35aq0/internal/error"
	"xrf197ilz35aq0/server/http/response"
)

const XrfAuthToken = "xrf-auth-token"

type AuthenticationMiddleware struct {
	logger      internal.Logger
	authService service.AuthService
}

func (m *AuthenticationMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get(XrfAuthToken)
		if authToken == "" {
			externalErr := &xrfErr.External{Message: "missing xrf-auth-token authToken", Code: 401}
			response.WriteErrorResponse(externalErr, w, m.logger)
			return
		}

		ctx := context.WithValue(r.Context(), XrfAuthToken, authToken)
		userId, err := m.authService.VerifyToken(authToken, ctx)
		if err != nil {
			response.WriteErrorResponse(err, w, m.logger)
			return
		}
		m.logger.Debug(fmt.Sprintf("event=authenticatingToken :: userId=%s", userId))

		next.ServeHTTP(w, r)
	})
}

func NewAuthenticationMiddleware(logger internal.Logger, authService service.AuthService) AuthenticationMiddleware {
	return AuthenticationMiddleware{
		logger:      logger,
		authService: authService,
	}
}
