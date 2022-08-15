package in_memory_tests

import (
	"context"

	_core "websocket-in-go-boilerplate/src/core"
	_user "websocket-in-go-boilerplate/src/domains/user"
)

const _collectionName = "UsersInMemory"

var dbDataUser map[string]_user.User = make(map[string]_user.User)

type userRepository struct {
	db _core.IDatabase
}

func NewUserRepositoryFromConfig() _user.IUserRepository {
	var db _core.IDatabase

	return newUserRepository(db)
}

func newUserRepository(db _core.IDatabase) _user.IUserRepository {
	return &userRepository{db: db}
}

func (repository userRepository) FindOneById(ctx context.Context, id string) (*_user.User, error) {
	var user _user.User = dbDataUser[id]

	userId := string(user.Id[:])

	if userId == "" {
		return nil, _core.NotFoundError("User")
	}

	return &user, nil
}

func (repository userRepository) Create(ctx context.Context, user *_user.User) (*_user.User, error) {
	var u *_user.User

	u, _ = repository.FindOneByEmail(ctx, user.Email)
	if u != nil {
		return nil, _core.AlreadyExistsError("User")
	}

	userId := string(user.Id[:])
	dbDataUser[userId] = *user

	return user, nil
}

func (repository userRepository) FindOneByEmail(ctx context.Context, email string) (*_user.User, error) {
	var user *_user.User

	for _, u := range dbDataUser {
		if u.Email == email {
			user = &u
		}
	}

	if user == nil {
		return nil, _core.NotFoundError("User")
	}

	return user, nil
}
