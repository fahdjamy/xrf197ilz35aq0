package handlers

import (
	"context"
	"net/http"
	"time"
	"xrf197ilz35aq0/core/service"
	xrf "xrf197ilz35aq0/internal"
	"xrf197ilz35aq0/internal/exchange"
	"xrf197ilz35aq0/server/http/decoder"
	"xrf197ilz35aq0/server/http/response"
)

type AuthHandler struct {
	logger      xrf.Logger
	userService service.UserService
	authService service.AuthService
}

func (h *AuthHandler) getAuthToken(w http.ResponseWriter, r *http.Request) {
	var request exchange.AuthRequest
	err := decoder.DecodeJSONBody(r, &request)
	if err != nil {
		response.WriteErrorResponse(err, w, h.logger)
		return
	}

	tokenResp, err := h.authService.Authenticate(&request, context.Background())
	if err != nil {
		response.WriteErrorResponse(err, w, h.logger)
		return
	}
	resp := response.DataResponse{
		Code: 200,
		Data: struct {
			Token string `json:"token"`
		}{
			Token: tokenResp,
		},
	}
	response.WriteResponse(resp, w, h.logger)
}

func (h *AuthHandler) revokeToken(w http.ResponseWriter, r *http.Request) {
	var request exchange.RevokeTokenRequest
	err := decoder.DecodeJSONBody(r, &request)
	if err != nil {
		response.WriteErrorResponse(err, w, h.logger)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = h.authService.RevokeToken(request.Token, ctx)
	if err != nil {
		response.WriteErrorResponse(err, w, h.logger)
		return
	}
	resp := response.DataResponse{
		Code: 200,
	}
	response.WriteResponse(resp, w, h.logger)
}

func (h *AuthHandler) RegisterRoutes(serveMux *http.ServeMux) {
	serveMux.Handle("POST /api/v1/auth", http.Handler(http.HandlerFunc(h.getAuthToken)))
	serveMux.Handle("POST /api/v1/auth/revoke", http.Handler(http.HandlerFunc(h.revokeToken)))
}

func NewAuthHandler(logger xrf.Logger, services service.Services) *AuthHandler {
	return &AuthHandler{
		logger:      logger,
		userService: services.UserService,
		authService: services.AuthService,
	}
}
