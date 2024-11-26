package repository

import (
	"context"
	"xrf197ilz35aq0/core/model/user"
)

type UserRepository interface {
	CreateUser(user *user.User, ctx context.Context) (any, error)
}
