package handlers

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"xrf197ilz35aq0/core/exchange"
	"xrf197ilz35aq0/core/service"
	xrf "xrf197ilz35aq0/internal"
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
		writeErrorResponse(errors.New("invalid id"), w, user.logger)
		return
	}

	data := map[string]int64{
		UserIdKey: userId,
	}

	resp := dataResponse{Data: data, Code: http.StatusOK}
	writeResponse(resp, w, user.logger)

}

func getAndValidateId(req *http.Request) (int64, bool) {
	vars := mux.Vars(req)
	userId, ok := vars[UserIdKey]
	if !ok {
		return 0, false
	}
	id, err := strconv.ParseInt(userId, 10, 64)

	if err != nil || id <= 0 {
		return 0, false
	}

	return id, true
}

func (user *User) RegisterAndListen() {
	user.router.HandleFunc("/user", user.createUser).Methods("POST")
	user.router.HandleFunc(fmt.Sprintf("/user/%s", UserIdKey), user.getUserById).Methods("GET")
}
