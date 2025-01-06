package service

import (
	"context"
	"unicode"
	"xrf197ilz35aq0/core/exchange"
	"xrf197ilz35aq0/core/model/org"
	"xrf197ilz35aq0/core/repository"
	"xrf197ilz35aq0/internal"
	xrfErr "xrf197ilz35aq0/internal/error"
)

type RoleService interface {
	CreateRole(req *exchange.RoleRequest, ctx context.Context) (string, error)
}

type roleService struct {
	log      internal.Logger
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
	// Role should be all letters or _
	for _, char := range name {
		if !unicode.IsLetter(char) && char != '_' {
			externalErr.Message = "role name should all letter characters"
			return externalErr
		}
	}
	countUnderScore := 0
	for _, char := range name {
		if char == '_' {
			countUnderScore++
		}
	}
	if countUnderScore > 2 {
		externalErr.Message = "role name shouldn't contain more than 3 underscore characters"
		return externalErr
	}
	lettersCount := 0
	for _, char := range name {
		if unicode.IsLetter(char) {
			lettersCount++
		}
	}
	if lettersCount < 3 {
		externalErr.Message = "role name should at least contain 3 letter characters"
		return externalErr
	}
	return nil
}

func NewRoleService(log internal.Logger, roleRepo repository.RoleRepo) RoleService {
	return &roleService{
		log:      log,
		roleRepo: roleRepo,
	}
}
