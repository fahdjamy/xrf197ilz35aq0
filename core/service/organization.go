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
	config  xrf.Security
	log     internal.Logger
	ctx     context.Context
	orgRepo repository.OrganizationRepository
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

func validateOrgName(name string) error {
	externalErr := &xrfErr.External{Source: "service/organization#validateOrgName"}
	if name == "" || len(name) < 3 || len(name) > 255 {
		externalErr.Message = "Org name must be between 3 and 255 characters"
		return externalErr
	}
	return nil
}

func NewOrganizationService(config xrf.Security, logger internal.Logger, orgRepo repository.OrganizationRepository) OrgService {
	return &organizationService{
		config:  config,
		log:     logger,
		orgRepo: orgRepo,
	}
}
