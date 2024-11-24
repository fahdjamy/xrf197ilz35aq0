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

var externalClientErr *xrfErr.External

func (u *UserRequest) UnmarshalJSON(bytes []byte) error {
	externalClientErr = &xrfErr.External{}
	externalClientErr.Source = "core/exchange/UserRequest#UnmarshalJSON"
	type Alias UserRequest
	aux := &struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,password,min=18,max=55"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	if err := json.Unmarshal(bytes, &aux); err != nil {
		externalClientErr.Message = "Failed to unmarshal JSON"
		return externalClientErr
	}

	if aux.Email == "" {
		externalClientErr.Message = "Invalid or missing email address"
		return externalClientErr
	}
	if aux.Password == "" {
		externalClientErr.Message = "Invalid or missing password"
		return externalClientErr
	}

	u.Email = *custom.NewSecret(aux.Email)
	u.Password = *custom.NewSecret(aux.Password)
	return nil
}

func (u *UserRequest) MarshalJSON() ([]byte, error) {
	if u == nil {
		externalClientErr.Message = "UserRequest is nil"
		return nil, externalClientErr
	}
	userObj := *u
	if userObj.Email.Data() == "" {
		externalClientErr.Message = "invalid user email"
		return nil, externalClientErr
	}
	if userObj.Password.Data() == "" {
		externalClientErr.Message = "invalid user password"
		return nil, externalClientErr
	}

	type Alias UserRequest

	auxUser := (Alias)(userObj)
	return json.Marshal(auxUser)
}

func (u *UserRequest) String() string {
	return fmt.Sprintf("{firstName:%s, lastName%s, anonymous=%t}", u.FirstName, u.LastName, u.Anonymous)
}

type UserResponse struct {
	UserId    int64                 `json:"userId"`
	FirstName string                `json:"firstName,omitempty"`
	LastName  string                `json:"lastName,omitempty"`
	Email     custom.Secret[string] `json:"email"`
	Anonymous bool                  `json:"anonymous"`
	CreatedAt model.Time            `json:"createdAt"`
	UpdatedAt model.Time            `json:"updatedAt"`
	Settings  SettingResponse       `json:"settings,omitempty"`
}

func (u *UserResponse) String() string {
	return fmt.Sprintf("{UserId: %d, CreatedAt: %s, UpdatedAt: %s}", u.UserId, u.CreatedAt, u.UpdatedAt)
}

type SettingRequest struct {
	RotateKey     bool   `json:"rotateKey"`
	RotateAfter   int    `json:"rotateAfter"`
	EncryptionKey string `json:"encryptionKey"`
}

func (s *SettingRequest) String() string {
	return fmt.Sprintf("{rotateKey: %t, rotateAfter: %d}", s.RotateKey, s.RotateAfter)
}

type SettingResponse struct {
	CreatedAt     time.Time             `json:"createdAt"`
	UpdatedAt     time.Time             `json:"updatedAt"`
	EncryptionKey custom.Secret[string] `json:"encryptionKey"`
	RotateKey     bool                  `json:"rotateKey"`
}

func (s *SettingResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
		RotateKey bool      `json:"rotateKey"`
	}{
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
		RotateKey: s.RotateKey,
	})
}

func (s *SettingResponse) String() string {
	return fmt.Sprintf("{rotateKey: %t: createdAt: %s, updatedAt: %s}", s.RotateKey, s.CreatedAt, s.UpdatedAt)
}
