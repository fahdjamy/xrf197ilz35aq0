package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"xrf197ilz35aq0/core/service"
	xrf "xrf197ilz35aq0/internal"
)

type AuthHandler struct {
	logger      xrf.Logger
	router      *mux.Router
	userService service.UserService
}

func (h *AuthHandler) getAuthToken(w http.ResponseWriter, r *http.Request) {
	resp := dataResponse{
		Code: 200,
		Data: struct {
			Token string `json:"token"`
		}{
			Token: "most-secure-token",
		},
	}
	writeResponse(resp, w, h.logger)
}

func (h *AuthHandler) RegisterAndListen() {
	h.router.HandleFunc("/api/v1/auth", h.getAuthToken).Methods(POST)
}

func NewAuthHandler(logger xrf.Logger, userService service.UserService, router *mux.Router) *AuthHandler {
	return &AuthHandler{
		logger:      logger,
		router:      router,
		userService: userService,
	}
}
