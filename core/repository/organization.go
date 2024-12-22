package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"xrf197ilz35aq0/core/model/org"
	"xrf197ilz35aq0/internal"
	xrfErr "xrf197ilz35aq0/internal/error"
)

const OrgCollection = "organization"

type OrganizationRepository interface {
	Create(organization *org.Organization, ctx context.Context) (any, error)
}

type orgRepo struct {
	db  *mongo.Database
	log internal.Logger
}

func (or *orgRepo) Create(organization *org.Organization, ctx context.Context) (any, error) {
	internalError = &xrfErr.Internal{}
	document, err := or.db.Collection(OrgCollection).InsertOne(ctx, organization)
	if err != nil {
		or.log.Error(fmt.Sprintf("event=mongoDBFailure :: action=saveOrg :: err=%s", err))
		internalError.Message = "Creating new org in mongodb failed"
		internalError.Err = err
		return nil, err
	}
	or.log.Debug(fmt.Sprintf("event=saveOrg :: success=true :: objectID=%v", document.InsertedID))

	return document, nil
}

func NewOrganizationRepository(db *mongo.Database, log internal.Logger) OrganizationRepository {
	return &orgRepo{db: db, log: log}
}
