package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"xrf197ilz35aq0/core/model/user"
	"xrf197ilz35aq0/internal"
	xrfErr "xrf197ilz35aq0/internal/error"
)

const UserCollection = "user"

type UserRepository interface {
	GetUserById(id int64, ctx context.Context) (*user.User, error)
	CreateUser(user *user.User, ctx context.Context) (any, error)
	UpdatePassword(userFPrint string, newPassword string, ctx context.Context) (bool, error)
}

type userRepo struct {
	db  *mongo.Database
	log internal.Logger
}

var internalError *xrfErr.Internal
var externalError *xrfErr.External

func (up *userRepo) CreateUser(user *user.User, ctx context.Context) (any, error) {
	internalError = &xrfErr.Internal{}
	document, err := up.db.Collection(UserCollection).InsertOne(ctx, user)
	if err != nil {
		up.log.Error(fmt.Sprintf("event=mongoDBFailure :: action=saveUser :: err=%s", err))
		internalError.Err = err
		internalError.Message = "Saving new user failed"
		return nil, err
	}
	up.log.Debug(fmt.Sprintf("event=createUser :: success=true :: objectID=%v", document.InsertedID))

	return document, nil
}

func (up *userRepo) UpdatePassword(userFPrint string, newPassword string, ctx context.Context) (bool, error) {
	internalError = &xrfErr.Internal{}
	internalError.Source = "core/repository/user#updateUser"
	filter := bson.D{{"fingerprint", userFPrint}}
	update := bson.D{{"$set", bson.D{{"password", newPassword}}}}

	resp, err := up.db.Collection(UserCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		internalError.Err = err
		internalError.Message = "Updating user failed"
		return false, internalError
	}

	return resp.ModifiedCount == 1, nil
}

func (up *userRepo) GetUserById(id int64, ctx context.Context) (*user.User, error) {
	externalError = &xrfErr.External{}
	internalError = &xrfErr.Internal{}
	internalError.Source = "core/repository/user#getUserById"

	filter := bson.D{{"_id", id}}

	var userResponse user.User
	resp := up.db.Collection(UserCollection).FindOne(ctx, filter)

	if err := resp.Decode(&userResponse); err != nil {
		internalError.Err = err
		internalError.Message = "Failed to decode userResponse object"
		return nil, internalError
	}
	return &userResponse, nil
}

func NewUserRepository(db *mongo.Database, log internal.Logger) UserRepository {
	return &userRepo{
		db:  db,
		log: log,
	}
}
