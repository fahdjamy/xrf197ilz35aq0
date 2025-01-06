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
	xrfErr "xrf197ilz35aq0/internal/error"
)

const RoleCollection = "role"

type RoleRepo interface {
	UpdateRole(role *org.Role, ctx context.Context) error
	SaveRole(role *org.Role, ctx context.Context) (string, error)
	FindRoleById(id string, ctx context.Context) (*org.Role, error)
	FindRoleByName(name string, ctx context.Context) (*org.Role, error)
	FindRoleByMongoId(mongoId string, ctx context.Context) (*org.Role, error)
}

type roleRepo struct {
	db  *mongo.Database
	log internal.Logger
}

func (repo *roleRepo) SaveRole(role *org.Role, ctx context.Context) (string, error) {
	internalError = &xrfErr.Internal{}
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
		internalError.Message = "Creating new role in mongodb failed"
		internalError.Err = err
		return "", err
	}
	repo.log.Debug(fmt.Sprintf("event=saveRole :: success=true :: objectID=%v", document.InsertedID))

	return document.InsertedID.(string), nil
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
	internalError = &xrfErr.Internal{}

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
		internalError.Err = err
		internalError.Message = "Failed to decode role object"
		repo.log.Error(fmt.Sprintf("event=mongoDBFailure :: action=FindRoleByName :: err=%s", err))
		return nil, internalError
	}
	return &result, nil
}

func (repo *roleRepo) FindRoleByMongoId(mongoId string, ctx context.Context) (*org.Role, error) {
	var result org.Role
	internalError = &xrfErr.Internal{}

	filter := bson.M{"_id": mongoId}
	resp := repo.db.Collection(RoleCollection).FindOne(ctx, filter)

	if resp.Err() != nil {
		if errors.Is(resp.Err(), mongo.ErrNoDocuments) {
			externalError.Message = "Role not found"
			return nil, externalError
		}
		return nil, resp.Err()
	}

	if err := resp.Decode(&result); err != nil {
		internalError.Err = err
		internalError.Message = "Failed to decode role object"
		repo.log.Error(fmt.Sprintf("event=mongoDBFailure :: action=FindRoleByMongoId :: err=%s", err))
		return nil, internalError
	}
	return &result, nil
}

func NewRoleRepo(db *mongo.Database, log internal.Logger) (RoleRepo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexModel := mongo.IndexModel{
		// create an index on the name field in ascending order (1)
		Keys: bson.D{{Key: "name", Value: 1}},
		// sets the unique option to true, enforcing uniqueness for role name.
		Options: options.Index().SetUnique(true),
	}

	// creates the index on the 'Role' collection
	_, err := db.Collection(RoleCollection).Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return nil, &xrfErr.Internal{
			Err:     err,
			Message: "Failed to create index",
			Source:  "core/repository/role#NewRoleRepo",
		}
	}
	return &roleRepo{
		db:  db,
		log: log,
	}, nil
}
