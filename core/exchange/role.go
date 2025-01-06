package exchange

import "xrf197ilz35aq0/core/model"

type RoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RoleResponse struct {
	Name        string     `json:"name"`
	UpdatedAt   model.Time `json:"updatedAt"`
	Description string     `json:"description"`
}
