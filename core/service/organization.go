package service

import (
	"context"
	"fmt"
	"strings"
	"sync"
	xrf "xrf197ilz35aq0"
	"xrf197ilz35aq0/core/exchange"
	"xrf197ilz35aq0/core/model/org"
	"xrf197ilz35aq0/core/model/user"
	"xrf197ilz35aq0/core/repository"
	"xrf197ilz35aq0/internal"
	xrfErr "xrf197ilz35aq0/internal/error"
)

type OrgService interface {
	FindOrgMembers(orgId string, ctx context.Context) ([]exchange.OrgMemberResponse, error)
	CreateOrg(request exchange.OrgRequest, ctx context.Context) (string, error)
	GetOrgById(orgId string, ctx context.Context) (*exchange.OrgResponse, error)
}

type organizationService struct {
	config   xrf.Security
	log      internal.Logger
	roleRepo repository.RoleRepository
	userRepo repository.UserRepository
	orgRepo  repository.OrganizationRepository
}

func (os *organizationService) CreateOrg(request exchange.OrgRequest, ctx context.Context) (string, error) {
	request.Name = strings.TrimSpace(request.Name)
	err := validateOrgName(request.Name)
	if err != nil {
		return "", err
	}

	orgMembers, err := os.validateAndCreateMembers(request.Members, ctx)
	if err != nil {
		return "", err
	}

	newOrg, err := org.CreateOrganization(request.Name, request.Category, request.Description, request.IsAnonymous, orgMembers)
	if err != nil {
		os.log.Error(fmt.Sprintf("event=creatOrg:: name=%s :: err=%v", request.Name, err))
		return "", err
	}

	orgId, err := os.orgRepo.Create(newOrg, ctx)
	if err != nil {
		return "", err
	}
	return orgId, nil
}

func (os *organizationService) GetOrgById(orgId string, ctx context.Context) (*exchange.OrgResponse, error) {
	savedOrg, err := os.findOrg(orgId, ctx)
	if err != nil {
		os.log.Error(fmt.Sprintf("event=getOrgIdFailure :: orgId=%s :: err=%v", orgId, err))
		return nil, err
	}
	return toOrgResponse(savedOrg), nil
}

