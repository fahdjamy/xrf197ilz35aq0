package handlers

import (
	"context"
	"net/http"
	"time"
	"xrf197ilz35aq0/core/service"
	xrf "xrf197ilz35aq0/internal"
	"xrf197ilz35aq0/internal/exchange"
	"xrf197ilz35aq0/server/http/decoder"
	"xrf197ilz35aq0/server/http/middleware"
	"xrf197ilz35aq0/server/http/response"
)

type AuthHandler struct {
	logger      xrf.Logger
	userService service.UserService
	authService service.AuthService
}

func (auth *AuthHandler) getAuthToken(w http.ResponseWriter, r *http.Request) {
	var request exchange.AuthRequest
	err := decoder.DecodeJSONBody(r, &request)
	if err != nil {
		response.WriteErrorResponse(err, w, auth.logger)
		return
	}

	tokenResp, err := auth.authService.GetAuthToken(&request, context.Background())
	if err != nil {
		response.WriteErrorResponse(err, w, auth.logger)
		return
	}
	resp := response.DataResponse{
		Code: 200,
		Data: tokenResp,
	}
	response.WriteResponse(resp, w, auth.logger)
}

func (auth *AuthHandler) revokeToken(w http.ResponseWriter, r *http.Request) {
	var request exchange.RevokeTokenRequest
	err := decoder.DecodeJSONBody(r, &request)
	if err != nil {
		response.WriteErrorResponse(err, w, auth.logger)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = auth.authService.RevokeToken(request.Token, ctx)
	if err != nil {
		response.WriteErrorResponse(err, w, auth.logger)
		return
	}
	resp := response.DataResponse{
		Code: 200,
	}
	response.WriteResponse(resp, w, auth.logger)
}

func (auth *AuthHandler) RegisterRoutes(serveMux *http.ServeMux) {
	serveMux.Handle("POST /api/v1/auth", middleware.EnforceJSONMiddleware(auth.logger, http.HandlerFunc(auth.getAuthToken)))
	serveMux.Handle("POST /api/v1/auth/revoke", middleware.EnforceJSONMiddleware(auth.logger, http.HandlerFunc(auth.revokeToken)))
}

func NewAuthHandler(logger xrf.Logger, services service.Services) *AuthHandler {
	return &AuthHandler{
		logger:      logger,
		userService: services.UserService,
		authService: services.AuthService,
	}
}
