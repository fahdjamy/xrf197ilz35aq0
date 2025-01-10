package repository

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	"xrf197ilz35aq0/core/model/org"
	"xrf197ilz35aq0/internal"
	xrfErr "xrf197ilz35aq0/internal/error"
)

const RoleCollection = "role"

type RoleRepository interface {
	UpdateRole(role *org.Role, ctx context.Context) error
	SaveRole(role *org.Role, ctx context.Context) (string, error)
	FindRoleById(id string, ctx context.Context) (*org.Role, error)
	FindRoleByName(name string, ctx context.Context) (*org.Role, error)
	FindRolesByNames(names []string, ctx context.Context) ([]*org.Role, error)
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

func (repo *roleRepo) FindRolesByNames(names []string, ctx context.Context) ([]*org.Role, error) {
	if names == nil || len(names) == 0 {
		return []*org.Role{}, nil
	}
	internalErr = &xrfErr.Internal{}
	// 1. Build query filter
	filter := bson.M{"name": bson.M{"$in": names}}

	// 2. Query mongoDB
	cursor, err := repo.db.Collection(RoleCollection).Find(ctx, filter)
	if err != nil {
		internalErr.Message = "Failed to query roles"
		internalErr.Err = err
		return nil, internalErr
	}

	defer cursor.Close(ctx)

	// 3. Decode the results into a slice of Role structs
	var orgRoles []*org.Role

	if err := cursor.All(ctx, &orgRoles); err != nil {
		internalErr.Err = err
		internalErr.Message = "Failed to decode role objects"
		return nil, internalErr
	}

	return orgRoles, nil
}

func NewRoleRepo(db *mongo.Database, log internal.Logger) (RoleRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := createUniqueIndex(db, log, ctx, "name", RoleCollection); err != nil {
		return nil, err
	}

	return &roleRepo{
		db:  db,
		log: log,
	}, nil
}
