package service

import (
	"xrf197ilz35aq0/core/repository"
	"xrf197ilz35aq0/internal"
	"xrf197ilz35aq0/internal/exchange"
)

type AuthService interface {
	Authenticate(request *exchange.AuthRequest) (string, error)
}

type authService struct {
	log      internal.Logger
	userRepo repository.UserRepository
}

func (service *authService) Authenticate(request *exchange.AuthRequest) (string, error) {

	return "", nil
}

func NewAuthService(log internal.Logger, userRepo repository.UserRepository) AuthService {
	return &authService{
		log:      log,
		userRepo: userRepo,
	}
}
