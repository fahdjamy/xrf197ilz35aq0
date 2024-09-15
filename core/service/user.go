package service

import (
	"fmt"
	"net/mail"
	"xrf197ilz35aq0"
	"xrf197ilz35aq0/core"
	"xrf197ilz35aq0/core/exchange"
	"xrf197ilz35aq0/core/model/user"
	"xrf197ilz35aq0/internal/custom"
)

// UserManager is a port (Driven side)
type UserManager interface {
	CreateUser(request *exchange.UserRequest) (*exchange.UserResponse, error)
}

type UserService struct {
	logger xrf197ilz35aq0.Logger
}

func (uc *UserService) CreateUser(request *exchange.UserRequest) (*exchange.UserResponse, error) {
	userName := request.FirstName + " " + request.LastName
	uc.logger.Info(fmt.Sprintf("event=%s, name=%s", "creatUser", userName))

	newUser := user.NewUser(
		request.FirstName,
		request.LastName,
		request.Email.Data(),
		request.Password.Data())

	err := uc.validateUser(request)
	if err != nil {
		return nil, err
	}

	return toUserResponse(*newUser, request.Email.Data()), nil
}

func (uc *UserService) validateUser(user *exchange.UserRequest) error {
	// is validEmail
	_, err := mail.ParseAddress(user.Email.Data())
	if err != nil {
		return core.InvalidRequest{Message: "Invalid email address"}
	}
	lastNameLen := len(user.LastName)
	if lastNameLen != 0 && lastNameLen < 3 {
		return core.InvalidRequest{Message: "If last name is specified, it should be at least 3 characters long"}
	}
	firstNameLen := len(user.FirstName)
	if firstNameLen != 0 && firstNameLen < 3 {
		return core.InvalidRequest{Message: "If first name is specified, it should be at least 3 characters long"}
	}
	return nil
}

func toUserResponse(user user.User, email string) *exchange.UserResponse {
	secretEmail := custom.NewSecret[string](email)
	return &exchange.UserResponse{
		UserId:    user.Id,
		CreatedAt: user.Joined,
		Email:     *secretEmail,
		LastName:  user.LastName,
		FirstName: user.FirstName,
		UpdatedAt: user.UpdatedAt,
		Anonymous: user.IsAnonymous(),
	}
}

func NewUserService(logger xrf197ilz35aq0.Logger) UserManager {
	return &UserService{logger: logger}
}
