package handlers

import (
	"fmt"
	"net/http"
	"xrf197ilz35aq0/core/service"
	xrf "xrf197ilz35aq0/internal"
	xrfErr "xrf197ilz35aq0/internal/error"
	"xrf197ilz35aq0/internal/exchange"
	"xrf197ilz35aq0/server/http/decoder"
	"xrf197ilz35aq0/server/http/response"
)

const (
	UserIdKey = "userId"
)

type UserHandler struct {
	logger      xrf.Logger
	userService service.UserService
}

func NewUserHandler(logger xrf.Logger, userManager service.UserService) *UserHandler {
	return &UserHandler{
		logger:      logger,
		userService: userManager,
	}
}

func (user *UserHandler) createUser(w http.ResponseWriter, req *http.Request) {
	var userReq exchange.UserRequest

	err := decoder.DecodeJSONBody(req, &userReq)
	if err != nil {
		response.WriteErrorResponse(err, w, user.logger)
		return
	}

	// create a user
	userResp, err := user.userService.CreateUser(&userReq)
	if err != nil {
		response.WriteErrorResponse(err, w, user.logger)
		return
	}

	resp := response.DataResponse{Data: userResp, Code: http.StatusCreated}
	response.WriteResponse(resp, w, user.logger)
}

func (user *UserHandler) getUserById(w http.ResponseWriter, req *http.Request) {
	userId, isValid := getAndValidateId(req, UserIdKey)
	if !isValid {
		externalError := &xrfErr.External{
			Message: "invalid user id",
		}
		response.WriteErrorResponse(externalError, w, user.logger)
		return
	}
	user.logger.Debug(fmt.Sprintf("event=getUserBy id :: userId=%s", userId))

	userResp, err := user.userService.GetUserById(userId)

	if err != nil {
		response.WriteErrorResponse(err, w, user.logger)
		return
	}

	resp := response.DataResponse{Data: userResp, Code: http.StatusOK}
	response.WriteResponse(resp, w, user.logger)
}

func (user *UserHandler) RegisterAndListen() {
	//user.router.HandleFunc("/api/v1/user", user.createUser).Methods(POST)
	//user.router.HandleFunc(fmt.Sprintf("/api/v1/user/{%s}", UserIdKey), user.getUserById).Methods(GET)
}
