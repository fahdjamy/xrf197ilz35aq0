package service

import (
	"context"
	xrf "xrf197ilz35aq0"
	"xrf197ilz35aq0/core/exchange"
	"xrf197ilz35aq0/core/repository"
	"xrf197ilz35aq0/internal"
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
	//TODO implement me
	panic("implement me")
}

func (os *organizationService) GetOrgById(orgId string, ctx context.Context) (exchange.OrgResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewOrganizationService(config xrf.Security, logger internal.Logger, orgRepo repository.OrganizationRepository) OrgService {
	return &organizationService{
		config:  config,
		log:     logger,
		orgRepo: orgRepo,
	}
}
