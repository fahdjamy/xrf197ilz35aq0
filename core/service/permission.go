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

type PermissionService interface {
	CreatePermission(req *exchange.PermissionRequest, ctx context.Context) (string, error)
}

type permissionService struct {
	log            internal.Logger
	permissionRepo repository.PermissionRepository
}

func (svc *permissionService) CreatePermission(req *exchange.PermissionRequest, ctx context.Context) (string, error) {
	err := validatePermissionName(req.Name)
	if err != nil {
		return "", err
	}

	newRole := org.CreateRole(req.Name, req.Name)
	_, err = svc.permissionRepo.CreatePermission(newRole, ctx)
	if err != nil {
		svc.log.Error(fmt.Sprintf("event=CreatePermission :: action=savePermissionToDB :: err=%v", err))
		return "", err
	}

	savedRole, err := svc.permissionRepo.FindPermissionByName(strings.ToUpper(req.Name), ctx)
	if err != nil {
		internalErr := &xrfErr.Internal{Err: err, Message: "internal error", Source: "core/service/permission#createPermission"}
		svc.log.Error(fmt.Sprintf("event=CreatePermission :: action=savingPermissionToDB :: err=%v", err))
		return "", internalErr
	}

	return savedRole.Id, nil
}

func validatePermissionName(name string) error {
	externalErr := &xrfErr.External{Source: "core/service/permission#validatePermissionName"}
	if name == "" || len(name) < 3 || len(name) > 63 {
		externalErr.Message = "permission name should be between 3 and 63 characters"
		return externalErr
	}
	// Permission should be all letters or _
	for _, char := range name {
		if !unicode.IsLetter(char) && char != '_' {
			externalErr.Message = "permission name should all letter characters"
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
		externalErr.Message = "permission name shouldn't contain more than 3 underscore characters"
		return externalErr
	}
	lettersCount := 0
	for _, char := range name {
		if unicode.IsLetter(char) {
			lettersCount++
		}
	}
	if lettersCount < 3 {
		externalErr.Message = "permission name should at least contain 3 letter characters"
		return externalErr
	}
	return nil
}

func NewRoleService(log internal.Logger, permissionRepo repository.PermissionRepository) PermissionService {
	return &permissionService{
		log:            log,
		permissionRepo: permissionRepo,
	}
}
