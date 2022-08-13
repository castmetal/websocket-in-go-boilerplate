package repositories

import (
	"context"

	_core "websocket-in-go-boilerplate/src/core"
	_user "websocket-in-go-boilerplate/src/domains/user"
)

const _collectionName = "Users"

type userRepository struct {
	db *_core.IDatabase
}

func NewUserRepositoryFromConfig() _user.IUserRepository {
	var db _core.IDatabase

	// TO DO - Get database connection from config

	return newUserRepository(&db)
}

func newUserRepository(db *_core.IDatabase) _user.IUserRepository {
	return &userRepository{db: db}
}

func (repository userRepository) FindOneById(ctx context.Context, id string) (*_user.User, error) {
	var user *_user.User

	// your implementation here

	return user, nil
}

func (repository userRepository) Create(ctx context.Context, user *_user.User) (*_user.User, error) {
	var u *_user.User = user

	// your implementation here

	// TO DO - Domain events

	return u, nil
}
