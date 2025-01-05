package repository

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"xrf197ilz35aq0/core/model/org"
	"xrf197ilz35aq0/internal"
	"xrf197ilz35aq0/internal/constants"
	xrfErr "xrf197ilz35aq0/internal/error"
)

const OrgCollection = "organization"

type OrganizationRepository interface {
	Create(organization *org.Organization, ctx context.Context) (any, error)
	GetOrgById(id string, ctx context.Context) (*org.Organization, error)
	FindByMongoId(id string, ctx context.Context) (*org.Organization, error)
}

type orgRepo struct {
	db  *mongo.Database
	log internal.Logger
}

func (repo *orgRepo) Create(organization *org.Organization, ctx context.Context) (any, error) {
	internalError = &xrfErr.Internal{}
	document, err := repo.db.Collection(OrgCollection).InsertOne(ctx, organization)
	if err != nil {
		repo.log.Error(fmt.Sprintf("event=mongoDBFailure :: action=saveOrg :: err=%s", err))
		internalError.Message = "Creating new org in mongodb failed"
		internalError.Err = err
		return nil, err
	}
	repo.log.Debug(fmt.Sprintf("event=saveOrg :: success=true :: objectID=%v", document.InsertedID))

	return document.InsertedID, nil
}

func (repo *orgRepo) GetOrgById(id string, ctx context.Context) (*org.Organization, error) {
	internalError = &xrfErr.Internal{}
	externalError = &xrfErr.External{}
	internalError.Source = "core/repository/organization#getOrgById"

	filter := bson.M{constants.OrgId: id}

	var result org.Organization
	resp := repo.db.Collection(OrgCollection).FindOne(ctx, filter)

	if resp.Err() != nil {
		if errors.Is(resp.Err(), mongo.ErrNoDocuments) {
			externalError.Message = "Org not found"
			return nil, externalError
		}
		return nil, resp.Err()
	}

	if err := resp.Decode(&result); err != nil {
		internalError.Err = err
		internalError.Message = "Failed to decode org object"
		repo.log.Error(fmt.Sprintf("event=mongoDBFailure :: action=getOrgById :: err=%s", err))
		return nil, internalError
	}
	return &result, nil
}

func (repo *orgRepo) FindByMongoId(id string, ctx context.Context) (*org.Organization, error) {
	internalError = &xrfErr.Internal{}

	filter := bson.M{"_id": id}
	var result org.Organization
	resp := repo.db.Collection(OrgCollection).FindOne(ctx, filter)

	if resp.Err() != nil {
		if errors.Is(resp.Err(), mongo.ErrNoDocuments) {
			externalError.Message = "Org not found"
			return nil, externalError
		}
		return nil, resp.Err()
	}

	if err := resp.Decode(&result); err != nil {
		internalError.Err = err
		internalError.Message = "Failed to decode org object"
		repo.log.Error(fmt.Sprintf("event=mongoDBFailure :: action=getOrgById :: err=%s", err))
		return nil, internalError
	}
	return &result, nil
}

func NewOrganizationRepository(db *mongo.Database, log internal.Logger) OrganizationRepository {
	return &orgRepo{db: db, log: log}
}
