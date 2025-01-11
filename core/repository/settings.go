package repository

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"xrf197ilz35aq0/core/model/user"
	"xrf197ilz35aq0/internal"
	xrfErr "xrf197ilz35aq0/internal/error"
)

const SettingsCollection = "settings"

type SettingsRepository interface {
	CreateSettings(settings *user.Settings, ctx context.Context) (any, error)
	FetchUserSettings(ctx context.Context, userFP string) (settings *user.Settings, err error)
}

type settingsRepo struct {
	db  *mongo.Database
	log internal.Logger
}

func (sr *settingsRepo) FetchUserSettings(ctx context.Context, userFP string) (settings *user.Settings, err error) {
	internalErr := &xrfErr.Internal{}
	externalError := &xrfErr.External{}
	internalErr.Source = "core/repository/settings#fetchUserSettings"

	filter := bson.D{{"fingerPrint", userFP}}
	var userSettings user.Settings
	resp := sr.db.Collection(SettingsCollection).FindOne(ctx, filter)

	if resp.Err() != nil {
		if errors.Is(resp.Err(), mongo.ErrNoDocuments) {
			externalError.Message = "Settings for user not found"
			return nil, externalError
		}
		return nil, resp.Err()
	}

	if err := resp.Decode(&userSettings); err != nil {
		internalErr.Err = err
		internalErr.Message = "Failed to decode user settings"
		sr.log.Error(fmt.Sprintf("event=mongoDBFailure :: action=fetchUserSettings :: err=%s", err))
		return nil, internalErr
	}

	return &userSettings, nil
}

func (sr *settingsRepo) CreateSettings(settings *user.Settings, ctx context.Context) (any, error) {
	internalErr := &xrfErr.Internal{} // defined in the repository/user file
	internalErr.Source = "core/repository/user#createSettings"
	document, err := sr.db.Collection(SettingsCollection).InsertOne(ctx, settings)
	if err != nil {
		internalErr.Err = err
		internalErr.Message = "Saving new user-settings failed"
		return nil, err
	}

	return document.InsertedID, nil
}

func NewSettingsRepository(db *mongo.Database, log internal.Logger) SettingsRepository {
	return &settingsRepo{
		db:  db,
		log: log,
	}
}
