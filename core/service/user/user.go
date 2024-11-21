package user

import (
	"fmt"
	"net/mail"
	xrf "xrf197ilz35aq0"
	"xrf197ilz35aq0/core/exchange"
	"xrf197ilz35aq0/core/model/user"
	"xrf197ilz35aq0/internal/encryption"
	xrfErr "xrf197ilz35aq0/internal/error"
	"xrf197ilz35aq0/storage"
)

// Manager is a port (Driven side)
type Manager interface {
	NewUser(request *exchange.UserRequest) (*exchange.UserResponse, error)
}

type service struct {
	log             xrf.Logger
	settingsService SettingsManager
	store           storage.Store
}

func (uc *service) NewUser(request *exchange.UserRequest) (*exchange.UserResponse, error) {
	err := uc.validateUser(request)
	if err != nil {
		return nil, err
	}

	userName := request.FirstName + " " + request.LastName
	uc.log.Info(fmt.Sprintf("event=creatUser :: action=creatingUser :: user=%s", userName))

	newUser := user.NewUser(request.FirstName, request.LastName, request.Email.Data(), request.Password.Data())
	settings, err := uc.settingsService.NewSettings(request.Settings, *newUser)
	if err != nil {
		return nil, err
	}

	passCode, err := encryption.Encrypt([]byte(request.Password.Data()), []byte(settings.EncryptionKey.Data()))
	if err != nil {
		return nil, &xrfErr.Internal{
			Message: "Encrypt user password failed",
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
		return &xrfErr.External{Message: "Invalid email address"}
	}
	lastNameLen := len(request.LastName)
	if lastNameLen != 0 && lastNameLen < 3 {
		return &xrfErr.External{Message: "If last name is specified, it should be at least 3 characters long"}
	}
	firstNameLen := len(request.FirstName)
	if firstNameLen != 0 && firstNameLen < 3 {
		return &xrfErr.External{Message: "If first name is specified, it should be at least 3 characters long"}
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

func NewUserManager(logger xrf.Logger, settingsService SettingsManager, store storage.Store) Manager {

	return &service{
		store:           store,
		log:             logger,
		settingsService: settingsService,
	}
}
