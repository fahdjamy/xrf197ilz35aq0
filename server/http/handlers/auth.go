package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"xrf197ilz35aq0/core/service"
	xrf "xrf197ilz35aq0/internal"
	"xrf197ilz35aq0/internal/exchange"
	"xrf197ilz35aq0/server/http/decoder"
	"xrf197ilz35aq0/server/http/response"
)

type AuthHandler struct {
	logger         xrf.Logger
	router         *mux.Router
	userService    service.UserService
	authService    service.AuthService
	contextTimeout time.Duration
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

	ctx, cancel := context.WithTimeout(context.Background(), h.contextTimeout*time.Second)
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

func (h *AuthHandler) RegisterAndListen() {
	h.router.HandleFunc("/api/v1/auth", h.getAuthToken).Methods(POST)
	h.router.HandleFunc("/api/v1/auth/revoke", h.revokeToken).Methods(POST)
}

func NewAuthHandler(logger xrf.Logger, services service.Services, router *mux.Router) *AuthHandler {
	return &AuthHandler{
		logger:      logger,
		router:      router,
		userService: services.UserService,
		authService: services.AuthService,
	}
}
