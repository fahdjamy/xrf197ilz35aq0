package repository

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"xrf197ilz35aq0/core/model/user"
	"xrf197ilz35aq0/internal"
	"xrf197ilz35aq0/internal/constants"
	xrfErr "xrf197ilz35aq0/internal/error"
)

const UserCollection = "user"

type UserRepository interface {
	CreateUser(user *user.User, ctx context.Context) (string, error)
	GetUserById(userId string, ctx context.Context) (*user.User, error)
	FindUsersByEmails(emails []string, ctx context.Context) ([]user.User, error)
	UpdatePassword(userFPrint string, newPassword string, ctx context.Context) (bool, error)
	FindUsersByFingerPrints(fingerPrints []string, ctx context.Context) ([]user.User, error)
}

type userRepo struct {
	db  *mongo.Database
	log internal.Logger
}

func (up *userRepo) CreateUser(newUser *user.User, ctx context.Context) (string, error) {
	internalErr := &xrfErr.Internal{}
	if newUser == nil {
		internalErr.Message = "user is nil"
		return "", internalErr
	}
	document, err := up.db.Collection(UserCollection).InsertOne(ctx, newUser)
	if err != nil {
		up.log.Error(fmt.Sprintf("event=mongoDBFailure :: action=saveUser :: err=%s", err))
		internalErr.Err = err
		internalErr.Message = "Saving new user failed"
		return "", err
	}
	up.log.Debug(fmt.Sprintf("event=createUser :: success=true :: objectID=%v", document.InsertedID))

	return newUser.Id, nil
}

func (up *userRepo) UpdatePassword(userFPrint string, newPassword string, ctx context.Context) (bool, error) {
	if newPassword == "" || userFPrint == "" {
		return false, nil
	}
	internalErr := &xrfErr.Internal{}
	internalErr.Source = "core/repository/user#updateUser"
	filter := bson.D{{"fingerprint", userFPrint}}
	update := bson.D{{"$set", bson.D{{constants.PASSWORD, newPassword}}}}

	resp, err := up.db.Collection(UserCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		internalErr.Err = err
		internalErr.Message = "Updating user failed"
		return false, internalErr
	}

	return resp.ModifiedCount == 1, nil
}

func (up *userRepo) GetUserById(userId string, ctx context.Context) (*user.User, error) {
	internalErr := &xrfErr.Internal{}
	externalError := &xrfErr.External{}
	internalErr.Source = "core/repository/user#getUserById"

	filter := bson.D{{constants.USERID, userId}}

	var userResponse user.User
	resp := up.db.Collection(UserCollection).FindOne(ctx, filter)

	if resp.Err() != nil {
		if errors.Is(resp.Err(), mongo.ErrNoDocuments) {
			externalError.Message = "User not found"
			return nil, externalError
		}
		return nil, resp.Err()
	}

	if err := resp.Decode(&userResponse); err != nil {
		internalErr.Err = err
		internalErr.Message = "Failed to decode userResponse object"
		up.log.Error(fmt.Sprintf("event=mongoDBFailure :: action=getUserById :: err=%s", err))
		return nil, internalErr
	}
	return &userResponse, nil
}

func (up *userRepo) FindUsersByEmails(emails []string, ctx context.Context) ([]user.User, error) {
	return up.findUsersByFilter(emails, "email", ctx)
}

func (up *userRepo) FindUsersByFingerPrints(fingerPrints []string, ctx context.Context) ([]user.User, error) {
	return up.findUsersByFilter(fingerPrints, "fingerPrint", ctx)
}

func (up *userRepo) findUsersByFilter(values []string, filterBy string, ctx context.Context) ([]user.User, error) {
	if values == nil || len(values) == 0 {
		return []user.User{}, nil
	}

	internalErr := &xrfErr.Internal{}
	internalErr.Source = "core/repository/user#findUsersByFilter"

	filter := bson.D{{filterBy, bson.M{"$in": values}}}

	var userResponse []user.User
	cursor, err := up.db.Collection(UserCollection).Find(ctx, filter)

	if err != nil {
		internalErr.Err = err
		internalErr.Message = fmt.Sprintf("Error finding users by filterBy: ('%s')", filterBy)
		up.log.Error(fmt.Sprintf("event=mongoDBFailure :: action=findUsersBy ('%s') :: err=%s", filterBy, err))
		return nil, internalErr
	}

	if err := cursor.All(ctx, &userResponse); err != nil {
		internalErr.Err = err
		internalErr.Message = "Failed to decode userResponse object"
		return nil, internalErr
	}

	return userResponse, nil
}

func NewUserRepository(db *mongo.Database, log internal.Logger) UserRepository {
	return &userRepo{
		db:  db,
		log: log,
	}
}
