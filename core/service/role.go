package service

import (
	"context"
	"unicode"
	"xrf197ilz35aq0/core/exchange"
	"xrf197ilz35aq0/core/model/org"
	"xrf197ilz35aq0/core/repository"
	xrfErr "xrf197ilz35aq0/internal/error"
)

type RoleService interface {
	CreateRole(req *exchange.RoleRequest, ctx context.Context) (string, error)
}

type roleService struct {
	roleRepo repository.RoleRepo
}

func (svc *roleService) CreateRole(req *exchange.RoleRequest, ctx context.Context) (string, error) {
	err := validateName(req.Name)
	if err != nil {
		return "", err
	}

	newRole := org.CreateRole(req.Name, req.Name)
	roleMongoId, err := svc.roleRepo.SaveRole(newRole, ctx)
	if err != nil {
		return "", err
	}

	savedRole, err := svc.roleRepo.FindRoleByMongoId(roleMongoId, ctx)
	if err != nil {
		internalErr := &xrfErr.Internal{Err: err, Message: "internal error", Source: "core/service/role#createRole"}
		return "", internalErr
	}

	return savedRole.RoleId, nil
}

func validateName(name string) error {
	externalErr := &xrfErr.External{Source: "core/service/role#validateName"}
	if name == "" || len(name) < 3 || len(name) > 63 {
		externalErr.Message = "role name should be between 3 and 63 characters"
		return externalErr
	}
	for _, char := range name {
		if !unicode.IsLetter(char) {
			externalErr.Message = "role name should all letter characters"
			return externalErr
		}
	}
	return nil
}

func NewRoleService() RoleService {
	return &roleService{}
}
