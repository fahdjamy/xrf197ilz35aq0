package model

import "xrf197ilz35aq0/core/model/user"

type Model interface {
	user.User | user.Settings
}
