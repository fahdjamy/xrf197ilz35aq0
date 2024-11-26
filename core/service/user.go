package service

import (
	"context"
	"fmt"
	"net/mail"
	"xrf197ilz35aq0/core/exchange"
	"xrf197ilz35aq0/core/model/user"
	xrf "xrf197ilz35aq0/internal"
	"xrf197ilz35aq0/internal/encryption"
	xrfErr "xrf197ilz35aq0/internal/error"
	"xrf197ilz35aq0/storage"
)

const Collection = "user"

var internalError *xrfErr.Internal

// UserService is a port (Driven side)
type UserService interface {
	CreateUser(request *exchange.UserRequest) (*exchange.UserResponse, error)
}

type service struct {
	log             xrf.Logger
	settingsService SettingsService
	store           storage.Store
	ctx             context.Context
}

func (uc *service) CreateUser(request *exchange.UserRequest) (*exchange.UserResponse, error) {
	internalError = &xrfErr.Internal{}
	internalError.Source = "core/service/user/user#createUser"
	err := uc.validateUser(request)
	if err != nil {
		return nil, err
	}

	userName := "Unknown name"
	if request.FirstName != "" && request.LastName != "" {
		userName = request.FirstName + " " + request.LastName
	}
	uc.log.Info(fmt.Sprintf("event=creatUser :: action=creatingUser :: username=%s", userName))

	newUser := user.NewUser(request.FirstName, request.LastName, request.Email.Data(), request.Password.Data())

	// save user to DB
	createdUser, err := uc.store.Save(Collection, newUser)
	if err != nil {
		internalError.Err = err
		internalError.Message = "Saving new user failed"
		return nil, internalError
	}

	uc.log.Debug(fmt.Sprintf("message=saved user successfully :: userId=%s", createdUser))

	settingRequest := request.Settings
	if settingRequest == nil {
		settingRequest = &exchange.SettingRequest{
			RotateKey:   false,
			RotateAfter: 13,
		}
	}
	settings, err := uc.settingsService.NewSettings(settingRequest, *newUser)
	if err != nil {
		return nil, err
	}

	passCode, err := encryption.Encrypt([]byte(request.Password.Data()), []byte(settings.EncryptionKey.Data()))
	if err != nil {
		internalError.Err = err
		internalError.Message = "Encrypting password failed"
		return nil, internalError
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

func NewUserService(logger xrf.Logger, settingsService SettingsService, store storage.Store, ctx context.Context) UserService {
	return &service{
		ctx:             ctx,
		store:           store,
		log:             logger,
		settingsService: settingsService,
	}
}
