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

	tokenResp, err := auth.authService.GetAuthToken(r.Context(), &request)
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
	var request exchange.VerifyRevokeTokenReq
	err := decoder.DecodeJSONBody(r, &request)
	if err != nil {
		response.WriteErrorResponse(err, w, auth.logger)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	err = auth.authService.RevokeToken(ctx, request.Token)
	if err != nil {
		response.WriteErrorResponse(err, w, auth.logger)
		return
	}
	resp := response.DataResponse{
		Code: 200,
	}
	response.WriteResponse(resp, w, auth.logger)
}

func (auth *AuthHandler) verifAuthToken(w http.ResponseWriter, r *http.Request) {
	var request exchange.VerifyRevokeTokenReq
	err := decoder.DecodeJSONBody(r, &request)
	if err != nil {
		response.WriteErrorResponse(err, w, auth.logger)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	userId, err := auth.authService.VerifyToken(ctx, request.Token)
	if err != nil {
		response.WriteErrorResponse(err, w, auth.logger)
		return
	}

	resp := response.DataResponse{
		Code: 200,
		Data: struct {
			UserId string `json:"userId"`
		}{
			UserId: userId,
		},
	}
	response.WriteResponse(resp, w, auth.logger)
}

func (auth *AuthHandler) RegisterRoutes(serveMux *http.ServeMux) {
	serveMux.Handle("POST /api/v1/auth/token", middleware.EnforceJSONMiddleware(auth.logger, http.HandlerFunc(auth.getAuthToken)))
	serveMux.Handle("POST /api/v1/auth/token/revoke", middleware.EnforceJSONMiddleware(auth.logger, http.HandlerFunc(auth.revokeToken)))
	serveMux.Handle("POST /api/v1/auth/token/verify", middleware.EnforceJSONMiddleware(auth.logger, http.HandlerFunc(auth.verifAuthToken)))
}

func NewAuthHandler(logger xrf.Logger, services service.Services) *AuthHandler {
	return &AuthHandler{
		logger:      logger,
		userService: services.UserService,
		authService: services.AuthService,
	}
}
