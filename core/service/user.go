package service

import (
	"fmt"
	"xrf197ilz35aq0"
	"xrf197ilz35aq0/core/exchange"
	"xrf197ilz35aq0/core/model/user"
)

// UserOrchestrator is a port (Driven side)
type UserOrchestrator interface {
	CreateUser(request *exchange.UserRequest) (exchange.UserResponse, error)
}

type UserService struct {
	logger xrf197ilz35aq0.Logger
}

func (uc *UserService) CreateUser(request *exchange.UserRequest) (exchange.UserResponse, error) {
	uc.logger.Info(fmt.Sprintf("event=%s, name=%s", "creatUser", request.FirstName+" "+request.LastName))
	_ = user.NewUser(request.FirstName, request.LastName, request.Email, request.Password)

	err := uc.validateUser(request)
	if err != nil {
		return exchange.UserResponse{}, err
	}

	return exchange.UserResponse{}, nil
}

func (uc *UserService) validateUser(user *exchange.UserRequest) error {
	return nil
}

func NewUserService() UserOrchestrator {
	return &UserService{}
}
