package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
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
}

type roleRepo struct {
	db  *mongo.Database
	log internal.Logger
}

func (repo *roleRepo) SaveRole(role *org.Role, ctx context.Context) (string, error) {
	internalError = &xrfErr.Internal{}
	document, err := repo.db.Collection(RoleCollection).InsertOne(ctx, role)
	if err != nil {
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
	//TODO implement me
	panic("implement me")
}

func NewRoleRepo(db *mongo.Database, log internal.Logger) RoleRepo {
	return &roleRepo{
		db:  db,
		log: log,
	}
}
