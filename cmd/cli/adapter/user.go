package adapter

import (
	"fmt"
	"xrf197ilz35aq0"
	"xrf197ilz35aq0/core/exchange"
	"xrf197ilz35aq0/core/service/user"
	"xrf197ilz35aq0/storage"
)

type UserAdapter struct {
	store storage.Store
	log   xrf197ilz35aq0.Logger
}

func (ua UserAdapter) CreateUser(userReq *exchange.UserRequest) {
	// create the services
	settingsService := user.NewSettingManager(ua.log)
	userManager := user.NewUserManager(ua.log, settingsService, ua.store)

	userResp, err := userManager.NewUser(userReq)

	if err != nil {
		ua.log.Error(err.Error())
		return
	}
	ua.log.Info(fmt.Sprintf("user created with id '%d'", userResp.UserId))
}

func NewUserAdapter(store storage.Store, log xrf197ilz35aq0.Logger) *UserAdapter {
	return &UserAdapter{
		store: store,
		log:   log,
	}
}
