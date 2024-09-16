package user

import (
	"fmt"
	"net/mail"
	"time"
	"xrf197ilz35aq0"
	"xrf197ilz35aq0/core"
	"xrf197ilz35aq0/core/exchange"
	"xrf197ilz35aq0/core/model/user"
	"xrf197ilz35aq0/internal/encryption"
)

// Manager is a port (Driven side)
type Manager interface {
	NewUser(request *exchange.UserRequest) (*exchange.UserResponse, error)
}

type service struct {
	log             xrf197ilz35aq0.Logger
	settingsService SettingsManager
}

func (uc *service) NewUser(request *exchange.UserRequest) (*exchange.UserResponse, error) {
	now := time.Now()
	userName := request.FirstName + " " + request.LastName
	uc.log.Info(fmt.Sprintf("event=creatUser :: action=creatingUser :: user=%s", userName))

	err := uc.validateUser(request)
	if err != nil {
		return nil, err
	}

	newUser := user.NewUser(request.FirstName, request.LastName, request.Email.Data(), request.Password.Data())
	settings, err := uc.settingsService.NewSettings(request.Settings, *newUser)
	if err != nil {
		return nil, err
	}

	passCode, err := encryption.Encrypt([]byte(request.Password.Data()), []byte(settings.EncryptionKey.Data()))
	if err != nil {
		return nil, core.InternalError{
			Message: "Encrypt user password failed",
			Time:    now,
			Source:  "createUser",
			Err:     err,
		}
	}
	newUser.UpdatePassword(string(passCode))

	userResponse := toUserResponse(newUser, request)
	userResponse.Settings = *settings

	return userResponse, nil
}

func (uc *service) validateUser(request *exchange.UserRequest) error {
	// is validEmail
	_, err := mail.ParseAddress(request.Email.Data())
	if err != nil {
		return core.InvalidRequest{Message: "Invalid email address"}
	}
	lastNameLen := len(request.LastName)
	if lastNameLen != 0 && lastNameLen < 3 {
		return core.InvalidRequest{Message: "If last name is specified, it should be at least 3 characters long"}
	}
	firstNameLen := len(request.FirstName)
	if firstNameLen != 0 && firstNameLen < 3 {
		return core.InvalidRequest{Message: "If first name is specified, it should be at least 3 characters long"}
	}
	return nil
}

func toUserResponse(newUser *user.User, request *exchange.UserRequest) *exchange.UserResponse {
	return &exchange.UserResponse{
		UserId:    newUser.Id,
		CreatedAt: newUser.Joined,
		Email:     request.Email,
		LastName:  newUser.LastName,
		FirstName: newUser.FirstName,
		UpdatedAt: newUser.UpdatedAt,
		Anonymous: newUser.IsAnonymous(),
	}
}

func NewUserManager(logger xrf197ilz35aq0.Logger, settingsService SettingsManager) Manager {
	return &service{
		log:             logger,
		settingsService: settingsService,
	}
}
