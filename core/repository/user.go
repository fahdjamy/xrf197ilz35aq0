package repository

import (
	"context"
	"xrf197ilz35aq0/core/model/user"
	xrfErr "xrf197ilz35aq0/internal/error"
	"xrf197ilz35aq0/storage"
)

const UserCollection = "user"

type UserRepository interface {
	CreateUser(user *user.User, ctx context.Context) (any, error)
}

type userRepo struct {
	store storage.Store
}

var internalError *xrfErr.Internal

func (up *userRepo) CreateUser(user *user.User, ctx context.Context) (any, error) {
	internalError = &xrfErr.Internal{}
	internalError.Source = "core/repository/user#createUser"
	// save user to DB
	objectId, err := up.store.Save(UserCollection, user, ctx)
	if err != nil {
		internalError.Err = err
		internalError.Message = "Saving new user failed"
		return nil, internalError
	}
	return objectId, nil
}

func NewUserRepository(store storage.Store) UserRepository {
	return &userRepo{
		store: store,
	}
}
