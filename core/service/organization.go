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
	config         xrf.Security
	log            internal.Logger
	userRepo       repository.UserRepository
	permissionRepo repository.PermissionRepository
	orgRepo        repository.OrganizationRepository
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
	uniquePermissionIds := make(map[string]string)

	userRoleMap := make(map[string][]string)
	userFps := make([]string, 0)

	for key, value := range savedOrg.Members { // key is user's fingerPrint
		userFps = append(userFps, key)
		// add everyUsers unique permissionId
		for _, permissionId := range value.Permissions {
			uniquePermissionIds[permissionId] = ""
		}
		userRoleMap[key] = value.Permissions
	}

	allPermission := make([]string, 0)
	for _, permId := range uniquePermissionIds {
		allPermission = append(allPermission, permId)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	var usersErr error
	var foundUsers []user.User

	// Call DB to find all users info asynchronously
	//dbCtx, dbCancel := context.WithTimeout(ctx, 5*time.Second)
	//defer dbCancel() // defer the Cancelling the dbCtx context after goroutines are done.
	go func() {
		defer wg.Done()
		foundUsers, usersErr = os.userRepo.FindUsersByFingerPrints(userFps, ctx)
		if err != nil {
			os.log.Error(fmt.Sprintf("event=findOrgMembers :: action=findUsersByFingerPrints :: err=%v", err))
		}
	}()

	var permissionErr error
	var foundRoles []org.Permission
	// Call DB to get detailed info about each permission asynchronously
	go func() {
		defer wg.Done()
		foundRoles, permissionErr = os.permissionRepo.FindPermissionsByIds(allPermission, ctx)
		if err != nil {
			os.log.Error(fmt.Sprintf("event=findOrgMembers :: action=findPermissions :: err=%v", err))
		}
	}()

	// Wait for both goroutines to finish.
	wg.Wait()

	// Check for errors from either goroutine
	if usersErr != nil {
		return nil, usersErr
	} else if permissionErr != nil {
		return nil, permissionErr
	}

	for _, foundRole := range foundRoles {
		uniquePermissionIds[foundRole.Id] = foundRole.Name
	}

	response := make([]exchange.OrgMemberResponse, 0)

	for _, foundUser := range foundUsers {
		userPermissions := make([]string, 0)
		for _, userRole := range userRoleMap[foundUser.FingerPrint] {
			userPermissions = append(userPermissions, uniquePermissionIds[userRole])
		}
		response = append(response, exchange.OrgMemberResponse{
			Permissions: userPermissions,
			UserId:      foundUser.Id,
			Email:       foundUser.Email,
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

func (os *organizationService) validateAndCreateMembers(req []exchange.OrgMemberRequest, ctx context.Context) (map[string]org.Member, error) {
	externalErr := &xrfErr.External{Source: "service/organization#validateAndCreateMembers"}
	if req == nil || len(req) == 0 {
		externalErr.Message = "an org should have at least one member"
		return nil, externalErr
	}
	hasOwner := false
	userEmails := make([]string, 0)
	seenMembers := make(map[string]bool) // avoid duplicates
	for _, member := range req {
		if !seenMembers[member.Email] {
			userEmails = append(userEmails, member.Email)
			if member.Owner {
				hasOwner = true
			}
			seenMembers[member.Email] = true
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
	dbUserMap := make(map[string]user.User)
	// convert users to user map, {userEmail : userObject}
	for _, savedUser := range foundUsers {
		dbUserMap[savedUser.Email] = savedUser
	}

	memberMap := make(map[string]struct {
		isOwner       bool
		userFp        string
		permissionIds []string
	})

	missingUsers := make([]string, 0)

	for _, member := range req {
		permissionMap, err := os.validatePermissions(member.Permissions, ctx)
		if err != nil {
			return nil, err
		}

		// gets the user from the request to the dbUserMap
		userObj, ok := dbUserMap[member.Email]
		if !ok {
			missingUsers = append(missingUsers, member.Email)
		} else {
			memberMap[userObj.FingerPrint] = struct {
				isOwner       bool
				userFp        string
				permissionIds []string
			}{
				isOwner:       member.Owner,
				userFp:        userObj.FingerPrint,
				permissionIds: getUserPermissions(permissionMap, member.Permissions),
			}
		}
	}

	if len(missingUsers) > 0 {
		externalErr.Message = fmt.Sprintf("you provided unknown user emails, please check your request")
		os.log.Error(fmt.Sprintf("event=validateAndCreateMembers :: action=userProvidedInvalidEmails :: invalidEmails=[%v]", missingUsers))
		return nil, externalErr
	}

	orgMembers := make(map[string]org.Member)
	for _, value := range memberMap {
		orgMembers[value.userFp] = *org.CreateMember(value.userFp, value.isOwner, value.permissionIds)
	}
	return orgMembers, nil
}

func (os *organizationService) validatePermissions(permissions []string, ctx context.Context) (map[string]string, error) {
	for _, permission := range permissions {
		if err := validatePermissionName(permission); err != nil {
			return nil, err
		}
	}

	savedRoles, err := os.permissionRepo.FindPermissionsByNames(permissions, ctx)
	if err != nil {
		os.log.Error(fmt.Sprintf("event=validatePermissions:: name=%s :: err=%v", permissions, err))
		return nil, err
	}
	permissionLen := len(permissions)
	savedRolesLen := len(savedRoles)
	permissionMap := make(map[string]string)
	// map permissions to their ids
	for _, permission := range savedRoles {
		permissionMap[permission.Name] = permission.Id
	}

	if permissionLen != savedRolesLen {
		missingRoles := make([]string, permissionLen-savedRolesLen)

		for _, permission := range permissions {
			_, ok := permissionMap[permission]
			if !ok {
				missingRoles = append(missingRoles, permission)
			}
		}

		return nil, &xrfErr.External{
			Source:  "service/organization#validatePermissions",
			Message: fmt.Sprintf("unknown permisions ['%v']", missingRoles),
		}
	}
	return permissionMap, nil
}

func getUserPermissions(permissionMap map[string]string, permissionIds []string) []string {
	result := make([]string, 0)
	for _, permission := range permissionIds {
		result = append(result, permissionMap[permission])
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
		config:         config,
		log:            logger,
		orgRepo:        allRepos.OrgRepo,
		userRepo:       allRepos.UserRepo,
		permissionRepo: allRepos.PermissionRepo,
	}
}
