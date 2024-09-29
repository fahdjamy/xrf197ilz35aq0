package exchange

import (
	"time"
	"xrf197ilz35aq0/core/model"
	"xrf197ilz35aq0/internal/custom"
)

type UserRequest struct {
	FirstName string                `json:"firstName"`
	LastName  string                `json:"lastName"`
	Email     custom.Secret[string] `json:"email"`
	Password  custom.Secret[string] `json:"password"`
	Anonymous bool                  `json:"anonymous"`
	Settings  *SettingRequest       `json:"settings"`
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
