package repository

import (
	"context"
	"xrf197ilz35aq0/core/model/user"
	xrfErr "xrf197ilz35aq0/internal/error"
	"xrf197ilz35aq0/storage"
)

const SettingsCollection = "settings"

type SettingsRepository interface {
	CreateSettings(settings *user.Settings, ctx context.Context) (any, error)
}

type settingsRepo struct {
	store storage.Store
}

func (sr *settingsRepo) CreateSettings(settings *user.Settings, ctx context.Context) (any, error) {
	internalError = &xrfErr.Internal{} // defined in the repository/user file
	internalError.Source = "core/repository/user#createSettings"
	// save user to DB
	objectId, err := sr.store.Save(SettingsCollection, settings, ctx)
	if err != nil {
		internalError.Err = err
		internalError.Message = "Saving user settings failed"
		return nil, internalError
	}
	return objectId, nil
}

func NewSettingsRepository(store storage.Store) SettingsRepository {
	return &settingsRepo{
		store: store,
	}
}
