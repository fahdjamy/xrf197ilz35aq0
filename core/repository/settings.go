package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"xrf197ilz35aq0/core/model/user"
	"xrf197ilz35aq0/internal"
	xrfErr "xrf197ilz35aq0/internal/error"
)

const SettingsCollection = "settings"

type SettingsRepository interface {
	CreateSettings(settings *user.Settings, ctx context.Context) (any, error)
}

type settingsRepo struct {
	db  *mongo.Database
	log internal.Logger
}

func (sr *settingsRepo) CreateSettings(settings *user.Settings, ctx context.Context) (any, error) {
	internalError = &xrfErr.Internal{} // defined in the repository/user file
	internalError.Source = "core/repository/user#createSettings"

	sr.log.Debug(fmt.Sprintf("event=createUserSettings :: message=saving new user response"))
	document, err := sr.db.Collection(SettingsCollection).InsertOne(ctx, settings)
	if err != nil {
		internalError.Err = err
		internalError.Message = "Saving new user-settings failed"
		return nil, err
	}
	sr.log.Debug(fmt.Sprintf("event=createUserSettings :: success=true :: objectID=%v", document.InsertedID))

	return document.InsertedID, nil
}

func NewSettingsRepository(db *mongo.Database) SettingsRepository {
	return &settingsRepo{
		db: db,
	}
}
