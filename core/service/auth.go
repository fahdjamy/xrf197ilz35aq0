package service

import (
	"context"
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
	RevokeToken(token string, ctx context.Context) error
	VerifyToken(token string, ctx context.Context) (string, error)
	GetAuthToken(request *exchange.AuthRequest, ctx context.Context) (*exchange.AuthResponse, error)
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
	Salt    string `json:"salt"`
	UserId  string `json:"userId"`
	Revoked bool   `json:"revoked"`
	UserFp  string `json:"userFp"`
}

func (at authTokenCache) MarshalBinary() ([]byte, error) {
	return json.Marshal(at)
}

func (service *authService) GetAuthToken(request *exchange.AuthRequest, ctx context.Context) (*exchange.AuthResponse, error) {
	email := request.Email
	password := request.Password
	internalErr := &xrfErr.Internal{Source: "service/auth#GetAuthToken"}
	externalErr := &xrfErr.External{Code: 400, Message: "invalid credentials"}

	savedUsers, err := service.userRepo.FindUsersByEmails([]string{email}, ctx)
	if err != nil {
		return nil, err
	}
	if len(savedUsers) == 0 || savedUsers == nil {
		return nil, externalErr
	}

	user := savedUsers[0]
	userSettings, err := service.settingsRepo.FetchUserSettings(ctx, user.FingerPrint)

	userId := user.Id
	if err != nil {
		service.log.Error(fmt.Sprintf("event=authenticate :: userId=%s :: err=%s", userId, err))
		return nil, err
	}
	isValid, err := verifyPassword(
		userSettings.Threads, userSettings.Memory, uint32(userSettings.Time), password, user.Password)
	if err != nil {
		service.log.Error(fmt.Sprintf("event=authenticate :: err=%s", err))
		return nil, err
	}

	if !isValid {
		service.log.Debug(fmt.Sprintf("event=authenticate :: userId=%s :: validPassword=%t", userId, isValid))
		return nil, externalErr
	}

	tokenExpiration := 6 * time.Hour
	tokenPayload := security.UserTokenPayload{UserId: userId, ExpiresAt: tokenExpiration, ServerId: service.serverId}
	authToken, err := security.GenerateAuthToken(tokenPayload, []byte(service.secretKey))
	if err != nil {
		service.log.Error(fmt.Sprintf("event=authenticate :: action=GenerateTokenFailure :: err=%s", err))
		internalErr.Message = "failed to generate auth token"
		internalErr.Err = err
		return nil, internalErr
	}

	authTokenCachePayload := authTokenCache{
		Revoked: false,
		UserId:  userId,
		Salt:    authToken.Salt,
		UserFp:  user.FingerPrint,
	}

	err = service.cache.Set(authToken.Token, authTokenCachePayload, tokenExpiration, ctx)
	if err != nil {
		service.log.Error(fmt.Sprintf("event=authenticate :: action=cacheUserAuthToken :: err=%s", err))
		internalErr.Message = "failed to store auth token in cache"
		internalErr.Err = err
		return nil, internalErr
	}

	return &exchange.AuthResponse{
		Token:  authToken.Token,
		Expiry: int64(tokenExpiration),
	}, nil
}

func (service *authService) VerifyToken(token string, ctx context.Context) (string, error) {
	externalErr := &xrfErr.External{Code: 401}
	cachedTokenData, err := service.cache.Get(token, ctx)
	internalErr := &xrfErr.Internal{Source: "service/auth#GetAuthToken"}
	if err != nil {
		if err.Error() == "redis: nil" {
			externalErr.Message = "invalid / expired token"
			return "", externalErr
		}
		internalErr.Err = err
		internalErr.Message = "failed to fetch auth token"
		service.log.Error(fmt.Sprintf("event=verifyToken :: action=fetchCachedTokenData :: err=%s", err.Error()))
		return "", internalErr
	}

	// Use a bytes.Reader for efficient reading from the byte slice.
	cacheToByte := cachedTokenData.(string)
	if len(cacheToByte) == 0 {
		externalErr.Message = "invalid/expired token"
		return "", externalErr
	}

	data := &authTokenCache{}
	err = json.Unmarshal([]byte(cacheToByte), data)
	if err != nil {
		service.log.Error(fmt.Sprintf("event=verifyToken :: action=binaryRead :: err=%s", err.Error()))
		internalErr.Message = "failed to read auth token binary data from cache"
		internalErr.Err = err
		return "", internalErr
	}

	if data.Revoked {
		externalErr.Message = "token revoked"
		return "", externalErr
	}

	return data.UserId, nil
}

func (service *authService) RevokeToken(token string, ctx context.Context) error {
	internalErr := &xrfErr.Internal{Source: "service/auth#RevokeToken"}
	deletedValCount, err := service.cache.Delete(token, ctx)
	if err != nil {
		internalErr.Message = "failed to delete auth token from cache"
		internalErr.Err = err
		return internalErr
	}
	service.log.Debug(fmt.Sprintf("event=revokeToken :: tokensDeleted=%d", deletedValCount))
	return nil
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
