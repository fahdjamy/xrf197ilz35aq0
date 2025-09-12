package repository

import (
	"context"
	"errors"
	"fmt"
	"time"
	"xrf197ilz35aq0/core/model/org"
	"xrf197ilz35aq0/internal"
	"xrf197ilz35aq0/internal/constants"
	xrfErr "xrf197ilz35aq0/internal/error"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PermissionRepository interface {
	FindPermissionById(id string, ctx context.Context) (*org.Permission, error)
	FindPermissionByName(name string, ctx context.Context) (*org.Permission, error)
	UpdatePermission(permission *org.Permission, ctx context.Context) (bool, error)
	CreatePermission(permission *org.Permission, ctx context.Context) (string, error)
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
	document, err := repo.db.Collection(constants.PermissionsCol).InsertOne(ctx, permission)
	if err != nil {
		// Check for the duplicate key error
		if mongo.IsDuplicateKeyError(err) {
			repo.log.Error(fmt.Sprintf("event=mongoDBFailure :: action=createPermission :: err=duplicateName :: name=%s", permission.Name))
			externalError.Message = "permission name already exists"
			return "", externalError
		}
		repo.log.Error(fmt.Sprintf("event=mongoDBFailure :: action=createPermission :: err=%s", err))
		internalErr.Message = "Creating new permission in mongodb failed"
		internalErr.Err = err
		return "", err
	}
	repo.log.Debug(fmt.Sprintf("event=savePermission :: success=true :: objectID=%v", document.InsertedID))

	return document.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (repo *permissionsRepo) UpdatePermission(permission *org.Permission, ctx context.Context) (bool, error) {
	internalErr := &xrfErr.Internal{}
	externalError := &xrfErr.External{}
	filter := bson.M{constants.PermissionId: permission.Id}
	update := bson.D{{"$set", bson.D{
		{"name", permission.Name},
		{"updatedAt", permission.UpdatedAt},
		{"description", permission.Description},
	}}}

	resp, err := repo.db.Collection(constants.PermissionsCol).UpdateOne(ctx, filter, update)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			repo.log.Error(fmt.Sprintf("event=mongoDBFailure :: action=updatePermission :: err=%s", err))
			externalError.Message = "permission name already exists"
			return false, externalError
		}
		repo.log.Error(fmt.Sprintf("event=mongoDBFailure :: action=updatePermission :: err=%s", err))
		internalErr.Message = "Updating permission in mongodb failed"
		internalErr.Err = err
		return false, err
	}

	if (resp.MatchedCount == 1 && resp.ModifiedCount == 0) ||
		(resp.UpsertedCount == 0 && resp.ModifiedCount == 0) {
		repo.log.Warn(fmt.Sprintf("event=updatePermission :: success=true :: permissionId=%s :: modified=%d :: upSerted=%d matched=%d",
			permission.Id, resp.ModifiedCount, resp.UpsertedCount, resp.MatchedCount))
	}

	// ModifiedCount: The number of documents modified by the operation
	// Upsert edCount: The number of documents upsert ed by the operation
	return resp.UpsertedCount == 1 && resp.ModifiedCount == 1, nil
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
	resp := repo.db.Collection(constants.PermissionsCol).FindOne(ctx, filter)

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
	return repo.findPermissionsByFilter(names, constants.NAME, ctx)
}

func (repo *permissionsRepo) FindPermissionsByIds(ids []string, ctx context.Context) ([]org.Permission, error) {
	return repo.findPermissionsByFilter(ids, constants.PermissionId, ctx)
}

func (repo *permissionsRepo) findPermissionsByFilter(values []string, filterBy string, ctx context.Context) ([]org.Permission, error) {
	if values == nil || len(values) == 0 {
		return []org.Permission{}, nil
	}
	internalError := &xrfErr.Internal{}
	// 1. Build query filter
	filter := bson.M{filterBy: bson.M{"$in": values}}

	// 2. Query mongoDB
	cursor, err := repo.db.Collection(constants.PermissionsCol).Find(ctx, filter)
	if err != nil {
		internalError.Message = fmt.Sprintf("Failed to query permissions by filter: %s", filterBy)
		internalError.Err = err
		return nil, internalError
	}

	defer cursor.Close(ctx)

	// 3. Decode the results into a slice of Permission structs
	var orgPermissions []org.Permission

	if err := cursor.All(ctx, &orgPermissions); err != nil {
		internalError.Err = err
		internalError.Message = "Failed to decode permission objects"
		return nil, internalError
	}

	return orgPermissions, nil
}

func NewPermissionRepo(db *mongo.Database, log internal.Logger) (PermissionRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := createUniqueIndex(db, log, ctx, constants.PermissionsCol, constants.NAME); err != nil {
		return nil, err
	}

	return &permissionsRepo{
		db:  db,
		log: log,
	}, nil
}
