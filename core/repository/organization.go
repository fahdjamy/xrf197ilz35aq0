package repository

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"xrf197ilz35aq0/core/model/org"
	"xrf197ilz35aq0/internal"
	"xrf197ilz35aq0/internal/constants"
	xrfErr "xrf197ilz35aq0/internal/error"
)

const OrgCollection = "organization"

type OrganizationRepository interface {
	Create(organization *org.Organization, ctx context.Context) (string, error)
	GetOrgById(id string, ctx context.Context) (*org.Organization, error)
}

type orgRepo struct {
	db  *mongo.Database
	log internal.Logger
}

func (repo *orgRepo) Create(organization *org.Organization, ctx context.Context) (string, error) {
	internalErr = &xrfErr.Internal{}
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
	internalErr = &xrfErr.Internal{}
	externalError = &xrfErr.External{}
	internalErr.Source = "core/repository/organization#getOrgById"

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

	err := createIndexForDoc(db, log, ctx, OrgCollection, constants.OrgId)
	if err != nil {
		return nil, err
	}
	return &orgRepo{db: db, log: log}, nil
}

func createIndexForDoc(db *mongo.Database, log internal.Logger, ctx context.Context, docName, indexName string) error {
	internalErr = &xrfErr.Internal{}

	collection := db.Collection(docName)
	// 1. List existing indexes on the collection
	cursor, err := collection.Indexes().List(ctx)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to list indexes for collection '%s': %v", RoleCollection, err))
		// Handle the error appropriately (e.g., exit)
	}
	defer cursor.Close(ctx)

	indexExists := false
	for cursor.Next(ctx) {
		var index bson.M
		if err := cursor.Decode(&index); err != nil {
			log.Error(fmt.Sprintf("Failed to decode index data: %v", err))
			// Handle the error appropriately
			internalErr.Err = err
			internalErr.Message = "Failed to decode index data"
			return internalErr
		}

		// 2. Check if an index on 'name' exists
		keys, ok := index["key"].(bson.D)
		if !ok {
			continue // Skip if the index doesn't have a 'key' field
		}

		// Check if it's an index on 'name' in ascending order
		if _, found := keys.Map()[indexName]; !found {
			indexExists = true
			// 3. Check if the existing index is unique (if it exists)
			unique, ok := index["unique"].(bool)
			if ok && unique {
				log.Debug(fmt.Sprintf("Unique index on '%s' already exists in collection '%s'", indexName, docName))
			} else {
				log.Warn(fmt.Sprintf("Index on '%s' exists in collection '%s' but is not unique. Consider dropping and recreating it.", indexName, docName))
				// You might want to handle this case differently:
				// - Drop and recreate the index (be very careful in production)
				// - Exit the application with an error
			}
			break // No need to continue checking other indexes
		}
	}

	if !indexExists {
		indexModel := mongo.IndexModel{
			// create an index on the 'providedIndexName' field in ascending order (1)
			Keys: bson.D{{Key: indexName, Value: 1}},
			// sets the unique option to true, enforcing uniqueness.
			Options: options.Index().SetUnique(true),
		}

		// creates the index on the 'provided' collection
		_, err = db.Collection(docName).Indexes().CreateOne(ctx, indexModel)
		if err != nil {
			return &xrfErr.Internal{
				Err:     err,
				Message: "Failed to create index",
				Source:  "core/repository#createRoleDocIndex",
			}
		}
	}

	return nil
}
