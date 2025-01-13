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
	"xrf197ilz35aq0/internal/constants"
	xrfErr "xrf197ilz35aq0/internal/error"
)

const PermissionsCollection = "permission"

type PermissionRepository interface {
	UpdatePermission(permission *org.Permission, ctx context.Context) error
	CreatePermission(permission *org.Permission, ctx context.Context) (string, error)
	FindPermissionById(id string, ctx context.Context) (*org.Permission, error)
	FindPermissionByName(name string, ctx context.Context) (*org.Permission, error)
	FindPermissionsByIds(ids []string, ctx context.Context) ([]org.Permission, error)
	FindPermissionsByNames(names []string, ctx context.Context) ([]org.Permission, error)
}

type permissionsRepo struct {
	db  *mongo.Database
	log internal.Logger
}

func (repo *permissionsRepo) CreatePermission(permission *org.Permission, ctx context.Context) (string, error) {
	internalErr := &xrfErr.Internal{}
	externalError := &xrfErr.External{}
	document, err := repo.db.Collection(PermissionsCollection).InsertOne(ctx, permission)
	if err != nil {
		// Check for the duplicate key error
		if mongo.IsDuplicateKeyError(err) {
			repo.log.Error(fmt.Sprintf("event=mongoDBFailure :: action=createRole :: err=duplicateName :: name=%s", permission.Name))
			externalError.Message = "permission name already exists"
			return "", externalError
		}
		repo.log.Error(fmt.Sprintf("event=mongoDBFailure :: action=createRole :: err=%s", err))
		internalErr.Message = "Creating new permission in mongodb failed"
		internalErr.Err = err
		return "", err
	}
	repo.log.Debug(fmt.Sprintf("event=saveRole :: success=true :: objectID=%v", document.InsertedID))

	return document.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (repo *permissionsRepo) UpdatePermission(permission *org.Permission, ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (repo *permissionsRepo) FindPermissionById(id string, ctx context.Context) (*org.Permission, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *permissionsRepo) FindPermissionByName(name string, ctx context.Context) (*org.Permission, error) {
	var result org.Permission
	internalErr := &xrfErr.Internal{}
	externalError := &xrfErr.External{}

	filter := bson.M{"name": name}
	resp := repo.db.Collection(PermissionsCollection).FindOne(ctx, filter)

	if resp.Err() != nil {
		if errors.Is(resp.Err(), mongo.ErrNoDocuments) {
			externalError.Message = "Permission not found"
			return nil, externalError
		}
		return nil, resp.Err()
	}

	if err := resp.Decode(&result); err != nil {
		internalErr.Err = err
		internalErr.Message = "Failed to decode permission object"
		repo.log.Error(fmt.Sprintf("event=mongoDBFailure :: action=FindPermissionByName :: err=%s", err))
		return nil, internalErr
	}
	return &result, nil
}

func (repo *permissionsRepo) FindPermissionsByNames(names []string, ctx context.Context) ([]org.Permission, error) {
	return repo.findRolesByFilter(names, constants.NAME, ctx)
}

func (repo *permissionsRepo) FindPermissionsByIds(ids []string, ctx context.Context) ([]org.Permission, error) {
	return repo.findRolesByFilter(ids, constants.PermissionId, ctx)
}

func (repo *permissionsRepo) findRolesByFilter(values []string, filterBy string, ctx context.Context) ([]org.Permission, error) {
	if values == nil || len(values) == 0 {
		return []org.Permission{}, nil
	}
	internalError := &xrfErr.Internal{}
	// 1. Build query filter
	filter := bson.M{filterBy: bson.M{"$in": values}}

	// 2. Query mongoDB
	cursor, err := repo.db.Collection(PermissionsCollection).Find(ctx, filter)
	if err != nil {
		internalError.Message = fmt.Sprintf("Failed to query permissions by filter: %s", filterBy)
		internalError.Err = err
		return nil, internalError
	}

	defer cursor.Close(ctx)

	// 3. Decode the results into a slice of Permission structs
	var orgRoles []org.Permission

	if err := cursor.All(ctx, &orgRoles); err != nil {
		internalError.Err = err
		internalError.Message = "Failed to decode permission objects"
		return nil, internalError
	}

	return orgRoles, nil
}

func NewPermissionRepo(db *mongo.Database, log internal.Logger) (PermissionRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := createUniqueIndex(db, log, ctx, "name", PermissionsCollection); err != nil {
		return nil, err
	}

	return &permissionsRepo{
		db:  db,
		log: log,
	}, nil
}
