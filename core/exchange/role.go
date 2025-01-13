package exchange

import "xrf197ilz35aq0/core/model"

type PermissionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PermissionResponse struct {
	Name        string     `json:"name"`
	UpdatedAt   model.Time `json:"updatedAt"`
	Description string     `json:"description"`
}
