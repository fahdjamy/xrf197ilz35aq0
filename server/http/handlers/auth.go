package handlers

import (
	"context"
	"net/http"
	"time"
	"xrf197ilz35aq0/core/service"
	xrf "xrf197ilz35aq0/internal"
	xrfErr "xrf197ilz35aq0/internal/error"
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

func (auth *AuthHandler) verifyToken(w http.ResponseWriter, r *http.Request) {
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

func (auth *AuthHandler) verifyTokenAndGetEnrichedResponse(w http.ResponseWriter, r *http.Request) {
	appToAppToken := r.Header.Get("xrf-to-xrf-token")
	if appToAppToken == "" {
		auth.logger.Warn("event=verifyTokenAndGetEnrichedResponse id :: error=invalid XRF-TO-XRF-TOKEN")
		response.WriteErrorResponse(xrfErr.InvalidXrfToXrfTokenError, w, auth.logger)
		return
	}

	var request exchange.VerifyRevokeTokenReq
	err := decoder.DecodeJSONBody(r, &request)
	if err != nil {
		response.WriteErrorResponse(err, w, auth.logger)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	userResponse, err := auth.authService.VerifyTokenAndGetEnrichedResponse(ctx, request.Token)
	if err != nil {
		response.WriteErrorResponse(err, w, auth.logger)
		return
	}

	resp := response.DataResponse{
		Code: 200,
		Data: userResponse,
	}
	response.WriteResponse(resp, w, auth.logger)
}

func (auth *AuthHandler) RegisterRoutes(serveMux *http.ServeMux) {
	serveMux.Handle("POST /api/v1/auth/token", middleware.EnforceJSONMiddleware(auth.logger, http.HandlerFunc(auth.getAuthToken)))
	serveMux.Handle("POST /api/v1/auth/token/verify", middleware.EnforceJSONMiddleware(auth.logger, http.HandlerFunc(auth.verifyToken)))
	serveMux.Handle("POST /api/v1/auth/token/revoke", middleware.EnforceJSONMiddleware(auth.logger, http.HandlerFunc(auth.revokeToken)))
	serveMux.Handle("POST /api/v1/auth/token/verify-with-enriched", middleware.EnforceJSONMiddleware(auth.logger, http.HandlerFunc(auth.verifyTokenAndGetEnrichedResponse)))
}

func NewAuthHandler(logger xrf.Logger, services service.Services) *AuthHandler {
	return &AuthHandler{
		logger:      logger,
		userService: services.UserService,
		authService: services.AuthService,
	}
}
