package repository

import (
	"context"
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
	document, err := sr.db.Collection(SettingsCollection).InsertOne(ctx, settings)
	if err != nil {
		internalError.Err = err
		internalError.Message = "Saving new user-settings failed"
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
