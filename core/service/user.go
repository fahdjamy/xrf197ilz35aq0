package service

import (
	"fmt"
	"time"
	"xrf197ilz35aq0"
	"xrf197ilz35aq0/core/model/user"
)

type UserRequest struct {
	FirstName           string        `json:"first_name"`
	LastName            string        `json:"last_name"`
	Email               string        `json:"email"`
	Password            string        `json:"password"`
	RotateEncryptionKey bool          `json:"rotate_encryption_key"`
	EncryptAfter        time.Duration `json:"encrypt_after"`
	Anonymous           bool          `json:"anonymous"`
}

type UserOrchestrator interface {
	CreateUser(request UserRequest) (*user.User, error)
}

type UserService struct {
	logger xrf197ilz35aq0.Logger
}

func (uc *UserService) CreateUser(request UserRequest) (*user.User, error) {
	uc.logger.Info(fmt.Sprintf("event=%s, name=%s", "creatUser", request.FirstName+" "+request.LastName))
	newUser := user.NewUser(request.FirstName, request.LastName, request.Email, request.Password)

	err := uc.validateUser(newUser)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (uc *UserService) validateUser(user *user.User) error {
	return nil
}

func NewUserService() UserOrchestrator {
	return &UserService{}
}