func (os *organizationService) FindOrgMembers(orgId string, ctx context.Context) ([]exchange.OrgMemberResponse, error) {
	savedOrg, err := os.orgRepo.GetOrgById(orgId, ctx)
	if err != nil {
		os.log.Error(fmt.Sprintf("event=findOrgMembers action=findOrgFailed :: orgId=%s :: err=%v", orgId, err))
		return nil, err
	}
	orgMembers := savedOrg.Members
	uniqueRoleIds := make(map[string]string)

	userRoleMap := make(map[string][]string)
	userFps := make([]string, len(orgMembers))

	for _, orgMember := range orgMembers {
		userFps = append(userFps, orgMember.Fingerprint)
		// add everyUsers unique roleId
		for _, roleId := range orgMember.RoleIds {
			uniqueRoleIds[roleId] = ""
		}
		userRoleMap[orgMember.Fingerprint] = orgMember.RoleIds
	}

	allRoles := make([]string, len(orgMembers))
	for _, roleId := range uniqueRoleIds {
		allRoles = append(allRoles, roleId)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	var usersErr error
	var foundUsers []user.User

	// Call DB to find all users info asynchronously
	//dbCtx, dbCancel := context.WithTimeout(ctx, 40*time.Second)
	//defer dbCancel() // defer the Cancelling the dbCtx context after goroutines are done.
	go func() {
		defer wg.Done()
		foundUsers, usersErr = os.userRepo.FindUsersByFingerPrints(userFps, ctx)
		if err != nil {
			os.log.Error(fmt.Sprintf("event=findOrgMembers :: action=findUsersByFingerPrints :: err=%v", err))
		}
	}()

	var rolesErr error
	var foundRoles []org.Role
	// Call DB to get detailed info about each role asynchronously
	go func() {
		defer wg.Done()
		foundRoles, rolesErr = os.roleRepo.FindRolesByIds(allRoles, ctx)
		if err != nil {
			os.log.Error(fmt.Sprintf("event=findOrgMembers :: action=findRoles :: err=%v", err))
		}
	}()

	// Wait for both goroutines to finish.
	wg.Wait()

	// Check for errors from either goroutine
	if usersErr != nil {
		return nil, usersErr
	} else if rolesErr != nil {
		return nil, rolesErr
	}

	for _, foundRole := range foundRoles {
		uniqueRoleIds[foundRole.RoleId] = foundRole.Name
	}

	response := make([]exchange.OrgMemberResponse, 0)

	for _, foundUser := range foundUsers {
		userRoles := make([]string, 0)
		for _, userRole := range userRoleMap[foundUser.FingerPrint] {
			userRoles = append(userRoles, uniqueRoleIds[userRole])
		}
		response = append(response, exchange.OrgMemberResponse{
			Roles:  userRoles,
			UserId: foundUser.Id,
			Email:  foundUser.Email,
		})
	}

	return response, nil
}

func (os *organizationService) findOrg(orgId string, ctx context.Context) (*org.Organization, error) {
	if orgId == "" {
		return nil, &xrfErr.External{Source: "service/organizationService#findOrg", Message: "Invalid org id"}
	}
	savedOrg, err := os.orgRepo.GetOrgById(orgId, ctx)
	if err != nil {
		return nil, err
	}
	return savedOrg, nil
}

func (os *organizationService) validateAndCreateMembers(req []exchange.OrgMemberRequest, ctx context.Context) ([]org.Member, error) {
	externalErr := &xrfErr.External{Source: "service/organization#validateAndCreateMembers"}
	if req == nil || len(req) == 0 {
		externalErr.Message = "an org should have at least one member"
		return nil, externalErr
	}
	hasOwner := false
	userEmails := make([]string, 0)
	for _, member := range req {
		userEmails = append(userEmails, member.Email)
		if member.Owner {
			hasOwner = true
		}
	}

	if !hasOwner {
		externalErr.Message = "an org should have at least one owner"
		return nil, externalErr
	}

	foundUsers, err := os.userRepo.FindUsersByEmails(userEmails, ctx)
	if err != nil {
		os.log.Error(fmt.Sprintf("event=validateAndCreateMembers :: err=%v", err))
		return nil, err
	}
	// convert users to user map, {userEmail : userObject}
	dbUserMap := make(map[string]user.User)
	for _, savedUser := range foundUsers {
		dbUserMap[savedUser.Email] = savedUser
	}

	userMap := make(map[string]struct {
		isOwner bool
		userFp  string
		roleIds []string
	})

	missingUsers := make([]string, 0)

	for _, member := range req {
		roleMap, err := os.validateRoles(member.Roles, ctx)
		if err != nil {
			return nil, err
		}

		userObj, ok := dbUserMap[member.Email]
		if !ok {
			missingUsers = append(missingUsers, member.Email)
		} else {
			userMap[userObj.FingerPrint] = struct {
				isOwner bool
				userFp  string
				roleIds []string
			}{
				isOwner: member.Owner,
				userFp:  userObj.FingerPrint,
				roleIds: getUserRoles(roleMap, member.Roles),
			}
		}
	}

	if len(missingUsers) > 0 {
		externalErr.Message = fmt.Sprintf("invalid user emails '%v'", missingUsers)
		return nil, externalErr
	}

	orgMembers := make([]org.Member, 0)
	for _, value := range userMap {
		orgMembers = append(orgMembers, *org.CreateMember(value.userFp, value.isOwner, value.roleIds))
	}
	return orgMembers, nil
}

func (os *organizationService) validateRoles(roles []string, ctx context.Context) (map[string]string, error) {
	for _, role := range roles {
		if err := validateRoleName(role); err != nil {
			return nil, err
		}
	}

	savedRoles, err := os.roleRepo.FindRolesByNames(roles, ctx)
	if err != nil {
		os.log.Error(fmt.Sprintf("event=validateRoles:: name=%s :: err=%v", roles, err))
		return nil, err
	}
	rolesLen := len(roles)
	savedRolesLen := len(savedRoles)
	roleMap := make(map[string]string)
	// map roles to their ids
	for _, role := range savedRoles {
		roleMap[role.Name] = role.RoleId
	}

	if rolesLen != savedRolesLen {
		missingRoles := make([]string, rolesLen-savedRolesLen)

		for _, role := range roles {
			_, ok := roleMap[role]
			if !ok {
				missingRoles = append(missingRoles, role)
			}
		}

		return nil, &xrfErr.External{
			Source:  fmt.Sprintf("you have invalid roles '%v'", missingRoles),
			Message: "service/organization#validateRoles",
		}
	}
	return roleMap, nil
}

func getUserRoles(roleMap map[string]string, roleIds []string) []string {
	result := make([]string, 0)
	for _, role := range roleIds {
		result = append(result, roleMap[role])
	}
	return result
}

func validateOrgName(name string) error {
	externalErr := &xrfErr.External{Source: "service/organization#validateOrgName"}
	if name == "" || len(name) < 3 || len(name) > 255 {
		externalErr.Message = "Org name must be between 3 and 255 characters"
		return externalErr
	}
	return nil
}

func toOrgResponse(domainOrg *org.Organization) *exchange.OrgResponse {
	return &exchange.OrgResponse{
		OrgId:        domainOrg.Id,
		Category:     domainOrg.Category,
		CreatedAt:    domainOrg.CreatedAt,
		Description:  domainOrg.Description,
		Name:         domainOrg.DisplayName,
		MembersCount: len(domainOrg.Members),
		IsAnonymous:  domainOrg.IsAnonymous,
	}
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
