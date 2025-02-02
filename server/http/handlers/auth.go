package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"xrf197ilz35aq0/core/service"
	xrf "xrf197ilz35aq0/internal"
	"xrf197ilz35aq0/internal/exchange"
)

type AuthHandler struct {
	logger      xrf.Logger
	router      *mux.Router
	userService service.UserService
	authService service.AuthService
}

func (h *AuthHandler) getAuthToken(w http.ResponseWriter, r *http.Request) {
	var request exchange.AuthRequest
	err := decodeJSONBody(r, &request)
	if err != nil {
		writeErrorResponse(err, w, h.logger)
		return
	}

	tokenResp, err := h.authService.Authenticate(&request, context.Background())
	if err != nil {
		writeErrorResponse(err, w, h.logger)
		return
	}
	resp := dataResponse{
		Code: 200,
		Data: struct {
			Token string `json:"token"`
		}{
			Token: tokenResp,
		},
	}
	writeResponse(resp, w, h.logger)
}

func (h *AuthHandler) RegisterAndListen() {
	h.router.HandleFunc("/api/v1/auth", h.getAuthToken).Methods(POST)
}

func NewAuthHandler(logger xrf.Logger, services service.Services, router *mux.Router) *AuthHandler {
	return &AuthHandler{
		logger:      logger,
		router:      router,
		userService: services.UserService,
		authService: services.AuthService,
	}
}
