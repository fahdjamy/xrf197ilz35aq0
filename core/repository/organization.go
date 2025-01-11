package repository

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	"xrf197ilz35aq0/core/model/org"
	"xrf197ilz35aq0/internal"
	"xrf197ilz35aq0/internal/constants"
	xrfErr "xrf197ilz35aq0/internal/error"
)

const (
	OrgCollection = "organization"
)

type OrganizationRepository interface {
	Create(organization *org.Organization, ctx context.Context) (string, error)
	GetOrgById(id string, ctx context.Context) (*org.Organization, error)
}

type orgRepo struct {
	db  *mongo.Database
	log internal.Logger
}

func (repo *orgRepo) Create(organization *org.Organization, ctx context.Context) (string, error) {
	internalErr := &xrfErr.Internal{}
	document, err := repo.db.Collection(OrgCollection).InsertOne(ctx, organization)
	if err != nil {
		repo.log.Error(fmt.Sprintf("event=mongoDBFailure :: action=saveOrg :: err=%s", err))
		internalErr.Message = "Creating new org in mongodb failed"
		internalErr.Err = err
		return "", err
	}
	repo.log.Debug(fmt.Sprintf("event=saveOrg :: success=true :: objectID=%v", document.InsertedID))

	return organization.Id, nil
}

func (repo *orgRepo) GetOrgById(id string, ctx context.Context) (*org.Organization, error) {
	internalErr := &xrfErr.Internal{}
	externalError := &xrfErr.External{}
	internalErr.Source = "core/repository/organization#getOrgById"

	filter := bson.M{constants.OrgId: id}

	var result org.Organization
	resp := repo.db.Collection(OrgCollection).FindOne(ctx, filter)

	if resp.Err() != nil {
		if errors.Is(resp.Err(), mongo.ErrNoDocuments) {
			externalError.Message = constants.NotFoundOrgErrMsg
			externalError.Err = errors.New(constants.NotFoundOrgErrMsg)
			return nil, externalError
		}
		return nil, resp.Err()
	}

	if err := resp.Decode(&result); err != nil {
		internalErr.Err = err
		internalErr.Message = "Failed to decode org object"
		repo.log.Error(fmt.Sprintf("event=mongoDBFailure :: action=getOrgById :: err=%s", err))
		return nil, internalErr
	}
	return &result, nil
}

func NewOrganizationRepository(db *mongo.Database, log internal.Logger) (OrganizationRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := createUniqueIndex(db, log, ctx, OrgCollection, constants.OrgId)
	if err != nil {
		return nil, err
	}
	return &orgRepo{db: db, log: log}, nil
}
