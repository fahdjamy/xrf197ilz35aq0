package service

import (
	"context"
	"fmt"
	"xrf197ilz35aq0/core/repository"
	"xrf197ilz35aq0/internal"
	xrfErr "xrf197ilz35aq0/internal/error"
	"xrf197ilz35aq0/internal/exchange"
)

type AuthService interface {
	Authenticate(request *exchange.AuthRequest, ctx context.Context) (string, error)
}

type authService struct {
	log          internal.Logger
	userRepo     repository.UserRepository
	settingsRepo repository.SettingsRepository
}

func (service *authService) Authenticate(request *exchange.AuthRequest, ctx context.Context) (string, error) {
	email := request.Email
	password := request.Password
	externalErr := &xrfErr.External{Code: 400, Message: "invalid credentials"}

	savedUsers, err := service.userRepo.FindUsersByEmails([]string{email}, ctx)
	if err != nil {
		return "", err
	}
	if len(savedUsers) == 0 {
		return "", externalErr
	}

	user := savedUsers[0]
	userSettings, err := service.settingsRepo.FetchUserSettings(ctx, user.FingerPrint)

	if err != nil {
		service.log.Error(fmt.Sprintf("event=authenticate :: userId=%s :: err=%s", user.Id, err))
		return "", err
	}
	isValid, err := verifyPassword(userSettings.Threads, userSettings.Memory, uint32(userSettings.Time), password, user.Password)
	if err != nil {
		service.log.Error(fmt.Sprintf("event=authenticate :: err=%s", err))
		return "", err
	}

	if !isValid {
		return "", externalErr
	}

	return "", nil
}

func NewAuthService(log internal.Logger, userRepo repository.UserRepository) AuthService {
	return &authService{
		log:      log,
		userRepo: userRepo,
	}
}
