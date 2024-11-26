package service

import (
	"context"
	"fmt"
	"strconv"
	"time"
	xrf "xrf197ilz35aq0"
	"xrf197ilz35aq0/core/exchange"
	"xrf197ilz35aq0/core/model/user"
	"xrf197ilz35aq0/internal"
	"xrf197ilz35aq0/internal/custom"
	"xrf197ilz35aq0/internal/encryption"
	xrfErr "xrf197ilz35aq0/internal/error"
	"xrf197ilz35aq0/internal/random"
	"xrf197ilz35aq0/storage"
)

type SettingsService interface {
	NewSettings(request *exchange.SettingRequest, userModel user.User) (*exchange.SettingResponse, error)
}

type settingService struct {
	log xrf.Logger
	ctx context.Context
	db  storage.Store
}

func (s *settingService) NewSettings(request *exchange.SettingRequest, userModel user.User) (*exchange.SettingResponse, error) {
	userId := strconv.FormatInt(userModel.Id, 10)
	s.log.Info(fmt.Sprintf("event=creatUserSettings :: action=creatingSettings :: userId=%s", userId))

	now := time.Now()
	rotateAfter := internal.AddMonths(now, request.RotateAfter)
	if len(request.EncryptionKey) == 0 {
		request.EncryptionKey = s.generateEncryptionKey()
	} else {
		err := s.validateEncryptionKey(request)
		if err != nil {
			return nil, err
		}
	}

	err := s.validateSettings(request)
	if err != nil {
		return nil, err
	}

	settings := user.NewSettings(
		request.RotateKey,
		time.Since(rotateAfter),
		userModel.FingerPrint(),
		request.EncryptionKey,
	)
	return toSettingsResponse(settings), nil
}

func (s *settingService) validateEncryptionKey(request *exchange.SettingRequest) error {
	key := request.EncryptionKey
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return &xrfErr.External{Message: "Encryption must be 16, 24, or 32 bytes"}
	}
	return nil
}

func toSettingsResponse(settings *user.Settings) *exchange.SettingResponse {
	key := custom.NewSecret(settings.Key())
	return &exchange.SettingResponse{
		EncryptionKey: *key,
		CreatedAt:     settings.CreatedAt,
		UpdatedAt:     settings.LastModified,
		RotateKey:     settings.RotateEncryptionKey,
	}
}

func (s *settingService) generateEncryptionKey() string {
	key, err := encryption.GenerateKey(32)
	if err != nil {
		s.log.Debug(fmt.Sprintf("event=%v :: error=%v", "generateEncryptionKeyFailure", err))
		return strconv.FormatInt(random.PositiveInt64(), 10)
	}
	return string(key)
}

func (s *settingService) validateSettings(request *exchange.SettingRequest) error {
	if request.RotateKey {
		switch encryptAfter := request.RotateAfter; {
		case encryptAfter < 3:
			return &xrfErr.External{Message: "Rotation should at least be between 3 and 12 months"}
		case encryptAfter > 12:
			return &xrfErr.External{Message: "Key rotation should at least happen every year"}
		}
	}
	return nil
}

func NewSettingService(logger xrf.Logger, store storage.Store, ctx context.Context) SettingsService {
	return &settingService{
		ctx: ctx,
		db:  store,
		log: logger,
	}
}
