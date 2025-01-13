package service

import (
	"context"
	"fmt"
	"strings"
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
	roleRepo repository.PermissionRepository
}

func (svc *roleService) CreateRole(req *exchange.RoleRequest, ctx context.Context) (string, error) {
	err := validateRoleName(req.Name)
	if err != nil {
		return "", err
	}

	newRole := org.CreateRole(req.Name, req.Name)
	_, err = svc.roleRepo.CreatePermission(newRole, ctx)
	if err != nil {
		svc.log.Error(fmt.Sprintf("event=CreateRole :: action=saveRoleToDB :: err=%v", err))
		return "", err
	}

	savedRole, err := svc.roleRepo.FindPermissionByName(strings.ToUpper(req.Name), ctx)
	if err != nil {
		internalErr := &xrfErr.Internal{Err: err, Message: "internal error", Source: "core/service/role#createRole"}
		svc.log.Error(fmt.Sprintf("event=CreateRole :: action=saveRoleToDB :: err=%v", err))
		return "", internalErr
	}

	return savedRole.RoleId, nil
}

func validateRoleName(name string) error {
	externalErr := &xrfErr.External{Source: "core/service/role#validateRoleName"}
	if name == "" || len(name) < 3 || len(name) > 63 {
		externalErr.Message = "role name should be between 3 and 63 characters"
		return externalErr
	}
	// Permission should be all letters or _
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

func NewRoleService(log internal.Logger, roleRepo repository.PermissionRepository) RoleService {
	return &roleService{
		log:      log,
		roleRepo: roleRepo,
	}
}
