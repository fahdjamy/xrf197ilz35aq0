package service

import (
	"context"
	"fmt"
	"net/mail"
	"xrf197ilz35aq0/core/exchange"
	"xrf197ilz35aq0/core/model/user"
	"xrf197ilz35aq0/core/repository"
	xrf "xrf197ilz35aq0/internal"
	"xrf197ilz35aq0/internal/encryption"
	xrfErr "xrf197ilz35aq0/internal/error"
)

var internalError *xrfErr.Internal

// UserService is a port (Driven side)
type UserService interface {
	CreateUser(request *exchange.UserRequest) (*exchange.UserResponse, error)
}

type service struct {
	log             xrf.Logger
	settingsService SettingsService
	ctx             context.Context
	userRepo        repository.UserRepository
}

func (uc *service) CreateUser(request *exchange.UserRequest) (*exchange.UserResponse, error) {
	internalError = &xrfErr.Internal{}
	internalError.Source = "core/service/user/user#createUser"

	userName := "Unknown name"
	if request.FirstName != "" && request.LastName != "" {
		userName = request.FirstName + " " + request.LastName
	}

	uc.log.Info(fmt.Sprintf("event=creatUser :: action=creatingUser :: username=%s", userName))
	err := uc.validateUser(request)
	if err != nil {
		return nil, err
	}

	newUser := user.NewUser(request.FirstName, request.LastName, request.Email.Data(), "")

	// SAVE-USER/DB: ACTION 1 - save user and settings to database
	uc.log.Debug(fmt.Sprintf("event=creatUser :: action=saveUserINDB :: userFP=%s :: userId=%d", newUser.FingerPrint[:7], newUser.Id))
	_, err = uc.userRepo.CreateUser(newUser, uc.ctx)
	if err != nil {
		internalError.Err = err
		internalError.Message = "User creation failed"
		return nil, internalError
	}

	settingRequest := request.Settings
	if settingRequest == nil {
		settingRequest = &exchange.SettingRequest{
			RotateKey:   false,
			RotateAfter: 13,
		}
	}

	// SAVE-USER/DB: ACTION 2 - create user settings
	settings, err := uc.settingsService.NewSettings(settingRequest, newUser.FingerPrint)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := encryption.EncryptAndEncode([]byte(request.Password.Data()), []byte(settings.EncryptionKey.Data()))
	if err != nil {
		internalError.Err = err
		internalError.Message = "Encrypting password failed"
		return nil, internalError
	}

	// SAVE-USER/DB: ACTION 3 - Update user password
	uc.log.Debug(fmt.Sprintf("event=creatUser :: action=setUserPassword :: userFP=%s :: userId=%d", newUser.FingerPrint[:7], newUser.Id))
	passwordSet, err := uc.userRepo.UpdatePassword(newUser.FingerPrint, hashedPassword, uc.ctx)
	if err != nil {
		internalError.Err = err
		internalError.Message = "Update password failed"
		return nil, internalError
	}
	uc.log.Debug(fmt.Sprintf("event=creatUser :: action=setUserPassword :: userFP=%s :: userId=%d passwordSet=%t", newUser.FingerPrint[:7], newUser.Id, passwordSet))

	// Return userResponse
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

func NewUserService(log xrf.Logger, userSettings SettingsService, userRepo repository.UserRepository, ctx context.Context) UserService {

	return &service{
		ctx:             ctx,
		log:             log,
		userRepo:        userRepo,
		settingsService: userSettings,
	}
}
