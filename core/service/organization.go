package service

import (
	"context"
	"fmt"
	xrf "xrf197ilz35aq0"
	"xrf197ilz35aq0/core/exchange"
	"xrf197ilz35aq0/core/model/org"
	"xrf197ilz35aq0/core/repository"
	"xrf197ilz35aq0/internal"
	xrfErr "xrf197ilz35aq0/internal/error"
)

type OrgService interface {
	CreateOrg(request exchange.OrgRequest, ctx context.Context) (string, error)
	GetOrgById(orgId string, ctx context.Context) (exchange.OrgResponse, error)
}

type organizationService struct {
	config   xrf.Security
	log      internal.Logger
	roleRepo repository.RoleRepository
	userRepo repository.UserRepository
	orgRepo  repository.OrganizationRepository
}

func (os *organizationService) CreateOrg(request exchange.OrgRequest, ctx context.Context) (string, error) {
	err := validateOrgName(request.Name)
	if err != nil {
		return "", err
	}

	newOrg, err := org.CreateOrganization(request.Name, request.Category, request.Description, nil)
	if err != nil {
		os.log.Error(fmt.Sprintf("event=creatOrg:: name=%s :: err=%v", request.Name, err))
		return "", err
	}

	insertedOrgId, err := os.orgRepo.Create(newOrg, ctx)
	if err != nil {
		return "", err
	}

	savedOrg, err := os.orgRepo.FindByMongoId(insertedOrgId.(string), ctx)
	if err != nil {
		os.log.Error(fmt.Sprintf("event=createOrg:: name=%s :: err=%v", request.Name, err))
		return "", err
	}
	return savedOrg.Id, nil
}

func (os *organizationService) GetOrgById(orgId string, ctx context.Context) (exchange.OrgResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (os *organizationService) validateRoles(roles []string, ctx context.Context) error {
	for _, role := range roles {
		if err := validateRoleName(role); err != nil {
			return err
		}
	}

	savedRoles, err := os.roleRepo.FindRolesByNames(roles, ctx)
	if err != nil {
		os.log.Error(fmt.Sprintf("event=validateRoles:: name=%s :: err=%v", roles, err))
		return err
	}
	rolesLen := len(roles)
	savedRolesLen := len(savedRoles)
	if rolesLen != savedRolesLen {
		missingRoles := make([]string, rolesLen-savedRolesLen)
		roleMap := make(map[string]bool)
		for _, role := range savedRoles {
			roleMap[role.Name] = true
		}

		for _, role := range roles {
			exists := roleMap[role]
			if !exists {
				missingRoles = append(missingRoles, role)
			}
		}

		return &xrfErr.External{
			Source:  fmt.Sprintf("you have invalid roles '%v'", missingRoles),
			Message: "service/organization#validateRoles",
		}
	}
	return nil
}

func (os *organizationService) validateOrgMembers(members *[]exchange.OrgMemberRequest, ctx context.Context) error {
	externalErr := &xrfErr.External{Source: "service/organization#validateOrgMembers"}
	if members == nil || len(*members) == 0 {
		externalErr.Message = "an org should have at least one member"
		return externalErr
	}
	hasOwner := false
	for _, member := range *members {
		if member.Owner {
			hasOwner = true
		}
	}

	if !hasOwner {
		externalErr.Message = "an org should have at least one owner"
		return externalErr
	}

	for _, member := range *members {
		if err := os.validateRoles(member.Roles, ctx); err != nil {
			return err
		}
	}
	return nil
}

func validateOrgName(name string) error {
	externalErr := &xrfErr.External{Source: "service/organization#validateOrgName"}
	if name == "" || len(name) < 3 || len(name) > 255 {
		externalErr.Message = "Org name must be between 3 and 255 characters"
		return externalErr
	}
	return nil
}

func NewOrganizationService(config xrf.Security, logger internal.Logger, allRepos *repository.Repositories) OrgService {
	return &organizationService{
		config:   config,
		log:      logger,
		orgRepo:  allRepos.OrgRepo,
		userRepo: allRepos.UserRepo,
		roleRepo: allRepos.RoleRepo,
	}
}
