package exchange

import (
	"time"
)

type PermissionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PermissionResponse struct {
	Name        string    `json:"name"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Description string    `json:"description"`
}
