package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"
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
	ctx, cancel := context.WithTimeout(req.Context(), 3*time.Second)
	defer cancel()

	userResp, err := user.userService.CreateUser(ctx, &userReq)
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

	ctx, cancel := context.WithTimeout(req.Context(), 3*time.Second)
	defer cancel()
	userResp, err := user.userService.GetUserById(ctx, userId)

	if err != nil {
		response.WriteErrorResponse(err, w, user.logger)
		return
	}

	resp := response.DataResponse{Data: userResp, Code: http.StatusOK}
	response.WriteResponse(resp, w, user.logger)
}

func (user *UserHandler) RegisterRoutes(serveMux *http.ServeMux) {
	serveMux.HandleFunc("POST /api/v1/user", user.createUser)
	serveMux.HandleFunc("GET /api/v1/user/{userId}", user.getUserById)
}
