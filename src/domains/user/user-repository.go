package user

import (
	"context"
	_core "websocket-in-go-boilerplate/src/core"
)

type IUserRepository interface {
	_core.IRepository
	FindOneById(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, user *User) (*User, error)
}
