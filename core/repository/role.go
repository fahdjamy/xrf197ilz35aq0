package repository

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"xrf197ilz35aq0/core/model/org"
	"xrf197ilz35aq0/internal"
	xrfErr "xrf197ilz35aq0/internal/error"
)

const RoleCollection = "role"

type RoleRepo interface {
	UpdateRole(role *org.Role, ctx context.Context) error
	SaveRole(role *org.Role, ctx context.Context) (string, error)
	FindRoleById(id string, ctx context.Context) (*org.Role, error)
	FindRoleByName(name string, ctx context.Context) (*org.Role, error)
	FindRoleByNames(names []string, ctx context.Context) ([]*org.Role, int, error)
}

type roleRepo struct {
	db  *mongo.Database
	log internal.Logger
}

func (repo *roleRepo) SaveRole(role *org.Role, ctx context.Context) (string, error) {
	internalErr = &xrfErr.Internal{}
	externalError = &xrfErr.External{}
	document, err := repo.db.Collection(RoleCollection).InsertOne(ctx, role)
	if err != nil {
		// Check for the duplicate key error
		if mongo.IsDuplicateKeyError(err) {
			repo.log.Error(fmt.Sprintf("event=mongoDBFailure :: action=createRole :: err=duplicateName :: name=%s", role.Name))
			externalError.Message = "role name already exists"
			return "", externalError
		}
		repo.log.Error(fmt.Sprintf("event=mongoDBFailure :: action=createRole :: err=%s", err))
		internalErr.Message = "Creating new role in mongodb failed"
		internalErr.Err = err
		return "", err
	}
	repo.log.Debug(fmt.Sprintf("event=saveRole :: success=true :: objectID=%v", document.InsertedID))

	return document.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (repo *roleRepo) UpdateRole(role *org.Role, ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (repo *roleRepo) FindRoleById(id string, ctx context.Context) (*org.Role, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *roleRepo) FindRoleByName(name string, ctx context.Context) (*org.Role, error) {
	var result org.Role
	internalErr = &xrfErr.Internal{}

	filter := bson.M{"name": name}
	resp := repo.db.Collection(RoleCollection).FindOne(ctx, filter)

	if resp.Err() != nil {
		if errors.Is(resp.Err(), mongo.ErrNoDocuments) {
			externalError.Message = "Role not found"
			return nil, externalError
		}
		return nil, resp.Err()
	}

	if err := resp.Decode(&result); err != nil {
		internalErr.Err = err
		internalErr.Message = "Failed to decode role object"
		repo.log.Error(fmt.Sprintf("event=mongoDBFailure :: action=FindRoleByName :: err=%s", err))
		return nil, internalErr
	}
	return &result, nil
}

func (repo *roleRepo) FindRoleByNames(names []string, ctx context.Context) ([]*org.Role, int, error) {
	internalErr = &xrfErr.Internal{}
	// 1. Build query filter
	filter := bson.M{"name": bson.M{"$in": names}}

	// 2. Query mongoDB
	cursor, err := repo.db.Collection(RoleCollection).Find(ctx, filter)
	if err != nil {
		internalErr.Message = "Failed to query roles"
		internalErr.Err = err
		return nil, 0, internalErr
	}

	defer cursor.Close(ctx)

	// 3. Decode the results into a slice of Role structs
	var orgRoles []*org.Role

	if err := cursor.All(ctx, &orgRoles); err != nil {
		internalErr.Err = err
		internalErr.Message = "Failed to decode role objects"
		return nil, 0, internalErr
	}

	return orgRoles, len(orgRoles), nil
}

func NewRoleRepo(db *mongo.Database, log internal.Logger) (RoleRepo, error) {
	if err := createRoleDocIndex(db, log); err != nil {
		return nil, err
	}

	return &roleRepo{
		db:  db,
		log: log,
	}, nil
}

func createRoleDocIndex(db *mongo.Database, log internal.Logger) error {
	internalErr = &xrfErr.Internal{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(RoleCollection)
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
		if _, found := keys.Map()["name"]; !found {
			indexExists = true
			// 3. Check if the existing index is unique (if it exists)
			unique, ok := index["unique"].(bool)
			if ok && unique {
				log.Debug(fmt.Sprintf("Unique index on 'name' already exists in collection '%s'", RoleCollection))
			} else {
				log.Warn(fmt.Sprintf("Index on 'name' exists in collection '%s' but is not unique. Consider dropping and recreating it.", RoleCollection))
				// You might want to handle this case differently:
				// - Drop and recreate the index (be very careful in production)
				// - Exit the application with an error
			}
			break // No need to continue checking other indexes
		}
	}

	if !indexExists {
		indexModel := mongo.IndexModel{
			// create an index on the name field in ascending order (1)
			Keys: bson.D{{Key: "name", Value: 1}},
			// sets the unique option to true, enforcing uniqueness for role name.
			Options: options.Index().SetUnique(true),
		}

		// creates the index on the 'Role' collection
		_, err = db.Collection(RoleCollection).Indexes().CreateOne(ctx, indexModel)
		if err != nil {
			return &xrfErr.Internal{
				Err:     err,
				Message: "Failed to create index",
				Source:  "core/repository/role#createRoleDocIndex",
			}
		}
	}

	return nil
}
