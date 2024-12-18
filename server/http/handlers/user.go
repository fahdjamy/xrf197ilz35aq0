package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"xrf197ilz35aq0/core/exchange"
	"xrf197ilz35aq0/core/service"
	xrf "xrf197ilz35aq0/internal"
	xrfErr "xrf197ilz35aq0/internal/error"
)

const (
	UserIdKey = "userId"
)

type User struct {
	logger      xrf.Logger
	router      *mux.Router
	userService service.UserService
}

func NewUser(logger xrf.Logger, userManager service.UserService, router *mux.Router) *User {
	return &User{
		router:      router,
		logger:      logger,
		userService: userManager,
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
	userResp, err := user.userService.CreateUser(&userReq)
	if err != nil {
		writeErrorResponse(err, w, user.logger)
		return
	}

	resp := dataResponse{Data: userResp, Code: http.StatusCreated}
	writeResponse(resp, w, user.logger)
}

func (user *User) getUserById(w http.ResponseWriter, req *http.Request) {
	userId, isValid := getAndValidateId(req)
	if !isValid {
		externalError := &xrfErr.External{
			Message: "invalid user id",
		}
		writeErrorResponse(externalError, w, user.logger)
		return
	}
	user.logger.Debug(fmt.Sprintf("event=getUserBy id :: userId=%s", userId))

	userResp, err := user.userService.GetUserById(userId)

	if err != nil {
		writeErrorResponse(err, w, user.logger)
		return
	}

	resp := dataResponse{Data: userResp, Code: http.StatusOK}
	writeResponse(resp, w, user.logger)

}

func getAndValidateId(req *http.Request) (string, bool) {
	vars := mux.Vars(req)
	userId, ok := vars[UserIdKey]
	if !ok || userId == "" {
		return "", false
	}

	return userId, true
}

func (user *User) RegisterAndListen() {
	user.router.HandleFunc("/user", user.createUser).Methods("POST")
	user.router.HandleFunc(fmt.Sprintf("/user/{%s}", UserIdKey), user.getUserById).Methods("GET")
}
