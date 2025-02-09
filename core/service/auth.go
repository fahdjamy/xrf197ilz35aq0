package service

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"
	"xrf197ilz35aq0/core/repository"
	"xrf197ilz35aq0/core/security"
	"xrf197ilz35aq0/internal"
	xrfErr "xrf197ilz35aq0/internal/error"
	"xrf197ilz35aq0/internal/exchange"
	"xrf197ilz35aq0/storage"
)

type AuthService interface {
	VerifyToken(token string, ctx context.Context) (string, error)
	Authenticate(request *exchange.AuthRequest, ctx context.Context) (string, error)
}

type authService struct {
	secretKey    string
	serverId     string
	cache        storage.Cache
	log          internal.Logger
	userRepo     repository.UserRepository
	settingsRepo repository.SettingsRepository
}

type authTokenCache struct {
	UserId  string `json:"userId"`
	Revoked bool   `json:"revoked"`
}

func (at authTokenCache) MarshalBinary() ([]byte, error) {
	return json.Marshal(at)
}

func (service *authService) Authenticate(request *exchange.AuthRequest, ctx context.Context) (string, error) {
	email := request.Email
	password := request.Password
	internalErr := &xrfErr.Internal{Source: "service/auth#Authenticate"}
	externalErr := &xrfErr.External{Code: 400, Message: "invalid credentials"}

	savedUsers, err := service.userRepo.FindUsersByEmails([]string{email}, ctx)
	if err != nil {
		return "", err
	}
	if len(savedUsers) == 0 || savedUsers == nil {
		return "", externalErr
	}

	user := savedUsers[0]
	userSettings, err := service.settingsRepo.FetchUserSettings(ctx, user.FingerPrint)

	if err != nil {
		service.log.Error(fmt.Sprintf("event=authenticate :: userId=%s :: err=%s", user.Id, err))
		return "", err
	}
	isValid, err := verifyPassword(
		userSettings.Threads, userSettings.Memory, uint32(userSettings.Time), password, user.Password)
	if err != nil {
		service.log.Error(fmt.Sprintf("event=authenticate :: err=%s", err))
		return "", err
	}

	if !isValid {
		service.log.Debug(fmt.Sprintf("event=authenticate :: userId=%s :: validPassword=%t", user.Id, isValid))
		return "", externalErr
	}

	tokenExpiration := 6 * time.Hour
	tokenPayload := security.UserTokenPayload{
		UserId:    user.Id,
		ExpiresAt: tokenExpiration,
		ServerId:  service.serverId,
	}
	authToken, err := security.GenerateAuthToken(tokenPayload, []byte(service.secretKey))
	if err != nil {
		service.log.Error(fmt.Sprintf("event=authenticate :: action=GenerateTokenFailure :: err=%s", err))
		internalErr.Message = "failed to generate auth token"
		internalErr.Err = err
		return "", internalErr
	}

	authTokenCachePayload := authTokenCache{UserId: user.Id, Revoked: false}

	err = service.cache.Set(authToken, authTokenCachePayload, tokenExpiration, ctx)
	if err != nil {
		service.log.Error(fmt.Sprintf("event=authenticate :: action=cacheUserAuthToken :: err=%s", err))
		internalErr.Message = "failed to store auth token in cache"
		internalErr.Err = err
		return "", internalErr
	}

	return authToken, nil
}

func (service *authService) VerifyToken(token string, ctx context.Context) (string, error) {
	externalErr := &xrfErr.External{Code: 401}
	cachedTokenData, err := service.cache.Get(token, ctx)
	internalErr := &xrfErr.Internal{Source: "service/auth#Authenticate"}
	if err != nil {
		internalErr.Err = err
		internalErr.Message = "failed to fetch auth token"
		service.log.Error(fmt.Sprintf("event=verifyToken :: action=fetchCachedTokenData :: err=%s", err))
		return "", internalErr
	}

	// Use a bytes.Reader for efficient reading from the byte slice.
	cacheToByte := cachedTokenData.([]byte)
	if len(cacheToByte) == 0 {
		externalErr.Message = "invalid/expired token"
		return "", externalErr
	}
	cachedTokenPayload := bytes.NewReader(cacheToByte)

	data := &authTokenCache{}
	err = binary.Read(cachedTokenPayload, binary.BigEndian, data)
	if err != nil {
		service.log.Error(fmt.Sprintf("event=verifyToken :: action=binaryRead :: err=%s", err))
		internalErr.Message = "failed to read auth token binary data from cache"
		internalErr.Err = err
		return "", internalErr
	}

	if !data.Revoked {
		externalErr.Message = "token revoked"
		return "", externalErr
	}

	return data.UserId, nil
}

func NewAuthService(server string,
	log internal.Logger, authSecret string, cache storage.Cache, repos *repository.Repositories) AuthService {
	return &authService{
		log:          log,
		cache:        cache,
		serverId:     server,
		secretKey:    authSecret,
		userRepo:     repos.UserRepo,
		settingsRepo: repos.SettingsRepo,
	}
}
