package exchange

import "time"

type UserRequest struct {
	FirstName           string        `json:"first_name"`
	LastName            string        `json:"last_name"`
	Email               string        `json:"email"`
	Password            string        `json:"password"`
	RotateEncryptionKey bool          `json:"rotate_encryption_key"`
	EncryptAfter        time.Duration `json:"encrypt_after"`
	Anonymous           bool          `json:"anonymous"`
}

type UserResponse struct {
	FirstName           string        `json:"first_name"`
	LastName            string        `json:"last_name"`
	Email               string        `json:"email"`
	RotateEncryptionKey bool          `json:"rotate_encryption_key"`
	EncryptAfter        time.Duration `json:"encrypt_after"`
	Anonymous           bool          `json:"anonymous"`
	CreatedAt           time.Time     `json:"created_at"`
	UpdatedAt           time.Time     `json:"updated_at"`
}
