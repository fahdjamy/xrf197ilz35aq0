package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	xrf "xrf197ilz35aq0"
	"xrf197ilz35aq0/core/exchange"
	"xrf197ilz35aq0/core/service"
)

type User struct {
	router      *mux.Router
	logger      xrf.Logger
	userManager service.UserService
}

func NewUser(logger xrf.Logger, userManager service.UserService, router *mux.Router) *User {
	return &User{
		router:      router,
		logger:      logger,
		userManager: userManager,
	}
}

func (user *User) createUser(w http.ResponseWriter, req *http.Request) {
	var userReq exchange.UserRequest

	err := decodeJSONBody(req, &userReq)
	if err != nil {
		writeErrorResponse(err, w, user.logger)
		return
	}

	// create a user
	userResp, err := user.userManager.CreateUser(&userReq)
	if err != nil {
		writeErrorResponse(err, w, user.logger)
		return
	}

	resp := dataResponse{Data: userResp, Code: http.StatusCreated}
	writeResponse(resp, w, user.logger)
}

func (user *User) RegisterAndListen() {
	user.router.HandleFunc("/user", user.createUser).Methods("POST")
}
