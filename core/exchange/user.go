package exchange

import (
	"time"
	"xrf197ilz35aq0/internal/custom"
)

type UserRequest struct {
	FirstName           string                `json:"first_name"`
	LastName            string                `json:"last_name"`
	Email               custom.Secret[string] `json:"email"`
	Password            custom.Secret[string] `json:"password"`
	RotateEncryptionKey bool                  `json:"rotate_encryption_key"`
	EncryptAfter        time.Duration         `json:"encrypt_after"`
	Anonymous           bool                  `json:"anonymous"`
}

type UserResponse struct {
	UserId              int64                 `json:"userId"`
	FirstName           string                `json:"first_name"`
	LastName            string                `json:"last_name"`
	Email               custom.Secret[string] `json:"email"`
	RotateEncryptionKey bool                  `json:"rotate_encryption_key"`
	EncryptAfter        time.Duration         `json:"encrypt_after"`
	Anonymous           bool                  `json:"anonymous"`
	CreatedAt           time.Time             `json:"created_at"`
	UpdatedAt           time.Time             `json:"updated_at"`
}
