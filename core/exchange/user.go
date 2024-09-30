package exchange

import (
	"encoding/json"
	"fmt"
	"time"
	"xrf197ilz35aq0/core/model"
	"xrf197ilz35aq0/internal/custom"
	xrfErr "xrf197ilz35aq0/internal/error"
)

type UserRequest struct {
	FirstName string                `json:"firstName"`
	LastName  string                `json:"lastName"`
	Email     custom.Secret[string] `json:"email"`
	Password  custom.Secret[string] `json:"password"`
	Anonymous bool                  `json:"anonymous"`
	Settings  *SettingRequest       `json:"settings"`
}

func (u *UserRequest) UnmarshalJSON(bytes []byte) error {
	type Alias UserRequest
	aux := &struct {
		*Alias
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Alias: (*Alias)(u),
	}
	if err := json.Unmarshal(bytes, &aux); err != nil {
		return &xrfErr.External{
			Err:     err,
			Time:    time.Now(),
			Message: "Failed to unmarshal JSON",
			Source:  "core/exchange/UserRequest#UnmarshalJSON",
		}
	}

	u.Email = *custom.NewSecret(aux.Email)
	u.Password = *custom.NewSecret(aux.Password)
	return nil
}

func (u *UserRequest) String() string {
	return fmt.Sprintf("{firstName:%s, lastName%s, anonymous=%t}", u.FirstName, u.LastName, u.Anonymous)
}

type UserResponse struct {
	UserId    int64                 `json:"userId"`
	FirstName string                `json:"firstName"`
	LastName  string                `json:"lastName"`
	Email     custom.Secret[string] `json:"email"`
	Anonymous bool                  `json:"anonymous"`
	CreatedAt model.Time            `json:"createdAt"`
	UpdatedAt model.Time            `json:"updatedAt"`
	Settings  SettingResponse       `json:"settings"`
}

type SettingRequest struct {
	RotateKey     bool   `json:"rotateEncryptionKey"`
	RotateAfter   int    `json:"rotateAfter"`
	EncryptionKey string `json:"encryptionKey"`
}

type SettingResponse struct {
	CreatedAt     time.Time             `json:"createdAt"`
	UpdatedAt     time.Time             `json:"updatedAt"`
	EncryptionKey custom.Secret[string] `json:"encryptionKey"`
	RotateKey     bool                  `json:"rotateEncryptionKey"`
}
