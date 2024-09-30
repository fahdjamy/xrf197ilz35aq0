package user

import (
	"fmt"
	"strconv"
	"time"
	"xrf197ilz35aq0"
	"xrf197ilz35aq0/core/exchange"
	"xrf197ilz35aq0/core/model/user"
	"xrf197ilz35aq0/internal"
	"xrf197ilz35aq0/internal/custom"
	"xrf197ilz35aq0/internal/encryption"
	xrfErr "xrf197ilz35aq0/internal/error"
	"xrf197ilz35aq0/internal/random"
)

type SettingsManager interface {
	NewSettings(request *exchange.SettingRequest, userModel user.User) (*exchange.SettingResponse, error)
}

type settingService struct {
	log xrf197ilz35aq0.Logger
}

func (s *settingService) NewSettings(request *exchange.SettingRequest, userModel user.User) (*exchange.SettingResponse, error) {
	s.log.Info(
		fmt.Sprintf("event=creatUserSettings :: action=creatingSettings :: userId=%s",
			strconv.FormatInt(userModel.Id, 10)))

	now := time.Now()
	rotateAfter := internal.AddMonths(now, request.RotateAfter)
	if len(request.EncryptionKey) == 0 {
		request.EncryptionKey = s.generateEncryptionKey()
	} else {
		err := s.validateEncryptionKey(request, now)
		if err != nil {
			return nil, err
		}
	}

	err := s.validateSettings(request, now)
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

func (s *settingService) validateEncryptionKey(request *exchange.SettingRequest, now time.Time) error {
	if len(request.EncryptionKey) != 0 && len(request.EncryptionKey) < 31 {
		return &xrfErr.External{Message: "Encryption key must be at least 31 characters long", Time: now}
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
	key, err := encryption.GenerateKey(35)
	errMsgTemplate := "event=%v :: error=%v"
	if err != nil {
		s.log.Debug(fmt.Sprintf(errMsgTemplate, "generateEncryptionKeyFailure", err))
		return strconv.FormatInt(random.PositiveInt64(), 10)
	}
	return string(key)
}

func (s *settingService) validateSettings(request *exchange.SettingRequest, now time.Time) error {
	if request.RotateKey {
		switch encryptAfter := request.RotateAfter; {
		case encryptAfter < 3:
			return &xrfErr.External{Message: "Rotation should at least be between 3 and 12 months", Time: now}
		case encryptAfter > 12:
			return &xrfErr.External{Message: "Key rotation should at least happen every year", Time: now}
		}
	}
	return nil
}

func NewSettingManager(logger xrf197ilz35aq0.Logger) SettingsManager {
	return &settingService{log: logger}
}
